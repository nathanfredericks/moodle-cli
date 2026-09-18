package database

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

type fieldsResponse struct {
	Fields []fieldItem `json:"fields"`
}

type fieldItem struct {
	ID          int    `json:"id"`
	DataID      int    `json:"dataid"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

func newFieldsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fields <database-id>",
		Short: "List fields in a database",
		Long:  "List the field definitions for a database activity.",
		Example: `  # List fields for database 235
  moodle database fields 235

  # Output as JSON
  moodle database fields 235 -f json`,
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

			var result fieldsResponse
			params := map[string]any{
				"databaseid": dbID,
			}
			if err := client.Call(cmd.Context(), "mod_data_get_fields", params, &result); err != nil {
				return fmt.Errorf("failed to get fields: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(&result, opts)
			}

			if len(result.Fields) == 0 {
				fmt.Fprintln(f.IO.Out, "No fields found.")
				return nil
			}

			table := &output.TableData{
				Columns: []output.Column{
					{Name: "ID", Width: 6},
					{Name: "Name", Width: 20},
					{Name: "Type", Width: 12},
					{Name: "Required", Width: 10},
					{Name: "Description", Width: 30},
				},
				Rows: make([]map[string]string, 0, len(result.Fields)),
			}

			for _, field := range result.Fields {
				reqStr := "No"
				if field.Required {
					reqStr = "Yes"
				}
				table.Rows = append(table.Rows, map[string]string{
					"ID":          strconv.Itoa(field.ID),
					"Name":        field.Name,
					"Type":        field.Type,
					"Required":    reqStr,
					"Description": field.Description,
				})
			}

			return f.Output(table, opts)
		},
	}

	return cmd
}
