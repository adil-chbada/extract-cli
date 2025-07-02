# Extract CLI 🚀

A powerful Go-based CLI tool designed to help developers share their code with AI assistants efficiently. It intelligently scans your project directory, categorizes files based on configurable patterns, and generates organized markdown files that can be easily shared as a single consolidated document with AI tools like ChatGPT, Claude, or GitHub Copilot.

## ✨ Key Features

- **🤖 AI-Ready Output**: Generates clean, organized markdown files perfect for sharing with AI assistants
- **📊 File Size Analytics**: Displays individual file sizes and total size information for better project insights
- **🎯 Smart Categorization**: Automatically categorizes files into code, data, and configuration files
- **📋 Template Library**: Built-in templates for popular frameworks (Flutter, React, Vue, Node.js, Laravel, Python)
- **🔒 Universal Exclusions**: Automatically excludes `.git`, IDE files, OS files, and other common artifacts
- **🌍 i18n Support**: Intelligent handling of internationalization and localization files
- **🎨 Beautiful Output**: Generates well-structured markdown with metadata, file sizes, and summaries
- **⚡ Cross-Platform**: Works seamlessly on Linux, macOS, and Windows
- **🎨 Rich CLI**: Colorful terminal output with progress indicators and size information

## 📦 Installation

### From Source
```bash
git clone https://github.com/adil-chbada/extract-cli.git
cd extract-cli
make build
```

### Using Go Install
```bash
go install github.com/adil-chbada/extract-cli@latest
```

## 🚀 Quick Start

1. **Create a configuration file**:
   ```bash
   extract-cli init react -o my-config.yaml
   ```

2. **Generate AI-ready documentation**:
   ```bash
   extract-cli generate -c my-config.yaml
   ```

3. **Share with AI** - Copy the content from generated markdown files and paste directly into ChatGPT, Claude, or any AI assistant for code review, debugging, or development assistance:
   - `project-code.md` - Source code files with individual and total sizes
   - `project-data.md` - Data files with size analytics
   - `project-locals.md` - Configuration files with size breakdown

## 💡 Usage

### Commands

#### `init` - Generate Configuration Template
```bash
extract-cli init [template] [flags]

# Examples
extract-cli init react -o react-config.yaml
extract-cli init flutter -o flutter-config.yaml
extract-cli init common -o common-exclusions.yaml  # Universal exclusions only
extract-cli init --list  # Show available templates
```

#### `generate` - Create Documentation
```bash
extract-cli generate [flags]

# Examples
extract-cli generate                           # Uses default config file
extract-cli generate -c config.yaml           # Specify config file
extract-cli generate -o ./docs                # Custom output directory
```

### Available Templates

| Template | Description | Best For |
|----------|-------------|----------|
| `common` | Universal exclusions (`.git`, IDE files, etc.) | Any project type |
| `go` | Go projects with vendor/, bin/, build/ exclusions | Go applications |
| `react` | React projects with i18n support | React applications |
| `vue` | Vue.js projects with localization | Vue applications |
| `nodejs` | Node.js projects with i18n patterns | Node.js backends/APIs |
| `flutter` | Flutter/Dart projects with .arb files | Mobile applications |
| `laravel` | PHP Laravel with localization | Web applications |
| `python` | Python projects with locale support | Python applications |

## 📝 Configuration

The YAML configuration file defines file categorization patterns. Common exclusions (`.git`, IDE files, etc.) are automatically applied to all projects.

```yaml
# Project metadata
project_name: "My React App"
project_path: "."

# Data file patterns (highest priority)
data_patterns:
  - "src/data/**"
  - "public/data/**"
  - "data/**/*.json"
  - "**/*.csv"
  - "**/*.xml"

# Configuration file patterns
local_patterns:
  - "*.config.js"
  - ".env*"
  - "package.json"
  - "src/i18n/**"
  - "src/locales/**"

# Main files (included in code documentation)
main_local_files:
  - "src/index.js"
  - "src/App.js"
  - "src/i18n/en.json"  # Main English locale only

# Project-specific exclusions (common exclusions are automatic)
exclude_patterns:
  - "node_modules/**"  # Project-specific
  - "build/**"         # Project-specific
  - "package-lock.json" # Project-specific

use_regex: false
```

## 📊 Output Files with Size Information

Extract CLI generates three AI-optimized markdown files that can be easily shared with AI assistants:

### `project-code.md`
- **Content**: Source code files and main application files formatted for AI code review, debugging, and development assistance
- **Size Info**: Individual file sizes, directory totals, overall project size
- **Includes**: Main locale file (English only)
- **AI Use**: Perfect for code review, bug fixing, architecture analysis, and development guidance

### `project-data.md`
- **Content**: Data files from designated directories organized for AI analysis of project structure
- **Size Info**: Data file sizes, largest files highlighted
- **Includes**: JSON, CSV, XML, database files
- **AI Use**: Ideal for data structure analysis, schema review, and content understanding

### `project-locals.md`
- **Content**: Configuration and localization files prepared for AI translation and i18n assistance
- **Size Info**: Config file sizes, locale file breakdown
- **Includes**: All i18n files (except main English), environment files
- **AI Use**: Optimized for translation tasks, configuration review, and internationalization guidance

Each file includes metadata and well-formatted content optimized for AI consumption. Simply copy and paste the content into your preferred AI assistant for instant project understanding.

## 🔧 Development

### Prerequisites
- Go 1.21+
- Make (optional)

### Setup
```bash
git clone https://github.com/adil-chbada/extract-cli.git
cd extract-cli
make deps
make build
```

### Available Commands
```bash
make build       # Build binary
make test        # Run tests
make clean       # Clean artifacts
make build-all   # Cross-compile
make help        # Show all commands
```

## 📚 Best Practices

### Configuration Tips
1. **Use specific patterns**: `data/**/*.json` instead of `**/*.json`
2. **Leverage automatic exclusions**: Common files are excluded automatically
3. **Organize by purpose**: Separate data, code, and configuration clearly
4. **Include main files wisely**: Only essential files in code documentation

### Pattern Precedence
1. `exclude_patterns` (highest)
2. `data_patterns`
3. `local_patterns`
4. Default to code files

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

### Template Guidelines
- Follow pattern precedence rules
- Include comprehensive i18n support
- Use specific directory patterns
- Test with real projects

## 📜 License

MIT License - see [LICENSE](LICENSE) file for details.

## 💬 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/adil-chbada/extract-cli/issues)
- 💡 **Features**: [Feature Requests](https://github.com/adil-chbada/extract-cli/issues)
- 📖 **Docs**: [GitHub Wiki](https://github.com/adil-chbada/extract-cli/wiki)

---