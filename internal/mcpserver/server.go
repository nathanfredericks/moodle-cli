// Package mcpserver exposes every moodle-cli command as a Model Context
// Protocol tool over the Streamable HTTP transport, so an AI agent (e.g.
// n8n's MCP Client Tool node) can drive moodle-cli with the credentials
// configured in the process environment (MOODLE_URL / MOODLE_TOKEN).
package mcpserver

import (
	"bytes"
	"context"
	"crypto/subtle"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/root"
)

// Options configures the MCP server.
type Options struct {
	// Addr is the address to listen on, e.g. ":8080".
	Addr string
	// APIKey is required as a Bearer token on every request to /mcp.
	APIKey string
	// Version is reported to MCP clients and by the "version" tool.
	Version string
	// DisableLocalhostProtection permits an external Host header when the
	// server is reached through a trusted loopback reverse proxy.
	DisableLocalhostProtection bool
}

// Serve starts the MCP server and blocks until it exits.
func Serve(opts Options) error {
	return http.ListenAndServe(opts.Addr, NewHandler(opts))
}

// NewHandler constructs the authenticated HTTP handler used by both the local
// MCP server and the AWS Lambda Web Adapter runtime.
func NewHandler(opts Options) http.Handler {
	mcpServer := server.NewMCPServer("moodle-cli", opts.Version, server.WithToolCapabilities(false))

	for _, leaf := range collectLeaves(root.Root()) {
		mcpServer.AddTool(buildTool(leaf), makeHandler(opts.Version, leaf.path))
	}

	httpServer := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithEndpointPath("/mcp"),
		server.WithDisableLocalhostProtection(opts.DisableLocalhostProtection),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/mcp", requireBearer(opts.APIKey, httpServer))

	return mux
}

func requireBearer(apiKey string, next http.Handler) http.Handler {
	want := "Bearer " + apiKey
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("Authorization")
		if subtle.ConstantTimeCompare([]byte(got), []byte(want)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// leaf is a runnable, argument-taking command discovered by walking the
// command tree, along with the path used to reach it (e.g. ["course", "get"]).
type leaf struct {
	path []string
	cmd  *cobra.Command
}

// collectLeaves walks the full command tree and returns every runnable leaf
// command — i.e. every actual moodle-cli command, excluding group commands
// like "course" or "auth" that only hold subcommands.
func collectLeaves(cmd *cobra.Command) []leaf {
	var out []leaf
	var walk func(c *cobra.Command, path []string)
	walk = func(c *cobra.Command, path []string) {
		if c.Runnable() && len(c.Commands()) == 0 {
			out = append(out, leaf{path: path, cmd: c})
			return
		}
		for _, child := range c.Commands() {
			walk(child, append(append([]string{}, path...), child.Name()))
		}
	}
	walk(cmd, nil)
	return out
}

// collectFlags returns every flag available on cmd, including those inherited
// from persistent flags on its ancestors (e.g. --format, --no-color, --verbose).
func collectFlags(cmd *cobra.Command) []*pflag.Flag {
	seen := map[string]bool{}
	var flags []*pflag.Flag
	add := func(f *pflag.Flag) {
		if !seen[f.Name] {
			seen[f.Name] = true
			flags = append(flags, f)
		}
	}
	cmd.Flags().VisitAll(add)
	cmd.InheritedFlags().VisitAll(add)
	sort.Slice(flags, func(i, j int) bool { return flags[i].Name < flags[j].Name })
	return flags
}

func buildTool(l leaf) mcp.Tool {
	name := l.path[0]
	for _, p := range l.path[1:] {
		name += "_" + p
	}

	desc := l.cmd.Short
	if l.cmd.Long != "" {
		desc = l.cmd.Long
	}
	desc += "\n\nUsage: " + l.cmd.UseLine()
	if l.cmd.Example != "" {
		desc += "\n\n" + l.cmd.Example
	}

	opts := []mcp.ToolOption{
		mcp.WithDescription(desc),
		mcp.WithArray("args",
			mcp.WithStringItems(),
			mcp.Description(fmt.Sprintf(
				"Positional arguments, in CLI order, matching: %s", l.cmd.UseLine(),
			)),
		),
	}

	for _, fl := range collectFlags(l.cmd) {
		propOpts := []mcp.PropertyOption{mcp.Description(fl.Usage)}
		switch fl.Value.Type() {
		case "bool":
			opts = append(opts, mcp.WithBoolean(fl.Name, propOpts...))
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
			opts = append(opts, mcp.WithInteger(fl.Name, propOpts...))
		default:
			opts = append(opts, mcp.WithString(fl.Name, propOpts...))
		}
	}

	return mcp.NewTool(name, opts...)
}

// makeHandler returns an MCP tool handler that re-runs moodle-cli's own
// command tree with a fresh Factory per call — the same isolation a real
// `moodle` process invocation gets — translating MCP tool arguments back
// into the equivalent CLI flags/positional args before executing.
func makeHandler(version string, path []string) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		f, err := cmdutil.NewFactory()
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to initialize moodle-cli", err), nil
		}

		buf := &bytes.Buffer{}
		f.IO.Out = buf
		f.IO.ErrOut = buf

		argv := append([]string{}, path...)

		args := req.GetArguments()
		var flagNames []string
		for k := range args {
			if k != "args" {
				flagNames = append(flagNames, k)
			}
		}
		sort.Strings(flagNames)
		for _, name := range flagNames {
			argv = append(argv, flagArgs(name, args[name])...)
		}
		if posRaw, ok := args["args"].([]any); ok {
			for _, a := range posRaw {
				argv = append(argv, fmt.Sprintf("%v", a))
			}
		}

		cmd := root.New(f, version)
		cmd.SetArgs(argv)

		execErr := cmd.ExecuteContext(ctx)

		output := buf.String()
		if execErr != nil {
			if output != "" {
				output += "\n"
			}
			output += execErr.Error()
			return mcp.NewToolResultError(output), nil
		}
		if output == "" {
			output = "(no output)"
		}
		return mcp.NewToolResultText(output), nil
	}
}

// flagArgs renders a single MCP argument as CLI flag tokens, e.g.
// ("course", float64(42)) -> ["--course", "42"], ("verbose", true) -> ["--verbose"].
func flagArgs(name string, v any) []string {
	flag := "--" + name
	switch val := v.(type) {
	case bool:
		if val {
			return []string{flag}
		}
		return []string{flag + "=false"}
	case float64:
		if val == math.Trunc(val) {
			return []string{flag, strconv.FormatInt(int64(val), 10)}
		}
		return []string{flag, strconv.FormatFloat(val, 'f', -1, 64)}
	case string:
		return []string{flag, val}
	default:
		return []string{flag, fmt.Sprintf("%v", val)}
	}
}
