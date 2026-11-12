package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elvin-mark/linux-practice/embedded"
	"github.com/elvin-mark/linux-practice/internal/config"
	"github.com/elvin-mark/linux-practice/internal/container"
	"github.com/elvin-mark/linux-practice/internal/scenario"
	"github.com/elvin-mark/linux-practice/internal/state"
	"github.com/elvin-mark/linux-practice/internal/tui"
	"github.com/elvin-mark/linux-practice/internal/ui"
	"github.com/spf13/cobra"
)

var (
	registry *scenario.Registry
	appState *state.State
	manager  *container.Manager
)

// Execute initializes and runs the root cobra command.
func Execute() error {
	var err error
	appState, err = state.Load()
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	registry = scenario.NewRegistry()
	if err := registry.LoadAll(embedded.ScenariosFS); err != nil {
		return fmt.Errorf("loading scenarios: %w", err)
	}

	rootCmd := &cobra.Command{
		Use:   "lp [command]",
		Short: "Linux Practice - Master Linux administration hands-on with realistic scenarios",
		Long: `lp (linux-practice) is a hands-on training tool for Linux systems administration.
It provisions isolated Ubuntu systemd containers, injects real-world problems and misconfigurations,
and validates your solutions automatically.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := container.NewManager()
			if err != nil {
				return err
			}
			manager = mgr

			m := tui.NewModel(registry, appState, manager)
			p := tea.NewProgram(m, tea.WithAltScreen())
			_, err = p.Run()
			return err
		},
	}

	rootCmd.AddCommand(
		newListCmd(),
		newStartCmd(),
		newAttachCmd(),
		newVerifyCmd(),
		newHintCmd(),
		newResetCmd(),
		newStopCmd(),
		newProgressCmd(),
		newCreateCmd(),
		newValidateCmd(),
	)

	return rootCmd.Execute()
}

func getManager() (*container.Manager, error) {
	if manager != nil {
		return manager, nil
	}
	var err error
	manager, err = container.NewManager()
	return manager, err
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available practice scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			scenarios := registry.All()
			ui.PrintBanner("LINUX-PRACTICE", "Hands-on Linux Systems Administration Labs")
			ui.PrintHeader("AVAILABLE PRACTICE SCENARIOS")

			columns := []ui.TableColumn{
				{Title: "#", Width: 3},
				{Title: "SCENARIO ID", Width: 26},
				{Title: "TRACK", Width: 16},
				{Title: "DIFFICULTY", Width: 10},
				{Title: "STATUS", Width: 6},
				{Title: "TITLE", Width: 46},
			}

			var rows [][]string
			totalCount := len(scenarios)
			completedCount := 0

			for i, sc := range scenarios {
				isDone := false
				if p, ok := appState.Scenarios[sc.ID]; ok && p.Completed {
					isDone = true
					completedCount++
				}
				rows = append(rows, []string{
					fmt.Sprintf("%2d", i+1),
					ui.CyanText.Render(sc.ID),
					sc.Track,
					ui.FormatDifficulty(sc.Difficulty),
					ui.FormatStatus(isDone),
					sc.Title,
				})
			}

			ui.RenderTable(columns, rows)

			pct := 0.0
			if totalCount > 0 {
				pct = float64(completedCount) / float64(totalCount) * 100
			}

			fmt.Printf("\nTotal: %s scenarios | Completed: %s | Progress: %s\n",
				ui.CyanText.Render(fmt.Sprintf("%d", totalCount)),
				ui.GreenText.Render(fmt.Sprintf("%d/%d", completedCount, totalCount)),
				ui.YellowText.Render(fmt.Sprintf("%.1f%%", pct)),
			)
			fmt.Printf("\nTo start a scenario: %s (e.g. %s)\n\n",
				ui.CyanText.Render("lp start <scenario-id>"),
				ui.YellowText.Render("lp start net-dns-failure"))
			return nil
		},
	}
}

func newStartCmd() *cobra.Command {
	var noAttach bool
	cmd := &cobra.Command{
		Use:   "start <scenario-id>",
		Short: "Start a scenario and open an interactive terminal shell",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			sc, ok := registry.Get(id)
			if !ok {
				return fmt.Errorf("scenario %q not found. Run 'lp list' to see available scenarios", id)
			}

			mgr, err := getManager()
			if err != nil {
				return err
			}

			ctx := context.Background()
			if err := mgr.EnsureBaseImage(ctx, embedded.Dockerfile, config.DefaultImageName); err != nil {
				return fmt.Errorf("ensuring base image: %w", err)
			}

			ui.PrintBanner("LINUX-PRACTICE", "Hands-on Linux Systems Administration Labs")
			ui.PrintInfo("Starting scenario: %s (%s)", sc.Title, sc.ID)
			if err := mgr.StartScenario(ctx, sc); err != nil {
				return fmt.Errorf("starting scenario: %w", err)
			}

			appState.RecordStart(sc.ID)
			_ = appState.Save()

			ui.PrintBox("INSTRUCTIONS", sc.Description)

			if noAttach {
				ui.PrintInfo("Container lp-%s is running. Run 'lp attach' to connect or 'lp verify' to test.", sc.ID)
				return nil
			}

			ui.PrintInfo("Dropping into interactive shell as 'student' (sudo available)...")
			ui.PrintInfo("Type 'exit' when ready, then run 'lp verify' to test your solution.")
			fmt.Println(ui.SubtleText.Render(strings.Repeat("─", 60)))

			return mgr.Attach(sc.ID)
		},
	}
	cmd.Flags().BoolVar(&noAttach, "no-attach", false, "Start the scenario without attaching interactive terminal")
	return cmd
}

func newAttachCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "attach [scenario-id]",
		Short: "Re-attach to the active or specified scenario container",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := appState.ActiveScenario
			if len(args) > 0 {
				id = args[0]
			}
			if id == "" {
				return fmt.Errorf("no active scenario. Provide a scenario ID or run 'lp start <id>'")
			}

			mgr, err := getManager()
			if err != nil {
				return err
			}

			ctx := context.Background()
			if !mgr.IsRunning(ctx, id) {
				return fmt.Errorf("container for %q is not running. Start it with 'lp start %s'", id, id)
			}

			fmt.Printf("Re-attaching to %s shell...\n", id)
			return mgr.Attach(id)
		},
	}
}

func newVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify [scenario-id]",
		Short: "Test whether you have successfully resolved the scenario",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := appState.ActiveScenario
			if len(args) > 0 {
				id = args[0]
			}
			if id == "" {
				return fmt.Errorf("no active scenario specified. Use 'lp verify <scenario-id>'")
			}

			sc, ok := registry.Get(id)
			if !ok {
				return fmt.Errorf("scenario %q not found", id)
			}

			mgr, err := getManager()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ui.PrintHeader(fmt.Sprintf("VERIFYING: %s (%s)", sc.Title, sc.ID))
			ui.PrintInfo("Evaluating verification checks inside container...")
			passed, out, err := mgr.Verify(ctx, sc)
			if err != nil {
				ui.PrintFailure("Verification execution error: %v", err)
				return err
			}

			if out != "" {
				fmt.Println(out)
			}
			if passed {
				appState.RecordCompletion(sc.ID)
				_ = appState.Save()
				ui.PrintSuccess("Scenario completed successfully!")
			} else {
				ui.PrintFailure("Verification failed. Review the output above or run 'lp hint' for guidance.")
			}

			return nil
		},
	}
}

func newHintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hint [scenario-id]",
		Short: "Reveal progressive hints for the current challenge",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := appState.ActiveScenario
			if len(args) > 0 {
				id = args[0]
			}
			if id == "" {
				return fmt.Errorf("no active scenario. Specify scenario ID or run 'lp start <id>'")
			}

			sc, ok := registry.Get(id)
			if !ok {
				return fmt.Errorf("scenario %q not found", id)
			}

			if len(sc.Hints) == 0 {
				ui.PrintInfo("No hints available for this scenario.")
				return nil
			}

			ui.PrintHeader(fmt.Sprintf("HINTS: %s", sc.Title))
			p := appState.GetOrCreateProgress(sc.ID)
			hintIdx := p.HintsViewed
			if hintIdx >= len(sc.Hints) {
				ui.PrintInfo("All available hints have already been revealed:")
				for i, h := range sc.Hints {
					ui.PrintHint("Hint %d/%d: %s", i+1, len(sc.Hints), h)
				}
				return nil
			}

			ui.PrintHint("Hint %d/%d: %s", hintIdx+1, len(sc.Hints), sc.Hints[hintIdx])
			appState.RecordHint(sc.ID, hintIdx)
			_ = appState.Save()
			return nil
		},
	}
}

func newResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset [scenario-id]",
		Short: "Reset the scenario container back to its initial broken state",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := appState.ActiveScenario
			if len(args) > 0 {
				id = args[0]
			}
			if id == "" {
				return fmt.Errorf("no scenario specified to reset")
			}

			sc, ok := registry.Get(id)
			if !ok {
				return fmt.Errorf("scenario %q not found", id)
			}

			mgr, err := getManager()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ui.PrintInfo("Resetting scenario %s to initial state...", sc.ID)
			if err := mgr.StartScenario(ctx, sc); err != nil {
				ui.PrintFailure("Resetting scenario: %v", err)
				return err
			}

			ui.PrintSuccess("Scenario '%s' has been reset to starting state. Re-attach with 'lp attach'.", sc.ID)
			return nil
		},
	}
}

func newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [scenario-id]",
		Short: "Stop and tear down the practice container",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := appState.ActiveScenario
			if len(args) > 0 {
				id = args[0]
			}
			if id == "" {
				return fmt.Errorf("no scenario specified to stop")
			}

			mgr, err := getManager()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ui.PrintInfo("Stopping scenario container lp-%s...", id)
			if err := mgr.StopScenario(ctx, id); err != nil {
				ui.PrintFailure("Stopping scenario: %v", err)
				return err
			}

			if appState.ActiveScenario == id {
				appState.ClearActive()
				_ = appState.Save()
			}

			ui.PrintSuccess("Scenario container lp-%s stopped and removed.", id)
			return nil
		},
	}
}

func newProgressCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "progress",
		Short: "Display your overall practice progress and completion statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			scenarios := registry.All()
			total := len(scenarios)
			completed := 0

			trackStats := make(map[string]struct{ total, completed int })

			for _, sc := range scenarios {
				ts := trackStats[sc.Track]
				ts.total++
				if p, ok := appState.Scenarios[sc.ID]; ok && p.Completed {
					completed++
					ts.completed++
				}
				trackStats[sc.Track] = ts
			}

			pct := 0.0
			if total > 0 {
				pct = float64(completed) / float64(total) * 100
			}

			ui.PrintBanner("LINUX-PRACTICE", "Hands-on Linux Systems Administration Labs")
			ui.PrintHeader("PROGRESS SUMMARY")

			columns := []ui.TableColumn{
				{Title: "TRACK", Width: 20},
				{Title: "COMPLETED", Width: 12},
				{Title: "TOTAL", Width: 8},
				{Title: "PERCENTAGE", Width: 12},
			}

			var rows [][]string
			for track, ts := range trackStats {
				tPct := float64(ts.completed) / float64(ts.total) * 100
				rows = append(rows, []string{
					strings.ToUpper(track),
					fmt.Sprintf("%d", ts.completed),
					fmt.Sprintf("%d", ts.total),
					fmt.Sprintf("%.0f%%", tPct),
				})
			}

			ui.RenderTable(columns, rows)

			fmt.Printf("\nOverall Progress: %s | Completed: %s\n\n",
				ui.YellowText.Render(fmt.Sprintf("%.1f%%", pct)),
				ui.GreenText.Render(fmt.Sprintf("%d/%d", completed, total)),
			)
			return nil
		},
	}
}

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <scenario-id>",
		Short: "Scaffold a new practice scenario template directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			dir := filepath.Join(".", id)

			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("creating scenario directory: %w", err)
			}

			manifestTemplate := fmt.Sprintf(`id: %s
title: "Sample Scenario Title"
track: systemd # permissions | users | systemd | networking | troubleshooting
difficulty: intermediate # beginner | intermediate | advanced
description: |
  Describe the scenario, incident symptoms, and what the user is tasked to resolve.
tags: [sample, tag]
estimated_minutes: 10
author: Community Contributor
container:
  image: "linux-practice-ubuntu:24.04"
  privileged: true
hints:
  - "First hint to guide the learner."
  - "Second hint with more specific commands."
`, id)

			setupTemplate := `#!/bin/bash
set -euo pipefail

# Provision broken state or test environment here
echo "Running setup script for scenario..."
`

			verifyTemplate := `#!/bin/bash
set -eo pipefail

# Test assertions. Exit 0 for success, non-zero for failure.
echo "Testing solution..."
if true; then
  echo "✅ Success: Challenge resolved!"
  exit 0
else
  echo "❌ Error: Condition not met."
  exit 1
fi
`

			if err := os.WriteFile(filepath.Join(dir, "scenario.yaml"), []byte(manifestTemplate), 0644); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(dir, "setup.sh"), []byte(setupTemplate), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(dir, "verify.sh"), []byte(verifyTemplate), 0755); err != nil {
				return err
			}

			fmt.Printf("✨ Created new scenario scaffold in ./%s/\n", id)
			fmt.Println("Files created:")
			fmt.Printf("  • %s/scenario.yaml\n", id)
			fmt.Printf("  • %s/setup.sh\n", id)
			fmt.Printf("  • %s/verify.sh\n", id)
			return nil
		},
	}
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <path>",
		Short: "Validate the manifest and scripts of a scenario directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetPath := args[0]
			sc, err := scenario.LoadFromDisk(targetPath)
			if err != nil {
				return fmt.Errorf("validation error: %w", err)
			}

			if err := sc.Validate(); err != nil {
				return fmt.Errorf("invalid scenario: %w", err)
			}

			fmt.Printf("✅ Scenario %q (%s) is valid!\n", sc.ID, sc.Title)
			fmt.Printf("  Track: %s | Difficulty: %s\n", sc.Track, sc.Difficulty)
			fmt.Printf("  Setup Script: %d bytes | Verify Script: %d bytes\n", len(sc.SetupScript), len(sc.VerifyScript))
			return nil
		},
	}
}
