package forum

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

func newDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <post-id>",
		Short: "Delete a post",
		Long:  "Delete a post by its ID. Deleting the first post in a discussion deletes the entire discussion.",
		Example: `  # Delete a post
  moodle forum delete 200`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			postID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid post ID: %s", args[0])
			}

			client, err := f.Client()
			if err != nil {
				return err
			}

			var result any
			params := map[string]any{
				"postid": postID,
			}
			if err := client.Call(cmd.Context(), "mod_forum_delete_post", params, &result); err != nil {
				return fmt.Errorf("failed to delete post: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "Post %d deleted.\n", postID)
			return nil
		},
	}

	return cmd
}
