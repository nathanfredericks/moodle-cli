package forum

import (
	"github.com/spf13/cobra"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

// NewCmd creates the forum command group.
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "forum",
		Short: "Manage forums",
		Long:  "List forums, read discussions, and post replies.",
		Example: `  # List forums in a course
  moodle forum list --course 42

  # List discussions in a forum
  moodle forum discussions 5

  # Read a discussion thread
  moodle forum read 100

  # Create a new discussion
  moodle forum post 5 --subject "Question" --message "Hello"`,
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newDiscussionsCmd(f))
	cmd.AddCommand(newReadCmd(f))
	cmd.AddCommand(newPostCmd(f))
	cmd.AddCommand(newReplyCmd(f))
	cmd.AddCommand(newEditCmd(f))
	cmd.AddCommand(newDeleteCmd(f))

	return cmd
}
