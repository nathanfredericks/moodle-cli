package assignment

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/auth"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/text"
)

func newDownloadCmd(f *cmdutil.Factory) *cobra.Command {
	var outputDir string
	var force bool
	var resources bool

	cmd := &cobra.Command{
		Use:   "download <assignment-id>",
		Short: "Download submission files",
		Long:  "Download all files from your submission for an assignment. Use --resources to download instructor-attached resource files instead.",
		Example: `  # Download submission files
  moodle assignment download 101

  # Download to a specific directory
  moodle assignment download 101 -o ./submissions

  # Overwrite existing files
  moodle assignment download 101 -o ./submissions -F

  # Download instructor-attached resources
  moodle assignment download 101 --resources`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			assignID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid assignment ID: %s", args[0])
			}

			var files []fileEntry

			if resources {
				// Download intro attachments (instructor resources)
				var listResult assignmentListResponse
				if err := client.Call(cmd.Context(), "mod_assign_get_assignments", map[string]any{}, &listResult); err != nil {
					return fmt.Errorf("failed to get assignments: %w", err)
				}

				var found *assignmentItem
				for _, c := range listResult.Courses {
					for i, a := range c.Assignments {
						if a.ID == assignID {
							found = &c.Assignments[i]
							break
						}
					}
					if found != nil {
						break
					}
				}

				if found == nil {
					return fmt.Errorf("assignment %d not found", assignID)
				}

				files = found.IntroAttachments

				if len(files) == 0 {
					fmt.Fprintln(f.IO.Out, "No resource files found for this assignment.")
					return nil
				}
			} else {
				// Download submission files
				var result submissionStatusResponse
				params := map[string]any{
					"assignid": assignID,
				}
				if err := client.Call(cmd.Context(), "mod_assign_get_submission_status", params, &result); err != nil {
					return fmt.Errorf("failed to get submission status: %w", err)
				}

				if result.LastAttempt == nil || result.LastAttempt.Submission == nil {
					return fmt.Errorf("no submission found for assignment %d", assignID)
				}

				for _, p := range result.LastAttempt.Submission.Plugins {
					if p.Type == "file" {
						for _, fa := range p.FileAreas {
							files = append(files, fa.Files...)
						}
					}
				}

				if len(files) == 0 {
					fmt.Fprintln(f.IO.Out, "No files found in submission.")
					return nil
				}
			}

			// Get auth token for download URLs
			token, err := f.Auth.Get(auth.TokenKey)
			if err != nil {
				return fmt.Errorf("no token available: %w", err)
			}

			if outputDir != "" {
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("unable to create output directory: %w", err)
				}
			}

			fmt.Fprintf(f.IO.Out, "Downloading %d file(s)...\n", len(files))

			httpClient := &http.Client{Timeout: 5 * time.Minute}
			for _, file := range files {
				fileName := file.FileName
				if outputDir != "" {
					fileName = filepath.Join(outputDir, fileName)
				}

				if !force {
					if _, err := os.Stat(fileName); err == nil {
						fmt.Fprintf(f.IO.ErrOut, "Skipping %s (already exists, use --force to overwrite)\n", fileName)
						continue
					}
				}

				if err := downloadSubmissionFile(httpClient, file.FileURL, token, fileName, file.FileSize, f.IO); err != nil {
					fmt.Fprintf(f.IO.ErrOut, "Error downloading %s: %v\n", file.FileName, err)
					continue
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
	cmd.Flags().BoolVarP(&force, "force", "F", false, "Overwrite existing files")
	cmd.Flags().BoolVarP(&resources, "resources", "R", false, "Download instructor-attached resource files instead of submission files")

	return cmd
}

func downloadSubmissionFile(httpClient *http.Client, fileURL, token, fileName string, fileSize int64, streams cmdutil.IOStreams) error {
	sep := "?"
	if strings.Contains(fileURL, "?") {
		sep = "&"
	}
	downloadURL := fileURL + sep + "token=" + token

	resp, err := httpClient.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer outFile.Close()

	written, err := io.Copy(outFile, resp.Body)
	if err != nil {
		outFile.Close()
		os.Remove(fileName)
		return fmt.Errorf("download interrupted: %w", err)
	}

	fmt.Fprintf(streams.Out, "  %s (%s)\n", filepath.Base(fileName), text.FormatFileSize(written))
	return nil
}
