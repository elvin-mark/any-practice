package progress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// ItemRecord holds progress for a single Linux scenario.
type ItemRecord struct {
	Completed    bool      `json:"completed" yaml:"completed"`
	Attempts     int       `json:"attempts" yaml:"attempts"`
	FirstStarted time.Time `json:"first_started,omitempty" yaml:"first_started,omitempty"`
	CompletedAt  time.Time `json:"completed_at,omitempty" yaml:"completed_at,omitempty"`
	HintsViewed  int       `json:"hints_viewed,omitempty" yaml:"hints_viewed,omitempty"`
}

// Progress represents the persistent user progress state.
type Progress struct {
	Version   int                   `json:"version" yaml:"version"`
	Tool      string                `json:"tool" yaml:"tool"`
	UpdatedAt time.Time             `json:"updated_at" yaml:"updated_at"`
	ActiveID  string                `json:"active_id,omitempty" yaml:"active_id,omitempty"`
	Items     map[string]ItemRecord `json:"items" yaml:"items"`

	filePath string
	mu       sync.Mutex
}

// DefaultPath returns the standard progress storage file path.
func DefaultPath() string {
	// 1. Prefer local directory if .lp/ exists
	if _, err := os.Stat(".lp"); err == nil {
		return filepath.Join(".lp", "progress.json")
	}
	// 2. User config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, "any-practice", "linux-progress.json")
	}
	// 3. Fallback to local
	return filepath.Join(".lp", "progress.json")
}

// LegacyYAMLPath returns the path to any pre-existing state.yaml file.
func LegacyYAMLPath() string {
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, "linux-practice", "state.yaml")
	}
	return ""
}

// Load reads progress from path, migrating from legacy state.yaml if present.
func Load(path string) (*Progress, error) {
	if path == "" {
		path = DefaultPath()
	}

	p := &Progress{
		Version:   1,
		Tool:      "lp",
		UpdatedAt: time.Now(),
		Items:     make(map[string]ItemRecord),
		filePath:  path,
	}

	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, p); err == nil && p.Items != nil {
			p.filePath = path
			return p, nil
		}
	}

	// Attempt migration from legacy state.yaml only when loading default path
	if path == DefaultPath() {
		legacyPath := LegacyYAMLPath()
		if legacyPath != "" {
			if yamlData, yErr := os.ReadFile(legacyPath); yErr == nil {
				var legacy struct {
					ActiveScenario string `yaml:"active_scenario"`
					Scenarios      map[string]struct {
						Completed    bool      `yaml:"completed"`
						Attempts     int       `yaml:"attempts"`
						FirstStarted time.Time `yaml:"first_started"`
						CompletedAt  time.Time `yaml:"completed_at"`
						HintsViewed  int       `yaml:"hints_viewed"`
					} `yaml:"scenarios"`
				}
				if err := yaml.Unmarshal(yamlData, &legacy); err == nil {
					p.ActiveID = legacy.ActiveScenario
					for id, sc := range legacy.Scenarios {
						p.Items[id] = ItemRecord{
							Completed:    sc.Completed,
							Attempts:     sc.Attempts,
							FirstStarted: sc.FirstStarted,
							CompletedAt:  sc.CompletedAt,
							HintsViewed:  sc.HintsViewed,
						}
					}
					_ = p.Save() // persist migrated json
					return p, nil
				}
			}
		}
	}

	return p, nil
}

// Save atomically writes progress to disk via a temporary file.
func (p *Progress) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.UpdatedAt = time.Now()
	p.Version = 1
	p.Tool = "lp"

	dir := filepath.Dir(p.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory for progress %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing progress: %w", err)
	}

	tmpFile := fmt.Sprintf("%s.tmp.%d", p.filePath, time.Now().UnixNano())
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("writing temp progress file: %w", err)
	}

	if err := os.Rename(tmpFile, p.filePath); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("renaming temp progress file: %w", err)
	}

	return nil
}

// IsCompleted returns true if the scenario has been resolved.
func (p *Progress) IsCompleted(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	item, ok := p.Items[id]
	return ok && item.Completed
}

// GetRecord returns a copy of the scenario record.
func (p *Progress) GetRecord(id string) ItemRecord {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Items[id]
}

// RecordStart marks a scenario as started.
func (p *Progress) RecordStart(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.ActiveID = id
	rec := p.Items[id]
	if rec.FirstStarted.IsZero() {
		rec.FirstStarted = time.Now()
	}
	rec.Attempts++
	p.Items[id] = rec
}

// RecordCompletion marks a scenario as successfully completed.
func (p *Progress) RecordCompletion(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rec := p.Items[id]
	rec.Completed = true
	rec.CompletedAt = time.Now()
	p.Items[id] = rec
}

// RecordHint marks that a hint has been viewed.
func (p *Progress) RecordHint(id string, hintIndex int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rec := p.Items[id]
	if hintIndex+1 > rec.HintsViewed {
		rec.HintsViewed = hintIndex + 1
	}
	p.Items[id] = rec
}

// SetActive updates the active scenario ID.
func (p *Progress) SetActive(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ActiveID = id
}

// GetActive returns the active scenario ID.
func (p *Progress) GetActive() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ActiveID
}

// ClearActive marks no scenario as active.
func (p *Progress) ClearActive() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ActiveID = ""
}

// CompletedCount returns the total number of completed scenarios.
func (p *Progress) CompletedCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	count := 0
	for _, it := range p.Items {
		if it.Completed {
			count++
		}
	}
	return count
}
