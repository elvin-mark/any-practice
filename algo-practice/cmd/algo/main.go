package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/elvin-mark/any-practice/algo-practice/internal/exercise"
	"github.com/elvin-mark/any-practice/algo-practice/internal/progress"
	"github.com/elvin-mark/any-practice/algo-practice/internal/runner"
	"github.com/elvin-mark/any-practice/algo-practice/internal/scaffold"
	"github.com/elvin-mark/any-practice/algo-practice/internal/tui"
	"github.com/elvin-mark/any-practice/algo-practice/internal/ui"
)

var (
	exercisesDir string
	progressFile string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "algo",
		Short: "algo is an interactive Go algorithm practice platform",
		Long:  "Master essential algorithms in Go through hands-on practice, live test feedback, complexity benchmarks, and AST validation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatch()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&exercisesDir, "dir", "d", "exercises", "Directory containing exercises")
	rootCmd.PersistentFlags().StringVarP(&progressFile, "progress", "p", progress.DefaultPath(), "Path to progress file")

	// Subcommands
	rootCmd.AddCommand(watchCmd())
	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(verifyCmd())
	rootCmd.AddCommand(hintCmd())
	rootCmd.AddCommand(solutionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func watchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Start the interactive TUI workbook watcher (default)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatch()
		},
	}
}

func runWatch() error {
	repo, prog, err := loadRepoAndProgress()
	if err != nil {
		return err
	}

	model := tui.NewModel(repo, prog)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("running TUI: %w", err)
	}
	return nil
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init [target-dir]",
		Short: "Initialize exercises in the target directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir := "."
			if len(args) > 0 {
				targetDir = args[0]
			}

			ui.PrintBanner("ALGO-PRACTICE", "Interactive Go Algorithm Workbook & Benchmark Platform")
			ui.PrintInfo("Extracting algorithm exercises into %s/exercises ...", targetDir)
			if err := scaffold.UnpackExercises(targetDir); err != nil {
				return fmt.Errorf("unpacking exercises: %w", err)
			}
			ui.PrintSuccess("Successfully initialized exercises!")
			fmt.Println("\nTo begin practicing:")
			fmt.Printf("  %s\n", ui.CyanText.Render("algo watch"))
			return nil
		},
	}
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all algorithm tracks and exercises",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, prog, err := loadRepoAndProgress()
			if err != nil {
				return err
			}

			ui.PrintBanner("ALGO-PRACTICE", "Interactive Go Algorithm Workbook & Benchmark Platform")
			ui.PrintHeader("ALGORITHM PRACTICE CURRICULUM")

			columns := []ui.TableColumn{
				{Title: "#", Width: 3},
				{Title: "EXERCISE ID", Width: 24},
				{Title: "TRACK", Width: 26},
				{Title: "DIFFICULTY", Width: 10},
				{Title: "STATUS", Width: 6},
				{Title: "NAME", Width: 32},
			}

			var rows [][]string
			solvedCount := prog.Stats()
			totalCount := len(repo.Exercises)

			for i, ex := range repo.Exercises {
				rows = append(rows, []string{
					fmt.Sprintf("%2d", i+1),
					ui.CyanText.Render(ex.ID),
					ex.Track,
					ui.FormatDifficulty(ex.Difficulty),
					ui.FormatStatus(prog.IsCompleted(ex.ID)),
					ex.Name,
				})
			}

			ui.RenderTable(columns, rows)

			pct := 0.0
			if totalCount > 0 {
				pct = float64(solvedCount) / float64(totalCount) * 100
			}

			fmt.Printf("\nTotal: %s exercises | Solved: %s | Progress: %s\n",
				ui.CyanText.Render(fmt.Sprintf("%d", totalCount)),
				ui.GreenText.Render(fmt.Sprintf("%d/%d", solvedCount, totalCount)),
				ui.YellowText.Render(fmt.Sprintf("%.1f%%", pct)),
			)
			fmt.Printf("\nTo begin interactive practice: %s\n\n", ui.CyanText.Render("algo watch"))
			return nil
		},
	}
}

func verifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify [exercise-id]",
		Short: "Verify an exercise from the command line without opening the TUI",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, prog, err := loadRepoAndProgress()
			if err != nil {
				return err
			}

			var targetID string
			if len(args) > 0 {
				targetID = args[0]
			} else {
				targetID = prog.GetCurrent()
			}

			ex, ok := repo.FindByID(targetID)
			if !ok {
				return fmt.Errorf("exercise %q not found", targetID)
			}

			ui.PrintHeader(fmt.Sprintf("VERIFYING: %s (%s)", ex.Name, ex.ID))
			ui.PrintInfo("Running AST validation, unit tests, and complexity benchmarks...")
			res := runner.Run(context.Background(), ex)

			if !res.ASTValid {
				ui.PrintFailure("AST Validation Failed: Forbidden Standard Library Usage")
				for _, v := range res.ASTViolations {
					fmt.Printf("  • %s\n", v)
				}
				os.Exit(1)
			}

			if !res.TestsPassed {
				ui.PrintFailure("Unit Tests Failed:")
				fmt.Println(res.TestOutput)
				os.Exit(1)
			}

			ui.PrintSuccess("All checks passed!")
			if res.NsPerOp > 0 {
				ui.PrintInfo("Benchmark: %.2f ns/op (%d B/op, %d allocs/op)", res.NsPerOp, res.AllocBytes, res.AllocsPerOp)
			}

			prog.MarkCompleted(ex.ID, res.NsPerOp, res.AllocBytes, res.AllocsPerOp)
			_ = prog.Save()
			return nil
		},
	}
}

func hintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hint [exercise-id] [hint-number]",
		Short: "Display hints for an exercise",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, prog, err := loadRepoAndProgress()
			if err != nil {
				return err
			}

			var targetID string
			if len(args) > 0 {
				targetID = args[0]
			} else {
				targetID = prog.GetCurrent()
			}

			ex, ok := repo.FindByID(targetID)
			if !ok {
				return fmt.Errorf("exercise %q not found", targetID)
			}

			if len(ex.Hints) == 0 {
				ui.PrintInfo("No hints available for this exercise.")
				return nil
			}

			ui.PrintHeader(fmt.Sprintf("HINTS: %s", ex.Name))
			for i, h := range ex.Hints {
				ui.PrintHint("Hint %d/%d:\n%s", i+1, len(ex.Hints), h)
			}
			return nil
		},
	}
}

func solutionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "solution [exercise-id]",
		Short: "Show the reference solution for an exercise",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, prog, err := loadRepoAndProgress()
			if err != nil {
				return err
			}

			var targetID string
			if len(args) > 0 {
				targetID = args[0]
			} else {
				targetID = prog.GetCurrent()
			}

			ex, ok := repo.FindByID(targetID)
			if !ok {
				return fmt.Errorf("exercise %q not found", targetID)
			}

			solPath := ex.SolutionFilePath()
			data, err := os.ReadFile(solPath)
			if err != nil {
				return fmt.Errorf("reading solution file: %w", err)
			}

			ui.PrintHeader(fmt.Sprintf("REFERENCE SOLUTION: %s", ex.Name))
			ui.PrintBox(ex.Name, string(data))
			return nil
		},
	}
}

func resolveExercisesDir() string {
	if _, err := os.Stat(exercisesDir); err == nil {
		return exercisesDir
	}
	if exercisesDir == "exercises" {
		if _, err := os.Stat("algo-practice/exercises"); err == nil {
			return "algo-practice/exercises"
		}
		if exe, err := os.Executable(); err == nil {
			dir := filepath.Join(filepath.Dir(exe), "exercises")
			if _, err := os.Stat(dir); err == nil {
				return dir
			}
			parentDir := filepath.Join(filepath.Dir(filepath.Dir(exe)), "exercises")
			if _, err := os.Stat(parentDir); err == nil {
				return parentDir
			}
			repoDir := filepath.Join(filepath.Dir(filepath.Dir(exe)), "algo-practice", "exercises")
			if _, err := os.Stat(repoDir); err == nil {
				return repoDir
			}
		}
	}
	return exercisesDir
}

func loadRepoAndProgress() (*exercise.Repository, *progress.Progress, error) {
	effectiveDir := resolveExercisesDir()
	if _, err := os.Stat(effectiveDir); os.IsNotExist(err) {
		// Attempt to auto-unpack if exercises directory doesn't exist yet
		fmt.Println("Exercises directory not found. Initializing exercises from template...")
		if err := scaffold.UnpackExercises("."); err != nil {
			return nil, nil, fmt.Errorf("auto-initializing exercises: %w", err)
		}
		effectiveDir = "exercises"
	}

	repo, err := exercise.LoadRepository(effectiveDir)
	if err != nil {
		return nil, nil, fmt.Errorf("loading exercises from %s: %w", effectiveDir, err)
	}

	prog, err := progress.Load(progressFile)
	if err != nil {
		return nil, nil, fmt.Errorf("loading progress from %s: %w", progressFile, err)
	}

	return repo, prog, nil
}
