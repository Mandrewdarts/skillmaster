package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"skillmaster/pkg/github"
)

// Installer handles package installation
type Installer struct {
	githubClient *github.Client
}

// New creates a new Installer instance
func New(githubClient *github.Client) *Installer {
	return &Installer{
		githubClient: githubClient,
	}
}

// InstallPackage installs a package from GitHub to the specified directory
func (i *Installer) InstallPackage(owner, repo, installDir string) (int, error) {
	// Get repository information
	repoInfo, err := i.githubClient.GetRepository(owner, repo)
	if err != nil {
		return 0, err
	}

	// Get the reference to download from (default branch or latest tag)
	ref := repoInfo.DefaultBranch

	// Download all markdown files
	files, err := i.githubClient.DownloadMarkdownFiles(owner, repo, ref)
	if err != nil {
		return 0, err
	}

	// Create namespaced directory: installDir/owner-repo/
	namespace := fmt.Sprintf("%s-%s", owner, repo)
	targetDir := filepath.Join(installDir, namespace)

	// Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create installation directory: %w", err)
	}

	// Copy files maintaining directory structure
	fileCount := 0
	for _, file := range files {
		targetPath := filepath.Join(targetDir, file.Path)
		
		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fileCount, fmt.Errorf("failed to create directory for %s: %w", file.Path, err)
		}

		// Write file
		if err := os.WriteFile(targetPath, file.Content, 0644); err != nil {
			return fileCount, fmt.Errorf("failed to write file %s: %w", file.Path, err)
		}

		fileCount++
	}

	return fileCount, nil
}

// UninstallPackage removes a package from the installation directory
func (i *Installer) UninstallPackage(owner, repo, installDir string) error {
	namespace := fmt.Sprintf("%s-%s", owner, repo)
	targetDir := filepath.Join(installDir, namespace)

	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("package not installed: %s/%s", owner, repo)
	}

	// Remove directory
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to remove package directory: %w", err)
	}

	return nil
}

// CountInstalledFiles counts the number of files in an installed package
func CountInstalledFiles(installDir, owner, repo string) (int, error) {
	namespace := fmt.Sprintf("%s-%s", owner, repo)
	targetDir := filepath.Join(installDir, namespace)

	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return 0, nil
	}

	count := 0
	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			count++
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to count files: %w", err)
	}

	return count, nil
}
