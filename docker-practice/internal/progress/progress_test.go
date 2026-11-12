package progress

import (
	"path/filepath"
	"testing"
)

func TestProgressLoadSave(t *testing.T) {
	tmpDir := t.TempDir()
	progPath := filepath.Join(tmpDir, "test-progress.json")

	p, err := Load(progPath)
	if err != nil {
		t.Fatalf("failed to load progress: %v", err)
	}

	if p.Tool != "dockp" {
		t.Errorf("expected tool dockp, got %s", p.Tool)
	}

	p.RecordStart("img-01-multistage")
	if p.GetActive() != "img-01-multistage" {
		t.Errorf("expected active img-01-multistage, got %s", p.GetActive())
	}

	p.RecordCompletion("img-01-multistage", 10)
	if !p.IsCompleted("img-01-multistage") {
		t.Error("expected completed to be true")
	}

	rec := p.GetRecord("img-01-multistage")
	if rec.Score != 10 {
		t.Errorf("expected score 10, got %d", rec.Score)
	}

	if err := p.Save(); err != nil {
		t.Fatalf("failed to save progress: %v", err)
	}

	p2, err := Load(progPath)
	if err != nil {
		t.Fatalf("failed to reload progress: %v", err)
	}

	if !p2.IsCompleted("img-01-multistage") {
		t.Error("persisted state was not loaded correctly")
	}
	if p2.GetRecord("img-01-multistage").Score != 10 {
		t.Errorf("expected score 10 in loaded state, got %d", p2.GetRecord("img-01-multistage").Score)
	}
}
