package course

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type siteInfoResponse struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
}

type enrolledCourse struct {
	ID           int     `json:"id"`
	Shortname    string  `json:"shortname"`
	Fullname     string  `json:"fullname"`
	CategoryName string  `json:"coursecategory"`
	HasProgress  bool    `json:"hasprogress"`
	Progress     float64 `json:"progress"`
	IsFavourite  bool    `json:"isfavourite"`
	StartDate    int64   `json:"startdate"`
	EndDate      int64   `json:"enddate"`
}

type timelineResponse struct {
	Courses    []enrolledCourse `json:"courses"`
	NextOffset int              `json:"nextoffset"`
}

type recentCourse struct {
	ID        int    `json:"id"`
	Shortname string `json:"shortname"`
	Fullname  string `json:"fullname"`
}

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	var classification string
	var recent bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List enrolled courses",
		Long:  "List courses the authenticated user is enrolled in. Use --recent to show recently accessed courses instead.",
		Example: `  # List in-progress courses
  moodle course list

  # List all courses
  moodle course list --timeline all

  # List past courses
  moodle course list --timeline past

  # List recently accessed courses
  moodle course list --recent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if recent {
				return listRecentCourses(f, cmd)
			}
			return listEnrolledCourses(f, cmd, classification)
		},
	}

	cmd.Flags().StringVar(&classification, "timeline", "inprogress", "Timeline filter: all, inprogress, past, future, hidden, favourites")
	cmd.Flags().BoolVar(&recent, "recent", false, "Show recently accessed courses")

	return cmd
}

func listEnrolledCourses(f *cmdutil.Factory, cmd *cobra.Command, classification string) error {
	client, err := f.Client()
	if err != nil {
		return err
	}

	var result timelineResponse
	params := map[string]any{
		"classification": classification,
	}
	err = client.Call(cmd.Context(), "core_course_get_enrolled_courses_by_timeline_classification", params, &result)
	if err != nil {
		return fmt.Errorf("failed to list courses: %w", err)
	}

	formatStr, _ := cmd.Flags().GetString("format")
	opts := output.FormatOptions{
		Format: output.ParseFormat(formatStr),
		Writer: f.IO.Out,
	}

	if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
		return f.Output(result.Courses, opts)
	}

	if len(result.Courses) == 0 {
		fmt.Fprintln(f.IO.Out, "No courses found.")
		return nil
	}

	table := &output.TableData{
		Columns: []output.Column{
			{Name: "ID", Width: 6},
			{Name: "Shortname", Width: 12},
			{Name: "Name", Width: 40},
			{Name: "Category", Width: 20},
			{Name: "Progress", Width: 10},
		},
		Rows: make([]map[string]string, 0, len(result.Courses)),
	}

	for _, c := range result.Courses {
		progress := "--"
		if c.HasProgress {
			progress = fmt.Sprintf("%.0f%%", c.Progress)
		}
		table.Rows = append(table.Rows, map[string]string{
			"ID":        strconv.Itoa(c.ID),
			"Shortname": c.Shortname,
			"Name":      c.Fullname,
			"Category":  c.CategoryName,
			"Progress":  progress,
		})
	}

	return f.Output(table, opts)
}

func listRecentCourses(f *cmdutil.Factory, cmd *cobra.Command) error {
	client, err := f.Client()
	if err != nil {
		return err
	}

	var courses []recentCourse
	err = client.Call(cmd.Context(), "core_course_get_recent_courses", nil, &courses)
	if err != nil {
		return fmt.Errorf("failed to get recent courses: %w", err)
	}

	formatStr, _ := cmd.Flags().GetString("format")
	opts := output.FormatOptions{
		Format: output.ParseFormat(formatStr),
		Writer: f.IO.Out,
	}

	if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
		return f.Output(courses, opts)
	}

	if len(courses) == 0 {
		fmt.Fprintln(f.IO.Out, "No courses found.")
		return nil
	}

	table := &output.TableData{
		Columns: []output.Column{
			{Name: "ID", Width: 6},
			{Name: "Shortname", Width: 12},
			{Name: "Name", Width: 40},
		},
		Rows: make([]map[string]string, 0, len(courses)),
	}

	for _, c := range courses {
		table.Rows = append(table.Rows, map[string]string{
			"ID":        strconv.Itoa(c.ID),
			"Shortname": c.Shortname,
			"Name":      c.Fullname,
		})
	}

	return f.Output(table, opts)
}
