package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/adil-chbada/extract-cli/internal/config"
)

// ScanResult holds the results of scanning a project directory
type ScanResult struct {
	Code     []string
	Data     []string
	Locals   []string
	Total    int
	Excluded int
}

// Scan scans the project directory and categorizes files
func Scan(projectPath string, cfg *config.Config) (*ScanResult, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	result := &ScanResult{
		Code:   []string{},
		Data:   []string{},
		Locals: []string{},
	}

	// Load .gitignore patterns
	ignorer, err := loadGitignore(projectPath)
	if err != nil {
		// .gitignore is optional, continue without it
		ignorer = nil
	}

	err = filepath.WalkDir(projectPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == projectPath {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Convert to forward slashes for consistent pattern matching
		relPath = filepath.ToSlash(relPath)

		result.Total++

		// Skip directories (we only process files)
		if d.IsDir() {
			return nil
		}

		// Check if file should be ignored by .gitignore
		if ignorer != nil && ignorer.MatchesPath(relPath) {
			result.Excluded++
			return nil
		}

		// Check if file should be excluded by config patterns
		if cfg.IsExcluded(relPath) {
			result.Excluded++
			return nil
		}

		// Categorize the file
		switch {
		case cfg.IsDataFile(relPath):
			result.Data = append(result.Data, relPath)
		case cfg.IsLocalFile(relPath):
			// Check if it's a main local file (should go to code)
			if cfg.IsMainLocalFile(relPath) {
				result.Code = append(result.Code, relPath)
			} else {
				result.Locals = append(result.Locals, relPath)
			}
		default:
			result.Code = append(result.Code, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	return result, nil
}

// loadGitignore loads .gitignore patterns from the project directory
func loadGitignore(projectPath string) (*ignore.GitIgnore, error) {
	gitignorePath := filepath.Join(projectPath, ".gitignore")

	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		return nil, fmt.Errorf(".gitignore not found")
	}

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .gitignore: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Filter out empty lines and comments
	var patterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	if len(patterns) == 0 {
		return nil, fmt.Errorf(".gitignore is empty")
	}

	return ignore.CompileIgnoreLines(patterns...), nil
}