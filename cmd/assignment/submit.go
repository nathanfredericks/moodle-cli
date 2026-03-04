package assignment

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

func newSubmitCmd(f *cmdutil.Factory) *cobra.Command {
	var acceptStatement bool

	cmd := &cobra.Command{
		Use:   "submit <assignment-id>",
		Short: "Submit assignment for grading",
		Long:  "Submit an assignment for grading. The assignment must have content saved first.",
		Example: `  # Submit an assignment for grading
  moodle assignment submit 101

  # Submit and accept the submission statement
  moodle assignment submit 101 --accept-statement`,
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

			params := map[string]any{
				"assignmentid":               assignID,
				"acceptsubmissionstatement":  acceptStatement,
			}

			var result any
			if err := client.Call(cmd.Context(), "mod_assign_submit_for_grading", params, &result); err != nil {
				return fmt.Errorf("failed to submit assignment: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(result, opts)
			}

			fmt.Fprintf(f.IO.Out, "Assignment %d submitted successfully.\n", assignID)
			return nil
		},
	}

	cmd.Flags().BoolVar(&acceptStatement, "accept-statement", false, "Accept the submission statement")
	return cmd
}
