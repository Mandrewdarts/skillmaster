package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"skillmaster/pkg/github"
	"skillmaster/pkg/installer"
	"skillmaster/pkg/manifest"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed packages",
	Long:  `Display all packages installed in the current project with their versions.`,
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
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

	// Check if there are any dependencies
	if len(m.Dependencies) == 0 {
		color.Yellow("No packages installed")
		fmt.Println()
		fmt.Println("Install a package with:")
		fmt.Printf("  %s\n", color.CyanString("skillmaster install <owner/repo>"))
		return nil
	}

	// Print header
	fmt.Println()
	color.Cyan("Installed Packages")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("%-40s %-15s %s\n", "Package", "Version", "Files")
	fmt.Println(strings.Repeat("─", 70))

	// Get installation directory
	installDir := filepath.Join(cwd, m.Config.InstallDir)

	// List all dependencies
	for packageName, version := range m.Dependencies {
		// Parse package name
		owner, repo, err := github.ParseRepoURL(packageName)
		if err != nil {
			color.Red("✗ Invalid package name: %s", packageName)
			continue
		}

		// Count installed files
		fileCount, err := installer.CountInstalledFiles(installDir, owner, repo)
		if err != nil {
			color.Yellow("⚠ %s", err)
			fileCount = 0
		}

		// Print package info
		if fileCount > 0 {
			fmt.Printf("%-40s %-15s %d file(s)\n", packageName, version, fileCount)
		} else {
			fmt.Printf("%-40s %-15s %s\n", packageName, version, color.RedString("not installed"))
		}
	}

	fmt.Println(strings.Repeat("─", 70))
	fmt.Println()
	color.Blue("ℹ Installation directory: %s", m.Config.InstallDir)

	return nil
}
