package github

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub API client
type Client struct {
	client *github.Client
	ctx    context.Context
}

// RepositoryInfo contains information about a repository
type RepositoryInfo struct {
	Owner       string
	Name        string
	Description string
	Stars       int
	UpdatedAt   string
	DefaultBranch string
}

// FileContent represents a file downloaded from GitHub
type FileContent struct {
	Path    string
	Content []byte
}

// NewClient creates a new GitHub API client
func NewClient(token string) *Client {
	ctx := context.Background()
	var client *github.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	return &Client{
		client: client,
		ctx:    ctx,
	}
}

// GetRepository fetches repository information
func (c *Client) GetRepository(owner, repo string) (*RepositoryInfo, error) {
	repository, resp, err := c.client.Repositories.Get(c.ctx, owner, repo)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil, fmt.Errorf("repository not found: %s/%s", owner, repo)
		}
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}

	info := &RepositoryInfo{
		Owner:       owner,
		Name:        repo,
		DefaultBranch: repository.GetDefaultBranch(),
	}

	if repository.Description != nil {
		info.Description = *repository.Description
	}
	if repository.StargazersCount != nil {
		info.Stars = *repository.StargazersCount
	}
	if repository.UpdatedAt != nil {
		info.UpdatedAt = repository.UpdatedAt.Format("2006-01-02")
	}

	return info, nil
}

// GetLatestVersion attempts to get the latest version/tag for a repository
func (c *Client) GetLatestVersion(owner, repo string) (string, error) {
	// Try to get latest release
	release, _, err := c.client.Repositories.GetLatestRelease(c.ctx, owner, repo)
	if err == nil && release.TagName != nil {
		return *release.TagName, nil
	}

	// If no releases, try to get tags
	tags, _, err := c.client.Repositories.ListTags(c.ctx, owner, repo, &github.ListOptions{
		PerPage: 1,
	})
	if err == nil && len(tags) > 0 && tags[0].Name != nil {
		return *tags[0].Name, nil
	}

	// If no tags, get the default branch
	repoInfo, _, err := c.client.Repositories.Get(c.ctx, owner, repo)
	if err != nil {
		return "main", nil // Default fallback
	}

	if repoInfo.DefaultBranch != nil {
		return *repoInfo.DefaultBranch, nil
	}

	return "main", nil
}

// DownloadMarkdownFiles recursively downloads all markdown files from a repository
func (c *Client) DownloadMarkdownFiles(owner, repo, ref string) ([]FileContent, error) {
	var files []FileContent
	
	err := c.walkRepositoryTree(owner, repo, ref, "", &files)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no markdown files found in repository")
	}

	return files, nil
}

// walkRepositoryTree recursively walks through the repository tree
func (c *Client) walkRepositoryTree(owner, repo, ref, dirPath string, files *[]FileContent) error {
	_, contents, resp, err := c.client.Repositories.GetContents(c.ctx, owner, repo, dirPath, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return fmt.Errorf("path not found: %s", dirPath)
		}
		return fmt.Errorf("failed to get repository contents: %w", err)
	}

	for _, content := range contents {
		if content.Type == nil || content.Path == nil {
			continue
		}

		contentType := *content.Type
		contentPath := *content.Path

		// Skip hidden files and directories
		if strings.HasPrefix(path.Base(contentPath), ".") {
			continue
		}

		if contentType == "file" {
			// Only process markdown files
			if strings.HasSuffix(strings.ToLower(contentPath), ".md") {
				fileContent, err := c.downloadFile(owner, repo, ref, contentPath)
				if err != nil {
					return fmt.Errorf("failed to download file %s: %w", contentPath, err)
				}
				*files = append(*files, FileContent{
					Path:    contentPath,
					Content: fileContent,
				})
			}
		} else if contentType == "dir" {
			// Recursively walk subdirectories
			if err := c.walkRepositoryTree(owner, repo, ref, contentPath, files); err != nil {
				return err
			}
		}
	}

	return nil
}

// downloadFile downloads a single file from GitHub
func (c *Client) downloadFile(owner, repo, ref, filePath string) ([]byte, error) {
	fileContent, _, resp, err := c.client.Repositories.GetContents(c.ctx, owner, repo, filePath, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil, fmt.Errorf("file not found: %s", filePath)
		}
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}

	if fileContent == nil {
		return nil, fmt.Errorf("file content is nil: %s", filePath)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, fmt.Errorf("failed to decode file content: %w", err)
	}

	return []byte(content), nil
}

// SearchRepositories searches for repositories by topic
func (c *Client) SearchRepositories(query string, limit int) ([]*RepositoryInfo, error) {
	// Build search query with required topic
	searchQuery := fmt.Sprintf("topic:skillmaster-package %s", query)
	
	opts := &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	}

	result, resp, err := c.client.Search.Repositories(c.ctx, searchQuery, opts)
	if err != nil {
		if resp != nil && resp.StatusCode == 403 {
			return nil, fmt.Errorf("GitHub API rate limit exceeded. Please add a GitHub token to ~/.skillmaster/config.json")
		}
		return nil, fmt.Errorf("failed to search repositories: %w", err)
	}

	var repos []*RepositoryInfo
	for _, repo := range result.Repositories {
		info := &RepositoryInfo{
			Owner: repo.GetOwner().GetLogin(),
			Name:  repo.GetName(),
		}

		if repo.Description != nil {
			info.Description = *repo.Description
		}
		if repo.StargazersCount != nil {
			info.Stars = *repo.StargazersCount
		}
		if repo.UpdatedAt != nil {
			info.UpdatedAt = repo.UpdatedAt.Format("2006-01-02")
		}

		repos = append(repos, info)
	}

	return repos, nil
}

// ParseRepoURL parses a repository URL in the format "owner/repo"
func ParseRepoURL(repoURL string) (owner, repo string, err error) {
	parts := strings.Split(repoURL, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format. Expected: owner/repo")
	}
	
	owner = strings.TrimSpace(parts[0])
	repo = strings.TrimSpace(parts[1])
	
	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("owner and repository name cannot be empty")
	}
	
	return owner, repo, nil
}
