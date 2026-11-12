package challenge

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRegistryLoadAndValidate(t *testing.T) {
	tmpDir := t.TempDir()

	ch1YAML := `id: crypto-01
title: Caesar Shift
category: Crypto
difficulty: Easy
points: 50
description: Decrypt this shifted message.
hints:
  - Look into ROT13 or frequency analysis
flag: FLAG{julius_caesar_cipher}
`
	ch2JSON := `{
  "id": "web-01",
  "title": "SQL Injection",
  "category": "Web",
  "difficulty": "Medium",
  "points": 100,
  "description": "Bypass authentication.",
  "hints": ["Try SQL comment characters"],
  "flag": "FLAG{sql_bypass_login}"
}`

	if err := os.WriteFile(filepath.Join(tmpDir, "crypto-01.yaml"), []byte(ch1YAML), 0644); err != nil {
		t.Fatalf("failed to write ch1: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "web-01.json"), []byte(ch2JSON), 0644); err != nil {
		t.Fatalf("failed to write ch2: %v", err)
	}

	reg, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("unexpected NewRegistry error: %v", err)
	}

	list := reg.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 challenges, got %d", len(list))
	}

	ch, ok := reg.Get("crypto-01")
	if !ok {
		t.Fatalf("expected to find crypto-01")
	}
	if ch.Points != 50 {
		t.Errorf("expected 50 points, got %d", ch.Points)
	}

	// Validate Flag
	valid, err := reg.ValidateFlag("crypto-01", "FLAG{julius_caesar_cipher}")
	if err != nil || !valid {
		t.Errorf("expected valid flag match, got valid=%v, err=%v", valid, err)
	}

	validWrong, err := reg.ValidateFlag("crypto-01", "FLAG{wrong_flag}")
	if err != nil || validWrong {
		t.Errorf("expected wrong flag to fail, got valid=%v, err=%v", validWrong, err)
	}

	// Non-existent challenge
	_, err = reg.ValidateFlag("non-existent", "FLAG{test}")
	if err == nil {
		t.Errorf("expected error for non-existent challenge")
	}
}
