package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "skillmaster",
	Short: "AI Code Assistant Markdown Package Manager",
	Long: `SkillMaster is a CLI tool for managing, discovering, and sharing 
AI code assistant configuration files (prompts, instructions, context files) 
across projects using GitHub as the repository backend.

Think of it as npm/go modules for AI assistant markdown files.`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Register subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
}
