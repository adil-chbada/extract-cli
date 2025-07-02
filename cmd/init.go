package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed templates/*.yaml
var templatesFS embed.FS

var (
	outputFile string
	listTemplates bool
)

var initCmd = &cobra.Command{
	Use:   "init [template]",
	Short: "Generate a config file from a built-in template",
	Long: `Create a configuration file based on pre-defined templates for popular frameworks.

Available templates:
- common: Universal template with common exclusions (.git, IDE files, OS files)
- go: Go projects with vendor/, bin/, build/ exclusions
- flutter: Flutter/Dart projects with .data.dart, .g.dart exclusions
- laravel: PHP Laravel projects with vendor/, storage/ exclusions
- vue: Vue.js projects with node_modules/, dist/ exclusions
- react: React projects with build/, node_modules/ exclusions
- nodejs: Node.js projects with standard npm exclusions
- python: Python projects with __pycache__/, .pyc exclusions`,
	Example: `  extract-cli init common
extract-cli init go
extract-cli init go -o my-go-config.yaml
extract-cli init flutter -o my-flutter-config.yaml
extract-cli init react --output my-react-config.yaml
extract-cli init --list`,
	Args: func(cmd *cobra.Command, args []string) error {
		if listTemplates {
			return nil
		}
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one template name")
		}
		return nil
	},
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path (defaults to extract.config.yml)")
	initCmd.Flags().BoolVarP(&listTemplates, "list", "l", false, "list available templates")
}

func runInit(cmd *cobra.Command, args []string) error {
	if listTemplates {
		return showAvailableTemplates()
	}

	// Use default filename if not specified
	if outputFile == "" {
		outputFile = "extract.config.yml"
	}

	templateName := args[0]
	templatePath := fmt.Sprintf("templates/%s.yaml", templateName)

	logInfo(fmt.Sprintf("Loading template: %s", templateName))

	templateData, err := templatesFS.ReadFile(templatePath)
	if err != nil {
		logError(fmt.Sprintf("Template '%s' not found", templateName))
		fmt.Println("\nAvailable templates:")
		showAvailableTemplates()
		return err
	}

	// Create output directory if it doesn't exist
	if dir := filepath.Dir(outputFile); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logError(fmt.Sprintf("Failed to create output directory: %v", err))
			return err
		}
	}

	logInfo(fmt.Sprintf("Writing config to: %s", outputFile))

	if err := os.WriteFile(outputFile, templateData, 0644); err != nil {
		logError(fmt.Sprintf("Failed to write config file: %v", err))
		return err
	}

	logSuccess(fmt.Sprintf("Config file created: %s", outputFile))
	fmt.Printf("Template: %s\n", infoColor(templateName))
	fmt.Printf("You can now edit the config and run: %s\n",
		infoColor(fmt.Sprintf("extract-cli generate -c %s", outputFile)))

	return nil
}

func showAvailableTemplates() error {
	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return err
	}

	fmt.Println("Available templates:")
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".yaml") {
			templateName := strings.TrimSuffix(entry.Name(), ".yaml")
			fmt.Printf("  - %s\n", infoColor(templateName))
		}
	}

	return nil
}