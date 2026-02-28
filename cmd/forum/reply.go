package forum

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type addPostResponse struct {
	PostID int `json:"postid"`
}

func newReplyCmd(f *cmdutil.Factory) *cobra.Command {
	var message string
	var subject string

	cmd := &cobra.Command{
		Use:   "reply <post-id>",
		Short: "Reply to a post",
		Long:  "Reply to an existing post in a discussion.",
		Example: `  # Reply to a post
  moodle forum reply 200 --message "Thanks, that helped!"

  # Reply with a custom subject
  moodle forum reply 200 --subject "Re: Lab 3" --message "Here is my approach..."`,
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

			var result addPostResponse
			if subject == "" {
				subject = "Re:"
			}
			params := map[string]any{
				"postid":  postID,
				"subject": subject,
				"message": message,
			}
			if err := client.Call(cmd.Context(), "mod_forum_add_discussion_post", params, &result); err != nil {
				return fmt.Errorf("failed to reply: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(result, opts)
			}

			fmt.Fprintf(f.IO.Out, "Reply posted (Post ID: %d)\n", result.PostID)
			return nil
		},
	}

	cmd.Flags().StringVar(&message, "message", "", "Reply message (required)")
	cmd.Flags().StringVar(&subject, "subject", "", "Reply subject (optional, defaults to original subject)")
	_ = cmd.MarkFlagRequired("message")

	return cmd
}
