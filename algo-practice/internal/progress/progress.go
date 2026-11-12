package progress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Record holds completion details for an exercise.
type Record struct {
	CompletedAt string  `json:"completed_at"`
	Attempts    int     `json:"attempts"`
	BestNsPerOp float64 `json:"best_ns_per_op,omitempty"`
	AllocBytes  uint64  `json:"alloc_bytes,omitempty"`
	AllocsPerOp uint64  `json:"allocs_per_op,omitempty"`
}

// Progress tracks the user's overall learning journey.
type Progress struct {
	CurrentExerciseID string            `json:"current_exercise_id"`
	Completed         map[string]Record `json:"completed"`
	filePath          string
	mu                sync.Mutex
}

// Load loads progress from the given path, initializing a new one if it doesn't exist.
func Load(path string) (*Progress, error) {
	p := &Progress{
		Completed: make(map[string]Record),
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
	if p.Completed == nil {
		p.Completed = make(map[string]Record)
	}
	p.filePath = path
	return p, nil
}

// Save atomically writes progress to disk.
func (p *Progress) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

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
		return fmt.Errorf("renaming temp progress file: %w", err)
	}

	return nil
}

// IsCompleted returns true if the exercise has been solved.
func (p *Progress) IsCompleted(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.Completed[id]
	return ok
}

// MarkCompleted records an exercise as solved.
func (p *Progress) MarkCompleted(id string, nsPerOp float64, allocBytes, allocsPerOp uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rec := p.Completed[id]
	rec.CompletedAt = time.Now().Format(time.RFC3339)
	if rec.BestNsPerOp == 0 || (nsPerOp > 0 && nsPerOp < rec.BestNsPerOp) {
		rec.BestNsPerOp = nsPerOp
		rec.AllocBytes = allocBytes
		rec.AllocsPerOp = allocsPerOp
	}
	p.Completed[id] = rec
}

// RecordAttempt increments the attempts count for an exercise.
func (p *Progress) RecordAttempt(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	rec := p.Completed[id]
	rec.Attempts++
	p.Completed[id] = rec
}

// SetCurrent updates the active exercise.
func (p *Progress) SetCurrent(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.CurrentExerciseID = id
}

// GetCurrent returns the active exercise ID.
func (p *Progress) GetCurrent() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.CurrentExerciseID
}

// Stats returns the count of completed exercises.
func (p *Progress) Stats() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.Completed)
}
