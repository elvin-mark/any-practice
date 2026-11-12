package progress

import (
	"path/filepath"
	"testing"
)

func TestProgress(t *testing.T) {
	tmpDir := t.TempDir()
	progPath := filepath.Join(tmpDir, ".algo", "progress.json")

	p, err := Load(progPath)
	if err != nil {
		t.Fatalf("failed to load progress: %v", err)
	}

	if p.IsCompleted("ex1") {
		t.Fatalf("expected ex1 not completed initially")
	}

	p.RecordAttempt("ex1")
	p.MarkCompleted("ex1", 125.5, 32, 2)
	p.SetCurrent("ex1")

	if !p.IsCompleted("ex1") {
		t.Fatalf("expected ex1 to be completed")
	}

	if err := p.Save(); err != nil {
		t.Fatalf("failed to save progress: %v", err)
	}

	// Reload from disk
	p2, err := Load(progPath)
	if err != nil {
		t.Fatalf("failed to reload progress: %v", err)
	}

	if !p2.IsCompleted("ex1") {
		t.Fatalf("expected reloaded progress to have ex1 completed")
	}

	rec := p2.Completed["ex1"]
	if rec.Attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", rec.Attempts)
	}
	if rec.BestNsPerOp != 125.5 {
		t.Errorf("expected best ns/op 125.5, got %f", rec.BestNsPerOp)
	}
}
