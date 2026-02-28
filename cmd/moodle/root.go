package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

// NewRootCmd creates the root command for the Moodle CLI.
func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	var formatStr string
	var noColor bool
	var verbose bool

	cmd := &cobra.Command{
		Use:   "moodle",
		Short: "CLI for the Moodle LMS",
		Long:  "A command-line interface for the Moodle Learning Management System API.",
		Example: `  # List your enrolled courses
  moodle course list

  # Get details about a specific course
  moodle course get 42

  # Search for courses by name
  moodle course search "Introduction to Computing"

  # View your assignments
  moodle assignment list --course 42

  # Output as JSON for scripting
  moodle course list -f json`,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			f.NoColor = noColor || output.NoColorEnabled()
			f.Verbose = verbose
		},
	}

	cmd.PersistentFlags().StringVarP(&formatStr, "format", "f", "", "Output format: table, json, csv, yaml, plain")
	cmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Version command
	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(f.IO.Out, "moodle version %s\n", Version)
		},
	})

	return cmd
}

// Execute runs the root command.
func Execute() {
	f, err := cmdutil.NewFactory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	rootCmd := NewRootCmd(f)

	// Register all subcommand groups here.
	registerCommands(rootCmd, f)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
