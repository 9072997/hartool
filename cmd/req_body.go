package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"hartool/internal"
	"hartool/pretty"
)

var reqBodyCmd = &cobra.Command{
	Use:   "req-body <file.har> <index|url-substring>",
	Short: "Print request body for an entry",
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
			return internal.WriteRequestBody(os.Stdout, entry.Request.PostData)
		}
		body, mimeType, err := internal.DecodeRequestBody(entry.Request.PostData)
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
	reqBodyCmd.Flags().BoolP("pretty", "p", false, "Pretty-print the body (JSON, XML, HTML, YAML, URL-encoded, multipart/form-data)")
	rootCmd.AddCommand(reqBodyCmd)
}
