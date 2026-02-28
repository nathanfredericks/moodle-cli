package forum

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type addDiscussionResponse struct {
	DiscussionID int `json:"discussionid"`
}

func newPostCmd(f *cmdutil.Factory) *cobra.Command {
	var subject string
	var message string

	cmd := &cobra.Command{
		Use:   "post <forum-id>",
		Short: "Create a new discussion",
		Long:  "Create a new discussion thread in a forum.",
		Example: `  # Create a new discussion
  moodle forum post 5 --subject "Help with Lab 3" --message "I'm stuck on step 2..."`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			forumID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid forum ID: %s", args[0])
			}

			client, err := f.Client()
			if err != nil {
				return err
			}

			var result addDiscussionResponse
			params := map[string]any{
				"forumid": forumID,
				"subject": subject,
				"message": message,
			}
			if err := client.Call(cmd.Context(), "mod_forum_add_discussion", params, &result); err != nil {
				return fmt.Errorf("failed to create discussion: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(result, opts)
			}

			fmt.Fprintf(f.IO.Out, "Discussion created (ID: %d)\n", result.DiscussionID)
			return nil
		},
	}

	cmd.Flags().StringVar(&subject, "subject", "", "Discussion subject (required)")
	cmd.Flags().StringVar(&message, "message", "", "Discussion message (required)")
	_ = cmd.MarkFlagRequired("subject")
	_ = cmd.MarkFlagRequired("message")

	return cmd
}
