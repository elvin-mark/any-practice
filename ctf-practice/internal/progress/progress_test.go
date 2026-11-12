package progress

import (
	"path/filepath"
	"testing"
)

func TestCTFProgressLifecycle(t *testing.T) {
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
	if p.TotalPoints() != 0 {
		t.Errorf("expected 0 points, got %d", p.TotalPoints())
	}

	// Record start & attempt
	p.RecordStart("crypto-01")
	p.RecordAttempt("crypto-01")
	p.RecordAttempt("crypto-01")
	p.RecordHint("crypto-01", 0)

	rec := p.GetRecord("crypto-01")
	if rec.Attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", rec.Attempts)
	}
	if rec.HintsViewed != 1 {
		t.Errorf("expected 1 hint viewed, got %d", rec.HintsViewed)
	}
	if p.IsCompleted("crypto-01") {
		t.Errorf("challenge should not be completed yet")
	}

	// Complete challenge
	p.RecordCompletion("crypto-01", 100)
	if !p.IsCompleted("crypto-01") {
		t.Errorf("expected crypto-01 to be completed")
	}
	if p.CompletedCount() != 1 {
		t.Errorf("expected 1 completed, got %d", p.CompletedCount())
	}
	if p.TotalPoints() != 100 {
		t.Errorf("expected 100 points, got %d", p.TotalPoints())
	}

	// Active ID
	p.SetActive("crypto-01")
	if p.GetActive() != "crypto-01" {
		t.Errorf("expected active crypto-01, got %s", p.GetActive())
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
	if p2.TotalPoints() != 100 {
		t.Errorf("expected 100 points on reload, got %d", p2.TotalPoints())
	}
	if !p2.IsCompleted("crypto-01") {
		t.Errorf("expected crypto-01 to be completed on reload")
	}
}
