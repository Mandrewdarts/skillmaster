package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GitHubConfig holds GitHub-specific configuration
type GitHubConfig struct {
	Token string `json:"token"`
}

// GlobalConfig represents the global configuration file
type GlobalConfig struct {
	GitHub     GitHubConfig `json:"github"`
	InstallDir string       `json:"installDir"`
}

const (
	ConfigDirName  = ".skillmaster"
	ConfigFileName = "config.json"
)

// GetConfigPath returns the path to the global config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ConfigDirName)
	return filepath.Join(configDir, ConfigFileName), nil
}

// Load reads the global configuration file
func Load() (*GlobalConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &GlobalConfig{
			GitHub: GitHubConfig{
				Token: "",
			},
			InstallDir: ".ai",
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config GlobalConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults if not specified
	if config.InstallDir == "" {
		config.InstallDir = ".ai"
	}

	return &config, nil
}

// Save writes the global configuration file
func (c *GlobalConfig) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetGitHubToken returns the GitHub token from the config
func (c *GlobalConfig) GetGitHubToken() string {
	return c.GitHub.Token
}

// InitializeConfig creates a default config file if it doesn't exist
func InitializeConfig() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Don't overwrite existing config
	if _, err := os.Stat(configPath); err == nil {
		return nil
	}

	config := &GlobalConfig{
		GitHub: GitHubConfig{
			Token: "",
		},
		InstallDir: ".ai",
	}

	return config.Save()
}
