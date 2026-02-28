package user

import (
	"github.com/spf13/cobra"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

// NewCmd creates the user command group.
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
		Long:  "List, view, and manage Moodle users.",
		Example: `  # Show current user info
  moodle user whoami

  # List users in a course
  moodle user list --course 42

  # Get user details by ID
  moodle user get 7`,
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newWhoamiCmd(f))

	return cmd
}
