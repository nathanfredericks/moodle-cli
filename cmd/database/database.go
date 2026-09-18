package database

import (
	"github.com/spf13/cobra"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

// NewCmd creates the database command group.
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "database",
		Short: "Manage database activities",
		Long:  "List database activities, view fields, and browse entries.",
		Example: `  # List database activities in a course
  moodle database list --course 42

  # View fields defined in a database
  moodle database fields 235

  # Browse entries in a database
  moodle database entries 235`,
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newFieldsCmd(f))
	cmd.AddCommand(newEntriesCmd(f))

	return cmd
}
