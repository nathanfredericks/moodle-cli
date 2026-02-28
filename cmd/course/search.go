package course

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type searchCoursesResponse struct {
	Total   int              `json:"total"`
	Courses []searchedCourse `json:"courses"`
}

type searchedCourse struct {
	ID           int    `json:"id"`
	Shortname    string `json:"shortname"`
	Fullname     string `json:"fullname"`
	CategoryName string `json:"categoryname"`
}

func newSearchCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search courses",
		Long:  "Search for courses by name across the Moodle instance.",
		Example: `  # Search for courses by name
  moodle course search "Biology"

  # Search and output as JSON
  moodle course search "CS101" -f json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			var result searchCoursesResponse
			params := map[string]any{
				"criterianame":  "search",
				"criteriavalue": args[0],
			}
			err = client.Call(cmd.Context(), "core_course_search_courses", params, &result)
			if err != nil {
				return fmt.Errorf("failed to search courses: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(result, opts)
			}

			if len(result.Courses) == 0 {
				fmt.Fprintf(f.IO.Out, "No courses found matching %q.\n", args[0])
				return nil
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "ID", Width: 6},
					{Name: "Shortname", Width: 12},
					{Name: "Name", Width: 40},
					{Name: "Category", Width: 20},
				},
				Rows: make([]map[string]string, 0, len(result.Courses)),
			}

			for _, c := range result.Courses {
				table.Rows = append(table.Rows, map[string]string{
					"ID":        strconv.Itoa(c.ID),
					"Shortname": c.Shortname,
					"Name":      c.Fullname,
					"Category":  c.CategoryName,
				})
			}

			fmt.Fprintf(f.IO.Out, "Found %d courses matching %q\n\n", result.Total, args[0])
			return f.Output(table, opts)
		},
	}

	return cmd
}
