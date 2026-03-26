package internal

import (
	"fmt"
	"strconv"
	"strings"

	"hartool/har"
)

func ResolveEntry(entries []har.Entry, identifier string) (*har.Entry, error) {
	// Try integer index first
	if idx, err := strconv.Atoi(identifier); err == nil {
		if idx < 0 || idx >= len(entries) {
			return nil, fmt.Errorf("index %d out of range (0–%d)", idx, len(entries)-1)
		}
		return &entries[idx], nil
	}

	// Substring match
	var matches []int
	for i, e := range entries {
		if strings.Contains(e.Request.URL, identifier) {
			matches = append(matches, i)
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("no entries found matching %q", identifier)
	case 1:
		return &entries[matches[0]], nil
	default:
		var sb strings.Builder
		fmt.Fprintf(&sb, "multiple entries match %q; use an index instead:\n", identifier)
		for _, i := range matches {
			fmt.Fprintf(&sb, "  [%d] %s %s\n", i, entries[i].Request.Method, entries[i].Request.URL)
		}
		return nil, fmt.Errorf("%s", strings.TrimRight(sb.String(), "\n"))
	}
}
