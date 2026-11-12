package progress_test

import (
	"path/filepath"
	"testing"

	"github.com/elvin-mark/any-practice/linux-practice/internal/progress"
)

func TestProgress_Lifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	progPath := filepath.Join(tmpDir, "progress.json")

	p, err := progress.Load(progPath)
	if err != nil {
		t.Fatalf("failed to load initial progress: %v", err)
	}

	if p.IsCompleted("net-ufw") {
		t.Errorf("expected scenario not completed initially")
	}

	p.RecordStart("net-ufw")
	if p.GetActive() != "net-ufw" {
		t.Errorf("expected active scenario 'net-ufw', got %q", p.GetActive())
	}

	rec := p.GetRecord("net-ufw")
	if rec.Attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", rec.Attempts)
	}

	p.RecordHint("net-ufw", 0)
	rec = p.GetRecord("net-ufw")
	if rec.HintsViewed != 1 {
		t.Errorf("expected 1 hint viewed, got %d", rec.HintsViewed)
	}

	p.RecordCompletion("net-ufw")
	if !p.IsCompleted("net-ufw") {
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

	if !p2.IsCompleted("net-ufw") {
		t.Errorf("expected reloaded progress to have scenario completed")
	}
	if p2.CompletedCount() != 1 {
		t.Errorf("expected 1 completed scenario, got %d", p2.CompletedCount())
	}
}
