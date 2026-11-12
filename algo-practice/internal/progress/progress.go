package progress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ItemRecord holds completion and benchmark details for an exercise.
type ItemRecord struct {
	Completed    bool      `json:"completed"`
	Attempts     int       `json:"attempts"`
	FirstStarted time.Time `json:"first_started,omitempty"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
	HintsViewed  int       `json:"hints_viewed,omitempty"`
	BestNsPerOp  float64   `json:"best_ns_per_op,omitempty"`
	AllocBytes   uint64    `json:"alloc_bytes,omitempty"`
	AllocsPerOp  uint64    `json:"allocs_per_op,omitempty"`
}

// Progress tracks the user's overall learning journey.
type Progress struct {
	Version   int                   `json:"version"`
	Tool      string                `json:"tool"`
	UpdatedAt time.Time             `json:"updated_at"`
	ActiveID  string                `json:"active_id,omitempty"`
	Items     map[string]ItemRecord `json:"items"`
	Completed map[string]ItemRecord `json:"completed,omitempty"`

	filePath string
	mu       sync.Mutex
}

// DefaultPath returns the standard progress storage file path.
func DefaultPath() string {
	// 1. Prefer local directory if .algo/ exists
	if _, err := os.Stat(".algo"); err == nil {
		return filepath.Join(".algo", "progress.json")
	}
	// 2. User config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, "any-practice", "algo-progress.json")
	}
	// 3. Fallback to local .algo/progress.json
	return filepath.Join(".algo", "progress.json")
}

// Load loads progress from the given path, initializing a new one if it doesn't exist.
func Load(path string) (*Progress, error) {
	if path == "" {
		path = DefaultPath()
	}

	p := &Progress{
		Version:   1,
		Tool:      "algo",
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

	// Try unmarshaling modern schema
	if err := json.Unmarshal(data, p); err == nil && len(p.Items) > 0 {
		p.filePath = path
		p.Completed = p.Items
		return p, nil
	}

	// Backward compatibility fallback for legacy schema
	var legacy struct {
		CurrentExerciseID string `json:"current_exercise_id"`
		Completed         map[string]struct {
			CompletedAt string  `json:"completed_at"`
			Attempts    int     `json:"attempts"`
			BestNsPerOp float64 `json:"best_ns_per_op"`
			AllocBytes  uint64  `json:"alloc_bytes"`
			AllocsPerOp uint64  `json:"allocs_per_op"`
		} `json:"completed"`
	}
	if err := json.Unmarshal(data, &legacy); err == nil {
		p.ActiveID = legacy.CurrentExerciseID
		for id, rec := range legacy.Completed {
			t, _ := time.Parse(time.RFC3339, rec.CompletedAt)
			p.Items[id] = ItemRecord{
				Completed:   true,
				Attempts:    rec.Attempts,
				CompletedAt: t,
				BestNsPerOp: rec.BestNsPerOp,
				AllocBytes:  rec.AllocBytes,
				AllocsPerOp: rec.AllocsPerOp,
			}
		}
	}

	p.Completed = p.Items
	p.filePath = path
	return p, nil
}

// Save atomically writes progress to disk via a temporary file.
func (p *Progress) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.UpdatedAt = time.Now()
	p.Version = 1
	p.Tool = "algo"

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

// IsCompleted returns true if the exercise has been solved.
func (p *Progress) IsCompleted(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	item, ok := p.Items[id]
	return ok && item.Completed
}

// MarkCompleted records an exercise as solved with performance metrics.
func (p *Progress) MarkCompleted(id string, nsPerOp float64, allocBytes, allocsPerOp uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rec := p.Items[id]
	rec.Completed = true
	rec.CompletedAt = time.Now()
	if rec.BestNsPerOp == 0 || (nsPerOp > 0 && nsPerOp < rec.BestNsPerOp) {
		rec.BestNsPerOp = nsPerOp
		rec.AllocBytes = allocBytes
		rec.AllocsPerOp = allocsPerOp
	}
	p.Items[id] = rec
	p.Completed = p.Items
}

// RecordAttempt increments the attempts count for an exercise.
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

// SetCurrent updates the active exercise.
func (p *Progress) SetCurrent(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ActiveID = id
}

// GetCurrent returns the active exercise ID.
func (p *Progress) GetCurrent() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ActiveID
}

// Stats returns the count of completed exercises.
func (p *Progress) Stats() int {
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

// GetRecord returns a copy of the item record.
func (p *Progress) GetRecord(id string) ItemRecord {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Items[id]
}
