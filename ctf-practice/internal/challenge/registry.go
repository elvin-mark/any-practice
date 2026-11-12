package challenge

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Registry manages available CTF challenges.
type Registry struct {
	challenges []Challenge
	byID       map[string]Challenge
	dir        string
}

// ResolveChallengesDir finds the challenge directory using common lookup paths.
func ResolveChallengesDir(override string) (string, error) {
	if override != "" {
		if fi, err := os.Stat(override); err == nil && fi.IsDir() {
			return override, nil
		}
		return "", fmt.Errorf("specified challenges directory not found: %s", override)
	}

	candidates := []string{
		"challenges",
		"ctf-practice/challenges",
		"../challenges",
		"../../challenges",
	}

	// Also check relative to executable location
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "challenges"),
			filepath.Join(exeDir, "..", "challenges"),
			filepath.Join(exeDir, "..", "ctf-practice", "challenges"),
		)
	}

	for _, cand := range candidates {
		if fi, err := os.Stat(cand); err == nil && fi.IsDir() {
			return cand, nil
		}
	}

	return "", fmt.Errorf("challenges directory not found; specify via --challenges-dir or run from ctf-practice folder")
}

// NewRegistry loads all challenges from the specified directory.
func NewRegistry(dir string) (*Registry, error) {
	reg := &Registry{
		challenges: make([]Challenge, 0),
		byID:       make(map[string]Challenge),
		dir:        dir,
	}

	if err := reg.loadAll(); err != nil {
		return nil, err
	}

	return reg, nil
}

func (r *Registry) loadAll() error {
	entries, err := os.ReadDir(r.dir)
	if err != nil {
		return fmt.Errorf("reading challenges directory %s: %w", r.dir, err)
	}

	for _, entry := range entries {
		path := filepath.Join(r.dir, entry.Name())

		// If entry is a subdirectory, inspect files inside it (e.g. challenges/crypto-01/challenge.yaml)
		if entry.IsDir() {
			subEntries, err := os.ReadDir(path)
			if err != nil {
				continue
			}
			for _, sub := range subEntries {
				if isChallengeFile(sub.Name()) {
					if err := r.loadFile(filepath.Join(path, sub.Name())); err != nil {
						return err
					}
				}
			}
			continue
		}

		if isChallengeFile(entry.Name()) {
			if err := r.loadFile(path); err != nil {
				return err
			}
		}
	}

	// Sort challenges deterministically by Category then Points then ID
	sort.Slice(r.challenges, func(i, j int) bool {
		if r.challenges[i].Category != r.challenges[j].Category {
			return r.challenges[i].Category < r.challenges[j].Category
		}
		if r.challenges[i].Points != r.challenges[j].Points {
			return r.challenges[i].Points < r.challenges[j].Points
		}
		return r.challenges[i].ID < r.challenges[j].ID
	})

	return nil
}

func isChallengeFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".yaml" || ext == ".yml" || ext == ".json"
}

func (r *Registry) loadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading challenge file %s: %w", path, err)
	}

	var ch Challenge
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".json" {
		if err := json.Unmarshal(data, &ch); err != nil {
			return fmt.Errorf("parsing JSON challenge %s: %w", path, err)
		}
	} else {
		if err := yaml.Unmarshal(data, &ch); err != nil {
			return fmt.Errorf("parsing YAML challenge %s: %w", path, err)
		}
	}

	if ch.ID == "" {
		return fmt.Errorf("challenge in %s missing required 'id' field", path)
	}

	r.challenges = append(r.challenges, ch)
	r.byID[ch.ID] = ch
	return nil
}

// Get returns the challenge with the specified ID.
func (r *Registry) Get(id string) (Challenge, bool) {
	ch, ok := r.byID[id]
	return ch, ok
}

// List returns all loaded challenges.
func (r *Registry) List() []Challenge {
	result := make([]Challenge, len(r.challenges))
	copy(result, r.challenges)
	return result
}

// Categories returns unique sorted category names.
func (r *Registry) Categories() []string {
	catSet := make(map[string]struct{})
	for _, ch := range r.challenges {
		if ch.Category != "" {
			catSet[ch.Category] = struct{}{}
		}
	}
	var cats []string
	for cat := range catSet {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	return cats
}

// ValidateFlag checks if the provided flag matches the challenge's flag.
func (r *Registry) ValidateFlag(id, submittedFlag string) (bool, error) {
	ch, ok := r.Get(id)
	if !ok {
		return false, fmt.Errorf("challenge '%s' not found", id)
	}

	expected := strings.TrimSpace(ch.Flag)
	submitted := strings.TrimSpace(submittedFlag)

	if expected == "" {
		return false, fmt.Errorf("challenge '%s' has no flag configured", id)
	}

	return expected == submitted, nil
}
