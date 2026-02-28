package assignment

import (
	"github.com/spf13/cobra"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

// NewCmd creates the assignment command group.
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assignment",
		Short: "Manage assignments",
		Long:  "List, view, submit, and manage assignments.",
		Example: `  # List all assignments
  moodle assignment list

  # List assignments for a specific course
  moodle assignment list --course 42

  # Get assignment details
  moodle assignment get 101

  # Upload a file and submit
  moodle assignment upload 101 report.pdf
  moodle assignment submit 101 --accept-statement`,
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newDueCmd(f))
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newStatusCmd(f))
	cmd.AddCommand(newSubmitCmd(f))
	cmd.AddCommand(newUploadCmd(f))
	cmd.AddCommand(newDownloadCmd(f))
	cmd.AddCommand(newTextCmd(f))

	return cmd
}
