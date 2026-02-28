package assignment

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type assignmentListResponse struct {
	Courses []assignmentCourse `json:"courses"`
}

type assignmentCourse struct {
	ID          int              `json:"id"`
	FullName    string           `json:"fullname"`
	Assignments []assignmentItem `json:"assignments"`
}

type assignmentItem struct {
	ID                       int         `json:"id"`
	CMID                     int         `json:"cmid"`
	Name                     string      `json:"name"`
	DueDate                  int64       `json:"duedate"`
	AllowSubmissionsFromDate int64       `json:"allowsubmissionsfromdate"`
	Grade                    int         `json:"grade"`
	Intro                    string      `json:"intro"`
	Course                   int         `json:"course"`
	IntroAttachments         []fileEntry `json:"introattachments"`
}

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	var courseID int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List assignments",
		Long:  "List assignments across enrolled courses.",
		Example: `  # List all assignments
  moodle assignment list

  # List assignments for a specific course
  moodle assignment list --course 42

  # Output as JSON
  moodle assignment list -f json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			params := map[string]any{}
			if courseID > 0 {
				params["courseids"] = []int{courseID}
			}

			var result assignmentListResponse
			if err := client.Call(cmd.Context(), "mod_assign_get_assignments", params, &result); err != nil {
				return fmt.Errorf("failed to list assignments: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(&result, opts)
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "ID"},
					{Name: "Course"},
					{Name: "Name"},
					{Name: "Due Date"},
				},
				Rows: []map[string]string{},
			}

			for _, c := range result.Courses {
				for _, a := range c.Assignments {
					dueStr := "No due date"
					if a.DueDate > 0 {
						t := time.Unix(a.DueDate, 0)
						dueStr = t.Format("2006-01-02 15:04")
						if time.Now().After(t) {
							dueStr += " (overdue)"
						}
					}
					table.Rows = append(table.Rows, map[string]string{
						"ID":       strconv.Itoa(a.ID),
						"Course":   c.FullName,
						"Name":     a.Name,
						"Due Date": dueStr,
					})
				}
			}

			if len(table.Rows) == 0 {
				fmt.Fprintln(f.IO.Out, "No assignments found.")
				return nil
			}

			return f.Output(table, opts)
		},
	}

	cmd.Flags().IntVar(&courseID, "course", 0, "Filter by course ID")
	return cmd
}
