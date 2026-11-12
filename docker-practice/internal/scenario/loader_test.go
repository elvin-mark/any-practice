package scenario

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadScenario(t *testing.T) {
	tmpDir := t.TempDir()
	scDir := filepath.Join(tmpDir, "test-01-sample")
	if err := os.MkdirAll(scDir, 0o755); err != nil {
		t.Fatalf("failed to create scenario dir: %v", err)
	}

	metaContent := `id: test-01-sample
title: Sample Docker Test
domain: Dockerfile & Image Optimization
difficulty: Easy
weight: 5
hints:
  - "Look at the Dockerfile."
`
	if err := os.WriteFile(filepath.Join(scDir, "metadata.yaml"), []byte(metaContent), 0o644); err != nil {
		t.Fatalf("failed to write metadata.yaml: %v", err)
	}

	problemContent := "# Task\nBuild the image."
	if err := os.WriteFile(filepath.Join(scDir, "problem.md"), []byte(problemContent), 0o644); err != nil {
		t.Fatalf("failed to write problem.md: %v", err)
	}

	sc, err := LoadSingle(scDir)
	if err != nil {
		t.Fatalf("LoadSingle failed: %v", err)
	}

	if sc.ID != "test-01-sample" {
		t.Errorf("expected ID test-01-sample, got %s", sc.ID)
	}
	if sc.Title != "Sample Docker Test" {
		t.Errorf("expected title 'Sample Docker Test', got '%s'", sc.Title)
	}
	if len(sc.Hints) != 1 {
		t.Errorf("expected 1 hint, got %d", len(sc.Hints))
	}
	if sc.Description != problemContent {
		t.Errorf("expected description '%s', got '%s'", problemContent, sc.Description)
	}

	all, err := LoadAll(tmpDir)
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 scenario loaded, got %d", len(all))
	}

	found, err := FindByID(tmpDir, "test-01-sample")
	if err != nil || found == nil {
		t.Fatalf("FindByID failed: %v", err)
	}
}
