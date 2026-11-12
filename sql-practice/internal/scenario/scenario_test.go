package scenario

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSQLScenarioRegistryAndRunner(t *testing.T) {
	tmpDir := t.TempDir()

	scYAML := `id: query-01-basic
title: "Top Customers"
category: Queries
difficulty: Easy
description: Find customers with highest spending.
goal: Select name and total spending descending limit 2.
hints:
  - Use ORDER BY total DESC LIMIT 2
setup: |
  CREATE TABLE customers (id INT, name TEXT, total INT);
  INSERT INTO customers VALUES (1, 'Alice', 500), (2, 'Bob', 300), (3, 'Charlie', 700);
solution: |
  SELECT name, total FROM customers ORDER BY total DESC LIMIT 2;
`

	if err := os.WriteFile(filepath.Join(tmpDir, "query-01.yaml"), []byte(scYAML), 0644); err != nil {
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

	sc, ok := reg.Get("query-01-basic")
	if !ok {
		t.Fatalf("expected to find query-01-basic")
	}

	// Setup lab
	dbPath, err := SetupLab(sc)
	if err != nil {
		t.Fatalf("SetupLab failed: %v", err)
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("expected sqlite db file to exist at %s", dbPath)
	}

	// Tables check
	tables, err := GetTables(dbPath)
	if err != nil || len(tables) != 1 || tables[0] != "customers" {
		t.Errorf("expected table customers, got %v, err=%v", tables, err)
	}

	// Verify with wrong query (should fail)
	passed, _, err := VerifyLab(sc, "SELECT name, total FROM customers ORDER BY total ASC LIMIT 2;")
	if err != nil {
		t.Fatalf("unexpected error during verify: %v", err)
	}
	if passed {
		t.Errorf("expected wrong query to fail verification")
	}

	// Verify with correct query (should pass)
	passedCorrect, msg, err := VerifyLab(sc, "SELECT name, total FROM customers ORDER BY total DESC LIMIT 2;")
	if err != nil || !passedCorrect {
		t.Errorf("expected correct query to pass, got %v, msg=%s, err=%v", passedCorrect, msg, err)
	}
}
