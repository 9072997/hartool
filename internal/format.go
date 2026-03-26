package internal

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"hartool/har"
)

// PrintRequestHeaders prints request line + headers with cookie section.
func PrintRequestHeaders(w io.Writer, req har.Request) {
	fmt.Fprintf(w, "%s %s %s\n", req.Method, req.URL, req.HTTPVersion)
	for _, h := range req.Headers {
		lower := strings.ToLower(h.Name)
		if lower == "cookie" || lower == "set-cookie" {
			continue
		}
		fmt.Fprintf(w, "%s: %s\n", h.Name, h.Value)
	}
	if len(req.Cookies) > 0 {
		fmt.Fprintln(w, "\nCookies:")
		for _, c := range req.Cookies {
			fmt.Fprintf(w, "  %s=%s\n", c.Name, c.Value)
		}
	}
}

// PrintResponseHeaders prints status line + headers with Set-Cookie section.
func PrintResponseHeaders(w io.Writer, resp har.Response) {
	fmt.Fprintf(w, "%s %d %s\n", resp.HTTPVersion, resp.Status, resp.StatusText)
	for _, h := range resp.Headers {
		lower := strings.ToLower(h.Name)
		if lower == "cookie" || lower == "set-cookie" {
			continue
		}
		fmt.Fprintf(w, "%s: %s\n", h.Name, h.Value)
	}
	if len(resp.Cookies) > 0 {
		fmt.Fprintln(w, "\nSet-Cookie:")
		for _, c := range resp.Cookies {
			fmt.Fprintf(w, "  %s", formatSetCookie(c))
		}
	}
}

func formatSetCookie(c har.Cookie) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s=%s", c.Name, c.Value)
	if c.Path != "" {
		fmt.Fprintf(&sb, "; Path=%s", c.Path)
	}
	if c.Domain != "" {
		fmt.Fprintf(&sb, "; Domain=%s", c.Domain)
	}
	if c.Expires != "" {
		fmt.Fprintf(&sb, "; Expires=%s", c.Expires)
	}
	if c.Secure {
		sb.WriteString("; Secure")
	}
	if c.HTTPOnly {
		sb.WriteString("; HttpOnly")
	}
	sb.WriteString("\n")
	return sb.String()
}

// WriteResponseBody writes the response body, decoding base64 if needed.
func WriteResponseBody(w io.Writer, content har.Content) error {
	if content.Text == "" {
		fmt.Fprintln(w, "(no response body)")
		return nil
	}
	if strings.EqualFold(content.Encoding, "base64") {
		decoded, err := base64.StdEncoding.DecodeString(content.Text)
		if err != nil {
			return fmt.Errorf("base64 decode: %w", err)
		}
		_, err = w.Write(decoded)
		return err
	}
	_, err := w.Write([]byte(content.Text))
	return err
}

// WriteRequestBody writes the request body.
func WriteRequestBody(w io.Writer, postData *har.PostData) error {
	if postData == nil || postData.Text == "" {
		fmt.Fprintln(w, "(no request body)")
		return nil
	}
	_, err := w.Write([]byte(postData.Text))
	return err
}

// DecodeResponseBody returns decoded body bytes + MIME type (handles base64).
// Returns error if content.Text is empty.
func DecodeResponseBody(content har.Content) (body []byte, mimeType string, err error) {
	if content.Text == "" {
		return nil, "", fmt.Errorf("(no response body)")
	}
	mimeType = content.MimeType
	if strings.EqualFold(content.Encoding, "base64") {
		decoded, err := base64.StdEncoding.DecodeString(content.Text)
		if err != nil {
			return nil, "", fmt.Errorf("base64 decode: %w", err)
		}
		return decoded, mimeType, nil
	}
	return []byte(content.Text), mimeType, nil
}

// DecodeRequestBody returns body bytes + MIME type.
// Returns error if postData is nil or Text is empty.
func DecodeRequestBody(postData *har.PostData) (body []byte, mimeType string, err error) {
	if postData == nil || postData.Text == "" {
		return nil, "", fmt.Errorf("(no request body)")
	}
	return []byte(postData.Text), postData.MimeType, nil
}
