package exam

import (
	"testing"
	"time"

	"github.com/elvin-mark/any-practice/k8s-practice/internal/scenario"
)

func TestExamSession(t *testing.T) {
	scenarios := []*scenario.Scenario{
		{
			Metadata: scenario.Metadata{
				ID:     "q1",
				Domain: scenario.DomainStorage,
				Weight: 5,
			},
		},
		{
			Metadata: scenario.Metadata{
				ID:     "q2",
				Domain: scenario.DomainTrouble,
				Weight: 10,
			},
		},
	}

	session, err := StartExamSession(scenarios, 60)
	if err != nil {
		t.Fatalf("StartExamSession failed: %v", err)
	}

	if session.Mode != ModeExam {
		t.Errorf("expected ModeExam, got %s", session.Mode)
	}

	if len(session.ExamQuestions) != 2 {
		t.Errorf("expected 2 questions, got %d", len(session.ExamQuestions))
	}

	if session.IsExpired() {
		t.Errorf("new session should not be expired")
	}

	timeStr := session.TimeRemaining()
	if timeStr == "" || timeStr == "Untimed" {
		t.Errorf("expected formatted remaining time, got %s", timeStr)
	}

	// Update score
	session.Scores["q1"] = QuestionStatus{
		ScenarioID: "q1",
		Domain:     scenario.DomainStorage,
		Weight:     5,
		Earned:     5,
		Passed:     true,
		Attempted:  true,
		LastTested: time.Now(),
	}

	if session.Scores["q1"].Earned != 5 {
		t.Errorf("expected 5 earned points")
	}

	_ = ClearSession()
}
