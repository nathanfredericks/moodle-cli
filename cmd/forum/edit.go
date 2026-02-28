package forum

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

func newEditCmd(f *cmdutil.Factory) *cobra.Command {
	var subject string
	var message string

	cmd := &cobra.Command{
		Use:   "edit <post-id>",
		Short: "Edit a post",
		Long:  "Edit the subject and/or message of an existing post.",
		Example: `  # Edit a post's message
  moodle forum edit 200 --message "Updated content"

  # Edit both subject and message
  moodle forum edit 200 --subject "New Title" --message "New content"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			postID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid post ID: %s", args[0])
			}

			if subject == "" && message == "" {
				return fmt.Errorf("provide --subject and/or --message to update")
			}

			client, err := f.Client()
			if err != nil {
				return err
			}

			params := map[string]any{
				"postid": postID,
			}
			if subject != "" {
				params["subject"] = subject
			}
			if message != "" {
				params["message"] = message
			}

			var result any
			if err := client.Call(cmd.Context(), "mod_forum_update_discussion_post", params, &result); err != nil {
				return fmt.Errorf("failed to edit post: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "Post %d updated.\n", postID)
			return nil
		},
	}

	cmd.Flags().StringVar(&subject, "subject", "", "New subject")
	cmd.Flags().StringVar(&message, "message", "", "New message")

	return cmd
}
