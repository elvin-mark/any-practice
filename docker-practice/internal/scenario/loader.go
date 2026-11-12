package scenario

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultScenariosDir is the default location for scenario folders
const DefaultScenariosDir = "scenarios"

// FindScenariosDir attempts to locate the scenarios directory across common search paths
func FindScenariosDir() string {
	candidates := []string{
		"scenarios",
		"docker-practice/scenarios",
		"../scenarios",
		"../../scenarios",
	}

	for _, c := range candidates {
		if stat, err := os.Stat(c); err == nil && stat.IsDir() {
			return c
		}
	}
	return DefaultScenariosDir
}

// LoadAll scans the directory and loads all valid scenarios
func LoadAll(baseDir string) ([]*Scenario, error) {
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("scenarios directory '%s' does not exist", baseDir)
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("reading scenarios directory: %w", err)
	}

	var scenarios []*Scenario
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		scenarioPath := filepath.Join(baseDir, entry.Name())
		metaPath := filepath.Join(scenarioPath, "metadata.yaml")
		if _, err := os.Stat(metaPath); os.IsNotExist(err) {
			continue
		}

		sc, err := LoadSingle(scenarioPath)
		if err != nil {
			continue
		}
		scenarios = append(scenarios, sc)
	}

	// Sort scenarios by ID for consistent ordering
	sort.Slice(scenarios, func(i, j int) bool {
		return scenarios[i].ID < scenarios[j].ID
	})

	return scenarios, nil
}

// LoadSingle loads a single scenario from its directory path
func LoadSingle(scenarioDir string) (*Scenario, error) {
	metaPath := filepath.Join(scenarioDir, "metadata.yaml")
	metaBytes, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata at %s: %w", metaPath, err)
	}

	var meta Metadata
	if err := yaml.Unmarshal(metaBytes, &meta); err != nil {
		return nil, fmt.Errorf("parsing metadata at %s: %w", metaPath, err)
	}

	if meta.ID == "" {
		meta.ID = filepath.Base(scenarioDir)
	}

	var description string
	problemPath := filepath.Join(scenarioDir, "problem.md")
	if descBytes, err := os.ReadFile(problemPath); err == nil {
		description = strings.TrimSpace(string(descBytes))
	}

	var solution string
	solutionPath := filepath.Join(scenarioDir, "solution.md")
	if solBytes, err := os.ReadFile(solutionPath); err == nil {
		solution = strings.TrimSpace(string(solBytes))
	}

	return &Scenario{
		Metadata:    meta,
		Path:        scenarioDir,
		Description: description,
		Solution:    solution,
	}, nil
}

// FindByID searches for a scenario by ID in the base directory
func FindByID(baseDir, id string) (*Scenario, error) {
	scenarios, err := LoadAll(baseDir)
	if err != nil {
		return nil, err
	}

	for _, s := range scenarios {
		if strings.EqualFold(s.ID, id) {
			return s, nil
		}
	}

	return nil, fmt.Errorf("scenario with ID '%s' not found", id)
}
