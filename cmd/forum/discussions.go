package forum

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type discussion struct {
	Discussion   int    `json:"discussion"`
	Name         string `json:"name"`
	UserFullname string `json:"userfullname"`
	NumReplies   int    `json:"numreplies"`
	NumUnread    int    `json:"numunread"`
	Pinned       bool   `json:"pinned"`
	Locked       bool   `json:"locked"`
	TimeModified int64  `json:"timemodified"`
}

type discussionsResponse struct {
	Discussions []discussion `json:"discussions"`
}

func newDiscussionsCmd(f *cmdutil.Factory) *cobra.Command {
	var page int
	var perPage int

	cmd := &cobra.Command{
		Use:   "discussions <forum-id>",
		Short: "List discussions in a forum",
		Long:  "List discussions in a forum with pagination.",
		Example: `  # List discussions in forum 5
  moodle forum discussions 5

  # Paginate through discussions
  moodle forum discussions 5 --page 1 --per-page 10`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			forumID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid forum ID: %s", args[0])
			}

			client, err := f.Client()
			if err != nil {
				return err
			}

			var result discussionsResponse
			params := map[string]any{
				"forumid": forumID,
				"page":    page,
				"perpage": perPage,
			}
			if err := client.Call(cmd.Context(), "mod_forum_get_forum_discussions", params, &result); err != nil {
				return fmt.Errorf("failed to list discussions: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(result.Discussions, opts)
			}

			if len(result.Discussions) == 0 {
				fmt.Fprintln(f.IO.Out, "No discussions found.")
				return nil
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "ID", Width: 8},
					{Name: "Subject", Width: 40},
					{Name: "Author", Width: 20},
					{Name: "Replies", Width: 8},
					{Name: "Unread", Width: 8},
					{Name: "Date", Width: 16},
					{Name: "Status", Width: 10},
				},
				Rows: make([]map[string]string, 0, len(result.Discussions)),
			}

			for _, d := range result.Discussions {
				status := ""
				if d.Pinned {
					status = "pinned"
				}
				if d.Locked {
					if status != "" {
						status += ", "
					}
					status += "locked"
				}

				table.Rows = append(table.Rows, map[string]string{
					"ID":      strconv.Itoa(d.Discussion),
					"Subject": d.Name,
					"Author":  d.UserFullname,
					"Replies": strconv.Itoa(d.NumReplies),
					"Unread":  strconv.Itoa(d.NumUnread),
					"Date":    formatTimestamp(d.TimeModified),
					"Status":  status,
				})
			}

			return f.Output(table, opts)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number (0-indexed)")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Discussions per page")

	return cmd
}

func formatTimestamp(ts int64) string {
	if ts == 0 {
		return "--"
	}
	t := time.Unix(ts, 0)
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins <= 1 {
			return "just now"
		}
		return fmt.Sprintf("%dm ago", mins)
	case diff < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return t.Format("Mon 3:04 PM MST")
	default:
		if t.Year() == now.Year() {
			return t.Format("Jan 2, 3:04 PM MST")
		}
		return t.Format("Jan 2, 2006 MST")
	}
}
