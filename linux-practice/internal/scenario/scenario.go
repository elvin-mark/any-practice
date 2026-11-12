package scenario

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ContainerConfig holds container execution settings for a scenario.
type ContainerConfig struct {
	Image      string   `yaml:"image,omitempty"`
	Privileged bool     `yaml:"privileged,omitempty"`
	Ports      []string `yaml:"ports,omitempty"`
	Env        []string `yaml:"env,omitempty"`
}

// Scenario represents a practice lab challenge with metadata and execution scripts.
type Scenario struct {
	ID               string          `yaml:"id"`
	Title            string          `yaml:"title"`
	Track            string          `yaml:"track"`
	Difficulty       string          `yaml:"difficulty"`
	Description      string          `yaml:"description"`
	Tags             []string        `yaml:"tags,omitempty"`
	EstimatedMinutes int             `yaml:"estimated_minutes,omitempty"`
	Author           string          `yaml:"author,omitempty"`
	Container        ContainerConfig `yaml:"container,omitempty"`
	Hints            []string        `yaml:"hints,omitempty"`

	// In-memory script contents
	SetupScript  string `yaml:"-"`
	VerifyScript string `yaml:"-"`
	IsCustom     bool   `yaml:"-"`
	SourcePath   string `yaml:"-"`
}

// Validate checks whether the scenario has all required fields and valid configurations.
func (s *Scenario) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return fmt.Errorf("scenario id is required")
	}
	if strings.TrimSpace(s.Title) == "" {
		return fmt.Errorf("scenario title is required")
	}
	if strings.TrimSpace(s.Track) == "" {
		return fmt.Errorf("scenario track is required")
	}
	if strings.TrimSpace(s.Description) == "" {
		return fmt.Errorf("scenario description is required")
	}
	if strings.TrimSpace(s.SetupScript) == "" {
		return fmt.Errorf("scenario setup.sh script is missing or empty")
	}
	if strings.TrimSpace(s.VerifyScript) == "" {
		return fmt.Errorf("scenario verify.sh script is missing or empty")
	}
	return nil
}

// LoadFromFS reads a single scenario directory from any fs.FS (e.g. embed.FS or os.DirFS).
func LoadFromFS(sys fs.FS, dirPath string) (*Scenario, error) {
	manifestPath := filepath.Join(dirPath, "scenario.yaml")
	data, err := fs.ReadFile(sys, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading scenario manifest at %s: %w", manifestPath, err)
	}

	var sc Scenario
	if err := yaml.Unmarshal(data, &sc); err != nil {
		return nil, fmt.Errorf("parsing scenario YAML at %s: %w", manifestPath, err)
	}

	setupPath := filepath.Join(dirPath, "setup.sh")
	setupData, err := fs.ReadFile(sys, setupPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", setupPath, err)
	}
	sc.SetupScript = string(setupData)

	verifyPath := filepath.Join(dirPath, "verify.sh")
	verifyData, err := fs.ReadFile(sys, verifyPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", verifyPath, err)
	}
	sc.VerifyScript = string(verifyData)

	sc.SourcePath = dirPath

	// Default container values
	if sc.Container.Image == "" {
		sc.Container.Image = "linux-practice-ubuntu:24.04"
	}
	// By default, systemd practice containers require privileged access
	sc.Container.Privileged = true

	return &sc, nil
}

// LoadFromDisk loads a scenario from a local directory path.
func LoadFromDisk(dirPath string) (*Scenario, error) {
	manifestPath := filepath.Join(dirPath, "scenario.yaml")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading scenario manifest at %s: %w", manifestPath, err)
	}

	var sc Scenario
	if err := yaml.Unmarshal(data, &sc); err != nil {
		return nil, fmt.Errorf("parsing scenario YAML at %s: %w", manifestPath, err)
	}

	setupPath := filepath.Join(dirPath, "setup.sh")
	setupData, err := os.ReadFile(setupPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", setupPath, err)
	}
	sc.SetupScript = string(setupData)

	verifyPath := filepath.Join(dirPath, "verify.sh")
	verifyData, err := os.ReadFile(verifyPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", verifyPath, err)
	}
	sc.VerifyScript = string(verifyData)

	sc.IsCustom = true
	sc.SourcePath = dirPath

	if sc.Container.Image == "" {
		sc.Container.Image = "linux-practice-ubuntu:24.04"
	}
	sc.Container.Privileged = true

	return &sc, nil
}
