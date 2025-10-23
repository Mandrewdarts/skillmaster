package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"skillmaster/pkg/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the current SkillMaster configuration",
	Long: `Display the current SkillMaster configuration including GitHub token status, 
installation directory, and config file location.`,
	RunE: runConfig,
}

func runConfig(cmd *cobra.Command, args []string) error {
	// Get config path
	configPath, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if config file exists
	_, err = os.Stat(configPath)
	configExists := err == nil

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Print header
	fmt.Println()
	color.Cyan("SkillMaster Configuration")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println()

	// Show config file location
	fmt.Printf("%-20s %s\n", color.BlueString("Config File:"), configPath)
	if configExists {
		fmt.Printf("%-20s %s\n", color.BlueString("Status:"), color.GreenString("✓ exists"))
	} else {
		fmt.Printf("%-20s %s\n", color.BlueString("Status:"), color.YellowString("⚠ using defaults (file not found)"))
	}
	fmt.Println()

	// Show configuration values
	fmt.Println(color.CyanString("Settings:"))
	fmt.Println("─────────────────────────────────────────")
	fmt.Printf("%-20s %s\n", "Install Directory:", cfg.InstallDir)

	// Show GitHub token status (masked)
	if cfg.GitHub.Token != "" {
		maskedToken := cfg.GitHub.Token[:min(4, len(cfg.GitHub.Token))] + "..." +
			cfg.GitHub.Token[max(0, len(cfg.GitHub.Token)-4):]
		fmt.Printf("%-20s %s %s\n", "GitHub Token:", color.GreenString("✓ configured"), color.GreenString("("+maskedToken+")"))
	} else {
		fmt.Printf("%-20s %s\n", "GitHub Token:", color.YellowString("⚠ not configured"))
		fmt.Println()
		color.Yellow("  Configure GitHub token to increase API rate limits:")
		fmt.Println("  Set the GITHUB_TOKEN environment variable or run:")
		fmt.Printf("    %s\n", color.CyanString("skillmaster init"))
	}

	fmt.Println()

	// Show raw JSON if raw flag is set
	raw, _ := cmd.Flags().GetBool("raw")
	if raw {
		fmt.Println(color.CyanString("Raw Configuration:"))
		fmt.Println("─────────────────────────────────────────")
		jsonData, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}
		fmt.Println(string(jsonData))
		fmt.Println()
	}

	return nil
}

func init() {
	configCmd.Flags().BoolP("raw", "r", false, "Show raw JSON configuration")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
