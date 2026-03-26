package pretty

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

func indentXML(body []byte, w io.Writer) error {
	decoder := xml.NewDecoder(bytes.NewReader(body))
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if cd, ok := tok.(xml.CharData); ok {
			trimmed := strings.TrimSpace(string(cd))
			if trimmed == "" {
				continue
			}
			tok = xml.CharData([]byte(trimmed))
		}

		if err := encoder.EncodeToken(tok); err != nil {
			return err
		}
	}

	if err := encoder.Flush(); err != nil {
		return err
	}
	fmt.Fprintln(w)
	return nil
}

var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"param": true, "source": true, "track": true, "wbr": true,
}

func indentHTML(body []byte, w io.Writer) error {
	z := html.NewTokenizer(bytes.NewReader(body))
	depth := 0
	inRawText := false

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		tok := z.Token()
		indent := strings.Repeat("  ", depth)

		switch tt {
		case html.DoctypeToken:
			fmt.Fprintf(w, "<!DOCTYPE %s>\n", tok.Data)

		case html.StartTagToken:
			tagName := tok.Data
			fmt.Fprintf(w, "%s%s\n", indent, renderHTMLTag(tok))
			if tagName == "script" || tagName == "style" {
				inRawText = true
			}
			if !voidElements[tagName] {
				depth++
			}

		case html.EndTagToken:
			if tok.Data == "script" || tok.Data == "style" {
				inRawText = false
			}
			if depth > 0 {
				depth--
			}
			fmt.Fprintf(w, "%s</%s>\n", strings.Repeat("  ", depth), tok.Data)

		case html.SelfClosingTagToken:
			fmt.Fprintf(w, "%s%s\n", indent, renderHTMLTag(tok))

		case html.TextToken:
			if inRawText {
				lines := strings.Split(tok.Data, "\n")
				lineIndent := strings.Repeat("  ", depth)
				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					if trimmed != "" {
						fmt.Fprintf(w, "%s%s\n", lineIndent, trimmed)
					}
				}
			} else {
				text := strings.TrimSpace(tok.Data)
				if text != "" {
					fmt.Fprintf(w, "%s%s\n", indent, html.EscapeString(text))
				}
			}

		case html.CommentToken:
			fmt.Fprintf(w, "%s<!-- %s -->\n", indent, strings.TrimSpace(tok.Data))
		}
	}

	return nil
}

func renderHTMLTag(tok html.Token) string {
	var sb strings.Builder
	sb.WriteString("<")
	sb.WriteString(tok.Data)
	for _, attr := range tok.Attr {
		if attr.Val == "" {
			fmt.Fprintf(&sb, " %s", attr.Key)
		} else {
			fmt.Fprintf(&sb, ` %s="%s"`, attr.Key, html.EscapeString(attr.Val))
		}
	}
	sb.WriteString(">")
	return sb.String()
}
