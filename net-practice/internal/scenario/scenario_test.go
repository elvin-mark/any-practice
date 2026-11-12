package scenario

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNetScenarioRegistryAndRunner(t *testing.T) {
	tmpDir := t.TempDir()

	scYAML := `id: dns-01-test
title: "DNS Resolution Test"
category: DNS
difficulty: Easy
description: Test basic name resolution.
goal: Create hosts file with example.local mapped to 127.0.0.1.
hints:
  - Add 127.0.0.1 example.local to hosts
setup: |
  touch hosts
verify: |
  set -e
  grep -q "127.0.0.1 example.local" hosts
`

	if err := os.WriteFile(filepath.Join(tmpDir, "dns-01.yaml"), []byte(scYAML), 0644); err != nil {
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

	sc, ok := reg.Get("dns-01-test")
	if !ok {
		t.Fatalf("expected to find dns-01-test")
	}

	// Setup lab
	labDir, err := SetupLab(sc)
	if err != nil {
		t.Fatalf("SetupLab failed: %v", err)
	}
	defer CleanLab(sc)

	if _, err := os.Stat(filepath.Join(labDir, "hosts")); os.IsNotExist(err) {
		t.Errorf("expected hosts file to exist in %s", labDir)
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
	_ = os.WriteFile(filepath.Join(labDir, "hosts"), []byte("127.0.0.1 example.local\n"), 0644)

	// Verify after work
	passedCorrect, msg, err := VerifyLab(sc)
	if err != nil || !passedCorrect {
		t.Errorf("expected verification to pass, got %v, msg=%s, err=%v", passedCorrect, msg, err)
	}
}
