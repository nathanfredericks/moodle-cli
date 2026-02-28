package text

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reBlockElements = regexp.MustCompile(`(?i)<br\s*/?>|</p>|</div>|</li>|</tr>`)
	reListItem      = regexp.MustCompile(`(?i)<li[^>]*>`)
	reHTMLTag       = regexp.MustCompile(`<[^>]*>`)
	reSpaces        = regexp.MustCompile(`[ \t]+`)
	reNewlines      = regexp.MustCompile(`\n{3,}`)
)

// StripHTML removes HTML tags and decodes common entities, preserving
// basic structure (newlines for block elements, bullets for list items).
func StripHTML(s string) string {
	// Replace block-level closing tags with newlines
	s = reBlockElements.ReplaceAllString(s, "\n")
	// Replace <li> with bullet
	s = reListItem.ReplaceAllString(s, "- ")
	// Strip remaining tags
	s = reHTMLTag.ReplaceAllString(s, "")
	// Decode common HTML entities
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	// Clean up whitespace
	s = reSpaces.ReplaceAllString(s, " ")
	s = reNewlines.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}

// FormatFileSize formats a byte count as a human-readable string.
func FormatFileSize(bytes int64) string {
	switch {
	case bytes >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(1<<30))
	case bytes >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(1<<10))
	default:
		return strconv.FormatInt(bytes, 10) + " B"
	}
}

// Truncate truncates a string to max characters, appending "..." if truncated.
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
