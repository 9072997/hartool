package internal

import (
	"fmt"
	"regexp"
	"strings"

	"hartool/har"
)

// BuildMatcher returns a function that matches strings against the given term.
func BuildMatcher(term string, useRegex, caseInsensitive bool) (func(string) bool, error) {
	if useRegex {
		pattern := term
		if caseInsensitive {
			pattern = "(?i)" + term
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex: %w", err)
		}
		return re.MatchString, nil
	}
	if caseInsensitive {
		lowerTerm := strings.ToLower(term)
		return func(s string) bool {
			return strings.Contains(strings.ToLower(s), lowerTerm)
		}, nil
	}
	return func(s string) bool {
		return strings.Contains(s, term)
	}, nil
}

// SearchEntry returns the fields in which the match function found a hit.
func SearchEntry(entry har.Entry, match func(string) bool) []string {
	var fields []string

	if match(entry.Request.URL) {
		fields = append(fields, "url")
	}

	for _, h := range entry.Request.Headers {
		if match(h.Name) || match(h.Value) {
			fields = append(fields, "req-headers")
			break
		}
	}

	for _, h := range entry.Response.Headers {
		if match(h.Name) || match(h.Value) {
			fields = append(fields, "resp-headers")
			break
		}
	}

	if entry.Request.PostData != nil && match(entry.Request.PostData.Text) {
		fields = append(fields, "req-body")
	}

	body, _, err := DecodeResponseBody(entry.Response.Content)
	if err == nil && match(string(body)) {
		fields = append(fields, "resp-body")
	}

	return fields
}
