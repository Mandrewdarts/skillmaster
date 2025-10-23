package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the configuration section in the manifest
type Config struct {
	InstallDir string `json:"installDir"`
	AutoMerge  bool   `json:"autoMerge"`
}

// Manifest represents the skillmaster.json file
type Manifest struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
	Config       Config            `json:"config"`
}

const ManifestFileName = "skillmaster.json"

// New creates a new manifest with default values
func New(name string) *Manifest {
	return &Manifest{
		Name:         name,
		Version:      "1.0.0",
		Dependencies: make(map[string]string),
		Config: Config{
			InstallDir: ".ai",
			AutoMerge:  true,
		},
	}
}

// Load reads and parses the manifest file from the given directory
func Load(dir string) (*Manifest, error) {
	manifestPath := filepath.Join(dir, ManifestFileName)
	
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("manifest file not found: %s (run 'skillmaster init' to create one)", manifestPath)
		}
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Initialize dependencies map if nil
	if manifest.Dependencies == nil {
		manifest.Dependencies = make(map[string]string)
	}

	// Set default config values if not specified
	if manifest.Config.InstallDir == "" {
		manifest.Config.InstallDir = ".ai"
	}

	return &manifest, nil
}

// Save writes the manifest to the given directory
func (m *Manifest) Save(dir string) error {
	manifestPath := filepath.Join(dir, ManifestFileName)
	
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// AddDependency adds or updates a dependency in the manifest
func (m *Manifest) AddDependency(name, version string) {
	if m.Dependencies == nil {
		m.Dependencies = make(map[string]string)
	}
	m.Dependencies[name] = version
}

// RemoveDependency removes a dependency from the manifest
func (m *Manifest) RemoveDependency(name string) {
	if m.Dependencies != nil {
		delete(m.Dependencies, name)
	}
}

// Exists checks if a manifest file exists in the given directory
func Exists(dir string) bool {
	manifestPath := filepath.Join(dir, ManifestFileName)
	_, err := os.Stat(manifestPath)
	return err == nil
}
