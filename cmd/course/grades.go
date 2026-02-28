package course

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type gradeEntry struct {
	CourseID int    `json:"courseid"`
	Grade    string `json:"grade"`
	RawGrade string `json:"rawgrade"`
	Rank     int    `json:"rank"`
}

type gradesResponse struct {
	Grades []gradeEntry `json:"grades"`
}

func newGradesCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grades",
		Short: "View your course grades",
		Long:  "Display an overview of grades across all enrolled courses.",
		Example: `  # View grades across all courses
  moodle course grades

  # Output grades as JSON
  moodle course grades -f json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			var result gradesResponse
			err = client.Call(cmd.Context(), "gradereport_overview_get_course_grades", nil, &result)
			if err != nil {
				return fmt.Errorf("failed to get grades: %w", err)
			}

			// Build course ID → name map
			var timeline timelineResponse
			params := map[string]any{"classification": "all"}
			err = client.Call(cmd.Context(), "core_course_get_enrolled_courses_by_timeline_classification", params, &timeline)
			if err != nil {
				return fmt.Errorf("failed to get course names: %w", err)
			}

			courseNames := make(map[int]string, len(timeline.Courses))
			for _, c := range timeline.Courses {
				courseNames[c.ID] = c.Fullname
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(result.Grades, opts)
			}

			if len(result.Grades) == 0 {
				fmt.Fprintln(f.IO.Out, "No grades found.")
				return nil
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "ID", Width: 6},
					{Name: "Course", Width: 40},
					{Name: "Grade", Width: 10},
					{Name: "Rank", Width: 6},
				},
				Rows: make([]map[string]string, 0, len(result.Grades)),
			}

			for _, g := range result.Grades {
				name := courseNames[g.CourseID]
				if name == "" {
					name = fmt.Sprintf("Course %d", g.CourseID)
				}
				rank := "--"
				if g.Rank > 0 {
					rank = strconv.Itoa(g.Rank)
				}
				table.Rows = append(table.Rows, map[string]string{
					"ID":     strconv.Itoa(g.CourseID),
					"Course": name,
					"Grade":  g.Grade,
					"Rank":   rank,
				})
			}

			return f.Output(table, opts)
		},
	}

	return cmd
}
