package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format represents an output format type.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
	FormatYAML  Format = "yaml"
	FormatPlain Format = "plain"
)

// ParseFormat converts a string to a Format, returning FormatTable if unknown.
func ParseFormat(s string) Format {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON
	case "csv":
		return FormatCSV
	case "yaml":
		return FormatYAML
	case "plain":
		return FormatPlain
	default:
		return FormatTable
	}
}

// FormatOptions controls output behavior.
type FormatOptions struct {
	Format  Format
	Fields  []string // field names to include
	Sort    string   // field to sort by
	Filter  string   // filter expression (field=value, field!=value, etc.)
	NoColor bool
	Wide    bool
	Writer  io.Writer // defaults to os.Stdout
}

// Column defines a table column.
type Column struct {
	Name     string
	Width    int
	Priority int // lower = more important, shown first
}

// TableData represents structured data for table output.
type TableData struct {
	Columns []Column
	Rows    []map[string]string
}

// Print formats and writes data to the configured writer.
func Print(data any, opts FormatOptions) error {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}

	switch opts.Format {
	case FormatJSON:
		return printJSON(w, data)
	case FormatCSV:
		return printCSV(w, data)
	case FormatYAML:
		return printYAML(w, data)
	case FormatPlain:
		return printPlain(w, data)
	default:
		return printTable(w, data, opts)
	}
}

func printJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func printYAML(w io.Writer, data any) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer enc.Close()
	return enc.Encode(data)
}

func printCSV(w io.Writer, data any) error {
	td, ok := data.(*TableData)
	if !ok {
		return printJSON(w, data)
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	headers := make([]string, len(td.Columns))
	for i, col := range td.Columns {
		headers[i] = col.Name
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, row := range td.Rows {
		record := make([]string, len(td.Columns))
		for i, col := range td.Columns {
			record[i] = row[col.Name]
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func printPlain(w io.Writer, data any) error {
	td, ok := data.(*TableData)
	if !ok {
		return printJSON(w, data)
	}

	for _, row := range td.Rows {
		for _, col := range td.Columns {
			if val, ok := row[col.Name]; ok && val != "" {
				fmt.Fprintf(w, "%s: %s\n", col.Name, val)
			}
		}
		fmt.Fprintln(w)
	}
	return nil
}

func printTable(w io.Writer, data any, opts FormatOptions) error {
	td, ok := data.(*TableData)
	if !ok {
		return printJSON(w, data)
	}

	if len(td.Rows) == 0 {
		fmt.Fprintln(w, "No results found.")
		return nil
	}

	// Apply field filtering
	columns := td.Columns
	if len(opts.Fields) > 0 {
		columns = filterColumns(columns, opts.Fields)
	}

	// Apply sorting
	if opts.Sort != "" {
		sortRows(td.Rows, opts.Sort)
	}

	// Calculate column widths
	widths := make(map[string]int)
	for _, col := range columns {
		widths[col.Name] = len(col.Name)
	}
	for _, row := range td.Rows {
		for _, col := range columns {
			if l := len(row[col.Name]); l > widths[col.Name] {
				widths[col.Name] = l
			}
		}
	}

	// Print header
	for i, col := range columns {
		if i > 0 {
			fmt.Fprint(w, "  ")
		}
		fmt.Fprintf(w, "%-*s", widths[col.Name], strings.ToUpper(col.Name))
	}
	fmt.Fprintln(w)

	// Print rows
	for _, row := range td.Rows {
		for i, col := range columns {
			if i > 0 {
				fmt.Fprint(w, "  ")
			}
			fmt.Fprintf(w, "%-*s", widths[col.Name], row[col.Name])
		}
		fmt.Fprintln(w)
	}

	return nil
}

func filterColumns(columns []Column, fields []string) []Column {
	fieldSet := make(map[string]bool)
	for _, f := range fields {
		fieldSet[strings.ToLower(f)] = true
	}
	var result []Column
	for _, col := range columns {
		if fieldSet[strings.ToLower(col.Name)] {
			result = append(result, col)
		}
	}
	return result
}

func sortRows(rows []map[string]string, field string) {
	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i][field] < rows[j][field]
	})
}

// IsTerminal returns true if the writer is a terminal.
func IsTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		stat, err := f.Stat()
		if err != nil {
			return false
		}
		return (stat.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// DefaultFormat returns the appropriate default format based on the writer.
func DefaultFormat(w io.Writer) Format {
	if IsTerminal(w) {
		return FormatTable
	}
	return FormatJSON
}

// NoColorEnabled checks if color should be disabled.
func NoColorEnabled() bool {
	return os.Getenv("NO_COLOR") != ""
}
