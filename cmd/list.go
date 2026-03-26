package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list <file.har>",
	Short: "List all entries in the HAR file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := loadHAR(args[0])
		if err != nil {
			return err
		}
		for i, e := range h.Log.Entries {
			fmt.Printf("[%d] %s %s (%d)\n", i, e.Request.Method, e.Request.URL, e.Response.Status)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
