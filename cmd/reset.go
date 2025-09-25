package cmd

import (
	"fmt"
	"os"

	"color/internal"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset terminal to default dark theme",
	Long: `Reset the terminal background to a default dark theme.
	
This command restores the terminal to a standard dark background
color suitable for general terminal use.`,
	Run: func(cmd *cobra.Command, args []string) {
		cm := internal.NewColorManager()
		
		// Default dark theme
		defaultColor := internal.RGB{R: 30, G: 30, B: 30}
		
		if err := cm.SetITermColor(defaultColor); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting color: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("🔄 Reset to default dark theme: RGB(%d, %d, %d)\n", 
			defaultColor.R, defaultColor.G, defaultColor.B)
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}