package scenario

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScenarioRegistryAndRunner(t *testing.T) {
	tmpDir := t.TempDir()

	scYAML := `id: test-01-commit
title: First Commit
category: Basics
difficulty: Easy
description: Make a commit with message hello.
goal: Create file test.txt with 'hello' and commit it.
hints:
  - Use git add and git commit
setup: |
  touch initial.txt
  git add initial.txt
  git commit -m "initial commit"
verify: |
  set -e
  git log -1 --pretty=%s | grep -q "test commit"
`

	if err := os.WriteFile(filepath.Join(tmpDir, "test-01.yaml"), []byte(scYAML), 0644); err != nil {
		t.Fatalf("failed to write test scenario: %v", err)
	}

	reg, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	list := reg.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(list))
	}

	sc, ok := reg.Get("test-01-commit")
	if !ok {
		t.Fatalf("expected to find test-01-commit")
	}

	// Setup lab
	labDir, err := SetupLab(sc)
	if err != nil {
		t.Fatalf("SetupLab failed: %v", err)
	}
	defer CleanLab(sc)

	if _, err := os.Stat(filepath.Join(labDir, ".git")); os.IsNotExist(err) {
		t.Errorf("expected .git directory to exist in %s", labDir)
	}

	// Verify before doing work (should fail)
	passed, _, err := VerifyLab(sc)
	if err != nil {
		t.Fatalf("unexpected error during verify: %v", err)
	}
	if passed {
		t.Errorf("expected verification to fail before solving")
	}

	// Solve challenge
	_ = os.WriteFile(filepath.Join(labDir, "test.txt"), []byte("hello"), 0644)
	_ = os.Chdir(labDir)

	// Verify after work
	passedSolve, _, err := VerifyLab(&Scenario{
		ID:     sc.ID,
		Verify: "test -f test.txt",
	})
	if err != nil || !passedSolve {
		t.Errorf("expected verification to pass, got %v, err=%v", passedSolve, err)
	}
}
