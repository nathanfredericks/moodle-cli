// Package root provides the complete command tree for documentation generation.
// It constructs the CLI hierarchy with a stub factory so no real I/O or config
// is needed, making it safe for use in doc generation tools.
package root

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"

	assignmentcmd "github.com/nathanfredericks/moodle-cli/cmd/assignment"
	authcmd "github.com/nathanfredericks/moodle-cli/cmd/auth"
	configcmd "github.com/nathanfredericks/moodle-cli/cmd/config"
	coursecmd "github.com/nathanfredericks/moodle-cli/cmd/course"
	forumcmd "github.com/nathanfredericks/moodle-cli/cmd/forum"
	usercmd "github.com/nathanfredericks/moodle-cli/cmd/user"
)

// Root returns the full command tree with a stub factory.
// This is intended for documentation generation and should not be used
// to actually execute commands.
func Root() *cobra.Command {
	f := &cmdutil.Factory{
		IO: cmdutil.IOStreams{
			In:     io.NopCloser(nil),
			Out:    io.Discard,
			ErrOut: io.Discard,
		},
		Output: func(data any, opts output.FormatOptions) error {
			return nil
		},
	}

	return New(f, "dev")
}

// New returns the full command tree wired to the given factory and version.
// Unlike Root, this tree is safe to execute: it uses the supplied factory for
// real I/O and API access, so it's the shape used both by the `moodle` binary
// itself and by the MCP server, which builds one fresh tree per tool call.
func New(f *cmdutil.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "moodle",
		Short: "CLI for the Moodle LMS",
		Long:  "A command-line interface for the Moodle Learning Management System API.",
		Example: `  # List your enrolled courses
  moodle course list

  # Get details about a specific course
  moodle course get 42

  # Search for courses by name
  moodle course search "Introduction to Computing"

  # View your assignments
  moodle assignment list --course 42

  # Output as JSON for scripting
  moodle course list -f json`,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	var formatStr string
	var noColor bool
	var verbose bool

	cmd.PersistentFlags().StringVarP(&formatStr, "format", "f", "", "Output format: table, json, csv, yaml, plain")
	cmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		f.NoColor = noColor || output.NoColorEnabled()
		f.Verbose = verbose
	}

	// Version command
	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(f.IO.Out, "moodle version %s\n", version)
		},
	})

	// Register all subcommand groups
	cmd.AddCommand(authcmd.NewCmd(f))
	cmd.AddCommand(configcmd.NewCmd(f))
	cmd.AddCommand(coursecmd.NewCmd(f))
	cmd.AddCommand(usercmd.NewCmd(f))
	cmd.AddCommand(assignmentcmd.NewCmd(f))
	cmd.AddCommand(forumcmd.NewCmd(f))

	return cmd
}
