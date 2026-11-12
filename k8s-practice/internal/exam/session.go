package exam

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/elvin-mark/k8s-practice/internal/scenario"
)

type SessionMode string

const (
	ModePractice SessionMode = "practice"
	ModeExam     = "exam"
)

// ExamQuestionResult tracks status of a question in an exam session
type QuestionStatus struct {
	ScenarioID string    `json:"scenarioId"`
	Domain     string    `json:"domain"`
	Weight     int       `json:"weight"`
	Earned     int       `json:"earned"`
	Passed     bool      `json:"passed"`
	Attempted  bool      `json:"attempted"`
	LastTested time.Time `json:"lastTested"`
}

// Session state maintained during practice or exam
type Session struct {
	Mode            SessionMode               `json:"mode"`
	Active          bool                      `json:"active"`
	StartTime       time.Time                 `json:"startTime"`
	DurationMinutes int                       `json:"durationMinutes"`
	ActiveScenario  string                    `json:"activeScenario"`
	ExamQuestions   []string                  `json:"examQuestions"`
	CurrentIndex    int                       `json:"currentIndex"`
	Scores          map[string]QuestionStatus `json:"scores"`
}

func getSessionFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".cka-session.json", nil
	}
	dir := filepath.Join(home, ".cka")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return ".cka-session.json", nil
	}
	return filepath.Join(dir, "session.json"), nil
}

// LoadSession loads active session from disk
func LoadSession() (*Session, error) {
	path, err := getSessionFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No active session
		}
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// SaveSession persists session to disk
func SaveSession(s *Session) error {
	path, err := getSessionFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ClearSession removes the stored session
func ClearSession() error {
	path, err := getSessionFilePath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// StartPracticeSession initializes a single-scenario practice session
func StartPracticeSession(scenarioID string) (*Session, error) {
	s := &Session{
		Mode:           ModePractice,
		Active:         true,
		StartTime:      time.Now(),
		ActiveScenario: scenarioID,
		Scores:         make(map[string]QuestionStatus),
	}
	return s, SaveSession(s)
}

// StartExamSession initializes a full timed exam session with selected scenarios
func StartExamSession(allScenarios []*scenario.Scenario, durationMinutes int) (*Session, error) {
	if len(allScenarios) == 0 {
		return nil, fmt.Errorf("no scenarios available to generate exam")
	}

	// Shuffle scenarios slightly for randomized exam experience
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	indices := r.Perm(len(allScenarios))

	var questionIDs []string
	scores := make(map[string]QuestionStatus)

	for _, idx := range indices {
		sc := allScenarios[idx]
		questionIDs = append(questionIDs, sc.ID)
		scores[sc.ID] = QuestionStatus{
			ScenarioID: sc.ID,
			Domain:     sc.Domain,
			Weight:     sc.Weight,
			Earned:     0,
			Passed:     false,
			Attempted:  false,
		}
	}

	s := &Session{
		Mode:            ModeExam,
		Active:          true,
		StartTime:       time.Now(),
		DurationMinutes: durationMinutes,
		ActiveScenario:  questionIDs[0],
		ExamQuestions:   questionIDs,
		CurrentIndex:    0,
		Scores:          scores,
	}

	return s, SaveSession(s)
}

// TimeRemaining returns formatted time left in exam
func (s *Session) TimeRemaining() string {
	if s.Mode != ModeExam || s.DurationMinutes <= 0 {
		return "Untimed"
	}

	elapsed := time.Since(s.StartTime)
	total := time.Duration(s.DurationMinutes) * time.Minute
	remaining := total - elapsed

	if remaining <= 0 {
		return "00:00:00 (TIME EXPIRED)"
	}

	hours := int(remaining.Hours())
	minutes := int(remaining.Minutes()) % 60
	seconds := int(remaining.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d remaining", hours, minutes, seconds)
}

// IsExpired returns true if exam duration has elapsed
func (s *Session) IsExpired() bool {
	if s.Mode != ModeExam || s.DurationMinutes <= 0 {
		return false
	}
	return time.Since(s.StartTime) >= time.Duration(s.DurationMinutes)*time.Minute
}
