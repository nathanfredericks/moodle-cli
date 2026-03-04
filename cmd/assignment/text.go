package assignment

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

func newTextCmd(f *cmdutil.Factory) *cobra.Command {
	var save string
	var stdin bool

	cmd := &cobra.Command{
		Use:   "text <assignment-id>",
		Short: "View or update online text submission",
		Long:  "View the current online text submission, or update it with --save or --stdin.",
		Example: `  # View the current online text submission
  moodle assignment text 101

  # Save text directly
  moodle assignment text 101 --save "My submission text"

  # Pipe text from a file
  cat essay.txt | moodle assignment text 101 --stdin`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			assignID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid assignment ID: %s", args[0])
			}

			// If saving, update the online text
			if save != "" || stdin {
				found, _, err := lookupAssignment(cmd.Context(), client, assignID)
				if err != nil {
					return err
				}
				if found.AllowSubmissionsFromDate > 0 {
					openTime := time.Unix(found.AllowSubmissionsFromDate, 0)
					if time.Now().Before(openTime) {
						return fmt.Errorf("assignment %d is not open for submission until %s", assignID, openTime.Format("2006-01-02 15:04 MST"))
					}
				}

				text := save
				if stdin {
					data, err := io.ReadAll(os.Stdin)
					if err != nil {
						return fmt.Errorf("failed to read stdin: %w", err)
					}
					text = string(data)
				}

				params := map[string]any{
					"assignmentid": assignID,
					"plugindata": map[string]any{
						"onlinetext_editor": map[string]any{
							"text":   text,
							"format": 1,
							"itemid": 0,
						},
					},
				}

				var result any
				if err := client.Call(cmd.Context(), "mod_assign_save_submission", params, &result); err != nil {
					return fmt.Errorf("failed to save online text: %w", err)
				}

				fmt.Fprintf(f.IO.Out, "Online text saved for assignment %d.\n", assignID)
				return nil
			}

			// Otherwise, display the current online text
			var status submissionStatusResponse
			params := map[string]any{
				"assignid": assignID,
			}
			if err := client.Call(cmd.Context(), "mod_assign_get_submission_status", params, &status); err != nil {
				return fmt.Errorf("failed to get submission status: %w", err)
			}

			if status.LastAttempt == nil || status.LastAttempt.Submission == nil {
				return fmt.Errorf("no submission found for assignment %d", assignID)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			for _, p := range status.LastAttempt.Submission.Plugins {
				if p.Type == "onlinetext" {
					for _, ef := range p.EditorFields {
						if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
							return f.Output(map[string]any{
								"assignmentid": assignID,
								"text":         ef.Text,
								"format":       ef.Format,
							}, opts)
						}
						fmt.Fprintln(f.IO.Out, ef.Text)
						return nil
					}
				}
			}

			fmt.Fprintln(f.IO.Out, "No online text submission found for this assignment.")
			return nil
		},
	}

	cmd.Flags().StringVar(&save, "save", "", "Save this text as the online text submission")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Read online text from stdin")

	return cmd
}
