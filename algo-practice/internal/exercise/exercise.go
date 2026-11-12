package exercise

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Complexity specifies target time and space bounds.
type Complexity struct {
	Time  string `json:"time"`
	Space string `json:"space"`
}

// Exercise represents a single algorithm problem definition.
type Exercise struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Track       string     `json:"track"`
	TrackOrder  int        `json:"track_order"`
	Order       int        `json:"order"`
	Difficulty  string     `json:"difficulty"` // Easy, Medium, Hard
	Description string     `json:"description"`
	Complexity  Complexity `json:"complexity"`
	Forbidden   []string   `json:"forbidden"`
	EntryFile   string     `json:"entry_file"`
	TestFile    string     `json:"test_file"`
	Hints       []string   `json:"hints"`

	// Runtime fields populated on load
	Dir string `json:"-"`
}

// LoadFromDir reads exercise.json from the specified directory.
func LoadFromDir(dir string) (*Exercise, error) {
	manifestPath := filepath.Join(dir, "exercise.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading exercise manifest %s: %w", manifestPath, err)
	}

	var ex Exercise
	if err := json.Unmarshal(data, &ex); err != nil {
		return nil, fmt.Errorf("parsing exercise manifest %s: %w", manifestPath, err)
	}
	ex.Dir = dir
	if ex.EntryFile == "" {
		ex.EntryFile = "solution.go"
	}
	if ex.TestFile == "" {
		ex.TestFile = "solution_test.go"
	}

	return &ex, nil
}

// EntryFilePath returns the absolute path to the editable source file.
func (e *Exercise) EntryFilePath() string {
	return filepath.Join(e.Dir, e.EntryFile)
}

// TestFilePath returns the absolute path to the test file.
func (e *Exercise) TestFilePath() string {
	return filepath.Join(e.Dir, e.TestFile)
}

// SolutionFilePath returns the path to the reference solution file if present.
func (e *Exercise) SolutionFilePath() string {
	return filepath.Join(e.Dir, "solution", "solution.go")
}
