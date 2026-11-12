package progress_test

import (
	"path/filepath"
	"testing"

	"github.com/elvin-mark/any-practice/k8s-practice/internal/progress"
)

func TestProgress_Lifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	progPath := filepath.Join(tmpDir, "progress.json")

	p, err := progress.Load(progPath)
	if err != nil {
		t.Fatalf("failed to load initial progress: %v", err)
	}

	if p.IsCompleted("arch-01-rbac") {
		t.Errorf("expected scenario not completed initially")
	}

	p.RecordStart("arch-01-rbac")
	if p.GetActive() != "arch-01-rbac" {
		t.Errorf("expected active scenario 'arch-01-rbac', got %q", p.GetActive())
	}

	rec := p.GetRecord("arch-01-rbac")
	if rec.Attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", rec.Attempts)
	}

	p.RecordHint("arch-01-rbac", 0)
	rec = p.GetRecord("arch-01-rbac")
	if rec.HintsViewed != 1 {
		t.Errorf("expected 1 hint viewed, got %d", rec.HintsViewed)
	}

	p.RecordCompletion("arch-01-rbac", 6)
	if !p.IsCompleted("arch-01-rbac") {
		t.Errorf("expected scenario to be marked completed")
	}

	if err := p.Save(); err != nil {
		t.Fatalf("failed to save progress: %v", err)
	}

	// Reload from disk
	p2, err := progress.Load(progPath)
	if err != nil {
		t.Fatalf("failed to reload progress: %v", err)
	}

	if !p2.IsCompleted("arch-01-rbac") {
		t.Errorf("expected reloaded progress to have scenario completed")
	}
	if p2.CompletedCount() != 1 {
		t.Errorf("expected 1 completed scenario, got %d", p2.CompletedCount())
	}
	if p2.EarnedScore() != 6 {
		t.Errorf("expected 6 points earned, got %d", p2.EarnedScore())
	}
}
