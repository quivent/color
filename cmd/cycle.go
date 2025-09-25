package cmd

import (
	"fmt"
	"os"

	"color/internal"

	"github.com/spf13/cobra"
)

var cycleCmd = &cobra.Command{
	Use:   "cycle [mode]",
	Short: "Cycle through color variations",
	Long: `Generate color variations based on current terminal color.
	
Available modes:
  hue_shift    - Shift the hue while keeping saturation/value (default)
  brightness   - Adjust brightness/value
  saturation   - Adjust color saturation
  complement   - Use complementary color
  random       - Random mode selection`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mode := "hue_shift"
		if len(args) > 0 {
			mode = args[0]
		}
		
		if err := cycleColors(mode); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// cycleColors implements the color cycling functionality
func cycleColors(mode ...string) error {
	selectedMode := "hue_shift"
	if len(mode) > 0 {
		selectedMode = mode[0]
	}
	
	cm := internal.NewColorManager()
	
	// Get current color
	current, err := cm.GetCurrentColor()
	if err != nil {
		return fmt.Errorf("failed to get current color: %w", err)
	}
	
	// Generate variant
	newColor := cm.GenerateVariant(current, selectedMode)
	
	// Set new color
	if err := cm.SetITermColor(newColor); err != nil {
		return fmt.Errorf("failed to set color: %w", err)
	}
	
	// More human-friendly messages
	var message string
	switch selectedMode {
	case "hue_shift":
		message = "🌈 Shifted to a new hue"
	case "brightness":
		message = "💡 Adjusted brightness"
	case "saturation":
		message = "🎨 Changed color intensity"
	case "complement":
		message = "🔄 Switched to complementary color"
	default:
		message = "✨ Applied color variation"
	}
	
	fmt.Printf("%s: RGB(%d, %d, %d)\n", 
		message, newColor.R, newColor.G, newColor.B)
	
	return nil
}

func init() {
	rootCmd.AddCommand(cycleCmd)
}