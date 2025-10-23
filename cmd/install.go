package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"skillmaster/pkg/config"
	"skillmaster/pkg/github"
	"skillmaster/pkg/installer"
	"skillmaster/pkg/manifest"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [owner/repo]",
	Short: "Install packages from GitHub",
	Long: `Install SkillMaster packages from GitHub repositories.

When run without arguments, installs all packages listed in skillmaster.json.
When run with a package name, installs that specific package.
	
Examples:
  skillmaster install                            # Install all packages from manifest
  skillmaster install anthropic/claude-best-practices  # Install specific package
  skillmaster install --force                    # Force reinstall all packages`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInstall,
}

func init() {
	installCmd.Flags().BoolP("force", "f", false, "Force reinstall even if package is already installed")
}

func runInstall(cmd *cobra.Command, args []string) error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Load manifest
	m, err := manifest.Load(cwd)
	if err != nil {
		return err
	}

	// Get force flag
	force, _ := cmd.Flags().GetBool("force")

	// If no arguments, install all packages from manifest
	if len(args) == 0 {
		return installAll(m, cwd, force)
	}

	// Otherwise, install specific package
	return installPackage(args[0], m, cwd, force)
}

// installAll installs all packages from the manifest
func installAll(m *manifest.Manifest, cwd string, force bool) error {
	if len(m.Dependencies) == 0 {
		color.Yellow("No packages to install")
		fmt.Println()
		fmt.Println("Add packages with:")
		fmt.Printf("  %s\n", color.CyanString("skillmaster install <owner/repo>"))
		return nil
	}

	// Load global config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Show warning if no GitHub token
	if cfg.GetGitHubToken() == "" {
		color.Yellow("⚠ No GitHub token configured. API rate limits will be lower.")
		color.Blue("ℹ Add a token to ~/.skillmaster/config.json for higher rate limits")
		fmt.Println()
	}

	// Create GitHub client
	githubClient := github.NewClient(cfg.GetGitHubToken())

	// Create installer
	inst := installer.New(githubClient)

	// Install directory (absolute path)
	installDir := filepath.Join(cwd, m.Config.InstallDir)

	fmt.Println()
	color.Cyan("Installing packages...")
	fmt.Println()

	installedCount := 0
	skippedCount := 0

	// Install each package
	for packageName := range m.Dependencies {
		// Parse package name
		owner, repo, err := github.ParseRepoURL(packageName)
		if err != nil {
			color.Red("✗ Invalid package name: %s", packageName)
			continue
		}

		// Check if already installed (unless force flag is set)
		fileCount, err := installer.CountInstalledFiles(installDir, owner, repo)
		if !force && err == nil && fileCount > 0 {
			color.Green("✓ %s (already installed, %d files)", packageName, fileCount)
			skippedCount++
			continue
		}

		// Install package
		if force && fileCount > 0 {
			fmt.Printf("→ Reinstalling %s...\n", color.CyanString(packageName))
		} else {
			fmt.Printf("→ Installing %s...\n", color.CyanString(packageName))
		}
		fileCount, err = inst.InstallPackage(owner, repo, installDir)
		if err != nil {
			color.Red("✗ Failed to install %s: %v", packageName, err)
			continue
		}

		color.Green("✓ %s (%d files)", packageName, fileCount)
		installedCount++
	}

	// Summary
	fmt.Println()
	if installedCount > 0 {
		color.Green("✓ Installed %d package(s)", installedCount)
	}
	if skippedCount > 0 {
		color.Blue("ℹ Skipped %d already installed package(s)", skippedCount)
	}

	return nil
}

// installPackage installs a specific package
func installPackage(repoURL string, m *manifest.Manifest, cwd string, force bool) error {
	// Parse repository URL
	owner, repo, err := github.ParseRepoURL(repoURL)
	if err != nil {
		return err
	}

	// Load global config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Show warning if no GitHub token
	if cfg.GetGitHubToken() == "" {
		color.Yellow("⚠ No GitHub token configured. API rate limits will be lower.")
		color.Blue("ℹ Add a token to ~/.skillmaster/config.json for higher rate limits")
		fmt.Println()
	}

	// Create GitHub client
	githubClient := github.NewClient(cfg.GetGitHubToken())

	// Install directory (absolute path)
	installDir := filepath.Join(cwd, m.Config.InstallDir)

	// Check if already installed (unless force flag is set)
	packageName := fmt.Sprintf("%s/%s", owner, repo)
	existingFileCount, err := installer.CountInstalledFiles(installDir, owner, repo)
	if !force && err == nil && existingFileCount > 0 {
		color.Yellow("⚠ Package %s is already installed (%d files)", packageName, existingFileCount)
		fmt.Print("Reinstall? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			color.Blue("ℹ Installation cancelled")
			return nil
		}
	}

	// Get latest version
	color.Blue("→ Fetching repository information...")
	version, err := githubClient.GetLatestVersion(owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get repository version: %w", err)
	}

	// Create installer
	inst := installer.New(githubClient)

	// Install package
	if force || existingFileCount > 0 {
		color.Blue("→ Reinstalling markdown files...")
	} else {
		color.Blue("→ Downloading markdown files...")
	}
	fileCount, err := inst.InstallPackage(owner, repo, installDir)
	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Add to manifest dependencies
	m.AddDependency(packageName, version)

	// Save manifest
	if err := m.Save(cwd); err != nil {
		return fmt.Errorf("failed to update manifest: %w", err)
	}

	// Success message
	color.Green("✓ Successfully installed %s@%s", packageName, version)
	color.Blue("ℹ Installed %d markdown file(s) to %s/%s-%s/", fileCount, m.Config.InstallDir, owner, repo)

	return nil
}
