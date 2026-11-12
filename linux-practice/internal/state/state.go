package state

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/elvin-mark/linux-practice/internal/config"
	"gopkg.in/yaml.v3"
)

// ScenarioProgress records user progress on an individual scenario.
type ScenarioProgress struct {
	Completed    bool      `yaml:"completed"`
	Attempts     int       `yaml:"attempts"`
	FirstStarted time.Time `yaml:"first_started,omitempty"`
	CompletedAt  time.Time `yaml:"completed_at,omitempty"`
	HintsViewed  int       `yaml:"hints_viewed"`
}

// State represents the overall user progress and active challenge state.
type State struct {
	ActiveScenario string                      `yaml:"active_scenario,omitempty"`
	Scenarios      map[string]*ScenarioProgress `yaml:"scenarios"`
}

// Load reads and parses state.yaml, returning an empty state if it does not yet exist.
func Load() (*State, error) {
	statePath, err := config.StateFilePath()
	if err != nil {
		return nil, err
	}

	st := &State{
		Scenarios: make(map[string]*ScenarioProgress),
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return st, nil
		}
		return nil, fmt.Errorf("reading state file at %s: %w", statePath, err)
	}

	if err := yaml.Unmarshal(data, st); err != nil {
		return nil, fmt.Errorf("parsing state file: %w", err)
	}

	if st.Scenarios == nil {
		st.Scenarios = make(map[string]*ScenarioProgress)
	}

	return st, nil
}

// Save writes the state to state.yaml atomically via a temporary file.
func (s *State) Save() error {
	statePath, err := config.StateFilePath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling state to yaml: %w", err)
	}

	dir := filepath.Dir(statePath)
	tmpFile, err := os.CreateTemp(dir, "state.*.yaml.tmp")
	if err != nil {
		return fmt.Errorf("creating temp file for state: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("writing state data: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("closing temp state file: %w", err)
	}

	if err := os.Rename(tmpFile.Name(), statePath); err != nil {
		return fmt.Errorf("saving state file to %s: %w", statePath, err)
	}

	return nil
}

// GetOrCreateProgress returns the progress for a scenario, initializing if necessary.
func (s *State) GetOrCreateProgress(scenarioID string) *ScenarioProgress {
	p, exists := s.Scenarios[scenarioID]
	if !exists {
		p = &ScenarioProgress{}
		s.Scenarios[scenarioID] = p
	}
	return p
}

// RecordStart marks a scenario as currently started.
func (s *State) RecordStart(scenarioID string) {
	s.ActiveScenario = scenarioID
	p := s.GetOrCreateProgress(scenarioID)
	if p.FirstStarted.IsZero() {
		p.FirstStarted = time.Now()
	}
	p.Attempts++
}

// RecordCompletion marks a scenario as successfully verified.
func (s *State) RecordCompletion(scenarioID string) {
	p := s.GetOrCreateProgress(scenarioID)
	p.Completed = true
	p.CompletedAt = time.Now()
}

// RecordHint marks that a hint has been viewed.
func (s *State) RecordHint(scenarioID string, hintIndex int) {
	p := s.GetOrCreateProgress(scenarioID)
	if hintIndex+1 > p.HintsViewed {
		p.HintsViewed = hintIndex + 1
	}
}

// ClearActive marks no scenario as active.
func (s *State) ClearActive() {
	s.ActiveScenario = ""
}
