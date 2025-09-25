package cmd

import (
	"fmt"

	"color/internal"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show color persistence and Redis connection status",
	Long: `Display the current status of the color persistence system.
	
This command shows:
- Redis connection status
- Number of stored directory colors
- Last Claude theme usage
- Persistence configuration details`,
	Run: func(cmd *cobra.Command, args []string) {
		cm := internal.NewColorManager()
		
		// Get persistence status
		status := cm.GetPersistenceStatus()
		fmt.Println(status)
		
		// Show color history if available
		if history, err := cm.GetColorHistory(5); err == nil && len(history) > 0 {
			fmt.Println("\n📊 Recent Color History:")
			for i, entry := range history {
				fmt.Printf("%d. RGB(%d, %d, %d) - %s (%s)\n", 
					i+1, entry.Color.R, entry.Color.G, entry.Color.B, 
					entry.Source, entry.Timestamp.Format("2006-01-02 15:04"))
			}
		}
	},
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all stored color data",
	Long: `Remove all stored colors from the Redis cache.
	
This will clear:
- All directory color associations
- Last Claude theme color
- Color history data
	
Colors will be regenerated on next use.`,
	Run: func(cmd *cobra.Command, args []string) {
		cm := internal.NewColorManager()
		
		if err := cm.ClearColorCache(); err != nil {
			fmt.Printf("❌ Error clearing color cache: %v\n", err)
			return
		}
		
		fmt.Println("🗑️ Cleared all stored color data")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(clearCmd)
}