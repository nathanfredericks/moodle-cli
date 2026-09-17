package mcp

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/mcpserver"
)

// NewCmd creates the mcp command group.
func NewCmd(f *cmdutil.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Run moodle-cli as a Model Context Protocol server",
		Long:  "Expose every moodle-cli command as an MCP tool, so an AI agent (e.g. n8n's MCP Client Tool node) can call moodle-cli using the credentials configured in this environment.",
	}

	cmd.AddCommand(newServeCmd(version))

	return cmd
}

func newServeCmd(version string) *cobra.Command {
	var addr string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server",
		Long: `Start an HTTP server speaking the Model Context Protocol (Streamable HTTP
transport) at /mcp, exposing one MCP tool per moodle-cli command.

Requires MOODLE_URL and MOODLE_TOKEN to be set in the environment — see
"moodle auth login" and "moodle auth token" to obtain a token. Requires
MCP_API_KEY to be set; every request to /mcp must send it as a Bearer token
in the Authorization header.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := os.Getenv("MCP_API_KEY")
			if apiKey == "" {
				return fmt.Errorf("MCP_API_KEY environment variable is required")
			}
			if os.Getenv("MOODLE_URL") == "" {
				return fmt.Errorf("MOODLE_URL environment variable is required")
			}
			if os.Getenv("MOODLE_TOKEN") == "" {
				return fmt.Errorf("MOODLE_TOKEN environment variable is required")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "moodle-cli MCP server listening on %s\n", addr)
			return mcpserver.Serve(mcpserver.Options{
				Addr:    addr,
				APIKey:  apiKey,
				Version: version,
			})
		},
	}

	cmd.Flags().StringVar(&addr, "addr", ":8080", "Address to listen on")

	return cmd
}
