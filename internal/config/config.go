package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration for file extraction
type Config struct {
	ProjectName     string   `yaml:"project_name"`
	ProjectPath     string   `yaml:"project_path"`
	DataPatterns    []string `yaml:"data_patterns"`
	LocalPatterns   []string `yaml:"local_patterns"`
	ExcludePatterns []string `yaml:"exclude_patterns"`
	MainLocalFiles  []string `yaml:"main_local_files"`
	UseRegex        bool     `yaml:"use_regex"`
}

// getCommonExclusions returns common exclusion patterns that should be applied to all projects
func getCommonExclusions() []string {
	return []string{
		// Version control
		".git/**",
		".svn/**",
		".hg/**",
		".bzr/**",
		// IDE and editor files
		".vscode/**",
		".idea/**",
		"*.swp",
		"*.swo",
		"*~",
		".DS_Store",
		"Thumbs.db",
		// OS generated files
		"desktop.ini",
		".Spotlight-V100",
		".Trashes",
		"ehthumbs.db",
		// Temporary files
		"*.tmp",
		"*.temp",
		"*.bak",
		"*.backup",
		// Logs
		"*.log",
		"logs/**",
		// Cache and build artifacts
		"*.cache",
		".cache/**",
		"build/**",
		"dist/**",
		"out/**",
		"target/**",
		// Environment and sensitive files
		".env",
		".env.*",
		"*.key",
		"*.pem",
		"*.p12",
		"*.pfx",
		// Archives and binaries
		"*.zip",
		"*.tar",
		"*.tar.gz",
		"*.rar",
		"*.7z",
		"*.exe",
		"*.dll",
		"*.so",
		"*.dylib",
		// Coverage reports
		"coverage/**",
		"*.coverage",
		"lcov.info",
	}
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Load and merge common exclusions
	commonExclusions := getCommonExclusions()
	// Merge common exclusions with project-specific ones
	// Add common exclusions first, then project-specific ones
	cfg.ExcludePatterns = append(commonExclusions, cfg.ExcludePatterns...)

	// Set defaults
	if cfg.ProjectPath == "" {
		cfg.ProjectPath = "."
	}
	if cfg.MainLocalFiles == nil {
		cfg.MainLocalFiles = []string{"main.*", "index.*", "app.*"}
	}

	// Convert relative path to absolute
	if !filepath.IsAbs(cfg.ProjectPath) {
		abs, err := filepath.Abs(cfg.ProjectPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve project path: %w", err)
		}
		cfg.ProjectPath = abs
	}

	return &cfg, nil
}

// IsDataFile checks if a file matches data patterns
func (c *Config) IsDataFile(path string) bool {
	return c.matchesPatterns(path, c.DataPatterns)
}

// IsLocalFile checks if a file matches local patterns
func (c *Config) IsLocalFile(path string) bool {
	return c.matchesPatterns(path, c.LocalPatterns)
}

// IsMainLocalFile checks if a file is a main local file
func (c *Config) IsMainLocalFile(path string) bool {
	return c.matchesPatterns(path, c.MainLocalFiles)
}

// IsExcluded checks if a file should be excluded
func (c *Config) IsExcluded(path string) bool {
	return c.matchesPatterns(path, c.ExcludePatterns)
}

// matchesPatterns checks if a path matches any of the given patterns
func (c *Config) matchesPatterns(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if c.matchesPattern(path, pattern) {
			return true
		}
	}
	return false
}

// matchesPattern checks if a path matches a single pattern
func (c *Config) matchesPattern(path, pattern string) bool {
	// Handle regex patterns (prefix with 're:')
	if c.UseRegex && strings.HasPrefix(pattern, "re:") {
		regexPattern := strings.TrimPrefix(pattern, "re:")
		matched, err := regexp.MatchString(regexPattern, path)
		if err != nil {
			return false
		}
		return matched
	}

	// Handle glob patterns
	matched, err := filepath.Match(pattern, filepath.Base(path))
	if err != nil {
		return false
	}
	if matched {
		return true
	}

	// Handle directory patterns (e.g., "data/**")
	if strings.Contains(pattern, "**") {
		// Convert glob pattern to regex-like matching
		parts := strings.Split(pattern, "**")
		if len(parts) == 2 {
			prefix := strings.TrimSuffix(parts[0], "/")
			suffix := strings.TrimPrefix(parts[1], "/")

			// Check if path starts with prefix and ends with suffix
			if prefix != "" && !strings.HasPrefix(path, prefix) {
				return false
			}
			if suffix != "" && !strings.HasSuffix(path, suffix) {
				return false
			}
			return true
		}
	}

	// Check if the pattern matches the full path
	matched, err = filepath.Match(pattern, path)
	return err == nil && matched
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.ProjectPath == "" {
		return fmt.Errorf("project_path is required")
	}

	if _, err := os.Stat(c.ProjectPath); os.IsNotExist(err) {
		return fmt.Errorf("project_path does not exist: %s", c.ProjectPath)
	}

	return nil
}