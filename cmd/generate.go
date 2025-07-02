package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/adil-chbada/extract-cli/internal/config"
	"github.com/adil-chbada/extract-cli/internal/markdown"
	"github.com/adil-chbada/extract-cli/internal/scanner"
)

var (
	configPath string
	outputDir  string
)

// Default config file names to search for (in order of preference)
var defaultConfigFiles = []string{
	"extract.config.yml",
	"extract.config.yaml",
	"extract-config.yaml",
	"extract-config.yml",
	".extract-config.yaml",
	".extract-config.yml",
	"extract.yaml",
	"extract.yml",
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate markdown files from project based on config",
	Long: `Scan your project directory and generate three markdown files:
- project-code.md: All code files and main local files
- project-data.md: Data files only (*.data.dart, *.json, /data/**, etc.)
- project-locals.md: All other local/configuration files

The tool respects .gitignore patterns and custom exclude patterns from your config.

If no config file is specified, the tool will automatically search for default
config files in the following order: extract.config.yml, extract.config.yaml,
extract-config.yaml, extract-config.yml, .extract-config.yaml, .extract-config.yml,
extract.yaml, extract.yml`,
	Example: `  extract-cli generate
  extract-cli generate -c config.yaml
  extract-cli generate -c flutter-config.yaml -o ./output
  extract-cli generate --config myproject.yaml --output-dir ./docs`,
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file (if not specified, searches for default config files)")
	generateCmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "output directory for markdown files")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// If no config path specified, search for default config files
	if configPath == "" {
		foundConfig, err := findDefaultConfig()
		if err != nil {
			logError(fmt.Sprintf("No config file found. Please specify one with -c flag or create one of: %v", defaultConfigFiles))
			return err
		}
		configPath = foundConfig
	}

	logInfo(fmt.Sprintf("Loading config from: %s", configPath))

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logError(fmt.Sprintf("Failed to load config: %v", err))
		return err
	}

	logInfo(fmt.Sprintf("Scanning project directory: %s", cfg.ProjectPath))

	result, err := scanner.Scan(cfg.ProjectPath, cfg)
	if err != nil {
		logError(fmt.Sprintf("Failed to scan project: %v", err))
		return err
	}

	// Write markdown files
	files := map[string]struct {
		title string
		items []string
	}{
		"project-code.md": {
			title: "Project Code Files",
			items: result.Code,
		},
		"project-data.md": {
			title: "Project Data Files",
			items: result.Data,
		},
		"project-locals.md": {
			title: "Project Local Files",
			items: result.Locals,
		},
	}

	for filename, data := range files {
		outputPath := filepath.Join(outputDir, filename)
		logInfo(fmt.Sprintf("Writing %s (%d files)", outputPath, len(data.items)))

		if err := markdown.WriteMarkdown(outputPath, data.title, data.items, cfg); err != nil {
			logError(fmt.Sprintf("Failed to write %s: %v", filename, err))
			return err
		}
	}

	// Calculate total sizes for each category
	codeSize := calculateTotalSize(result.Code, cfg.ProjectPath)
	dataSize := calculateTotalSize(result.Data, cfg.ProjectPath)
	localsSize := calculateTotalSize(result.Locals, cfg.ProjectPath)
	totalSize := codeSize + dataSize + localsSize
	
	// Print summary
	fmt.Printf("\n%s\n", successColor("✓ Generation completed successfully!"))
	fmt.Printf("Total files scanned: %d (%s)\n", result.Total, formatFileSize(totalSize))
	fmt.Printf("├─ Code files: %d (%s)\n", len(result.Code), formatFileSize(codeSize))
	fmt.Printf("├─ Data files: %d (%s)\n", len(result.Data), formatFileSize(dataSize))
	fmt.Printf("├─ Local files: %d (%s)\n", len(result.Locals), formatFileSize(localsSize))
	fmt.Printf("└─ Excluded files: %d\n", result.Excluded)
	fmt.Printf("\nMarkdown files written to: %s\n", outputDir)

	return nil
}

// formatFileSize formats file size in human readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// calculateTotalSize calculates the total size of a list of files
func calculateTotalSize(files []string, projectPath string) int64 {
	totalSize := int64(0)
	for _, filePath := range files {
		info, err := os.Stat(filepath.Join(projectPath, filePath))
		if err == nil {
			totalSize += info.Size()
		}
	}
	return totalSize
}

// findDefaultConfig searches for default config files in the current directory
func findDefaultConfig() (string, error) {
	for _, filename := range defaultConfigFiles {
		if _, err := os.Stat(filename); err == nil {
			return filename, nil
		}
	}
	return "", fmt.Errorf("no default config file found")
}