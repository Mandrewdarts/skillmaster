package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"skillmaster/pkg/manifest"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new SkillMaster project",
	Long:  `Creates a skillmaster.json manifest file and .ai directory in the current directory.`,
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if manifest already exists
	if manifest.Exists(cwd) {
		color.Yellow("⚠ skillmaster.json already exists in this directory")
		
		// Prompt to overwrite
		fmt.Print("Overwrite? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			color.Blue("ℹ Initialization cancelled")
			return nil
		}
	}

	// Get project name from directory
	projectName := filepath.Base(cwd)

	// Create new manifest
	m := manifest.New(projectName)

	// Save manifest
	if err := m.Save(cwd); err != nil {
		return err
	}

	// Create installation directory
	installDir := filepath.Join(cwd, m.Config.InstallDir)
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create installation directory: %w", err)
	}

	// Add to .gitignore if it exists
	gitignorePath := filepath.Join(cwd, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		// Read existing gitignore
		content, err := os.ReadFile(gitignorePath)
		if err != nil {
			color.Yellow("⚠ Warning: Could not read .gitignore: %v", err)
		} else {
			// Check if .ai/ is already in gitignore
			if !strings.Contains(string(content), m.Config.InstallDir) {
				// Append to gitignore
				f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					color.Yellow("⚠ Warning: Could not update .gitignore: %v", err)
				} else {
					defer f.Close()
					
					// Add newline if file doesn't end with one
					if len(content) > 0 && content[len(content)-1] != '\n' {
						f.WriteString("\n")
					}
					
					f.WriteString(fmt.Sprintf("\n# SkillMaster installation directory\n%s/\n", m.Config.InstallDir))
					color.Green("✓ Added %s/ to .gitignore", m.Config.InstallDir)
				}
			}
		}
	}

	color.Green("✓ Initialized SkillMaster project: %s", projectName)
	color.Blue("ℹ Created skillmaster.json and %s/ directory", m.Config.InstallDir)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  • Search for packages: %s\n", color.CyanString("skillmaster search <query>"))
	fmt.Printf("  • Install a package: %s\n", color.CyanString("skillmaster install <owner/repo>"))

	return nil
}
