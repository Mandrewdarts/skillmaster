# SkillMaster - AI Code Assistant Markdown Package Manager

## Overview

SkillMaster is a CLI tool for managing, discovering, and sharing AI code assistant configuration files (prompts, instructions, context files) across projects using GitHub as the repository backend. Think of it as npm/go modules for AI assistant markdown files.

## Core Concept

Many developers use AI code assistants (Claude, GitHub Copilot, Cursor, etc.) with custom markdown files containing:
- System prompts
- Project-specific instructions
- Code style guides
- Architecture decision records
- Common patterns and examples

Currently, these files are duplicated across projects or manually copied. SkillMaster provides a centralized way to:
1. Publish packages of AI assistant markdown files to GitHub
2. Install and version them in projects
3. Keep them updated across all projects
4. Share them with teams and the community

---

## High-Impact Features (Priority Order)

### 1. **Package Installation & Management**
**Impact: Critical - Core functionality**

- `skillmaster init` - Initialize a new project with a `skillmaster.json` manifest
- `skillmaster install <github-user/repo>` - Install a package from GitHub
- `skillmaster install <github-user/repo>@<version>` - Install specific version
- `skillmaster update [package-name]` - Update packages to latest versions
- `skillmaster remove <package-name>` - Uninstall a package
- `skillmaster list` - List all installed packages with versions

**Technical Details:**
- Manifest file: `skillmaster.json` (similar to package.json)
- Lock file: `skillmaster.lock` (ensures reproducible installs)
- Default installation directory: `.ai/` or `.skillmaster/`
- GitHub API integration for fetching releases/tags
- Semantic versioning support

```json
// skillmaster.json example
{
  "name": "my-project",
  "version": "1.0.0",
  "dependencies": {
    "anthropic/claude-best-practices": "^2.1.0",
    "company/internal-style-guide": "1.0.0",
    "community/react-patterns": "~3.2.0"
  },
  "config": {
    "installDir": ".ai",
    "autoMerge": true
  }
}
```

### 2. **GitHub Repository Discovery & Search**
**Impact: High - Essential for package ecosystem**

- `skillmaster search <query>` - Search for packages on GitHub
- `skillmaster info <package-name>` - Show package details, README, stats
- Topic-based discovery using GitHub topics: `ai-assistant`, `skillmaster-package`, `claude-prompts`
- Star count, last updated, description display
- Filter by language, framework, use-case

**Technical Details:**
- GitHub REST API or GraphQL API
- Search by topics, keywords in README
- Cache search results locally
- Display compatible version ranges

### 3. **Package Publishing & Versioning**
**Impact: High - Enables ecosystem growth**

- `skillmaster publish` - Publish current directory as a package
- `skillmaster version <major|minor|patch>` - Bump version following semver
- Automatic git tagging on publish
- Validation of package structure before publish
- Package template generation

**Technical Details:**
- Requires GitHub authentication (personal access token)
- Creates Git tags for versions
- Validates required files (README.md, skillmaster-package.json)
- Optional: GitHub Releases integration

```json
// skillmaster-package.json (in published packages)
{
  "name": "react-patterns",
  "version": "3.2.0",
  "description": "Best practices for React development with AI assistants",
  "author": "username",
  "license": "MIT",
  "keywords": ["react", "patterns", "best-practices"],
  "files": [
    "prompts/*.md",
    "patterns/*.md",
    "examples/*.md"
  ]
}
```

### 4. **Intelligent File Merging & Conflict Resolution**
**Impact: High - Critical for usability**

- Merge multiple packages without file conflicts
- Namespace/prefix support for avoiding collisions
- Conflict detection and resolution strategies
- Option to keep packages separate vs. merged

**Technical Details:**
- Configurable merge strategies:
  - `merge`: Combine all files into shared directories
  - `namespace`: Prefix each package's files with package name
  - `separate`: Keep each package in its own subdirectory
- Conflict resolution prompts during install
- Ability to override specific files locally

### 5. **Local Development & Package Creation**
**Impact: High - Developer experience**

- `skillmaster create <package-name>` - Scaffold a new package
- Package templates for common use-cases:
  - Language-specific (Python, JavaScript, Go, etc.)
  - Framework-specific (React, Next.js, Django, etc.)
  - General best practices
  - Code review guidelines
- Local package linking for development: `skillmaster link`
- Validation and linting of markdown files

### 6. **Version Resolution & Dependency Management**
**Impact: Medium-High - Package ecosystem stability**

- Semantic versioning support (^, ~, >=, etc.)
- Dependency resolution algorithm
- Conflict detection when multiple packages require different versions
- Lock file for reproducible builds
- `skillmaster outdated` - Check for newer versions

**Technical Details:**
- Implement semver resolution (similar to npm)
- Handle transitive dependencies if packages depend on other packages
- Generate and maintain lock file

### 7. **Configuration & Profiles**
**Impact: Medium - Flexibility**

- Global configuration: `~/.skillmaster/config.json`
- Project-level configuration in `skillmaster.json`
- Profiles for different AI assistants (Claude, Copilot, Cursor)
- Custom installation directories
- GitHub token management
- Registry configuration (default: GitHub, extensible)

```json
// Global config example
{
  "github": {
    "token": "ghp_xxx"
  },
  "profiles": {
    "claude": {
      "installDir": ".claude"
    },
    "cursor": {
      "installDir": ".cursorrules"
    }
  },
  "registry": "github.com"
}
```

### 8. **Package Templates & Scaffolding**
**Impact: Medium - Accelerates adoption**

- Built-in templates for common scenarios
- Community template registry
- `skillmaster create --template <template-name>`
- Interactive package creation wizard

### 9. **Sync & Multi-Project Management**
**Impact: Medium - Enterprise/team usage**

- `skillmaster sync` - Sync packages across multiple projects
- Workspace support for monorepos
- Global packages that apply to all projects
- Team package presets

### 10. **Analytics & Usage Insights**
**Impact: Low-Medium - Community building**

- Anonymous usage statistics
- Popular packages dashboard
- Download counts (via GitHub API stars/clones)
- Trending packages

---

## Technical Architecture

### Core Components

1. **CLI Layer**
   - Command parser (use: Commander.js or Cobra for Go)
   - Interactive prompts (use: Inquirer.js or survey for Go)
   - Colorized output (use: Chalk.js or color for Go)

2. **Package Manager Core**
   - Manifest parser (JSON)
   - Lock file generator
   - Version resolver
   - Dependency graph builder

3. **GitHub Integration Layer**
   - REST API client
   - Authentication handler
   - Rate limiting
   - Caching layer

4. **File System Manager**
   - Installation handler
   - File merging logic
   - Conflict resolution
   - Directory structure manager

5. **Configuration Manager**
   - Global config reader/writer
   - Project config handler
   - Environment variable support

### Tech Stack Recommendations

**Option A: Node.js/TypeScript**
- Pros: Rich ecosystem, fast development, familiar to web devs
- Cons: Requires Node.js runtime
- Key libraries: Commander, Inquirer, Octokit, fs-extra

**Option B: Go**
- Pros: Single binary, fast, easy distribution
- Cons: Slightly more verbose, smaller ecosystem for CLI tools
- Key libraries: Cobra, survey, go-github, afero

**Option C: Rust**
- Pros: Blazing fast, single binary, memory safe
- Cons: Steeper learning curve, longer compile times
- Key libraries: clap, dialoguer, octocrab

**Recommendation: Go** - Best balance of performance, distribution simplicity, and developer experience

### File Structure

```
skillmaster/
├── cmd/
│   ├── init.go
│   ├── install.go
│   ├── search.go
│   ├── publish.go
│   └── ...
├── pkg/
│   ├── manifest/
│   ├── github/
│   ├── resolver/
│   ├── installer/
│   └── config/
├── internal/
│   └── utils/
├── testdata/
└── main.go
```

### Data Structures

```go
// Manifest represents skillmaster.json
type Manifest struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Dependencies map[string]string `json:"dependencies"`
    Config       Config            `json:"config"`
}

// Package represents a SkillMaster package
type Package struct {
    Name        string
    Version     string
    Description string
    Repository  string
    Files       []string
    Metadata    PackageMetadata
}

// LockFile ensures reproducible installs
type LockFile struct {
    Version  string
    Packages map[string]LockedPackage
}
```

---

## Package Repository Structure

Packages are standard GitHub repositories with specific conventions:

```
react-best-practices/          # Repository name
├── skillmaster-package.json   # Package manifest
├── README.md                  # Package documentation
├── prompts/
│   ├── component-design.md
│   └── hooks-guidelines.md
├── patterns/
│   ├── state-management.md
│   └── error-handling.md
└── examples/
    └── sample-component.md
```

GitHub Topics Required:
- `skillmaster-package` (primary identifier)
- Additional: `ai-assistant`, `claude`, `copilot`, etc.

---

## MVP Feature Set (Phase 1)

For initial release, focus on:

1. ✅ `skillmaster init` - Initialize project
2. ✅ `skillmaster install <github-user/repo>` - Install from GitHub
3. ✅ `skillmaster list` - List installed packages
4. ✅ Basic `skillmaster.json` manifest support
5. ✅ Simple file copying to installation directory
6. ✅ GitHub API integration (read-only)
7. ✅ `skillmaster search` - Basic search by topics

**MVP Deferred:**
- Complex version resolution
- Dependency graphs
- Publishing workflow
- Merge strategies (just copy files)
- Lock files

---

## Future Enhancements (Phase 2+)

- **AI-Powered Features:**
  - Automatic package recommendations based on project type
  - Smart merging of conflicting instructions
  - Package quality scoring using AI analysis

- **Community Features:**
  - Package ratings and reviews
  - Curated collections
  - Official verified packages

- **IDE Integrations:**
  - VS Code extension
  - JetBrains plugin
  - Direct Cursor/Claude integration

- **Advanced Features:**
  - Private package registries
  - Package namespaces/orgs
  - Automated testing of prompts
  - Package analytics dashboard

---

## Success Metrics

- Number of published packages in ecosystem
- Active users/installations
- GitHub stars on CLI tool
- Community contributions
- Package install frequency
- User retention rate

---

## Open Questions

1. Should packages support dependencies on other packages?
2. How to handle breaking changes in packages?
3. Licensing considerations for shared prompts?
4. Rate limiting strategy for GitHub API?
5. Should we support private GitHub repos?
6. Package naming conventions and collision handling?
7. How to version the CLI tool itself for backward compatibility?

---

## Getting Started (Post-Implementation)

```bash
# Install SkillMaster
go install github.com/yourusername/skillmaster@latest

# Initialize a project
skillmaster init

# Search for packages
skillmaster search react

# Install a package
skillmaster install anthropic/claude-best-practices

# Your AI assistant files are now in .ai/
```

---

## References & Inspiration

- npm (JavaScript package manager)
- Go modules (Go dependency management)
- Homebrew (macOS package manager)
- asdf (runtime version manager)
- GitHub Package Registry
