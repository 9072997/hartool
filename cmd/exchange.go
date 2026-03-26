package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"hartool/internal"
	"hartool/pretty"
)

var exchangeCmd = &cobra.Command{
	Use:   "exchange <file.har> <index|url-substring>",
	Short: "Print the full request/response exchange for an entry",
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

		// 1. Request headers
		internal.PrintRequestHeaders(os.Stdout, entry.Request)
		fmt.Println()

		// 2. Request body
		if prettyFlag {
			if body, mimeType, err := internal.DecodeRequestBody(entry.Request.PostData); err == nil {
				if err2 := pretty.Print(mimeType, body, os.Stdout); err2 != nil {
					fmt.Fprintln(os.Stderr, "pretty-print:", err2)
					os.Exit(1)
				}
			} else {
				internal.WriteRequestBody(os.Stdout, entry.Request.PostData)
			}
		} else {
			internal.WriteRequestBody(os.Stdout, entry.Request.PostData)
		}
		fmt.Println()

		// 3. Response headers
		internal.PrintResponseHeaders(os.Stdout, entry.Response)
		fmt.Println()

		// 4. Response body
		if prettyFlag {
			if body, mimeType, err := internal.DecodeResponseBody(entry.Response.Content); err == nil {
				if err2 := pretty.Print(mimeType, body, os.Stdout); err2 != nil {
					fmt.Fprintln(os.Stderr, "pretty-print:", err2)
					os.Exit(1)
				}
			} else {
				internal.WriteResponseBody(os.Stdout, entry.Response.Content)
			}
		} else {
			internal.WriteResponseBody(os.Stdout, entry.Response.Content)
		}
		return nil
	},
}

func init() {
	exchangeCmd.Flags().BoolP("pretty", "p", false, "Pretty-print request and response bodies")
	rootCmd.AddCommand(exchangeCmd)
}
