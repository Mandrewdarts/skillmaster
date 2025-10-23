package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"skillmaster/pkg/config"
	"skillmaster/pkg/github"
	"skillmaster/pkg/installer"
	"skillmaster/pkg/manifest"
)

var installCmd = &cobra.Command{
	Use:   "install <owner/repo>",
	Short: "Install a package from GitHub",
	Long: `Install a SkillMaster package from a GitHub repository.
	
Example:
  skillmaster install anthropic/claude-best-practices`,
	Args: cobra.ExactArgs(1),
	RunE: runInstall,
}

func runInstall(cmd *cobra.Command, args []string) error {
	repoURL := args[0]

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

	// Get latest version
	color.Blue("→ Fetching repository information...")
	version, err := githubClient.GetLatestVersion(owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get repository version: %w", err)
	}

	// Create installer
	inst := installer.New(githubClient)

	// Install directory (absolute path)
	installDir := filepath.Join(cwd, m.Config.InstallDir)

	// Install package
	color.Blue("→ Downloading markdown files...")
	fileCount, err := inst.InstallPackage(owner, repo, installDir)
	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Add to manifest dependencies
	packageName := fmt.Sprintf("%s/%s", owner, repo)
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
