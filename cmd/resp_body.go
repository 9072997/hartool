package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"hartool/internal"
	"hartool/pretty"
)

var respBodyCmd = &cobra.Command{
	Use:   "resp-body <file.har> <index|url-substring>",
	Short: "Print response body for an entry",
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
		prettyFlag, _ := cmd.Flags().GetBool("pretty")
		if !prettyFlag {
			return internal.WriteResponseBody(os.Stdout, entry.Response.Content)
		}
		body, mimeType, err := internal.DecodeResponseBody(entry.Response.Content)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := pretty.Print(mimeType, body, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, "pretty-print:", err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	respBodyCmd.Flags().BoolP("pretty", "p", false, "Pretty-print the body (JSON, XML, HTML, YAML, URL-encoded, multipart/form-data)")
	rootCmd.AddCommand(respBodyCmd)
}
