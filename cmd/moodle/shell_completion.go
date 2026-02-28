package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newShellCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "shell-completion [bash|zsh|fish|powershell]",
		Short:  "Generate shell completion scripts",
		Long:   "Generate shell completion scripts for bash, zsh, fish, or powershell.",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			switch args[0] {
			case "bash":
				return root.GenBashCompletionV2(os.Stdout, true)
			case "zsh":
				return root.GenZshCompletion(os.Stdout)
			case "fish":
				return root.GenFishCompletion(os.Stdout, true)
			case "powershell":
				return root.GenPowerShellCompletion(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell %q: use bash, zsh, fish, or powershell", args[0])
			}
		},
	}

	return cmd
}
