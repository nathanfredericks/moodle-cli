package forum

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type forumItem struct {
	ID               int    `json:"id"`
	Course           int    `json:"course"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	NumDiscussions   int    `json:"numdiscussions"`
	UnreadPostsCount int    `json:"unreadpostscount"`
}

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	var courseID int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List forums in a course",
		Long:  "List all forums in a course.",
		Example: `  # List forums in course 42
  moodle forum list --course 42

  # Output as JSON
  moodle forum list --course 42 -f json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			var forums []forumItem
			params := map[string]any{
				"courseids": []int{courseID},
			}
			if err := client.Call(cmd.Context(), "mod_forum_get_forums_by_courses", params, &forums); err != nil {
				return fmt.Errorf("failed to list forums: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(forums, opts)
			}

			if len(forums) == 0 {
				fmt.Fprintln(f.IO.Out, "No forums found.")
				return nil
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "ID", Width: 6},
					{Name: "Name", Width: 30},
					{Name: "Type", Width: 10},
					{Name: "Discussions", Width: 12},
					{Name: "Unread", Width: 8},
				},
				Rows: make([]map[string]string, 0, len(forums)),
			}

			for _, forum := range forums {
				table.Rows = append(table.Rows, map[string]string{
					"ID":          strconv.Itoa(forum.ID),
					"Name":        forum.Name,
					"Type":        forum.Type,
					"Discussions": strconv.Itoa(forum.NumDiscussions),
					"Unread":      strconv.Itoa(forum.UnreadPostsCount),
				})
			}

			return f.Output(table, opts)
		},
	}

	cmd.Flags().IntVar(&courseID, "course", 0, "Course ID (required)")
	_ = cmd.MarkFlagRequired("course")

	return cmd
}
