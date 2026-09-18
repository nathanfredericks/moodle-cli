package main

import (
	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"

	assignmentcmd "github.com/nathanfredericks/moodle-cli/cmd/assignment"
	authcmd "github.com/nathanfredericks/moodle-cli/cmd/auth"
	configcmd "github.com/nathanfredericks/moodle-cli/cmd/config"
	coursecmd "github.com/nathanfredericks/moodle-cli/cmd/course"
	databasecmd "github.com/nathanfredericks/moodle-cli/cmd/database"
	forumcmd "github.com/nathanfredericks/moodle-cli/cmd/forum"
	mcpcmd "github.com/nathanfredericks/moodle-cli/cmd/mcp"
	usercmd "github.com/nathanfredericks/moodle-cli/cmd/user"
)

func registerCommands(rootCmd *cobra.Command, f *cmdutil.Factory) {
	rootCmd.AddCommand(authcmd.NewCmd(f))
	rootCmd.AddCommand(configcmd.NewCmd(f))
	rootCmd.AddCommand(coursecmd.NewCmd(f))
	rootCmd.AddCommand(usercmd.NewCmd(f))
	rootCmd.AddCommand(assignmentcmd.NewCmd(f))
	rootCmd.AddCommand(databasecmd.NewCmd(f))
	rootCmd.AddCommand(forumcmd.NewCmd(f))
	rootCmd.AddCommand(mcpcmd.NewCmd(f, Version))
	rootCmd.AddCommand(newShellCompletionCmd())
}
