package config

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/config"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration settings",
		Long:  "Display all configuration keys and their current values in a table.",
		Example: `  # List all settings
  moodle config list

  # List settings as JSON
  moodle config list -f json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, ok := f.Config.(*config.FileConfigManager)
			if !ok {
				return fmt.Errorf("unsupported config manager type")
			}
			settings := mgr.AllSettings()
			if len(settings) == 0 {
				fmt.Fprintln(f.IO.Out, "No configuration settings.")
				return nil
			}

			keys := make([]string, 0, len(settings))
			for k := range settings {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			td := &output.TableData{
				Columns: []output.Column{
					{Name: "Key"},
					{Name: "Value"},
				},
			}
			for _, k := range keys {
				td.Rows = append(td.Rows, map[string]string{
					"Key":   k,
					"Value": settings[k],
				})
			}
			return f.Output(td, output.FormatOptions{Writer: f.IO.Out})
		},
	}
	return cmd
}
