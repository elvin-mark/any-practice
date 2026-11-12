package scenario

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/elvin-mark/k8s-practice/internal/ui"
)

// Setup prepares the scenario environment on the cluster
func Setup(ctx context.Context, s *Scenario, contextOverride string) error {
	targetContext := s.TargetContext
	if contextOverride != "" {
		targetContext = contextOverride
	}

	setupScript := filepath.Join(s.Path, "setup.sh")
	if _, err := os.Stat(setupScript); os.IsNotExist(err) {
		// If no setup.sh, ensure the target namespace exists if specified
		if s.TargetNamespace != "" && s.TargetNamespace != "default" && s.TargetNamespace != "kube-system" {
			ensureNsCmd := exec.CommandContext(ctx, "kubectl", "--context", targetContext, "create", "namespace", s.TargetNamespace, "--dry-run=client", "-o", "yaml")
			applyCmd := exec.CommandContext(ctx, "kubectl", "--context", targetContext, "apply", "-f", "-")
			pipe, _ := ensureNsCmd.StdoutPipe()
			applyCmd.Stdin = pipe
			_ = ensureNsCmd.Start()
			_ = applyCmd.Run()
		}
		return nil
	}

	cmd := exec.CommandContext(ctx, "bash", "setup.sh")
	cmd.Dir = s.Path
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("TARGET_CONTEXT=%s", targetContext),
		fmt.Sprintf("TARGET_NAMESPACE=%s", s.TargetNamespace),
		fmt.Sprintf("SCENARIO_ID=%s", s.ID),
		fmt.Sprintf("SCENARIO_DIR=%s", s.Path),
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	ui.PrintInfo("Applying scenario setup environment...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scenario setup failed: %s\nOutput:\n%s", err, output.String())
	}

	return nil
}

// Verify runs the verification script to grade the scenario
func Verify(ctx context.Context, s *Scenario, contextOverride string) (*VerificationResult, error) {
	targetContext := s.TargetContext
	if contextOverride != "" {
		targetContext = contextOverride
	}

	verifyScript := filepath.Join(s.Path, "verify.sh")
	if _, err := os.Stat(verifyScript); os.IsNotExist(err) {
		return nil, fmt.Errorf("verification script 'verify.sh' not found in %s", s.Path)
	}

	// Give a timeout for verification execution
	vCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	cmd := exec.CommandContext(vCtx, "bash", "verify.sh")
	cmd.Dir = s.Path
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("TARGET_CONTEXT=%s", targetContext),
		fmt.Sprintf("TARGET_NAMESPACE=%s", s.TargetNamespace),
		fmt.Sprintf("SCENARIO_ID=%s", s.ID),
		fmt.Sprintf("SCENARIO_DIR=%s", s.Path),
	)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()
	rawOutput := output.String()

	result := &VerificationResult{
		Output: rawOutput,
		Passed: err == nil,
	}

	// Parse individual check lines if script outputs [PASS] or [FAIL]
	lines := strings.Split(rawOutput, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "[PASS]") {
			result.Checks = append(result.Checks, CheckResult{
				Name:   strings.TrimPrefix(l, "[PASS]"),
				Passed: true,
			})
		} else if strings.HasPrefix(l, "[FAIL]") {
			result.Checks = append(result.Checks, CheckResult{
				Name:   strings.TrimPrefix(l, "[FAIL]"),
				Passed: false,
			})
		}
	}

	if result.Passed {
		result.Message = fmt.Sprintf("All checks passed successfully! (+%d points)", s.Weight)
	} else {
		result.Message = "One or more verification checks failed. Review details below."
	}

	return result, nil
}
