package progress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ItemRecord tracks completion and attempt details for a challenge.
type ItemRecord struct {
	Completed    bool      `json:"completed"`
	Attempts     int       `json:"attempts"`
	FirstStarted time.Time `json:"first_started,omitempty"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
	HintsViewed  int       `json:"hints_viewed,omitempty"`
	PointsEarned int       `json:"points_earned,omitempty"`
}

// Progress tracks the user's overall CTF progress and score.
type Progress struct {
	Version    int                   `json:"version"`
	Tool       string                `json:"tool"`
	UpdatedAt  time.Time             `json:"updated_at"`
	ActiveID   string                `json:"active_id,omitempty"`
	TotalScore int                   `json:"total_score,omitempty"`
	Items      map[string]ItemRecord `json:"items"`
	Completed  map[string]ItemRecord `json:"completed,omitempty"`

	filePath string
	mu       sync.Mutex
}

// DefaultPath returns the standard progress storage file path.
func DefaultPath() string {
	// 1. Prefer local directory if .ctf/ exists
	if _, err := os.Stat(".ctf"); err == nil {
		return filepath.Join(".ctf", "progress.json")
	}
	// 2. User config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, "any-practice", "ctf-progress.json")
	}
	// 3. Fallback to local .ctf/progress.json
	return filepath.Join(".ctf", "progress.json")
}

// Load loads progress from the given path, initializing a new one if it doesn't exist.
func Load(path string) (*Progress, error) {
	if path == "" {
		path = DefaultPath()
	}

	p := &Progress{
		Version:   1,
		Tool:      "ctf",
		UpdatedAt: time.Now(),
		Items:     make(map[string]ItemRecord),
		Completed: make(map[string]ItemRecord),
		filePath:  path,
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return p, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading progress file %s: %w", path, err)
	}

	if err := json.Unmarshal(data, p); err != nil {
		return nil, fmt.Errorf("parsing progress file %s: %w", path, err)
	}

	if p.Items == nil {
		p.Items = make(map[string]ItemRecord)
	}
	p.Completed = p.Items
	p.filePath = path
	p.recalcScore()

	return p, nil
}

// Save atomically writes progress to disk via a temporary file.
func (p *Progress) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.UpdatedAt = time.Now()
	p.Version = 1
	p.Tool = "ctf"
	p.recalcScore()

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

// recalcScore recalculates the total score earned across all completed challenges.
func (p *Progress) recalcScore() {
	score := 0
	for _, rec := range p.Items {
		if rec.Completed {
			score += rec.PointsEarned
		}
	}
	p.TotalScore = score
}

// IsCompleted returns true if the challenge has been solved.
func (p *Progress) IsCompleted(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	item, ok := p.Items[id]
	return ok && item.Completed
}

// RecordStart marks when work on a challenge begins.
func (p *Progress) RecordStart(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rec := p.Items[id]
	if rec.FirstStarted.IsZero() {
		rec.FirstStarted = time.Now()
	}
	p.Items[id] = rec
	p.Completed = p.Items
}

// RecordAttempt increments the attempt counter for a challenge.
func (p *Progress) RecordAttempt(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rec := p.Items[id]
	rec.Attempts++
	if rec.FirstStarted.IsZero() {
		rec.FirstStarted = time.Now()
	}
	p.Items[id] = rec
	p.Completed = p.Items
}

// RecordCompletion marks a challenge as solved and awards points.
func (p *Progress) RecordCompletion(id string, points int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rec := p.Items[id]
	rec.Completed = true
	rec.CompletedAt = time.Now()
	rec.PointsEarned = points
	p.Items[id] = rec
	p.Completed = p.Items
	p.recalcScore()
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
	p.Completed = p.Items
}

// SetActive sets the active challenge ID.
func (p *Progress) SetActive(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ActiveID = id
}

// GetActive returns the active challenge ID.
func (p *Progress) GetActive() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ActiveID
}

// ClearActive clears the active challenge.
func (p *Progress) ClearActive() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ActiveID = ""
}

// CompletedCount returns the number of completed challenges.
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

// TotalPoints returns the total points accumulated.
func (p *Progress) TotalPoints() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.TotalScore
}

// GetRecord returns a copy of the item record for a challenge.
func (p *Progress) GetRecord(id string) ItemRecord {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Items[id]
}
