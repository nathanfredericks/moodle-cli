package course

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type courseModuleInfo struct {
	ID         int     `json:"id"`
	Course     int     `json:"course"`
	Module     int     `json:"module"`
	Name       string  `json:"name"`
	ModName    string  `json:"modname"`
	Instance   int     `json:"instance"`
	Section    int     `json:"section"`
	SectionNum int     `json:"sectionnum"`
	Visible    int     `json:"visible"`
	Grade      float64 `json:"grade"`
	GradePass  string  `json:"gradepass"`
	Completion int     `json:"completion"`
}

type courseModuleResponse struct {
	CM courseModuleInfo `json:"cm"`
}

func newModuleCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module <module-id>",
		Short: "Get course module details",
		Long:  "Display detailed information about a specific course module.",
		Example: `  # Get module details
  moodle course module 1001

  # Get module details as JSON
  moodle course module 1001 -f json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			cmid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid course module ID: %s", args[0])
			}

			var result courseModuleResponse
			params := map[string]any{"cmid": cmid}
			err = client.Call(cmd.Context(), "core_course_get_course_module", params, &result)
			if err != nil {
				return fmt.Errorf("failed to get course module: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(result.CM, opts)
			}

			cm := result.CM
			visible := "No"
			if cm.Visible == 1 {
				visible = "Yes"
			}
			completionStr := "None"
			switch cm.Completion {
			case 1:
				completionStr = "Manual"
			case 2:
				completionStr = "Automatic"
			}
			grade := "--"
			if cm.Grade != 0 {
				grade = fmt.Sprintf("%.2f", cm.Grade)
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "Field", Width: 12},
					{Name: "Value", Width: 40},
				},
				Rows: []map[string]string{
					{"Field": "ID", "Value": strconv.Itoa(cm.ID)},
					{"Field": "Name", "Value": cm.Name},
					{"Field": "Type", "Value": cm.ModName},
					{"Field": "Course", "Value": strconv.Itoa(cm.Course)},
					{"Field": "Section", "Value": strconv.Itoa(cm.SectionNum)},
					{"Field": "Visible", "Value": visible},
					{"Field": "Grade", "Value": grade},
					{"Field": "Completion", "Value": completionStr},
				},
			}

			return f.Output(table, opts)
		},
	}

	return cmd
}
