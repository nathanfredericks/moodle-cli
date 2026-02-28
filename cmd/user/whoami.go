package user

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type siteInfo struct {
	UserID           int    `json:"userid"`
	Username         string `json:"username"`
	Fullname         string `json:"fullname"`
	SiteName         string `json:"sitename"`
	SiteURL          string `json:"siteurl"`
	UserIsSiteAdmin  bool   `json:"userissiteadmin"`
	UserPictureURL   string `json:"userpictureurl"`
}

func newWhoamiCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Show current user info",
		Long:  "Display information about the currently authenticated user.",
		Example: `  # Show current user info
  moodle user whoami

  # Output as JSON
  moodle user whoami -f json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			var info siteInfo
			err = client.Call(cmd.Context(), "core_webservice_get_site_info", nil, &info)
			if err != nil {
				return fmt.Errorf("failed to get site info: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(info, opts)
			}

			admin := "No"
			if info.UserIsSiteAdmin {
				admin = "Yes"
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "Field"},
					{Name: "Value"},
				},
				Rows: []map[string]string{
					{"Field": "User ID", "Value": strconv.Itoa(info.UserID)},
					{"Field": "Username", "Value": info.Username},
					{"Field": "Full Name", "Value": info.Fullname},
					{"Field": "Site", "Value": info.SiteName},
					{"Field": "URL", "Value": info.SiteURL},
					{"Field": "Admin", "Value": admin},
				},
			}

			return f.Output(table, opts)
		},
	}

	return cmd
}
