package assignment

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type calendarEventsResponse struct {
	Events []calendarEvent `json:"events"`
}

type calendarEvent struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	ModuleName string         `json:"modulename"`
	Overdue    bool           `json:"overdue"`
	TimeSort   int64          `json:"timesort"`
	URL        string         `json:"url"`
	Course     calendarCourse `json:"course"`
}

type calendarCourse struct {
	ID       int    `json:"id"`
	FullName string `json:"fullname"`
}

// resolvePeriod converts a period flag value to a Unix timestamp range.
func resolvePeriod(period string) (from, to int64, err error) {
	now := time.Now()
	switch period {
	case "overdue":
		from = now.AddDate(0, 0, -180).Unix()
		to = now.Unix()
	case "7days":
		from = now.Unix()
		to = now.AddDate(0, 0, 7).Unix()
	case "30days":
		from = now.Unix()
		to = now.AddDate(0, 0, 30).Unix()
	case "3months":
		from = now.Unix()
		to = now.AddDate(0, 0, 90).Unix()
	case "6months":
		from = now.Unix()
		to = now.AddDate(0, 0, 180).Unix()
	default:
		return 0, 0, fmt.Errorf("invalid period %q: must be one of overdue, 7days, 30days, 3months, 6months", period)
	}
	return from, to, nil
}

func newDueCmd(f *cmdutil.Factory) *cobra.Command {
	var (
		period   string
		courseID int
	)

	cmd := &cobra.Command{
		Use:   "due",
		Short: "Show upcoming due dates",
		Long:  "Show upcoming assignment and activity due dates from the calendar timeline.",
		Example: `  # Show items due in the next 7 days (default)
  moodle assignment due

  # Show overdue items
  moodle assignment due --period overdue

  # Show items due in the next 30 days
  moodle assignment due --period 30days

  # Filter by course
  moodle assignment due --course 42

  # Output as JSON
  moodle assignment due -f json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			from, to, err := resolvePeriod(period)
			if err != nil {
				return err
			}

			isOverdue := period == "overdue"

			var allEvents []calendarEvent
			afterEventID := 0
			const limitNum = 50

			for {
				params := map[string]any{
					"timesortfrom": from,
					"timesortto":   to,
					"limitnum":     limitNum,
				}
				if afterEventID > 0 {
					params["aftereventid"] = afterEventID
				}

				var result calendarEventsResponse
				if err := client.Call(cmd.Context(), "core_calendar_get_action_events_by_timesort", params, &result); err != nil {
					return fmt.Errorf("failed to get calendar events: %w", err)
				}

				if len(result.Events) == 0 {
					break
				}

				for i := range result.Events {
					e := &result.Events[i]

					if isOverdue && !e.Overdue {
						continue
					}
					if courseID > 0 && e.Course.ID != courseID {
						continue
					}
					allEvents = append(allEvents, *e)
				}

				if len(result.Events) < limitNum {
					break
				}
				afterEventID = result.Events[len(result.Events)-1].ID
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(allEvents, opts)
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "ID"},
					{Name: "Course ID"},
					{Name: "Type"},
					{Name: "Course"},
					{Name: "Name"},
					{Name: "Due Date"},
				},
				Rows: []map[string]string{},
			}

			for _, e := range allEvents {
				dueStr := "No due date"
				if e.TimeSort > 0 {
					t := time.Unix(e.TimeSort, 0)
					dueStr = t.Format("2006-01-02 15:04 MST")
					if e.Overdue {
						dueStr += " (overdue)"
					}
				}
				table.Rows = append(table.Rows, map[string]string{
					"ID":        strconv.Itoa(e.ID),
					"Course ID": strconv.Itoa(e.Course.ID),
					"Type":      e.ModuleName,
					"Course":    e.Course.FullName,
					"Name":      strings.TrimSuffix(e.Name, " is due"),
					"Due Date":  dueStr,
				})
			}

			if len(table.Rows) == 0 {
				fmt.Fprintln(f.IO.Out, "No upcoming due dates found.")
				return nil
			}

			return f.Output(table, opts)
		},
	}

	cmd.Flags().StringVar(&period, "period", "7days", "Time period: overdue, 7days, 30days, 3months, 6months")
	cmd.Flags().IntVar(&courseID, "course", 0, "Filter by course ID")
	return cmd
}
