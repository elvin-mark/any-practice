package progress

import (
	"path/filepath"
	"testing"
)

func TestGitProgressLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	progPath := filepath.Join(tmpDir, "progress.json")

	// Load fresh
	p, err := Load(progPath)
	if err != nil {
		t.Fatalf("unexpected Load error: %v", err)
	}

	if p.CompletedCount() != 0 {
		t.Errorf("expected 0 completed, got %d", p.CompletedCount())
	}

	// Record start & attempt
	p.RecordStart("rebase-01")
	p.RecordAttempt("rebase-01")
	p.RecordHint("rebase-01", 0)

	rec := p.GetRecord("rebase-01")
	if rec.Attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", rec.Attempts)
	}
	if rec.HintsViewed != 1 {
		t.Errorf("expected 1 hint viewed, got %d", rec.HintsViewed)
	}
	if p.IsCompleted("rebase-01") {
		t.Errorf("scenario should not be completed yet")
	}

	// Complete scenario
	p.RecordCompletion("rebase-01")
	if !p.IsCompleted("rebase-01") {
		t.Errorf("expected rebase-01 to be completed")
	}
	if p.CompletedCount() != 1 {
		t.Errorf("expected 1 completed, got %d", p.CompletedCount())
	}

	// Active ID
	p.SetActive("rebase-01")
	if p.GetActive() != "rebase-01" {
		t.Errorf("expected active rebase-01, got %s", p.GetActive())
	}

	// Save and reload
	if err := p.Save(); err != nil {
		t.Fatalf("unexpected Save error: %v", err)
	}

	p2, err := Load(progPath)
	if err != nil {
		t.Fatalf("unexpected Load reload error: %v", err)
	}

	if p2.CompletedCount() != 1 {
		t.Errorf("expected 1 completed on reload, got %d", p2.CompletedCount())
	}
	if !p2.IsCompleted("rebase-01") {
		t.Errorf("expected rebase-01 to be completed on reload")
	}
}
