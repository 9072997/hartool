package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"hartool/har"
	"hartool/internal"
	"hartool/pretty"
)

var replayCmd = &cobra.Command{
	Use:   "replay <file.har> <index|url-substring>",
	Short: "Re-issue a captured request live and print the response",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := loadHAR(args[0])
		if err != nil {
			return err
		}
		entry, err := internal.ResolveEntry(h.Log.Entries, args[1])
		if err != nil {
			return err
		}

		req, err := buildRequest(entry.Request)
		if err != nil {
			return err
		}

		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Printf("%s %s\n", resp.Proto, resp.Status)
		for name, vals := range resp.Header {
			for _, v := range vals {
				fmt.Printf("%s: %s\n", name, v)
			}
		}
		fmt.Println()
		prettyFlag, _ := cmd.Flags().GetBool("pretty")
		if !prettyFlag {
			_, err = io.Copy(os.Stdout, resp.Body)
			return err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		mimeType := resp.Header.Get("Content-Type")
		if err := pretty.Print(mimeType, body, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, "pretty-print:", err)
			os.Exit(1)
		}
		return nil
	},
}

func buildRequest(r har.Request) (*http.Request, error) {
	var body io.Reader
	if r.PostData != nil && r.PostData.Text != "" {
		body = strings.NewReader(r.PostData.Text)
	}

	req, err := http.NewRequest(r.Method, r.URL, body)
	if err != nil {
		return nil, err
	}

	for _, h := range r.Headers {
		switch strings.ToLower(h.Name) {
		case "host", "content-length", "transfer-encoding":
			// skip headers managed by net/http
		default:
			req.Header.Add(h.Name, h.Value)
		}
	}

	return req, nil
}

func init() {
	replayCmd.Flags().BoolP("pretty", "p", false, "Pretty-print the response body")
	rootCmd.AddCommand(replayCmd)
}
