package cmd

import (
	"fmt"
	"os"

	"color/internal"

	"github.com/spf13/cobra"
)

var directoryCmd = &cobra.Command{
	Use:   "directory [path]",
	Short: "Apply directory-based consistent color theme",
	Long: `Apply a color theme based on the directory path hash.
	
Each directory gets a unique, consistent color based on its path.
The same directory will always produce the same color, providing
visual consistency for navigation.

If no path is provided, uses current working directory.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) > 0 {
			path = args[0]
		}
		
		cm := internal.NewColorManager()
		color := cm.GenerateDirectoryTheme(path)
		
		if err := cm.SetITermColor(color); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting color: %v\n", err)
			os.Exit(1)
		}
		
		actualPath := path
		if actualPath == "" {
			var err error
			actualPath, err = os.Getwd()
			if err != nil {
				actualPath = "current directory"
			}
		}
		
		fmt.Printf("📁 Applied color for %s: RGB(%d, %d, %d)\n", 
			actualPath, color.R, color.G, color.B)
	},
}

func init() {
	rootCmd.AddCommand(directoryCmd)
}