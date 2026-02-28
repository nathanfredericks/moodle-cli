package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

func newGetCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Long:  "Get the value of a specific configuration key. Returns an error if the key is not set.",
		Example: `  # Get the default output format
  moodle config get format`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			val := f.Config.Get(key)
			if val == "" {
				return fmt.Errorf("key %q is not set", key)
			}
			fmt.Fprintln(f.IO.Out, val)
			return nil
		},
	}
	return cmd
}
