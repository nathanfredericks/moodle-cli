package course

import (
	"github.com/spf13/cobra"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

// NewCmd creates the course command group.
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "course",
		Short: "Manage courses",
		Long:  "List, view, and manage Moodle courses.",
		Example: `  # List enrolled courses
  moodle course list

  # Get course details
  moodle course get 42

  # View course contents
  moodle course content 42

  # Search for a course
  moodle course search "Biology 101"`,
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newContentCmd(f))
	cmd.AddCommand(newSearchCmd(f))
	cmd.AddCommand(newDownloadCmd(f))
	cmd.AddCommand(newGradesCmd(f))
	cmd.AddCommand(newModuleCmd(f))

	return cmd
}
