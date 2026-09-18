package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type entriesResponse struct {
	Entries    []entryItem `json:"entries"`
	TotalCount int         `json:"totalcount"`
}

type entryItem struct {
	ID           int            `json:"id"`
	UserID       int            `json:"userid"`
	GroupID      int            `json:"groupid"`
	DataID       int            `json:"dataid"`
	TimeCreated  int64          `json:"timecreated"`
	TimeModified int64          `json:"timemodified"`
	Approved     bool           `json:"approved"`
	FullName     string         `json:"fullname"`
	Contents     []entryContent `json:"contents"`
}

type entryContent struct {
	ID       int         `json:"id"`
	FieldID  int         `json:"fieldid"`
	RecordID int         `json:"recordid"`
	Content  string      `json:"content"`
	Content1 string      `json:"content1"`
	Content2 string      `json:"content2"`
	Content3 string      `json:"content3"`
	Content4 string      `json:"content4"`
	Files    []entryFile `json:"files"`
}

type entryFile struct {
	FileName string `json:"filename"`
	FilePath string `json:"filepath"`
	FileSize int64  `json:"filesize"`
	FileURL  string `json:"fileurl"`
	MimeType string `json:"mimetype"`
}

func newEntriesCmd(f *cmdutil.Factory) *cobra.Command {
	var page int
	var perPage int

	cmd := &cobra.Command{
		Use:   "entries <database-id>",
		Short: "List entries in a database",
		Long:  "List entries in a database activity, including their field contents and file attachments.",
		Example: `  # List entries for database 235
  moodle database entries 235

  # Output as JSON with full contents
  moodle database entries 235 -f json

  # Paginate results
  moodle database entries 235 --page 1 --per-page 10`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			dbID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid database ID: %s", args[0])
			}

			params := map[string]any{
				"databaseid":     dbID,
				"returncontents": true,
			}
			if page > 0 {
				params["page"] = page
			}
			if perPage > 0 {
				params["perpage"] = perPage
			}

			var result entriesResponse
			if err := client.Call(cmd.Context(), "mod_data_get_entries", params, &result); err != nil {
				return fmt.Errorf("failed to get entries: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(&result, opts)
			}

			if len(result.Entries) == 0 {
				fmt.Fprintln(f.IO.Out, "No entries found.")
				return nil
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "ID", Width: 6},
					{Name: "User", Width: 20},
					{Name: "User ID", Width: 8},
					{Name: "Created", Width: 20},
					{Name: "Files", Width: 30},
				},
				Rows: make([]map[string]string, 0, len(result.Entries)),
			}

			for _, entry := range result.Entries {
				created := time.Unix(entry.TimeCreated, 0).Format("2006-01-02 15:04 MST")

				var fileNames []string
				for _, c := range entry.Contents {
					for _, file := range c.Files {
						fileNames = append(fileNames, file.FileName)
					}
				}
				filesStr := "None"
				if len(fileNames) > 0 {
					filesStr = ""
					for i, name := range fileNames {
						if i > 0 {
							filesStr += ", "
						}
						filesStr += name
					}
				}

				table.Rows = append(table.Rows, map[string]string{
					"ID":      strconv.Itoa(entry.ID),
					"User":    entry.FullName,
					"User ID": strconv.Itoa(entry.UserID),
					"Created": created,
					"Files":   filesStr,
				})
			}

			fmt.Fprintf(f.IO.Out, "Total entries: %d\n\n", result.TotalCount)
			return f.Output(table, opts)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number (0-indexed)")
	cmd.Flags().IntVar(&perPage, "per-page", 0, "Entries per page (0 for all)")

	return cmd
}
