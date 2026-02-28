package assignment

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

func newUploadCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload <assignment-id> <file>...",
		Short: "Upload files to an assignment",
		Long:  "Upload one or more files to an assignment's submission draft area and save the submission.",
		Example: `  # Upload a single file
  moodle assignment upload 101 report.pdf

  # Upload multiple files
  moodle assignment upload 101 report.pdf appendix.pdf data.csv`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			assignID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid assignment ID: %s", args[0])
			}

			files := args[1:]
			draftItemID := 0

			for _, filePath := range files {
				draft, err := client.UploadFile(cmd.Context(), filePath, draftItemID)
				if err != nil {
					return fmt.Errorf("failed to upload %s: %w", filepath.Base(filePath), err)
				}
				draftItemID = draft.ItemID
				fmt.Fprintf(f.IO.Out, "Uploaded %s (%d bytes)\n", draft.FileName, draft.FileSize)
			}

			params := map[string]any{
				"assignmentid": assignID,
				"plugindata": map[string]any{
					"files_filemanager": draftItemID,
				},
			}

			var result any
			if err := client.Call(cmd.Context(), "mod_assign_save_submission", params, &result); err != nil {
				return fmt.Errorf("failed to save submission: %w", err)
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(result, opts)
			}

			fmt.Fprintf(f.IO.Out, "Submission saved for assignment %d.\n", assignID)
			return nil
		},
	}

	return cmd
}
