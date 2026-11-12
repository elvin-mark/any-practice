package runner

import (
	"bytes"
	"context"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elvin-mark/any-practice/algo-practice/internal/astcheck"
	"github.com/elvin-mark/any-practice/algo-practice/internal/exercise"
)

// Result holds the combined validation, test, and benchmark outcome.
type Result struct {
	ExerciseID      string
	Passed          bool
	ASTValid        bool
	ASTViolations   []string
	ParseErr        error
	TestsPassed     bool
	BenchmarkPassed bool
	TestOutput      string
	BenchmarkOutput string
	NsPerOp         float64
	AllocBytes      uint64
	AllocsPerOp     uint64
	Duration        time.Duration
}

var benchRegex = regexp.MustCompile(`Benchmark\S+\s+\d+\s+([0-9.]+)\s+ns/op(?:\s+(\d+)\s+B/op)?(?:\s+(\d+)\s+allocs/op)?`)

// Run executes the full validation pipeline for an exercise:
// 1. AST checking
// 2. Unit testing (with timeout)
// 3. Benchmarking (if tests pass)
func Run(ctx context.Context, ex *exercise.Exercise) *Result {
	start := time.Now()
	res := &Result{
		ExerciseID: ex.ID,
	}

	// 1. AST Validation
	astRes := astcheck.CheckFile(ex.EntryFilePath(), ex.Forbidden)
	res.ASTValid = astRes.Valid
	res.ASTViolations = astRes.Violations
	res.ParseErr = astRes.ParseErr

	if !astRes.Valid {
		res.Duration = time.Since(start)
		return res
	}

	// 2. Unit Tests
	testCtx, testCancel := context.WithTimeout(ctx, 15*time.Second)
	defer testCancel()

	cmd := exec.CommandContext(testCtx, "go", "test", "-v", "-run", "Test", ".")
	cmd.Dir = ex.Dir

	var testBuf bytes.Buffer
	cmd.Stdout = &testBuf
	cmd.Stderr = &testBuf

	err := cmd.Run()
	res.TestOutput = cleanOutput(testBuf.String())

	if testCtx.Err() == context.DeadlineExceeded {
		res.TestsPassed = false
		res.TestOutput += "\n\n[TIMEOUT] Execution exceeded 15 seconds. Check for infinite loops or non-terminating recursion."
		res.Duration = time.Since(start)
		return res
	}

	if err != nil {
		res.TestsPassed = false
		res.Duration = time.Since(start)
		return res
	}

	res.TestsPassed = true

	// 3. Benchmarks (optional, run only if tests pass)
	benchCtx, benchCancel := context.WithTimeout(ctx, 15*time.Second)
	defer benchCancel()

	bCmd := exec.CommandContext(benchCtx, "go", "test", "-run", "^$", "-bench", ".", "-benchmem", ".")
	bCmd.Dir = ex.Dir

	var bBuf bytes.Buffer
	bCmd.Stdout = &bBuf
	bCmd.Stderr = &bBuf

	bErr := bCmd.Run()
	benchOut := bBuf.String()
	res.BenchmarkOutput = cleanOutput(benchOut)

	if bErr == nil {
		res.BenchmarkPassed = true
		parseBenchmarkMetrics(benchOut, res)
	} else {
		// If no benchmarks were defined in the test file, that's not a failure
		if strings.Contains(benchOut, "testing: warning: no tests to run") || !strings.Contains(benchOut, "Benchmark") {
			res.BenchmarkPassed = true
		} else {
			res.BenchmarkPassed = false
		}
	}

	res.Passed = res.ASTValid && res.TestsPassed && res.BenchmarkPassed
	res.Duration = time.Since(start)
	return res
}

func parseBenchmarkMetrics(out string, res *Result) {
	matches := benchRegex.FindStringSubmatch(out)
	if len(matches) > 1 {
		if ns, err := strconv.ParseFloat(matches[1], 64); err == nil {
			res.NsPerOp = ns
		}
	}
	if len(matches) > 2 && matches[2] != "" {
		if b, err := strconv.ParseUint(matches[2], 10, 64); err == nil {
			res.AllocBytes = b
		}
	}
	if len(matches) > 3 && matches[3] != "" {
		if allocs, err := strconv.ParseUint(matches[3], 10, 64); err == nil {
			res.AllocsPerOp = allocs
		}
	}
}

func cleanOutput(s string) string {
	s = strings.TrimSpace(s)
	// Remove common noisy go test artifacts if desired
	return s
}
