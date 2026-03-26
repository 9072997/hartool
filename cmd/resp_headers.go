package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"hartool/internal"
)

var respHeadersCmd = &cobra.Command{
	Use:   "resp-headers <file.har> <index|url-substring>",
	Short: "Print response headers for an entry",
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
		internal.PrintResponseHeaders(os.Stdout, entry.Response)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(respHeadersCmd)
}
