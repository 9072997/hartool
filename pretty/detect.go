package pretty

import (
	"bytes"
	"strings"
)

// Format represents a detected content format.
type Format int

const (
	FormatUnknown    Format = iota
	FormatJSON
	FormatXML
	FormatHTML
	FormatYAML
	FormatURLEncoded
	FormatMultipart
)

// Detect returns format from MIME type, falling back to content sniffing.
func Detect(mimeType string, body []byte) Format {
	if f := fromMIME(mimeType); f != FormatUnknown {
		return f
	}
	return sniff(body)
}

func fromMIME(mimeType string) Format {
	if idx := strings.Index(mimeType, ";"); idx >= 0 {
		mimeType = mimeType[:idx]
	}
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))

	switch {
	case mimeType == "application/json" || mimeType == "text/json" || strings.HasSuffix(mimeType, "+json"):
		return FormatJSON
	case mimeType == "application/xml" || mimeType == "text/xml" || strings.HasSuffix(mimeType, "+xml"):
		return FormatXML
	case mimeType == "text/html":
		return FormatHTML
	case mimeType == "application/yaml" || mimeType == "text/yaml" || mimeType == "text/x-yaml" || mimeType == "application/x-yaml":
		return FormatYAML
	case mimeType == "application/x-www-form-urlencoded":
		return FormatURLEncoded
	case strings.HasPrefix(mimeType, "multipart/form-data"):
		return FormatMultipart
	}
	return FormatUnknown
}

func sniff(body []byte) Format {
	b := bytes.TrimSpace(body)
	if len(b) == 0 {
		return FormatUnknown
	}

	switch b[0] {
	case '{', '[':
		return FormatJSON
	case '<':
		limit := 512
		if len(b) < limit {
			limit = len(b)
		}
		upper := strings.ToUpper(string(b[:limit]))
		if strings.Contains(upper, "<!DOCTYPE HTML") || strings.HasPrefix(upper, "<HTML") {
			return FormatHTML
		}
		return FormatXML
	}

	if bytes.HasPrefix(b, []byte("---")) {
		return FormatYAML
	}

	return FormatUnknown
}
