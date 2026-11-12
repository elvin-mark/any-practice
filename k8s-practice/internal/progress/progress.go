package progress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ItemRecord holds progress for a single CKA scenario.
type ItemRecord struct {
	Completed    bool      `json:"completed"`
	Attempts     int       `json:"attempts"`
	FirstStarted time.Time `json:"first_started,omitempty"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
	HintsViewed  int       `json:"hints_viewed,omitempty"`
	Score        int       `json:"score,omitempty"`
}

// Progress represents persistent user progress across CKA practice drills.
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
	// 1. Prefer local directory if .cka/ exists
	if _, err := os.Stat(".cka"); err == nil {
		return filepath.Join(".cka", "progress.json")
	}
	// 2. User config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, "any-practice", "cka-progress.json")
	}
	// 3. User home .cka directory
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".cka", "progress.json")
	}
	return filepath.Join(".cka", "progress.json")
}

// Load reads progress from path, migrating from legacy session.json if available.
func Load(path string) (*Progress, error) {
	if path == "" {
		path = DefaultPath()
	}

	p := &Progress{
		Version:   1,
		Tool:      "cka",
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

	// Migrate from legacy ~/.cka/session.json if loading default path
	if path == DefaultPath() {
		if home, err := os.UserHomeDir(); err == nil {
			sessPath := filepath.Join(home, ".cka", "session.json")
			if sessData, sErr := os.ReadFile(sessPath); sErr == nil {
				var sess struct {
					ActiveScenario string `json:"activeScenario"`
					Scores         map[string]struct {
						Passed     bool      `json:"passed"`
						Earned     int       `json:"earned"`
						Attempted  bool      `json:"attempted"`
						LastTested time.Time `json:"lastTested"`
					} `json:"scores"`
				}
				if err := json.Unmarshal(sessData, &sess); err == nil {
					p.ActiveID = sess.ActiveScenario
					for id, q := range sess.Scores {
						if q.Passed || q.Attempted {
							p.Items[id] = ItemRecord{
								Completed:    q.Passed,
								Attempts:     1,
								FirstStarted: q.LastTested,
								CompletedAt:  q.LastTested,
								Score:        q.Earned,
							}
						}
					}
					_ = p.Save()
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
	p.Tool = "cka"

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

// RecordCompletion marks a scenario as successfully completed with earned score.
func (p *Progress) RecordCompletion(id string, score int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rec := p.Items[id]
	rec.Completed = true
	rec.CompletedAt = time.Now()
	if score > rec.Score {
		rec.Score = score
	}
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

// EarnedScore returns total points earned from completed scenarios.
func (p *Progress) EarnedScore() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	total := 0
	for _, it := range p.Items {
		if it.Completed {
			total += it.Score
		}
	}
	return total
}
