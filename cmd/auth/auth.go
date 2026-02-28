package auth

import (
	"github.com/spf13/cobra"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

// NewCmd creates the auth command group.
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
		Long:  "Log in, log out, and manage authentication tokens for Moodle instances.",
		Example: `  # Log in to a Moodle instance
  moodle auth login --url https://moodle.example.com

  # Check authentication status
  moodle auth status

  # Print the stored token
  moodle auth token`,
	}

	cmd.AddCommand(newLoginCmd(f))
	cmd.AddCommand(newLogoutCmd(f))
	cmd.AddCommand(newStatusCmd(f))
	cmd.AddCommand(newTokenCmd(f))

	return cmd
}
