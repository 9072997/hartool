package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"hartool/internal"
)

var headersCmd = &cobra.Command{
	Use:   "headers <file.har> <index|url-substring>",
	Short: "Print both request and response headers for an entry",
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
		fmt.Println("=== Request ===")
		internal.PrintRequestHeaders(os.Stdout, entry.Request)
		fmt.Println("\n=== Response ===")
		internal.PrintResponseHeaders(os.Stdout, entry.Response)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(headersCmd)
}
