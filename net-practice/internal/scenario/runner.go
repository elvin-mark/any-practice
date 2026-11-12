package scenario

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// LabDir returns the workspace directory path for a given scenario.
func LabDir(scenarioID string) string {
	return filepath.Join(os.TempDir(), "netp-labs", scenarioID)
}

// SetupLab initializes the workspace directory and executes scenario setup commands.
func SetupLab(sc *Scenario) (string, error) {
	dir := LabDir(sc.ID)

	_ = os.RemoveAll(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating lab directory: %w", err)
	}

	preamble := "set -e\n"
	fullScript := preamble + sc.Setup

	cmd := exec.Command("bash", "-c", fullScript)
	cmd.Dir = dir
	cmd.Env = os.Environ()

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("setting up lab for %s: %w\nOutput:\n%s", sc.ID, err, out.String())
	}

	return dir, nil
}

// VerifyLab runs the scenario's verification commands inside the lab directory.
func VerifyLab(sc *Scenario) (bool, string, error) {
	dir := LabDir(sc.ID)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false, fmt.Sprintf("Lab directory %s does not exist. Run 'netp start %s' first.", dir, sc.ID), nil
	}

	cmd := exec.Command("bash", "-c", sc.Verify)
	cmd.Dir = dir
	cmd.Env = os.Environ()

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	output := strings.TrimSpace(out.String())

	if err == nil {
		if output == "" {
			output = "All network assertions and connectivity checks passed!"
		}
		return true, output, nil
	}

	if output == "" {
		output = fmt.Sprintf("Verification failed with exit code: %v", err)
	}
	return false, output, nil
}

// CleanLab removes the temporary lab workspace.
func CleanLab(sc *Scenario) error {
	dir := LabDir(sc.ID)
	return os.RemoveAll(dir)
}
