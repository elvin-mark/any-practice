package scenario_test

import (
	"testing"

	"github.com/elvin-mark/linux-practice/embedded"
	"github.com/elvin-mark/linux-practice/internal/scenario"
)

func TestRegistry_LoadEmbedded(t *testing.T) {
	reg := scenario.NewRegistry()
	if err := reg.LoadEmbedded(embedded.ScenariosFS, "scenarios"); err != nil {
		t.Fatalf("unexpected error loading embedded scenarios: %v", err)
	}

	all := reg.All()
	if len(all) < 6 {
		t.Errorf("expected at least 6 embedded scenarios, got %d", len(all))
	}

	tracks := reg.Tracks()
	if len(tracks) == 0 {
		t.Errorf("expected non-empty tracks list")
	}

	// Verify required starter scenarios exist
	expectedIDs := []string{
		"perm-suid-cleanup",
		"user-sudo-syntax",
		"systemd-failed-service",
		"systemd-custom-timer",
		"net-ufw-lockdown",
		"troubleshoot-port-conflict",
	}

	for _, id := range expectedIDs {
		sc, found := reg.Get(id)
		if !found {
			t.Errorf("expected scenario %q to be registered", id)
			continue
		}
		if sc.Title == "" {
			t.Errorf("scenario %q title is empty", id)
		}
		if sc.SetupScript == "" {
			t.Errorf("scenario %q setup script is empty", id)
		}
		if sc.VerifyScript == "" {
			t.Errorf("scenario %q verify script is empty", id)
		}
	}
}

func TestScenario_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sc      scenario.Scenario
		wantErr bool
	}{
		{
			name:    "empty scenario",
			sc:      scenario.Scenario{},
			wantErr: true,
		},
		{
			name: "missing setup script",
			sc: scenario.Scenario{
				ID:           "test-id",
				Title:        "Test Title",
				Track:        "systemd",
				Description:  "Desc",
				VerifyScript: "echo ok",
			},
			wantErr: true,
		},
		{
			name: "valid scenario",
			sc: scenario.Scenario{
				ID:           "test-id",
				Title:        "Test Title",
				Track:        "systemd",
				Description:  "Desc",
				SetupScript:  "echo setup",
				VerifyScript: "echo ok",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sc.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
