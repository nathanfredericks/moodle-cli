package user

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/api"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type enrolledUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	var courseID int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users enrolled in a course",
		Long:  "List all users enrolled in a specific course.",
		Example: `  # List users in course 42
  moodle user list --course 42

  # Output as JSON
  moodle user list --course 42 -f json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			return listCourseUsers(cmd, client, f, courseID, opts)
		},
	}

	cmd.Flags().IntVar(&courseID, "course", 0, "Course ID to list enrolled users for")
	_ = cmd.MarkFlagRequired("course")

	return cmd
}

func listCourseUsers(cmd *cobra.Command, client api.MoodleClient, f *cmdutil.Factory, courseID int, opts output.FormatOptions) error {
	var users []enrolledUser
	params := map[string]any{
		"courseid": courseID,
	}
	err := client.Call(cmd.Context(), "core_enrol_get_enrolled_users", params, &users)
	if err != nil {
		return fmt.Errorf("failed to list enrolled users: %w", err)
	}

	if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
		return f.Output(users, opts)
	}

	if len(users) == 0 {
		fmt.Fprintln(f.IO.Out, "No users found.")
		return nil
	}

	table := &output.TableData{
		Columns: []output.Column{
			{Name: "ID", Width: 6},
			{Name: "Username", Width: 15},
			{Name: "Name", Width: 30},
			{Name: "Email", Width: 30},
		},
		Rows: make([]map[string]string, 0, len(users)),
	}

	for _, u := range users {
		table.Rows = append(table.Rows, map[string]string{
			"ID":       strconv.Itoa(u.ID),
			"Username": u.Username,
			"Name":     u.Fullname,
			"Email":    u.Email,
		})
	}

	return f.Output(table, opts)
}
