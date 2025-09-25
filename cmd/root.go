package cmd

import (
	"fmt"
	"os"

	"color/internal"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "color",
	Short: "Terminal color management CLI",
	Long: `A fast and flexible terminal color management tool.
	
This CLI provides automatic terminal color changes based on:
- Directory changes (consistent colors per directory)
- Claude Code sessions (blue/purple themes)
- Manual color cycling and themes

Built with Go and Cobra for speed and reliability.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action - apply directory color and show help
		cm := internal.NewColorManager()
		color := cm.GenerateDirectoryTheme("")
		
		if err := cm.SetITermColor(color); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting color: %v\n", err)
		} else {
			cwd, _ := os.Getwd()
			fmt.Printf("📁 Applied color for %s: RGB(%d, %d, %d)\n", 
				cwd, color.R, color.G, color.B)
		}
		
		fmt.Println()
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
}