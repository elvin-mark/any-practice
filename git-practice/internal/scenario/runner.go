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
	return filepath.Join(os.TempDir(), "gitp-labs", scenarioID)
}

// SetupLab initializes an authentic temporary git workspace and executes setup commands.
func SetupLab(sc *Scenario) (string, error) {
	dir := LabDir(sc.ID)

	// Clean any previous lab
	_ = os.RemoveAll(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating lab directory: %w", err)
	}

	// Prepare baseline repository setup script
	preamble := `set -e
git init -b main
git config user.name "Student"
git config user.email "student@any-practice.local"
git config commit.gpgsign false
`
	fullScript := preamble + "\n" + sc.Setup

	cmd := exec.Command("bash", "-c", fullScript)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Student",
		"GIT_AUTHOR_EMAIL=student@any-practice.local",
		"GIT_COMMITTER_NAME=Student",
		"GIT_COMMITTER_EMAIL=student@any-practice.local",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("setting up lab for %s: %w\nOutput:\n%s", sc.ID, err, out.String())
	}

	return dir, nil
}

// VerifyLab runs the scenario's verification commands inside the lab repository.
func VerifyLab(sc *Scenario) (bool, string, error) {
	dir := LabDir(sc.ID)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false, fmt.Sprintf("Lab directory %s does not exist. Run 'gitp start %s' first.", dir, sc.ID), nil
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
			output = "All verification assertions passed!"
		}
		return true, output, nil
	}

	if output == "" {
		output = fmt.Sprintf("Verification failed with exit code: %v", err)
	}
	return false, output, nil
}

// CleanLab removes the temporary lab directory.
func CleanLab(sc *Scenario) error {
	dir := LabDir(sc.ID)
	return os.RemoveAll(dir)
}
