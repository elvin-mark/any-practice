package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/elvin-mark/any-practice/sql-practice/internal/progress"
	"github.com/elvin-mark/any-practice/sql-practice/internal/scenario"
	"github.com/elvin-mark/any-practice/sql-practice/internal/tui"
	"github.com/elvin-mark/any-practice/sql-practice/internal/ui"
)

var (
	flagScenariosDir string
	flagProgress     string
	appProgress      *progress.Progress
	appRegistry      *scenario.Registry
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sqlp",
		Short: "SQL Mastery, Query Optimization & Schema Labs CLI/TUI",
		Long: `Interactive CLI & TUI training platform for SQL mastery:
Window functions, recursive CTEs, anti-joins, aggregations, index optimization,
and query execution plan tuning using fast, embedded SQLite databases.`,
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
	rootCmd.AddCommand(queryCmd())
	rootCmd.AddCommand(verifyCmd())
	rootCmd.AddCommand(tablesCmd())
	rootCmd.AddCommand(resetCmd())
	rootCmd.AddCommand(hintCmd())
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
		Short: "Launch interactive SQL dashboard (default)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}
}

func listCmd() *cobra.Command {
	var filterCategory string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all SQL labs and completion status",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.PrintBanner("SQL-PRACTICE", "SQL Mastery, Query Optimization & Schema Labs")

			scenarios := appRegistry.List()
			if len(scenarios) == 0 {
				ui.PrintWarning("No scenarios found in %s", flagScenariosDir)
				return nil
			}

			columns := []ui.TableColumn{
				{Title: "#", Width: 3},
				{Title: "SCENARIO ID", Width: 26},
				{Title: "CATEGORY", Width: 18},
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
			fmt.Printf("\nProgress: %s labs completed (%.1f%%)\n\n",
				ui.GreenText.Render(fmt.Sprintf("%d/%d", solvedCount, len(scenarios))),
				pct,
			)

			return nil
		},
	}

	cmd.Flags().StringVar(&filterCategory, "category", "", "Filter scenarios by category (Queries, Joins, Window Functions, CTEs, Indexes)")
	return cmd
}

func startCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start <scenario-id>",
		Short: "Initialize and start a SQL practice lab",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			sc, ok := appRegistry.Get(id)
			if !ok {
				ui.PrintFailure("Scenario '%s' not found", id)
				return nil
			}

			dbPath, err := scenario.SetupLab(sc)
			if err != nil {
				ui.PrintFailure("Failed to setup lab: %v", err)
				return nil
			}

			appProgress.RecordStart(id)
			appProgress.SetActive(id)
			_ = appProgress.Save()

			ui.PrintBanner("SQL-PRACTICE", sc.Title)
			fmt.Printf("\n%s %s\n", ui.BoldText.Render("Category:  "), ui.CategoryBadgeStyle.Render(sc.Category))
			fmt.Printf("%s %s\n", ui.BoldText.Render("Difficulty:"), ui.FormatDifficulty(sc.Difficulty))
			fmt.Printf("%s %s\n", ui.BoldText.Render("Status:    "), ui.FormatStatus(appProgress.IsCompleted(id)))

			ui.PrintBox("MISSION BRIEFING", sc.Description)
			ui.PrintBox("TARGET GOAL", sc.Goal)

			tables, _ := scenario.GetTables(dbPath)
			fmt.Printf("%s %s\n\n", ui.BoldText.Render("Seeded Tables:"), ui.CyanText.Render(strings.Join(tables, ", ")))

			solPath := scenario.SolutionFilePath(id)
			ui.PrintSuccess("Database ready at: %s", dbPath)
			fmt.Printf("  %s %s\n", ui.YellowText.Render("▶ Write SQL in:"), ui.CyanText.Render(solPath))
			fmt.Printf("  %s %s\n\n", ui.YellowText.Render("▶ Verify:      "), ui.CyanText.Render("sqlp verify "+id))

			return nil
		},
	}
}

func queryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "query [scenario-id] <sql-statement>",
		Short: "Execute an ad-hoc query against a scenario database and display tabular results",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetID := ""
			sqlQuery := ""

			if len(args) == 1 {
				targetID = appProgress.GetActive()
				sqlQuery = args[0]
			} else {
				targetID = args[0]
				sqlQuery = args[1]
			}

			if targetID == "" {
				ui.PrintFailure("No scenario specified. Use 'sqlp query <scenario-id> \"<sql>\"'")
				return nil
			}

			dbPath := scenario.DBPath(targetID)
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				ui.PrintFailure("Lab database not found. Run 'sqlp start %s' first.", targetID)
				return nil
			}

			cols, rows, err := scenario.ExecuteQuery(dbPath, sqlQuery)
			if err != nil {
				ui.PrintFailure("SQL Error: %v", err)
				return nil
			}

			if len(cols) == 0 {
				ui.PrintSuccess("Statement executed successfully. (0 columns returned)")
				return nil
			}

			var tblCols []ui.TableColumn
			for _, col := range cols {
				tblCols = append(tblCols, ui.TableColumn{Title: col, Width: len(col) + 4})
			}

			fmt.Println()
			ui.RenderTable(tblCols, rows)
			fmt.Printf("\n(%d rows returned)\n\n", len(rows))
			return nil
		},
	}
}

func verifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify [scenario-id] [inline-sql]",
		Short: "Verify a SQL solution against expected results",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetID := ""
			inlineSQL := ""

			if len(args) == 0 {
				targetID = appProgress.GetActive()
			} else if len(args) == 1 {
				// Could be scenario ID or inline SQL if active is set
				if _, ok := appRegistry.Get(args[0]); ok {
					targetID = args[0]
				} else {
					targetID = appProgress.GetActive()
					inlineSQL = args[0]
				}
			} else {
				targetID = args[0]
				inlineSQL = args[1]
			}

			if targetID == "" {
				ui.PrintFailure("No scenario specified. Use 'sqlp verify <scenario-id>'")
				return nil
			}

			sc, ok := appRegistry.Get(targetID)
			if !ok {
				ui.PrintFailure("Scenario '%s' not found", targetID)
				return nil
			}

			passed, msg, err := scenario.VerifyLab(sc, inlineSQL)
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
				ui.PrintFailure("Scenario '%s' FAILED:\n%s", targetID, msg)
				fmt.Printf("\nNeed a hint? Run: sqlp hint %s\n\n", targetID)
			}

			return nil
		},
	}
}

func tablesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tables [scenario-id]",
		Short: "List tables and schemas in a scenario database",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetID := ""
			if len(args) == 1 {
				targetID = args[0]
			} else {
				targetID = appProgress.GetActive()
			}

			if targetID == "" {
				ui.PrintFailure("No scenario specified. Use 'sqlp tables <scenario-id>'")
				return nil
			}

			dbPath := scenario.DBPath(targetID)
			tables, err := scenario.GetTables(dbPath)
			if err != nil {
				ui.PrintFailure("Error fetching tables: %v", err)
				return nil
			}

			ui.PrintBanner("SQL-PRACTICE", fmt.Sprintf("Tables in %s", targetID))
			for _, tbl := range tables {
				_, rows, _ := scenario.ExecuteQuery(dbPath, fmt.Sprintf("SELECT sql FROM sqlite_master WHERE type='table' AND name='%s';", tbl))
				schema := ""
				if len(rows) > 0 && len(rows[0]) > 0 {
					schema = rows[0][0]
				}
				ui.PrintBox("TABLE: "+tbl, schema)
			}
			return nil
		},
	}
}

func resetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset [scenario-id]",
		Short: "Reset a scenario's database back to its initial seeded state",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetID := ""
			if len(args) == 1 {
				targetID = args[0]
			} else {
				targetID = appProgress.GetActive()
			}

			if targetID == "" {
				ui.PrintFailure("No scenario specified to reset. Use 'sqlp reset <scenario-id>'")
				return nil
			}

			sc, ok := appRegistry.Get(targetID)
			if !ok {
				ui.PrintFailure("Scenario '%s' not found", targetID)
				return nil
			}

			dbPath, err := scenario.SetupLab(sc)
			if err != nil {
				ui.PrintFailure("Failed to reset database: %v", err)
				return nil
			}

			ui.PrintWarning("Database for '%s' reset to initial state.", targetID)
			fmt.Printf("Database: %s\n\n", dbPath)
			return nil
		},
	}
}

func hintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hint [scenario-id] [number]",
		Short: "Reveal hints for a SQL lab",
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
				ui.PrintFailure("No scenario specified. Use 'sqlp hint <scenario-id>'")
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

func progressCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "progress",
		Short: "Display detailed SQL learning progress and domain mastery",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.PrintBanner("SQL-PRACTICE", "Learning Track & Skill Breakdown")

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
				{Title: "CATEGORY", Width: 20},
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
				"Labs Completed:  %s (%.1f%%)\nActive Scenario: %s\nProgress File:   %s",
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
