package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"hartool/har"
)

var rootCmd = &cobra.Command{
	Use:   "hartool",
	Short: "Inspect HAR (HTTP Archive) files",
	Long: `hartool inspects HAR files captured by browsers or proxies.

The command word may appear in any position:
  hartool capture.har 7 resp-body
  hartool resp-body capture.har 7
  hartool capture.har resp-body 7

Default commands (no explicit command word):
  hartool capture.har          -> list
  hartool capture.har 7        -> exchange

License: GPLv2 or later <https://www.gnu.org/licenses/gpl-2.0.html>`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func loadHAR(path string) (*har.HAR, error) {
	return har.Load(path)
}
