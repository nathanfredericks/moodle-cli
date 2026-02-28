package user

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type userProfile struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	IDNumber    string `json:"idnumber"`
	Institution string `json:"institution"`
	Department  string `json:"department"`
	City        string `json:"city"`
	Country     string `json:"country"`
	FirstAccess int64  `json:"firstaccess"`
	LastAccess  int64  `json:"lastaccess"`
	Suspended   bool   `json:"suspended"`
}

func newGetCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <user-id>",
		Short: "Get user by ID",
		Long:  "Get detailed information about a user by their numeric ID.",
		Example: `  # Get user details
  moodle user get 7

  # Output as JSON
  moodle user get 7 -f json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			userID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[0])
			}

			var users []userProfile
			params := map[string]any{
				"field":  "id",
				"values": []string{strconv.Itoa(userID)},
			}
			err = client.Call(cmd.Context(), "core_user_get_users_by_field", params, &users)
			if err != nil {
				return fmt.Errorf("failed to get user: %w", err)
			}

			if len(users) == 0 {
				return fmt.Errorf("user not found: %d", userID)
			}

			u := users[0]

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(u, opts)
			}

			suspended := "No"
			if u.Suspended {
				suspended = "Yes"
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "Field"},
					{Name: "Value"},
				},
				Rows: []map[string]string{
					{"Field": "ID", "Value": strconv.Itoa(u.ID)},
					{"Field": "Username", "Value": u.Username},
					{"Field": "Full Name", "Value": u.Fullname},
					{"Field": "Email", "Value": u.Email},
					{"Field": "Institution", "Value": u.Institution},
					{"Field": "Department", "Value": u.Department},
					{"Field": "City", "Value": u.City},
					{"Field": "Country", "Value": u.Country},
					{"Field": "Suspended", "Value": suspended},
				},
			}

			return f.Output(table, opts)
		},
	}

	return cmd
}
