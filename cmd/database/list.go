package database

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type databaseListResponse struct {
	Databases []databaseItem `json:"databases"`
}

type databaseItem struct {
	ID           int    `json:"id"`
	Course       int    `json:"course"`
	CourseModule int    `json:"coursemodule"`
	Name         string `json:"name"`
	Intro        string `json:"intro"`
	MaxEntries   int    `json:"maxentries"`
}

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	var courseID int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List database activities",
		Long:  "List database activities across enrolled courses.",
		Example: `  # List all database activities
  moodle database list

  # List database activities for a specific course
  moodle database list --course 42

  # Output as JSON
  moodle database list --course 42 -f json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			params := map[string]any{}
			if courseID > 0 {
				params["courseids"] = []int{courseID}
			}

			var result databaseListResponse
			if err := client.Call(cmd.Context(), "mod_data_get_databases_by_courses", params, &result); err != nil {
				return fmt.Errorf("failed to list databases: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(&result, opts)
			}

			if len(result.Databases) == 0 {
				fmt.Fprintln(f.IO.Out, "No database activities found.")
				return nil
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "ID", Width: 6},
					{Name: "Course", Width: 8},
					{Name: "CMID", Width: 8},
					{Name: "Name", Width: 30},
				},
				Rows: make([]map[string]string, 0, len(result.Databases)),
			}

			for _, db := range result.Databases {
				table.Rows = append(table.Rows, map[string]string{
					"ID":     strconv.Itoa(db.ID),
					"Course": strconv.Itoa(db.Course),
					"CMID":   strconv.Itoa(db.CourseModule),
					"Name":   db.Name,
				})
			}

			return f.Output(table, opts)
		},
	}

	cmd.Flags().IntVar(&courseID, "course", 0, "Filter by course ID")
	return cmd
}
