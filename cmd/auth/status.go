package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type siteInfo struct {
	SiteName  string `json:"sitename"`
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	SiteURL   string `json:"siteurl"`
	UserID    int    `json:"userid"`
	Functions []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"functions"`
}

func newStatusCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Long:  "Display the current authentication status by querying the Moodle instance.",
		Example: `  # Show current authentication status
  moodle auth status

  # Output status as JSON
  moodle auth status -f json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			var info siteInfo
			if err := client.Call(cmd.Context(), "core_webservice_get_site_info", nil, &info); err != nil {
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

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "Field"},
					{Name: "Value"},
				},
				Rows: []map[string]string{
					{"Field": "URL", "Value": info.SiteURL},
					{"Field": "Site", "Value": info.SiteName},
					{"Field": "User", "Value": info.Fullname},
					{"Field": "Username", "Value": info.Username},
					{"Field": "User ID", "Value": fmt.Sprintf("%d", info.UserID)},
				},
			}

			return f.Output(table, opts)
		},
	}

	return cmd
}
