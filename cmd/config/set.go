package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

func newSetCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long:  "Set a configuration key to the given value. The change is persisted to disk immediately.",
		Example: `  # Set the default output format to JSON
  moodle config set format json

  # Set the default output format to table
  moodle config set format table`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]
			if err := f.Config.Set(key, value); err != nil {
				return err
			}
			fmt.Fprintf(f.IO.Out, "Set %q to %q\n", key, value)
			return nil
		},
	}
	return cmd
}
