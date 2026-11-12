package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elvin-mark/linux-practice/internal/state"
)

func TestState_Lifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	origConfigHome := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", origConfigHome)

	st, err := state.Load()
	if err != nil {
		t.Fatalf("unexpected error loading initial state: %v", err)
	}

	st.RecordStart("test-scenario")
	if st.ActiveScenario != "test-scenario" {
		t.Errorf("expected active scenario 'test-scenario', got %q", st.ActiveScenario)
	}

	p := st.GetOrCreateProgress("test-scenario")
	if p.Attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", p.Attempts)
	}
	if p.Completed {
		t.Errorf("expected completed to be false")
	}

	st.RecordHint("test-scenario", 0)
	if p.HintsViewed != 1 {
		t.Errorf("expected 1 hint viewed, got %d", p.HintsViewed)
	}

	st.RecordCompletion("test-scenario")
	if !p.Completed {
		t.Errorf("expected completed to be true")
	}

	if err := st.Save(); err != nil {
		t.Fatalf("unexpected error saving state: %v", err)
	}

	// Verify file on disk
	expectedFile := filepath.Join(tmpDir, "linux-practice", "state.yaml")
	if _, err := os.Stat(expectedFile); err != nil {
		t.Fatalf("state file was not created: %v", err)
	}

	// Reload state
	loaded, err := state.Load()
	if err != nil {
		t.Fatalf("unexpected error reloading state: %v", err)
	}

	lp := loaded.GetOrCreateProgress("test-scenario")
	if !lp.Completed {
		t.Errorf("reloaded state expected completed=true")
	}
	if lp.Attempts != 1 {
		t.Errorf("reloaded state expected attempts=1, got %d", lp.Attempts)
	}
}
