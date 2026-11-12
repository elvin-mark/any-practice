package scenario

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAll(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "scenarios-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create dummy scenario
	scDir := filepath.Join(tmpDir, "test-01")
	if err := os.MkdirAll(scDir, 0755); err != nil {
		t.Fatalf("failed to create test sc dir: %v", err)
	}

	metaContent := `
id: test-01
title: Test Scenario Title
domain: Storage
difficulty: Easy
weight: 5
targetContext: kind-test
targetNamespace: test-ns
hints:
  - Check PV status
`
	if err := os.WriteFile(filepath.Join(scDir, "metadata.yaml"), []byte(metaContent), 0644); err != nil {
		t.Fatalf("failed to write metadata: %v", err)
	}
	if err := os.WriteFile(filepath.Join(scDir, "problem.md"), []byte("Problem instructions here"), 0644); err != nil {
		t.Fatalf("failed to write problem.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(scDir, "solution.md"), []byte("kubectl apply -f pv.yaml"), 0644); err != nil {
		t.Fatalf("failed to write solution.md: %v", err)
	}

	scenarios, err := LoadAll(tmpDir)
	if err != nil {
		t.Fatalf("LoadAll returned error: %v", err)
	}

	if len(scenarios) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(scenarios))
	}

	s := scenarios[0]
	if s.ID != "test-01" {
		t.Errorf("expected ID 'test-01', got '%s'", s.ID)
	}
	if s.Title != "Test Scenario Title" {
		t.Errorf("expected Title 'Test Scenario Title', got '%s'", s.Title)
	}
	if s.Weight != 5 {
		t.Errorf("expected Weight 5, got %d", s.Weight)
	}
	if len(s.Hints) != 1 {
		t.Errorf("expected 1 hint, got %d", len(s.Hints))
	}
	if s.Description != "Problem instructions here" {
		t.Errorf("expected description match, got '%s'", s.Description)
	}
}
