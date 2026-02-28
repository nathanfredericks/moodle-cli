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
		DisableAutoGenTag: true,
	}

	cmd.PersistentFlags().StringP("format", "f", "", "Output format: table, json, csv, yaml, plain")
	cmd.PersistentFlags().Bool("no-color", false, "Disable color output")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	// Version command
	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "moodle version dev")
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
