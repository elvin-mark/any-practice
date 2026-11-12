package astcheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckFile(t *testing.T) {
	tmpDir := t.TempDir()

	validCode := `package main
func BubbleSort(arr []int) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr)-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}
`
	validFile := filepath.Join(tmpDir, "valid.go")
	if err := os.WriteFile(validFile, []byte(validCode), 0644); err != nil {
		t.Fatal(err)
	}

	cheatingCode := `package main
import "sort"

func BubbleSort(arr []int) {
	sort.Ints(arr)
}
`
	cheatingFile := filepath.Join(tmpDir, "cheating.go")
	if err := os.WriteFile(cheatingFile, []byte(cheatingCode), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("Valid Code Passes", func(t *testing.T) {
		res := CheckFile(validFile, []string{"sort", "slices"})
		if !res.Valid {
			t.Fatalf("expected valid code to pass, got violations: %v", res.Violations)
		}
	})

	t.Run("Cheating Code Fails", func(t *testing.T) {
		res := CheckFile(cheatingFile, []string{"sort", "slices"})
		if res.Valid {
			t.Fatalf("expected cheating code to fail AST check")
		}
		if len(res.Violations) == 0 {
			t.Fatalf("expected violations to be recorded")
		}
	})
}
