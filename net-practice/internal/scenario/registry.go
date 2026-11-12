package scenario

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Registry manages available networking scenarios.
type Registry struct {
	scenarios []*Scenario
	byID      map[string]*Scenario
	dir       string
}

// ResolveScenariosDir finds the scenarios directory using common lookup paths.
func ResolveScenariosDir(override string) (string, error) {
	if override != "" {
		if fi, err := os.Stat(override); err == nil && fi.IsDir() {
			return override, nil
		}
		return "", fmt.Errorf("specified scenarios directory not found: %s", override)
	}

	candidates := []string{
		"scenarios",
		"net-practice/scenarios",
		"../scenarios",
		"../../scenarios",
	}

	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "scenarios"),
			filepath.Join(exeDir, "..", "scenarios"),
			filepath.Join(exeDir, "..", "net-practice", "scenarios"),
		)
	}

	for _, cand := range candidates {
		if fi, err := os.Stat(cand); err == nil && fi.IsDir() {
			return cand, nil
		}
	}

	return "", fmt.Errorf("scenarios directory not found; specify via --scenarios-dir or run from net-practice folder")
}

// NewRegistry loads all scenarios from the specified directory.
func NewRegistry(dir string) (*Registry, error) {
	reg := &Registry{
		scenarios: make([]*Scenario, 0),
		byID:      make(map[string]*Scenario),
		dir:       dir,
	}

	if err := reg.loadAll(); err != nil {
		return nil, err
	}

	return reg, nil
}

func (r *Registry) loadAll() error {
	entries, err := os.ReadDir(r.dir)
	if err != nil {
		return fmt.Errorf("reading scenarios directory %s: %w", r.dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		path := filepath.Join(r.dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading scenario %s: %w", path, err)
		}

		var sc Scenario
		if err := yaml.Unmarshal(data, &sc); err != nil {
			return fmt.Errorf("parsing scenario %s: %w", path, err)
		}

		if sc.ID == "" {
			return fmt.Errorf("scenario %s missing required 'id' field", path)
		}

		r.scenarios = append(r.scenarios, &sc)
		r.byID[sc.ID] = &sc
	}

	// Deterministic sort: Category then ID
	sort.Slice(r.scenarios, func(i, j int) bool {
		if r.scenarios[i].Category != r.scenarios[j].Category {
			return r.scenarios[i].Category < r.scenarios[j].Category
		}
		return r.scenarios[i].ID < r.scenarios[j].ID
	})

	return nil
}

// Get returns the scenario with the specified ID.
func (r *Registry) Get(id string) (*Scenario, bool) {
	sc, ok := r.byID[id]
	return sc, ok
}

// List returns all loaded scenarios.
func (r *Registry) List() []*Scenario {
	result := make([]*Scenario, len(r.scenarios))
	copy(result, r.scenarios)
	return result
}

// Categories returns unique sorted category names.
func (r *Registry) Categories() []string {
	catSet := make(map[string]struct{})
	for _, sc := range r.scenarios {
		if sc.Category != "" {
			catSet[sc.Category] = struct{}{}
		}
	}
	var cats []string
	for cat := range catSet {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	return cats
}
