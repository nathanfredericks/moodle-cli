package course

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type courseDetail struct {
	ID               int              `json:"id"`
	Shortname        string           `json:"shortname"`
	Fullname         string           `json:"fullname"`
	DisplayName      string           `json:"displayname"`
	Summary          string           `json:"summary"`
	CategoryID       int              `json:"categoryid"`
	CategoryName     string           `json:"categoryname"`
	Format           string           `json:"format"`
	StartDate        int64            `json:"startdate"`
	EndDate          int64            `json:"enddate"`
	Visible          int              `json:"visible"`
	EnableCompletion int              `json:"enablecompletion"`
	Contacts         []courseContact  `json:"contacts"`
}

type courseContact struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
}

type coursesByFieldResponse struct {
	Courses []courseDetail `json:"courses"`
}

func newGetCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <course-id>",
		Short: "Get course details",
		Long:  "Get detailed information about a specific course by ID.",
		Example: `  # Get course details
  moodle course get 42

  # Get course details as JSON
  moodle course get 42 -f json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			courseID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}

			var result coursesByFieldResponse
			params := map[string]any{
				"field": "id",
				"value": strconv.Itoa(courseID),
			}
			err = client.Call(cmd.Context(), "core_course_get_courses_by_field", params, &result)
			if err != nil {
				return fmt.Errorf("failed to get course: %w", err)
			}

			if len(result.Courses) == 0 {
				return fmt.Errorf("course not found: %d", courseID)
			}

			c := result.Courses[0]

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(c, opts)
			}

			visible := "No"
			if c.Visible == 1 {
				visible = "Yes"
			}
			completion := "Disabled"
			if c.EnableCompletion == 1 {
				completion = "Enabled"
			}

			contacts := "--"
			if len(c.Contacts) > 0 {
				contacts = ""
				for i, ct := range c.Contacts {
					if i > 0 {
						contacts += ", "
					}
					contacts += ct.Fullname
				}
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "Field"},
					{Name: "Value"},
				},
				Rows: []map[string]string{
					{"Field": "ID", "Value": strconv.Itoa(c.ID)},
					{"Field": "Shortname", "Value": c.Shortname},
					{"Field": "Fullname", "Value": c.Fullname},
					{"Field": "Category", "Value": c.CategoryName},
					{"Field": "Format", "Value": c.Format},
					{"Field": "Visible", "Value": visible},
					{"Field": "Completion", "Value": completion},
					{"Field": "Contacts", "Value": contacts},
				},
			}

			return f.Output(table, opts)
		},
	}

	return cmd
}
