package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"hartool/internal"
)

var (
	searchRegex           bool
	searchCaseInsensitive bool
)

var searchCmd = &cobra.Command{
	Use:   "search <file.har> <term>",
	Short: "Search across all entries' URLs, headers, and bodies",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := loadHAR(args[0])
		if err != nil {
			return err
		}

		match, err := internal.BuildMatcher(args[1], searchRegex, searchCaseInsensitive)
		if err != nil {
			return err
		}

		for i, entry := range h.Log.Entries {
			if fields := internal.SearchEntry(entry, match); len(fields) > 0 {
				fmt.Printf("[%d] %s (%s)\n", i, entry.Request.URL, strings.Join(fields, ","))
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVarP(&searchRegex, "regex", "r", false, "treat term as a regular expression")
	searchCmd.Flags().BoolVarP(&searchCaseInsensitive, "case-insensitive", "i", false, "case-insensitive matching")
}
