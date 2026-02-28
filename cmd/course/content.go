package course

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type courseSection struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Section int            `json:"section"`
	Visible int            `json:"visible"`
	Modules []courseModule `json:"modules"`
}

type courseModule struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	ModName     string           `json:"modname"`
	Visible     int              `json:"visible"`
	Description string           `json:"description"`
	Contents    []moduleContent  `json:"contents"`
}

func newContentCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content <course-id>",
		Short: "Get course contents",
		Long:  "Display the sections and activities of a course.",
		Example: `  # View course sections and activities
  moodle course content 42

  # Output as JSON for processing
  moodle course content 42 -f json`,
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

			var sections []courseSection
			params := map[string]any{
				"courseid": courseID,
			}
			err = client.Call(cmd.Context(), "core_course_get_contents", params, &sections)
			if err != nil {
				return fmt.Errorf("failed to get course contents: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(sections, opts)
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "Section"},
					{Name: "Module ID"},
					{Name: "Type"},
					{Name: "Name"},
				},
				Rows: make([]map[string]string, 0),
			}

			for _, s := range sections {
				sectionName := s.Name
				if sectionName == "" {
					sectionName = fmt.Sprintf("Section %d", s.Section)
				}
				if len(s.Modules) == 0 {
					table.Rows = append(table.Rows, map[string]string{
						"Section":   sectionName,
						"Module ID": "--",
						"Type":      "--",
						"Name":      "(no activities)",
					})
					continue
				}
				for _, m := range s.Modules {
					table.Rows = append(table.Rows, map[string]string{
						"Section":   sectionName,
						"Module ID": strconv.Itoa(m.ID),
						"Type":      m.ModName,
						"Name":      m.Name,
					})
				}
			}

			return f.Output(table, opts)
		},
	}

	return cmd
}
