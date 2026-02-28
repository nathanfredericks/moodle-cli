package assignment

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
	"github.com/nathanfredericks/moodle-cli/internal/text"
)

func newGetCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <assignment-id>",
		Short: "Get assignment details",
		Long:  "Get detailed information about a specific assignment including grade and feedback.",
		Example: `  # Get assignment details
  moodle assignment get 101

  # Get assignment details as JSON
  moodle assignment get 101 -f json

  # Download any listed resources
  moodle assignment download 101 --resources`,
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

			// Fetch assignment details from the assignments list
			var listResult assignmentListResponse
			if err := client.Call(cmd.Context(), "mod_assign_get_assignments", map[string]any{}, &listResult); err != nil {
				return fmt.Errorf("failed to get assignments: %w", err)
			}

			var found *assignmentItem
			var courseName string
			for _, c := range listResult.Courses {
				for i, a := range c.Assignments {
					if a.ID == assignID {
						found = &c.Assignments[i]
						courseName = c.FullName
						break
					}
				}
				if found != nil {
					break
				}
			}

			if found == nil {
				return fmt.Errorf("assignment %d not found", assignID)
			}

			// Fetch submission status
			var status submissionStatusResponse
			params := map[string]any{
				"assignid": assignID,
			}
			if err := client.Call(cmd.Context(), "mod_assign_get_submission_status", params, &status); err != nil {
				return fmt.Errorf("failed to get submission status: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				combined := map[string]any{
					"assignment": found,
					"course":     courseName,
					"status":     status,
				}
				return f.Output(combined, opts)
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "Field"},
					{Name: "Value"},
				},
				Rows: []map[string]string{
					{"Field": "ID", "Value": strconv.Itoa(found.ID)},
					{"Field": "Name", "Value": found.Name},
					{"Field": "Course", "Value": courseName},
				},
			}

			if found.DueDate > 0 {
				t := time.Unix(found.DueDate, 0)
				dueStr := t.Format("2006-01-02 15:04 MST")
				if time.Now().After(t) {
					dueStr += " (overdue)"
				}
				table.Rows = append(table.Rows, map[string]string{"Field": "Due Date", "Value": dueStr})
			}

			table.Rows = append(table.Rows, map[string]string{"Field": "Max Grade", "Value": strconv.Itoa(found.Grade)})

			// Show intro attachments (instructor resources)
			for _, file := range found.IntroAttachments {
				table.Rows = append(table.Rows, map[string]string{
					"Field": "Resource",
					"Value": fmt.Sprintf("%s (%s)", file.FileName, text.FormatFileSize(file.FileSize)),
				})
			}

			if status.LastAttempt != nil && status.LastAttempt.Submission != nil {
				sub := status.LastAttempt.Submission
				table.Rows = append(table.Rows,
					map[string]string{"Field": "Submission Status", "Value": sub.Status},
				)
				if sub.TimeModified > 0 {
					table.Rows = append(table.Rows, map[string]string{"Field": "Last Modified", "Value": time.Unix(sub.TimeModified, 0).Format("2006-01-02 15:04 MST")})
				}

				// Show submitted files
				for _, p := range sub.Plugins {
					if p.Type == "file" {
						for _, fa := range p.FileAreas {
							for _, file := range fa.Files {
								table.Rows = append(table.Rows, map[string]string{
									"Field": "File",
									"Value": fmt.Sprintf("%s (%s)", file.FileName, text.FormatFileSize(file.FileSize)),
								})
							}
						}
					}
					if p.Type == "onlinetext" {
						for _, ef := range p.EditorFields {
							if ef.Text != "" {
								table.Rows = append(table.Rows, map[string]string{
									"Field": "Online Text",
									"Value": text.Truncate(text.StripHTML(ef.Text), 80),
								})
							}
						}
					}
				}
			}

			if status.LastAttempt != nil {
				table.Rows = append(table.Rows,
					map[string]string{"Field": "Can Submit", "Value": fmt.Sprintf("%v", status.LastAttempt.CanSubmit)},
				)
			}

			// Show grade and feedback
			if status.Feedback != nil {
				if status.Feedback.GradeForDisplay != "" {
					table.Rows = append(table.Rows, map[string]string{"Field": "Grade", "Value": text.StripHTML(status.Feedback.GradeForDisplay)})
				}
				if status.Feedback.GradedDate > 0 {
					table.Rows = append(table.Rows, map[string]string{"Field": "Graded Date", "Value": time.Unix(status.Feedback.GradedDate, 0).Format("2006-01-02 15:04 MST")})
				}
				for _, p := range status.Feedback.Plugins {
					if p.Type == "comments" {
						for _, ef := range p.EditorFields {
							feedback := strings.TrimSpace(text.StripHTML(ef.Text))
							if feedback != "" {
								table.Rows = append(table.Rows, map[string]string{"Field": "Feedback", "Value": feedback})
							}
						}
					}
				}
			}

			return f.Output(table, opts)
		},
	}

	return cmd
}
