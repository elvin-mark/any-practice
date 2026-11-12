package runner

import (
	"bytes"
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/elvin-mark/algo-practice/internal/exercise"
)

var pkgRegex = regexp.MustCompile(`^package\s+([a-zA-Z0-9_]+)`)

func TestReferenceSolutions(t *testing.T) {
	repo, err := exercise.LoadRepository("../../exercises")
	if err != nil {
		t.Fatalf("failed to load repository: %v", err)
	}

	for _, ex := range repo.Exercises {
		t.Run(ex.ID, func(t *testing.T) {
			solFile := ex.SolutionFilePath()
			solCode, err := os.ReadFile(solFile)
			if err != nil {
				t.Fatalf("failed to read solution file %s: %v", solFile, err)
			}

			// Read original entry file
			entryFile := ex.EntryFilePath()
			origCode, err := os.ReadFile(entryFile)
			if err != nil {
				t.Fatalf("failed to read entry file %s: %v", entryFile, err)
			}

			// Match original package name
			match := pkgRegex.FindSubmatch(origCode)
			if len(match) > 1 {
				targetPkg := string(match[1])
				solCode = bytes.Replace(solCode, []byte("package solution"), []byte("package "+targetPkg), 1)
			}

			// Replace entry file with solution temporarily
			if err := os.WriteFile(entryFile, solCode, 0644); err != nil {
				t.Fatalf("failed to overwrite entry file with solution: %v", err)
			}
			defer func() {
				_ = os.WriteFile(entryFile, origCode, 0644)
			}()

			res := Run(context.Background(), ex)

			if !res.ASTValid {
				t.Fatalf("solution failed AST validation: %v", res.ASTViolations)
			}

			if !res.TestsPassed {
				t.Fatalf("solution failed tests:\n%s", res.TestOutput)
			}

			if !res.BenchmarkPassed {
				t.Fatalf("solution failed benchmark:\n%s", res.BenchmarkOutput)
			}
		})
	}
}
