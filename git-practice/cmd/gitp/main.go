package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/elvin-mark/any-practice/git-practice/internal/progress"
	"github.com/elvin-mark/any-practice/git-practice/internal/scenario"
	"github.com/elvin-mark/any-practice/git-practice/internal/tui"
	"github.com/elvin-mark/any-practice/git-practice/internal/ui"
)

var (
	flagScenariosDir string
	flagProgress     string
	appProgress      *progress.Progress
	appRegistry      *scenario.Registry
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gitp",
		Short: "Git Mastery, History Surgery & Disaster Recovery Practice CLI/TUI",
		Long: `Interactive CLI & TUI training platform for Git mastery:
Interactive rebases, 3-way merge conflict resolution, reflog disaster recovery,
git bisect bug hunting, and detached HEAD rescue in safe sandboxed environments.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			resolvedDir, err := scenario.ResolveScenariosDir(flagScenariosDir)
			if err != nil {
				return err
			}
			flagScenariosDir = resolvedDir

			reg, err := scenario.NewRegistry(flagScenariosDir)
			if err != nil {
				return fmt.Errorf("loading scenarios from %s: %w", flagScenariosDir, err)
			}
			appRegistry = reg

			prog, err := progress.Load(flagProgress)
			if err != nil {
				return fmt.Errorf("loading progress: %w", err)
			}
			appProgress = prog

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&flagScenariosDir, "scenarios-dir", "s", "", "Path to scenarios directory")
	rootCmd.PersistentFlags().StringVarP(&flagProgress, "progress", "p", progress.DefaultPath(), "Path to progress file")

	// Subcommands
	rootCmd.AddCommand(tuiCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(verifyCmd())
	rootCmd.AddCommand(resetCmd())
	rootCmd.AddCommand(hintCmd())
	rootCmd.AddCommand(cleanCmd())
	rootCmd.AddCommand(progressCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runTUI() error {
	m := tui.NewModel(appRegistry, appProgress)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func tuiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch interactive Git dashboard (default)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}
}

func listCmd() *cobra.Command {
	var filterCategory string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all Git drills and completion status",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.PrintBanner("GIT-PRACTICE", "Git Mastery, History Surgery & Recovery Labs")

			scenarios := appRegistry.List()
			if len(scenarios) == 0 {
				ui.PrintWarning("No scenarios found in %s", flagScenariosDir)
				return nil
			}

			columns := []ui.TableColumn{
				{Title: "#", Width: 3},
				{Title: "SCENARIO ID", Width: 26},
				{Title: "CATEGORY", Width: 14},
				{Title: "DIFFICULTY", Width: 12},
				{Title: "STATUS", Width: 8},
				{Title: "TITLE", Width: 38},
			}

			var rows [][]string
			solvedCount := 0
			idx := 1

			for _, sc := range scenarios {
				if filterCategory != "" && !strings.EqualFold(sc.Category, filterCategory) {
					continue
				}

				isDone := appProgress.IsCompleted(sc.ID)
				if isDone {
					solvedCount++
				}

				rows = append(rows, []string{
					fmt.Sprintf("%2d", idx),
					ui.CyanText.Render(sc.ID),
					ui.CategoryBadgeStyle.Render(sc.Category),
					ui.FormatDifficulty(sc.Difficulty),
					ui.FormatStatus(isDone),
					sc.Title,
				})
				idx++
			}

			ui.RenderTable(columns, rows)

			pct := 0.0
			if len(scenarios) > 0 {
				pct = float64(solvedCount) / float64(len(scenarios)) * 100
			}
			fmt.Printf("\nProgress: %s scenarios completed (%.1f%%)\n\n",
				ui.GreenText.Render(fmt.Sprintf("%d/%d", solvedCount, len(scenarios))),
				pct,
			)

			return nil
		},
	}

	cmd.Flags().StringVar(&filterCategory, "category", "", "Filter scenarios by category (Rebase, Merge, Recovery, Bisect, Branching)")
	return cmd
}

func startCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start <scenario-id>",
		Short: "Initialize and start a git practice lab",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			sc, ok := appRegistry.Get(id)
			if !ok {
				ui.PrintFailure("Scenario '%s' not found", id)
				return nil
			}

			labDir, err := scenario.SetupLab(sc)
			if err != nil {
				ui.PrintFailure("Failed to setup lab: %v", err)
				return nil
			}

			appProgress.RecordStart(id)
			appProgress.SetActive(id)
			_ = appProgress.Save()

			ui.PrintBanner("GIT-PRACTICE", sc.Title)
			fmt.Printf("\n%s %s\n", ui.BoldText.Render("Category:  "), ui.CategoryBadgeStyle.Render(sc.Category))
			fmt.Printf("%s %s\n", ui.BoldText.Render("Difficulty:"), ui.FormatDifficulty(sc.Difficulty))
			fmt.Printf("%s %s\n", ui.BoldText.Render("Status:    "), ui.FormatStatus(appProgress.IsCompleted(id)))

			ui.PrintBox("MISSION BRIEFING", sc.Description)
			ui.PrintBox("TARGET GOAL", sc.Goal)

			ui.PrintSuccess("Lab repository initialized successfully!")
			fmt.Printf("\n  %s %s\n", ui.YellowText.Render("▶ Switch to lab:"), ui.CyanText.Render("cd "+labDir))
			fmt.Printf("  %s %s\n\n", ui.YellowText.Render("▶ When done:    "), ui.CyanText.Render("gitp verify "+id))

			return nil
		},
	}
}

func verifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify [scenario-id]",
		Short: "Verify completion of a git lab",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetID := ""
			if len(args) == 1 {
				targetID = args[0]
			} else {
				targetID = appProgress.GetActive()
			}

			if targetID == "" {
				ui.PrintFailure("No active scenario specified. Use 'gitp verify <scenario-id>'")
				return nil
			}

			sc, ok := appRegistry.Get(targetID)
			if !ok {
				ui.PrintFailure("Scenario '%s' not found", targetID)
				return nil
			}

			passed, msg, err := scenario.VerifyLab(sc)
			if err != nil {
				ui.PrintFailure("Verification error: %v", err)
				return nil
			}

			if passed {
				appProgress.RecordCompletion(targetID)
				_ = appProgress.Save()
				ui.PrintSuccess("Scenario '%s' PASSED! %s", targetID, msg)
			} else {
				appProgress.RecordAttempt(targetID)
				_ = appProgress.Save()
				ui.PrintFailure("Scenario '%s' FAILED verification:\n%s", targetID, msg)
				fmt.Printf("\nNeed a hint? Run: gitp hint %s\n\n", targetID)
			}

			return nil
		},
	}
}

func resetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset [scenario-id]",
		Short: "Reset a lab workspace back to its initial state",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetID := ""
			if len(args) == 1 {
				targetID = args[0]
			} else {
				targetID = appProgress.GetActive()
			}

			if targetID == "" {
				ui.PrintFailure("No scenario specified to reset. Use 'gitp reset <scenario-id>'")
				return nil
			}

			sc, ok := appRegistry.Get(targetID)
			if !ok {
				ui.PrintFailure("Scenario '%s' not found", targetID)
				return nil
			}

			labDir, err := scenario.SetupLab(sc)
			if err != nil {
				ui.PrintFailure("Failed to reset lab: %v", err)
				return nil
			}

			ui.PrintWarning("Lab workspace for '%s' reset to initial state.", targetID)
			fmt.Printf("Workspace: %s\n\n", labDir)
			return nil
		},
	}
}

func hintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hint [scenario-id] [number]",
		Short: "Reveal hints for a git scenario",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetID := ""
			hintArg := ""

			if len(args) >= 1 {
				targetID = args[0]
			} else {
				targetID = appProgress.GetActive()
			}
			if len(args) == 2 {
				hintArg = args[1]
			}

			if targetID == "" {
				ui.PrintFailure("No scenario specified. Use 'gitp hint <scenario-id>'")
				return nil
			}

			sc, ok := appRegistry.Get(targetID)
			if !ok {
				ui.PrintFailure("Scenario '%s' not found", targetID)
				return nil
			}

			if len(sc.Hints) == 0 {
				ui.PrintInfo("No hints configured for '%s'", targetID)
				return nil
			}

			rec := appProgress.GetRecord(targetID)
			targetIndex := rec.HintsViewed
			if hintArg != "" {
				if parsed, err := strconv.Atoi(hintArg); err == nil && parsed > 0 {
					targetIndex = parsed - 1
				}
			}

			if targetIndex >= len(sc.Hints) {
				ui.PrintWarning("All %d hints have already been revealed for '%s'", len(sc.Hints), targetID)
				targetIndex = len(sc.Hints) - 1
			}

			appProgress.RecordHint(targetID, targetIndex)
			_ = appProgress.Save()

			ui.PrintHint("Hint %d/%d for '%s':", targetIndex+1, len(sc.Hints), targetID)
			fmt.Printf("  %s\n\n", sc.Hints[targetIndex])

			return nil
		},
	}
}

func cleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean [scenario-id]",
		Short: "Remove temporary lab directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetID := ""
			if len(args) == 1 {
				targetID = args[0]
			} else {
				targetID = appProgress.GetActive()
			}

			if targetID == "" {
				ui.PrintFailure("No scenario specified. Use 'gitp clean <scenario-id>'")
				return nil
			}

			sc, ok := appRegistry.Get(targetID)
			if !ok {
				ui.PrintFailure("Scenario '%s' not found", targetID)
				return nil
			}

			if err := scenario.CleanLab(sc); err != nil {
				ui.PrintFailure("Failed to clean lab: %v", err)
				return nil
			}

			ui.PrintSuccess("Cleaned up workspace for scenario '%s'.", targetID)
			return nil
		},
	}
}

func progressCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "progress",
		Short: "Display detailed Git practice learning progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.PrintBanner("GIT-PRACTICE", "Learning Track & Skill Breakdown")

			scenarios := appRegistry.List()
			if len(scenarios) == 0 {
				ui.PrintInfo("No scenarios loaded.")
				return nil
			}

			type catStat struct {
				total  int
				solved int
			}

			stats := make(map[string]*catStat)
			for _, cat := range appRegistry.Categories() {
				stats[cat] = &catStat{}
			}

			totalSolved := 0
			for _, sc := range scenarios {
				st := stats[sc.Category]
				if st == nil {
					st = &catStat{}
					stats[sc.Category] = st
				}
				st.total++
				if appProgress.IsCompleted(sc.ID) {
					st.solved++
					totalSolved++
				}
			}

			columns := []ui.TableColumn{
				{Title: "CATEGORY", Width: 18},
				{Title: "SOLVED / TOTAL", Width: 20},
				{Title: "COMPLETION", Width: 16},
			}

			var rows [][]string
			for _, cat := range appRegistry.Categories() {
				st := stats[cat]
				pct := 0.0
				if st.total > 0 {
					pct = float64(st.solved) / float64(st.total) * 100
				}
				rows = append(rows, []string{
					ui.CategoryBadgeStyle.Render(cat),
					fmt.Sprintf("%d / %d", st.solved, st.total),
					fmt.Sprintf("%.1f%%", pct),
				})
			}

			ui.RenderTable(columns, rows)

			overallPct := 0.0
			if len(scenarios) > 0 {
				overallPct = float64(totalSolved) / float64(len(scenarios)) * 100
			}

			fmt.Println()
			ui.PrintBox("OVERALL PROGRESS", fmt.Sprintf(
				"Scenarios Completed: %s (%.1f%%)\nActive Scenario:     %s\nProgress File:       %s",
				ui.GreenText.Render(fmt.Sprintf("%d / %d", totalSolved, len(scenarios))),
				overallPct,
				ui.CyanText.Render(appProgress.GetActive()),
				flagProgress,
			))
			fmt.Println()

			return nil
		},
	}
}
