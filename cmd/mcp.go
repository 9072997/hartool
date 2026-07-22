package cmd

import (
	"github.com/spf13/cobra"

	mcpserver "github.com/9072997/hartool/mcp"
)

var mcpCmd = &cobra.Command{
	Use:           "mcp",
	Short:         "Start MCP server (stdio transport)",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcpserver.Run()
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
