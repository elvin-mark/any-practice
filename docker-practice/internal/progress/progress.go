package progress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ItemRecord holds progress for a single Docker scenario.
type ItemRecord struct {
	Completed    bool      `json:"completed"`
	Attempts     int       `json:"attempts"`
	FirstStarted time.Time `json:"first_started,omitempty"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
	HintsViewed  int       `json:"hints_viewed,omitempty"`
	Score        int       `json:"score,omitempty"`
}

// Progress represents persistent user progress across Docker practice drills.
type Progress struct {
	Version   int                   `json:"version"`
	Tool      string                `json:"tool"`
	UpdatedAt time.Time             `json:"updated_at"`
	ActiveID  string                `json:"active_id,omitempty"`
	Items     map[string]ItemRecord `json:"items"`

	filePath string
	mu       sync.Mutex
}

// DefaultPath returns the standard progress storage file path.
func DefaultPath() string {
	// 1. Prefer local directory if .dockp/ exists
	if _, err := os.Stat(".dockp"); err == nil {
		return filepath.Join(".dockp", "progress.json")
	}
	// 2. User config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, "any-practice", "docker-progress.json")
	}
	// 3. User home .dockp directory
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".dockp", "progress.json")
	}
	return filepath.Join(".dockp", "progress.json")
}

// Load reads progress from path.
func Load(path string) (*Progress, error) {
	if path == "" {
		path = DefaultPath()
	}

	p := &Progress{
		Version:   1,
		Tool:      "dockp",
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

	return p, nil
}

// Save persists the current progress to disk.
func (p *Progress) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.UpdatedAt = time.Now()

	dir := filepath.Dir(p.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing progress: %w", err)
	}

	tmpFile := p.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0o644); err != nil {
		return fmt.Errorf("writing progress file: %w", err)
	}

	if err := os.Rename(tmpFile, p.filePath); err != nil {
		return fmt.Errorf("renaming temp progress file: %w", err)
	}

	return nil
}

// GetActive returns currently active scenario ID.
func (p *Progress) GetActive() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ActiveID
}

// SetActive sets the active scenario ID.
func (p *Progress) SetActive(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ActiveID = id
}

// ClearActive unsets the active scenario ID.
func (p *Progress) ClearActive() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ActiveID = ""
}

// RecordStart marks scenario as started and increments attempts.
func (p *Progress) RecordStart(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.ActiveID = id
	item := p.Items[id]
	if item.Attempts == 0 {
		item.FirstStarted = time.Now()
	}
	item.Attempts++
	p.Items[id] = item
}

// RecordCompletion records scenario completion and awarded points.
func (p *Progress) RecordCompletion(id string, score int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	item := p.Items[id]
	item.Completed = true
	item.CompletedAt = time.Now()
	if item.Score < score {
		item.Score = score
	}
	p.Items[id] = item
}

// RecordHint increments hints viewed counter.
func (p *Progress) RecordHint(id string) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	item := p.Items[id]
	item.HintsViewed++
	p.Items[id] = item
	return item.HintsViewed
}

// IsCompleted checks if an item is completed.
func (p *Progress) IsCompleted(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	item, exists := p.Items[id]
	return exists && item.Completed
}

// GetRecord returns record for an item.
func (p *Progress) GetRecord(id string) ItemRecord {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.Items[id]
}
