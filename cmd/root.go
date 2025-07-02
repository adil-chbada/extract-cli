package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	version = "1.0.0"

	// Color functions
	successColor = color.New(color.FgGreen).SprintFunc()
	errorColor   = color.New(color.FgRed).SprintFunc()
	warnColor    = color.New(color.FgYellow).SprintFunc()
	infoColor    = color.New(color.FgCyan).SprintFunc()
)

var rootCmd = &cobra.Command{
	Use:   "extract-cli",
	Short: "A CLI tool to extract and categorize project files into markdown",
	Long: `Extract CLI is a powerful tool that scans your project directory,
categorizes files based on configurable patterns, and generates organized
markdown files for code, data, and local files.

It supports multiple project templates (Flutter, Laravel, Vue, React) and
respects .gitignore patterns for intelligent file filtering.`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Add subcommands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(completionCmd)
}

// Logging helpers
func logInfo(msg string) {
	if verbose {
		fmt.Fprintf(os.Stderr, "%s %s\n", infoColor("[INFO]"), msg)
	}
}

func logSuccess(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", successColor("[SUCCESS]"), msg)
}

func logWarn(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", warnColor("[WARN]"), msg)
}

func logError(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", errorColor("[ERROR]"), msg)
}