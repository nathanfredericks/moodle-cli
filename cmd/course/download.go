package course

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

type moduleContent struct {
	Type         string `json:"type"`
	FileName     string `json:"filename"`
	FileSize     int64  `json:"filesize"`
	FileURL      string `json:"fileurl"`
	TimeModified int64  `json:"timemodified"`
}

func newDownloadCmd(f *cmdutil.Factory) *cobra.Command {
	var outputDir string
	var force bool

	cmd := &cobra.Command{
		Use:   "download <course-id> <module-id>",
		Short: "Download module files",
		Long:  "Download all files from a course module (e.g. a resource or folder activity).",
		Example: `  # Download files from a module
  moodle course download 42 1001

  # Download to a specific directory
  moodle course download 42 1001 -o ./downloads

  # Overwrite existing files
  moodle course download 42 1001 -o ./downloads -F`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.Client()
			if err != nil {
				return err
			}

			courseID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}
			moduleID, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid module ID: %s", args[1])
			}

			var sections []courseSection
			params := map[string]any{"courseid": courseID}
			if err := client.Call(cmd.Context(), "core_course_get_contents", params, &sections); err != nil {
				return fmt.Errorf("failed to get course contents: %w", err)
			}

			// Find the module
			var contents []moduleContent
			var moduleName string
			for _, s := range sections {
				for _, m := range s.Modules {
					if m.ID == moduleID {
						contents = m.Contents
						moduleName = m.Name
						break
					}
				}
			}
			if moduleName == "" {
				return fmt.Errorf("module %d not found in course %d", moduleID, courseID)
			}

			// Filter to actual files
			var files []moduleContent
			for _, c := range contents {
				if c.Type == "file" {
					files = append(files, c)
				}
			}
			if len(files) == 0 {
				fmt.Fprintf(f.IO.Out, "No downloadable files in module %q.\n", moduleName)
				return nil
			}

			// Get auth token
			token, err := f.Auth.Get(auth.TokenKey)
			if err != nil {
				return fmt.Errorf("no token available: %w", err)
			}

			// Ensure output directory exists
			if outputDir != "" {
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("unable to create output directory: %w", err)
				}
			}

			fmt.Fprintf(f.IO.Out, "Downloading %d file(s) from %q...\n", len(files), moduleName)

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

				if err := downloadFile(httpClient, file.FileURL, token, fileName, file.FileSize, f.IO); err != nil {
					fmt.Fprintf(f.IO.ErrOut, "Error downloading %s: %v\n", file.FileName, err)
					continue
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
	cmd.Flags().BoolVarP(&force, "force", "F", false, "Overwrite existing files")

	return cmd
}

func downloadFile(httpClient *http.Client, fileURL, token, fileName string, fileSize int64, streams cmdutil.IOStreams) error {
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

	totalSize := resp.ContentLength
	if totalSize <= 0 {
		totalSize = fileSize
	}

	if totalSize > 0 {
		fmt.Fprintf(streams.Out, "  %s (%s)\n", filepath.Base(fileName), text.FormatFileSize(totalSize))
	} else {
		fmt.Fprintf(streams.Out, "  %s\n", filepath.Base(fileName))
	}

	_, err = io.Copy(outFile, &progressReader{
		reader: resp.Body,
		total:  totalSize,
		w:      streams.ErrOut,
	})
	if err != nil {
		outFile.Close()
		os.Remove(fileName)
		return fmt.Errorf("download interrupted: %w", err)
	}

	return nil
}

type progressReader struct {
	reader  io.Reader
	total   int64
	current int64
	w       io.Writer
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.current += int64(n)
	if pr.total > 0 && pr.total > 1024*1024 {
		pct := float64(pr.current) / float64(pr.total) * 100
		fmt.Fprintf(pr.w, "\r  %s / %s (%.0f%%)",
			text.FormatFileSize(pr.current), text.FormatFileSize(pr.total), pct)
		if err == io.EOF {
			fmt.Fprintln(pr.w)
		}
	}
	return n, err
}
