package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"skillmaster/pkg/config"
	"skillmaster/pkg/github"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for packages on GitHub",
	Long: `Search for SkillMaster packages on GitHub by topic and keywords.
	
Example:
  skillmaster search react
  skillmaster search python best-practices`,
	Args: cobra.MinimumNArgs(1),
	RunE: runSearch,
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")

	// Load global config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Show warning if no GitHub token
	if cfg.GetGitHubToken() == "" {
		color.Yellow("⚠ No GitHub token configured. Search results may be limited.")
		color.Blue("ℹ Add a token to ~/.skillmaster/config.json for better results")
		fmt.Println()
	}

	// Create GitHub client
	githubClient := github.NewClient(cfg.GetGitHubToken())

	// Search repositories
	color.Blue("→ Searching GitHub repositories...")
	fmt.Println()

	repos, err := githubClient.SearchRepositories(query, 20)
	if err != nil {
		return err
	}

	// Check if no results
	if len(repos) == 0 {
		color.Yellow("No packages found matching: %s", query)
		fmt.Println()
		fmt.Println("Tips:")
		fmt.Println("  • Try different keywords")
		fmt.Println("  • Check package naming conventions")
		fmt.Println("  • Packages must have the 'skillmaster-package' topic")
		return nil
	}

	// Print results
	color.Cyan("Found %d package(s)", len(repos))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	for i, repo := range repos {
		// Package name with stars
		packageName := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
		stars := fmt.Sprintf("⭐ %d", repo.Stars)
		
		fmt.Printf("%s %s\n", color.GreenString(packageName), color.YellowString(stars))
		
		// Description
		if repo.Description != "" {
			fmt.Printf("  %s\n", color.WhiteString(repo.Description))
		}
		
		// Updated date
		if repo.UpdatedAt != "" {
			fmt.Printf("  %s %s\n", color.BlueString("Updated:"), repo.UpdatedAt)
		}
		
		// Install command
		installCmd := fmt.Sprintf("skillmaster install %s", packageName)
		fmt.Printf("  %s %s\n", color.BlueString("Install:"), color.CyanString(installCmd))
		
		// Separator between results
		if i < len(repos)-1 {
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("─", 80))

	return nil
}
