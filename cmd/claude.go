package cmd

import (
	"fmt"
	"os"

	"color/internal"

	"github.com/spf13/cobra"
)

var claudeCmd = &cobra.Command{
	Use:   "claude",
	Short: "Apply Claude Code session theme",
	Long: `Apply a Claude Code specific color theme with blue/purple hues.
	
This command sets terminal colors optimized for Claude Code sessions,
using a palette of blues and purples with appropriate contrast for
terminal readability.`,
	Run: func(cmd *cobra.Command, args []string) {
		cm := internal.NewColorManager()
		color := cm.GenerateClaudeTheme()
		
		if err := cm.SetITermColor(color); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting color: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("🤖 Applied Claude Code theme: RGB(%d, %d, %d)\n", color.R, color.G, color.B)
	},
}

func init() {
	rootCmd.AddCommand(claudeCmd)
}