package config

import (
	"github.com/spf13/cobra"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

// NewCmd creates the config command group.
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "Get, set, and manage CLI configuration.",
		Example: `  # List all configuration settings
  moodle config list

  # Get a specific configuration value
  moodle config get format

  # Set the default output format
  moodle config set format json`,
	}
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newSetCmd(f))
	cmd.AddCommand(newListCmd(f))
	return cmd
}
