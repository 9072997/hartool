package pretty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"sort"
	"unicode"

	"gopkg.in/yaml.v3"
)

// Print detects format and pretty-prints body to w.
// Returns error if format is unknown or body fails to parse.
func Print(mimeType string, body []byte, w io.Writer) error {
	format := Detect(mimeType, body)

	switch format {
	case FormatJSON:
		return printJSON(body, w)
	case FormatXML:
		return indentXML(body, w)
	case FormatHTML:
		return indentHTML(body, w)
	case FormatYAML:
		return printYAML(body, w)
	case FormatURLEncoded:
		return printURLEncoded(body, w)
	case FormatMultipart:
		return printMultipart(mimeType, body, w)
	default:
		return fmt.Errorf("cannot detect format")
	}
}

func printJSON(body []byte, w io.Writer) error {
	var v interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func printYAML(body []byte, w io.Writer) error {
	var v interface{}
	if err := yaml.Unmarshal(body, &v); err != nil {
		return err
	}
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	return enc.Encode(v)
}

func printURLEncoded(body []byte, w io.Writer) error {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range values[k] {
			fmt.Fprintf(w, "%s = %s\n", k, v)
		}
	}
	return nil
}

func printMultipart(mimeType string, body []byte, w io.Writer) error {
	_, params, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return fmt.Errorf("parse media type: %w", err)
	}
	boundary, ok := params["boundary"]
	if !ok {
		return fmt.Errorf("no boundary in multipart content-type")
	}

	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	for {
		part, err := reader.NextPart()
		if err != nil {
			break
		}

		for name, values := range part.Header {
			for _, v := range values {
				fmt.Fprintf(w, "%s: %s\n", name, v)
			}
		}
		fmt.Fprintln(w)

		content, readErr := io.ReadAll(part)
		if readErr != nil {
			return readErr
		}

		if isPrintable(content) {
			w.Write(content)
		} else {
			fmt.Fprintf(w, "[%d bytes binary data]\n", len(content))
		}
		fmt.Fprintln(w)
	}
	return nil
}

func isPrintable(b []byte) bool {
	for _, r := range string(b) {
		if !unicode.IsPrint(r) && r != '\n' && r != '\r' && r != '\t' {
			return false
		}
	}
	return true
}
