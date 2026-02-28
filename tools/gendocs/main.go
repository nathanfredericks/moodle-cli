// Command gendocs generates CLI documentation from the command tree.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/nathanfredericks/moodle-cli/internal/root"
)

func main() {
	outDir := flag.String("out", "docs/", "Output directory for generated docs")
	format := flag.String("format", "markdown", "Output format: markdown, man, or rest")
	flag.Parse()

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	cmd := root.Root()

	var err error
	switch *format {
	case "markdown":
		err = doc.GenMarkdownTree(cmd, *outDir)
	case "man":
		header := &doc.GenManHeader{
			Title:   "MOODLE",
			Section: "1",
		}
		err = doc.GenManTree(cmd, header, *outDir)
	case "rest":
		err = doc.GenReSTTree(cmd, *outDir)
	default:
		log.Fatalf("unsupported format: %s (use markdown, man, or rest)", *format)
	}

	if err != nil {
		log.Fatalf("failed to generate docs: %v", err)
	}

	fmt.Printf("Documentation generated in %s (format: %s)\n", *outDir, *format)
}
