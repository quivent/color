package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"color/internal"

	"github.com/spf13/cobra"
)

var wrapperCmd = &cobra.Command{
	Use:   "wrap [command] [args...]",
	Short: "Wrap a command with automatic color management",
	Long: `Wrap a command execution with automatic color management.
	
This command:
1. Sets Claude theme colors before execution
2. Runs the specified command with all arguments
3. Restores directory-based colors after execution
4. Preserves the exit code of the wrapped command

Example:
  color wrap claude --help
  color wrap claude code --session mysession`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := wrapCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			if exitError, ok := err.(*exec.ExitError); ok {
				if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
					os.Exit(status.ExitStatus())
				}
			}
			os.Exit(1)
		}
	},
}

// wrapCommand implements command wrapping with color management
func wrapCommand(args []string) error {
	cm := internal.NewColorManager()
	
	// Set Claude session colors
	claudeColor := cm.GenerateClaudeTheme()
	if err := cm.SetITermColor(claudeColor); err != nil {
		return fmt.Errorf("failed to set Claude theme: %w", err)
	}
	
	// Execute the wrapped command
	command := args[0]
	commandArgs := args[1:]
	
	execCmd := exec.Command(command, commandArgs...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin
	
	err := execCmd.Run()
	
	// Always try to restore directory colors, even if command failed
	cwd, cwdErr := os.Getwd()
	if cwdErr == nil {
		dirColor := cm.GenerateDirectoryTheme(cwd)
		if restoreErr := cm.SetITermColor(dirColor); restoreErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to restore directory colors: %v\n", restoreErr)
		}
	}
	
	return err
}

func init() {
	rootCmd.AddCommand(wrapperCmd)
}