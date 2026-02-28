package assignment

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type submissionStatusResponse struct {
	AssignmentData assignmentDetail `json:"assignmentdata"`
	LastAttempt    *lastAttemptData `json:"lastattempt"`
	Feedback       *feedbackData   `json:"feedback"`
}

type assignmentDetail struct {
	Activity       string `json:"activity"`
	ActivityFormat int    `json:"activityformat"`
}

type lastAttemptData struct {
	Submission    *submissionData `json:"submission"`
	CanSubmit     bool            `json:"cansubmit"`
	CanEdit       bool            `json:"canedit"`
	GradingStatus string          `json:"gradingstatus"`
	Locked        bool            `json:"locked"`
}

type submissionData struct {
	ID            int              `json:"id"`
	UserID        int              `json:"userid"`
	Status        string           `json:"status"`
	TimeCreated   int64            `json:"timecreated"`
	TimeModified  int64            `json:"timemodified"`
	AttemptNumber int              `json:"attemptnumber"`
	Plugins       []pluginData     `json:"plugins"`
}

type pluginData struct {
	Type         string           `json:"type"`
	Name         string           `json:"name"`
	FileAreas    []fileAreaData   `json:"fileareas"`
	EditorFields []editorField    `json:"editorfields"`
}

type fileAreaData struct {
	Area  string         `json:"area"`
	Files []fileEntry    `json:"files"`
}

type fileEntry struct {
	FileName string `json:"filename"`
	FilePath string `json:"filepath"`
	FileSize int64  `json:"filesize"`
	FileURL  string `json:"fileurl"`
	MimeType string `json:"mimetype"`
}

type editorField struct {
	Name   string `json:"name"`
	Text   string `json:"text"`
	Format int    `json:"format"`
}

type feedbackData struct {
	Grade            *gradeData   `json:"grade"`
	GradedDate       int64        `json:"gradeddate"`
	GradeForDisplay  string       `json:"gradefordisplay"`
	Plugins          []pluginData `json:"plugins"`
}

type gradeData struct {
	ID    int    `json:"id"`
	Grade string `json:"grade"`
}

func newStatusCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <assignment-id>",
		Short: "Check submission status",
		Long:  "View submission status for the current user on an assignment.",
		Example: `  # Check submission status
  moodle assignment status 101

  # Output status as JSON
  moodle assignment status 101 -f json`,
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

			var result submissionStatusResponse
			params := map[string]any{
				"assignid": assignID,
			}
			if err := client.Call(cmd.Context(), "mod_assign_get_submission_status", params, &result); err != nil {
				return fmt.Errorf("failed to get submission status: %w", err)
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
					{Name: "Field"},
					{Name: "Value"},
				},
				Rows: []map[string]string{
					{"Field": "Assignment ID", "Value": args[0]},
				},
			}

			if result.LastAttempt != nil {
				if result.LastAttempt.Submission != nil {
					sub := result.LastAttempt.Submission
					table.Rows = append(table.Rows,
						map[string]string{"Field": "Status", "Value": sub.Status},
						map[string]string{"Field": "Attempt", "Value": strconv.Itoa(sub.AttemptNumber + 1)},
					)
				}
				table.Rows = append(table.Rows,
					map[string]string{"Field": "Can Edit", "Value": fmt.Sprintf("%v", result.LastAttempt.CanEdit)},
					map[string]string{"Field": "Can Submit", "Value": fmt.Sprintf("%v", result.LastAttempt.CanSubmit)},
					map[string]string{"Field": "Locked", "Value": fmt.Sprintf("%v", result.LastAttempt.Locked)},
					map[string]string{"Field": "Grading Status", "Value": result.LastAttempt.GradingStatus},
				)
			}

			return f.Output(table, opts)
		},
	}

	return cmd
}
