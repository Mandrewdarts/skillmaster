# SkillMaster

AI Code Assistant Markdown Package Manager

## Overview

SkillMaster is a CLI tool for managing, discovering, and sharing AI code assistant configuration files (prompts, instructions, context files) across projects using GitHub as the repository backend. Think of it as npm/go modules for AI assistant markdown files.

## Features (Phase 1 MVP)

- ✅ **Project Initialization** - Set up SkillMaster in any project
- ✅ **Package Installation** - Install AI assistant markdown files from GitHub repositories
- ✅ **Package Listing** - View all installed packages with versions
- ✅ **Package Search** - Discover packages on GitHub by topic
- ✅ **Namespaced Installation** - Packages are isolated in their own directories
- ✅ **Version Tracking** - Track package versions in manifest file

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/skillmaster.git
cd skillmaster

# Build the binary
go build -o skillmaster

# Move to your PATH (optional)
sudo mv skillmaster /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/yourusername/skillmaster@latest
```

## Quick Start

### 1. Initialize a Project

```bash
# Navigate to your project directory
cd my-project

# Initialize SkillMaster
skillmaster init
```

This creates:

- `skillmaster.json` - Manifest file tracking dependencies
- `.ai/` - Directory where packages are installed
- Updates `.gitignore` to exclude the `.ai/` directory

### 2. Search for Packages

```bash
# Search for packages
skillmaster search react
skillmaster search python best-practices
```

### 3. Install a Package

```bash
# Install a package from GitHub
skillmaster install username/repository-name
```

The package will be installed to `.ai/username-repository-name/` with all markdown files preserved in their original directory structure.

### 4. List Installed Packages

```bash
skillmaster list
```

Example output:

```
Installed Packages
──────────────────────────────────────────────────────────────────────
Package                          Version         Files
──────────────────────────────────────────────────────────────────────
anthropic/claude-best-practices  main            12 file(s)
company/style-guide              v1.2.0          8 file(s)
──────────────────────────────────────────────────────────────────────

ℹ Installation directory: .ai
```

## Configuration

### Global Configuration

SkillMaster stores global configuration in `~/.skillmaster/config.json`:

```json
{
  "github": {
    "token": "ghp_your_token_here"
  },
  "installDir": ".ai"
}
```

#### Adding a GitHub Token

For higher API rate limits and access to private repositories:

1. Generate a [GitHub Personal Access Token](https://github.com/settings/tokens)
2. Add it to your config:

```bash
# Edit the config file
nano ~/.skillmaster/config.json

# Add your token:
{
  "github": {
    "token": "ghp_your_token_here"
  },
  "installDir": ".ai"
}
```

### Project Configuration

Each project has a `skillmaster.json` manifest:

```json
{
  "name": "my-project",
  "version": "1.0.0",
  "dependencies": {
    "anthropic/claude-best-practices": "main",
    "company/style-guide": "v1.2.0"
  },
  "config": {
    "installDir": ".ai",
    "autoMerge": true
  }
}
```

## Creating Packages

To create a package that others can install:

### 1. Create a GitHub Repository

Structure your repository with markdown files:

```
your-package/
├── README.md
├── prompts/
│   ├── component-design.md
│   └── code-review.md
├── patterns/
│   └── best-practices.md
└── examples/
    └── sample.md
```

### 2. Add GitHub Topics

Add these topics to your repository:

- `skillmaster-package` (required - identifies it as a SkillMaster package)
- Additional topics: `ai-assistant`, `claude`, `copilot`, `cursor`, etc.

### 3. Tag Versions (Optional)

Create Git tags for versions:

```bash
git tag -a v1.0.0 -m "Version 1.0.0"
git push origin v1.0.0
```

### 4. Others Can Install

```bash
skillmaster install your-username/your-package
```

## File Structure

After initialization and installing packages:

```
my-project/
├── .ai/                                    # Installation directory
│   ├── anthropic-claude-best-practices/   # Namespaced package
│   │   ├── prompts/
│   │   └── patterns/
│   └── company-style-guide/               # Another package
│       └── guidelines/
├── skillmaster.json                       # Manifest file
└── .gitignore                            # Updated to exclude .ai/
```

## Commands Reference

### `skillmaster init`

Initialize a SkillMaster project in the current directory.

```bash
skillmaster init
```

### `skillmaster install <owner/repo>`

Install a package from GitHub.

```bash
skillmaster install anthropic/claude-best-practices
```

### `skillmaster list`

List all installed packages with versions and file counts.

```bash
skillmaster list
```

### `skillmaster search <query>`

Search for packages on GitHub by topic and keywords.

```bash
skillmaster search react
skillmaster search python machine-learning
```

### `skillmaster --version`

Show the SkillMaster CLI version.

```bash
skillmaster --version
```

### `skillmaster --help`

Show help information.

```bash
skillmaster --help
skillmaster install --help
```

## How It Works

1. **GitHub as Registry**: SkillMaster uses GitHub as its package registry. Any GitHub repository with the `skillmaster-package` topic can be installed.

2. **Namespaced Installation**: Packages are installed to `.ai/owner-repo/` to prevent file conflicts.

3. **Version Tracking**: Package versions (tags, releases, or branches) are tracked in `skillmaster.json`.

4. **Markdown Focus**: Only `.md` files are downloaded and installed, keeping installations lightweight.

5. **GitHub API**: Uses GitHub's API to search, download, and discover packages.

## Troubleshooting

### Rate Limit Errors

**Problem**: Getting rate limit errors from GitHub API.

**Solution**: Add a GitHub personal access token to `~/.skillmaster/config.json`.

### Repository Not Found

**Problem**: `repository not found` error when installing.

**Solution**:

- Check the repository exists: `owner/repo`
- Verify it's public (or add a GitHub token for private repos)
- Ensure you're using the correct owner and repository name

### No Markdown Files Found

**Problem**: Installation fails with "no markdown files found".

**Solution**: The repository must contain at least one `.md` file to be installable.

## Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/skillmaster.git
cd skillmaster

# Install dependencies
go mod download

# Build
go build -o skillmaster

# Run tests
go test ./...
```

### Project Structure

```
skillmaster/
├── cmd/                  # Command implementations
│   ├── root.go          # Root command
│   ├── init.go          # init command
│   ├── install.go       # install command
│   ├── list.go          # list command
│   └── search.go        # search command
├── pkg/
│   ├── manifest/        # Manifest file handling
│   ├── github/          # GitHub API client
│   ├── installer/       # Installation logic
│   └── config/          # Configuration management
├── main.go
└── go.mod
```

## Roadmap

### Phase 2 Features (Planned)

- [ ] `skillmaster update` - Update packages to latest versions
- [ ] `skillmaster remove` - Uninstall packages
- [ ] Version resolution with semantic versioning (^, ~, >=)
- [ ] Lock file for reproducible installs
- [ ] Package dependencies (packages depending on other packages)
- [ ] `skillmaster publish` - Publish packages
- [ ] Advanced merge strategies
- [ ] Local package development with `skillmaster link`

### Phase 3+ Features (Future)

- [ ] Package templates and scaffolding
- [ ] Multi-project sync
- [ ] IDE integrations (VS Code, Cursor)
- [ ] Private package registries
- [ ] Package analytics and ratings
- [ ] AI-powered package recommendations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Credits

Inspired by:

- npm (JavaScript package manager)
- Go modules (Go dependency management)
- Homebrew (macOS package manager)

---

**Note**: This is Phase 1 MVP. More features are planned for future releases. See the [project roadmap](project.md) for details.
