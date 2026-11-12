package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/elvin-mark/any-practice/k8s-practice/internal/cluster"
	"github.com/elvin-mark/any-practice/k8s-practice/internal/exam"
	"github.com/elvin-mark/any-practice/k8s-practice/internal/progress"
	"github.com/elvin-mark/any-practice/k8s-practice/internal/scenario"
	"github.com/elvin-mark/any-practice/k8s-practice/internal/tui"
	"github.com/elvin-mark/any-practice/k8s-practice/internal/ui"
)

var (
	scenariosDir string
	flagContext  string
	flagDomain   string
	flagDuration int
	flagProgress string
	appProgress  *progress.Progress
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cka",
		Short: "CKA (Certified Kubernetes Administrator) Interactive Mock Exam & Practice CLI",
		Long: `An interactive CLI tool for CKA exam preparation using local Kind clusters.
Practice all 5 exam domains with instant verification, or simulate a full 2-hour timed exam.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			appProgress, err = progress.Load(flagProgress)
			if err != nil {
				return fmt.Errorf("loading progress: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&scenariosDir, "scenarios-dir", "s", findScenariosDir(), "Path to scenarios directory")
	rootCmd.PersistentFlags().StringVar(&flagContext, "context", "", "Override target Kubernetes context")
	rootCmd.PersistentFlags().StringVarP(&flagProgress, "progress", "p", progress.DefaultPath(), "Path to progress state file")

	// Cluster commands
	clusterCmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage local Kind clusters for CKA practice",
	}

	var withCalico bool
	initClusterCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize dedicated multi-node Kind clusters with Calico CNI",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			if err := cluster.CheckPrerequisites(); err != nil {
				return err
			}

			ui.PrintBanner()
			ui.PrintHeader("INITIALIZING CKA KIND CLUSTERS")

			// Cluster 1 (Control-plane + Worker)
			if err := cluster.CreateCluster(ctx, cluster.DefaultCluster1, true, withCalico); err != nil {
				return fmt.Errorf("creating cluster 1: %w", err)
			}

			// Cluster 2 (Control-plane only, for multi-context exam tasks)
			if err := cluster.CreateCluster(ctx, cluster.DefaultCluster2, false, false); err != nil {
				return fmt.Errorf("creating cluster 2: %w", err)
			}

			// Set active context to cluster 1 by default
			_ = cluster.SwitchContext(ctx, fmt.Sprintf("kind-%s", cluster.DefaultCluster1))

			ui.PrintSuccess("All Kind clusters successfully initialized and ready for CKA practice!")
			return nil
		},
	}
	initClusterCmd.Flags().BoolVar(&withCalico, "calico", true, "Install Calico CNI for NetworkPolicy support")

	statusClusterCmd := &cobra.Command{
		Use:   "status",
		Short: "Display status of Kind clusters and nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ui.PrintBanner()
			ui.PrintHeader("KUBERNETES CLUSTERS STATUS")

			clusters, err := cluster.ListClusters(ctx)
			if err != nil {
				return err
			}

			if len(clusters) == 0 {
				ui.PrintWarning("No Kind clusters found. Run 'cka cluster init' to create them.")
				return nil
			}

			currentCtx, _ := cluster.GetCurrentContext(ctx)
			fmt.Printf("Active kubectl context: %s\n\n", ui.CyanBold(currentCtx))

			for _, c := range clusters {
				contextName := fmt.Sprintf("kind-%s", c)
				activeTag := ""
				if contextName == currentCtx {
					activeTag = ui.GreenBold(" (ACTIVE)")
				}

				fmt.Printf("• Cluster: %s%s\n", ui.WhiteBold(c), activeTag)
				details, err := cluster.GetClusterDetails(ctx, contextName)
				if err != nil {
					fmt.Printf("   Error getting details: %v\n", err)
					continue
				}

				for _, n := range details.Nodes {
					fmt.Printf("   - Node: %-25s Role: %-15s Status: %s\n",
						n.Name, n.Role, ui.GreenBold(n.Status))
				}
			}
			return nil
		},
	}

	destroyClusterCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Tear down Kind clusters created for CKA",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ui.PrintHeader("TEARING DOWN CKA CLUSTERS")
			_ = cluster.DeleteCluster(ctx, cluster.DefaultCluster1)
			_ = cluster.DeleteCluster(ctx, cluster.DefaultCluster2)
			ui.PrintSuccess("Clusters destroyed.")
			return nil
		},
	}

	clusterCmd.AddCommand(initClusterCmd, statusClusterCmd, destroyClusterCmd)

	// List scenarios
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all available CKA scenarios in tabular format",
		RunE: func(cmd *cobra.Command, args []string) error {
			scenarios, err := scenario.LoadAll(scenariosDir)
			if err != nil {
				return err
			}

			ui.PrintBanner()
			ui.PrintHeader("CKA EXAM CURRICULUM SCENARIOS")

			columns := []ui.TableColumn{
				{Title: "#", Width: 2},
				{Title: "SCENARIO ID", Width: 29},
				{Title: "DOMAIN", Width: 12},
				{Title: "DIFFICULTY", Width: 10},
				{Title: "POINTS", Width: 8},
				{Title: "STATUS", Width: 6},
				{Title: "TITLE", Width: 46},
			}

			var rows [][]string
			totalPoints := 0
			earnedPoints := 0
			completedCount := 0
			count := 0

			for _, s := range scenarios {
				if flagDomain != "" && !strings.Contains(strings.ToLower(s.Domain), strings.ToLower(flagDomain)) {
					continue
				}
				count++
				totalPoints += s.Weight
				isDone := appProgress.IsCompleted(s.ID)
				if isDone {
					completedCount++
					earnedPoints += s.Weight
				}
				rows = append(rows, []string{
					fmt.Sprintf("%2d", count),
					ui.CyanBold(s.ID),
					ui.FormatDomain(s.Domain),
					ui.FormatDifficulty(s.Difficulty),
					fmt.Sprintf("%2d pts", s.Weight),
					ui.FormatStatus(isDone),
					s.Title,
				})
			}

			ui.RenderTable(columns, rows)

			pct := 0.0
			if count > 0 {
				pct = float64(completedCount) / float64(count) * 100
			}

			fmt.Printf("\nTotal: %s scenarios | Solved: %s (%s) | Points: %s / %d pts\n",
				ui.CyanBold(fmt.Sprintf("%d", count)),
				ui.GreenBold(fmt.Sprintf("%d/%d", completedCount, count)),
				ui.YellowBold(fmt.Sprintf("%.1f%%", pct)),
				ui.GreenBold(fmt.Sprintf("%d pts", earnedPoints)),
				totalPoints,
			)
			fmt.Println("\nTo start a scenario: " + ui.CyanBold("cka practice <scenario-id>") + " (e.g. " + ui.YellowBold("cka practice arch-01-rbac") + ")")
			return nil
		},
	}
	listCmd.Flags().StringVarP(&flagDomain, "domain", "d", "", "Filter scenarios by domain")

	// Practice command
	practiceCmd := &cobra.Command{
		Use:   "practice <scenario-id>",
		Short: "Start practicing an individual CKA scenario",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scenarioID := args[0]
			sc, err := scenario.FindByID(scenariosDir, scenarioID)
			if err != nil {
				return err
			}

			ctx := context.Background()
			targetCtx := resolveContext(flagContext, sc.TargetContext)

			// Switch context if available
			if targetCtx != "" {
				_ = cluster.SwitchContext(ctx, targetCtx)
			}

			// Setup cluster state
			if err := scenario.Setup(ctx, sc, targetCtx); err != nil {
				ui.PrintWarning("Setup warning: %v", err)
			}

			// Record practice session & persistent progress
			appProgress.RecordStart(sc.ID)
			_ = appProgress.Save()
			_, _ = exam.StartPracticeSession(sc.ID)

			displayScenario(sc, targetCtx, "")
			return nil
		},
	}

	// Verify command
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Grade the active scenario and display test assertion results",
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := exam.LoadSession()
			if err != nil || session == nil || !session.Active || session.ActiveScenario == "" {
				return fmt.Errorf("no active scenario found. Run 'cka practice <id>' or 'cka exam start' first")
			}

			sc, err := scenario.FindByID(scenariosDir, session.ActiveScenario)
			if err != nil {
				return err
			}

			ctx := context.Background()
			targetCtx := resolveContext(flagContext, sc.TargetContext)

			ui.PrintHeader(fmt.Sprintf("GRADING: %s (%s)", sc.Title, sc.ID))
			ui.PrintInfo("Evaluating verification checks against cluster...")

			result, err := scenario.Verify(ctx, sc, targetCtx)
			if err != nil {
				ui.PrintFailure("Verification execution error: %v", err)
				if result != nil && result.Output != "" {
					fmt.Printf("\nOutput:\n%s\n", result.Output)
				}
				return nil
			}

			// Show individual checks
			if len(result.Checks) > 0 {
				fmt.Println()
				for _, c := range result.Checks {
					if c.Passed {
						ui.PrintSuccess("%s", c.Name)
					} else {
						ui.PrintFailure("%s", c.Name)
					}
				}
				fmt.Println()
			}

			if result.Passed {
				appProgress.RecordCompletion(sc.ID, sc.Weight)
				_ = appProgress.Save()
				ui.PrintSuccess("%s", result.Message)
			} else {
				ui.PrintFailure("%s", result.Message)
				if len(result.Checks) == 0 && result.Output != "" {
					fmt.Printf("\nRaw output:\n%s\n", result.Output)
				}
				fmt.Println("\nTip: Run 'cka hint' for a clue, or 'cka solve' to view the complete solution.")
			}

			// Update exam score if in exam mode
			if session.Mode == exam.ModeExam {
				qs := session.Scores[sc.ID]
				qs.Attempted = true
				qs.Passed = result.Passed
				if result.Passed {
					qs.Earned = sc.Weight
				} else {
					qs.Earned = 0
				}
				qs.LastTested = time.Now()
				session.Scores[sc.ID] = qs
				_ = exam.SaveSession(session)
			}

			return nil
		},
	}

	// Hint command
	hintCmd := &cobra.Command{
		Use:   "hint",
		Short: "Display progressive hints for the active scenario",
		RunE: func(cmd *cobra.Command, args []string) error {
			session, _ := exam.LoadSession()
			if session == nil || session.ActiveScenario == "" {
				return fmt.Errorf("no active scenario. Run 'cka practice <id>' first")
			}

			sc, err := scenario.FindByID(scenariosDir, session.ActiveScenario)
			if err != nil {
				return err
			}

			ui.PrintHeader(fmt.Sprintf("HINTS: %s", sc.Title))
			if len(sc.Hints) == 0 {
				fmt.Println("No specific hints recorded for this scenario. Check official Kubernetes documentation.")
				return nil
			}

			for i, h := range sc.Hints {
				fmt.Printf("\n%s %s\n", ui.YellowBold(fmt.Sprintf("Hint %d:", i+1)), h)
			}
			appProgress.RecordHint(sc.ID, len(sc.Hints)-1)
			_ = appProgress.Save()
			fmt.Println()
			return nil
		},
	}

	// Solve command
	solveCmd := &cobra.Command{
		Use:   "solve",
		Short: "Display step-by-step solution walkthrough and YAML manifests",
		RunE: func(cmd *cobra.Command, args []string) error {
			session, _ := exam.LoadSession()
			if session == nil || session.ActiveScenario == "" {
				return fmt.Errorf("no active scenario. Run 'cka practice <id>' first")
			}

			sc, err := scenario.FindByID(scenariosDir, session.ActiveScenario)
			if err != nil {
				return err
			}

			ui.PrintHeader(fmt.Sprintf("SOLUTION WALKTHROUGH: %s", sc.Title))
			if sc.Solution == "" {
				fmt.Println("No solution file available.")
				return nil
			}

			fmt.Println(sc.Solution)
			return nil
		},
	}

	// Reset command
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset the current scenario to its initial broken state",
		RunE: func(cmd *cobra.Command, args []string) error {
			session, _ := exam.LoadSession()
			if session == nil || session.ActiveScenario == "" {
				return fmt.Errorf("no active scenario to reset")
			}

			sc, err := scenario.FindByID(scenariosDir, session.ActiveScenario)
			if err != nil {
				return err
			}

			ctx := context.Background()
			targetCtx := resolveContext(flagContext, sc.TargetContext)

			ui.PrintInfo("Resetting scenario environment...")
			if err := scenario.Setup(ctx, sc, targetCtx); err != nil {
				return fmt.Errorf("failed to reset scenario: %w", err)
			}

			ui.PrintSuccess("Scenario '%s' has been reset to starting state.", sc.ID)
			return nil
		},
	}

	// Exam commands
	examCmd := &cobra.Command{
		Use:   "exam",
		Short: "Manage full timed CKA mock exam sessions",
	}

	startExamCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a timed 2-hour CKA mock exam",
		RunE: func(cmd *cobra.Command, args []string) error {
			allScenarios, err := scenario.LoadAll(scenariosDir)
			if err != nil {
				return err
			}

			session, err := exam.StartExamSession(allScenarios, flagDuration)
			if err != nil {
				return err
			}

			ui.PrintBanner()
			ui.PrintHeader("CKA MOCK EXAM STARTED")
			fmt.Printf("Total Questions: %d\n", len(session.ExamQuestions))
			fmt.Printf("Duration: %d minutes\n", session.DurationMinutes)
			fmt.Printf("Time Remaining: %s\n\n", ui.YellowBold(session.TimeRemaining()))

			// Load first question
			firstSc, err := scenario.FindByID(scenariosDir, session.ActiveScenario)
			if err != nil {
				return err
			}

			ctx := context.Background()
			_ = cluster.SwitchContext(ctx, firstSc.TargetContext)
			_ = scenario.Setup(ctx, firstSc, firstSc.TargetContext)

			displayScenario(firstSc, firstSc.TargetContext, fmt.Sprintf("QUESTION 1 of %d", len(session.ExamQuestions)))
			return nil
		},
	}
	startExamCmd.Flags().IntVarP(&flagDuration, "duration", "d", 120, "Exam duration in minutes")

	statusExamCmd := &cobra.Command{
		Use:   "status",
		Short: "Show active exam timer, current question, and answered questions",
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := exam.LoadSession()
			if err != nil || session == nil || session.Mode != exam.ModeExam {
				return fmt.Errorf("no active exam found. Run 'cka exam start' to begin")
			}

			ui.PrintHeader("EXAM SESSION STATUS")
			fmt.Printf("Time Remaining: %s\n", ui.YellowBold(session.TimeRemaining()))
			fmt.Printf("Current Question: %d of %d (%s)\n\n",
				session.CurrentIndex+1, len(session.ExamQuestions), session.ActiveScenario)

			examCols := []ui.TableColumn{
				{Title: "", Width: 2},
				{Title: "#", Width: 2},
				{Title: "SCENARIO ID", Width: 29},
				{Title: "POINTS", Width: 8},
				{Title: "STATUS", Width: 12},
				{Title: "LAST TESTED", Width: 12},
			}

			var examRows [][]string
			for i, qID := range session.ExamQuestions {
				qs := session.Scores[qID]
				statusStr := ui.Gray("Unattempted")
				if qs.Attempted {
					if qs.Passed {
						statusStr = ui.GreenBold("PASSED")
					} else {
						statusStr = ui.RedBold("FAILED")
					}
				}
				activeMark := " "
				if i == session.CurrentIndex {
					activeMark = ui.CyanBold("▶")
				}
				testTime := "-"
				if !qs.LastTested.IsZero() {
					testTime = qs.LastTested.Format("15:04:05")
				}
				examRows = append(examRows, []string{
					activeMark,
					fmt.Sprintf("%2d", i+1),
					qID,
					fmt.Sprintf("%2d pts", qs.Weight),
					statusStr,
					testTime,
				})
			}
			ui.RenderTable(examCols, examRows)
			return nil
		},
	}

	nextExamCmd := &cobra.Command{
		Use:   "next",
		Short: "Move to the next exam question",
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := exam.LoadSession()
			if err != nil || session == nil || session.Mode != exam.ModeExam {
				return fmt.Errorf("no active exam")
			}

			if session.CurrentIndex >= len(session.ExamQuestions)-1 {
				ui.PrintWarning("Already at the last question. Run 'cka exam finish' when ready.")
				return nil
			}

			session.CurrentIndex++
			session.ActiveScenario = session.ExamQuestions[session.CurrentIndex]
			_ = exam.SaveSession(session)

			sc, err := scenario.FindByID(scenariosDir, session.ActiveScenario)
			if err != nil {
				return err
			}

			ctx := context.Background()
			_ = cluster.SwitchContext(ctx, sc.TargetContext)
			_ = scenario.Setup(ctx, sc, sc.TargetContext)

			displayScenario(sc, sc.TargetContext, fmt.Sprintf("QUESTION %d of %d [%s]",
				session.CurrentIndex+1, len(session.ExamQuestions), session.TimeRemaining()))
			return nil
		},
	}

	prevExamCmd := &cobra.Command{
		Use:   "prev",
		Short: "Move to the previous exam question",
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := exam.LoadSession()
			if err != nil || session == nil || session.Mode != exam.ModeExam {
				return fmt.Errorf("no active exam")
			}

			if session.CurrentIndex <= 0 {
				ui.PrintWarning("Already at the first question.")
				return nil
			}

			session.CurrentIndex--
			session.ActiveScenario = session.ExamQuestions[session.CurrentIndex]
			_ = exam.SaveSession(session)

			sc, err := scenario.FindByID(scenariosDir, session.ActiveScenario)
			if err != nil {
				return err
			}

			ctx := context.Background()
			_ = cluster.SwitchContext(ctx, sc.TargetContext)
			_ = scenario.Setup(ctx, sc, sc.TargetContext)

			displayScenario(sc, sc.TargetContext, fmt.Sprintf("QUESTION %d of %d [%s]",
				session.CurrentIndex+1, len(session.ExamQuestions), session.TimeRemaining()))
			return nil
		},
	}

	gotoExamCmd := &cobra.Command{
		Use:   "goto <index>",
		Short: "Jump to a specific question number (1-based)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			idx, err := strconv.Atoi(args[0])
			if err != nil || idx < 1 {
				return fmt.Errorf("invalid question number")
			}

			session, err := exam.LoadSession()
			if err != nil || session == nil || session.Mode != exam.ModeExam {
				return fmt.Errorf("no active exam")
			}

			if idx > len(session.ExamQuestions) {
				return fmt.Errorf("question number out of range (1-%d)", len(session.ExamQuestions))
			}

			session.CurrentIndex = idx - 1
			session.ActiveScenario = session.ExamQuestions[session.CurrentIndex]
			_ = exam.SaveSession(session)

			sc, err := scenario.FindByID(scenariosDir, session.ActiveScenario)
			if err != nil {
				return err
			}

			ctx := context.Background()
			_ = cluster.SwitchContext(ctx, sc.TargetContext)
			_ = scenario.Setup(ctx, sc, sc.TargetContext)

			displayScenario(sc, sc.TargetContext, fmt.Sprintf("QUESTION %d of %d [%s]",
				session.CurrentIndex+1, len(session.ExamQuestions), session.TimeRemaining()))
			return nil
		},
	}

	finishExamCmd := &cobra.Command{
		Use:   "finish",
		Short: "Finish the exam and calculate total score",
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := exam.LoadSession()
			if err != nil || session == nil || session.Mode != exam.ModeExam {
				return fmt.Errorf("no active exam to finish")
			}

			// Grade all questions that haven't been tested yet
			ctx := context.Background()
			ui.PrintInfo("Running final automated verification across all questions...")
			for _, qID := range session.ExamQuestions {
				qs := session.Scores[qID]
				if !qs.Passed {
					sc, err := scenario.FindByID(scenariosDir, qID)
					if err == nil {
						res, vErr := scenario.Verify(ctx, sc, sc.TargetContext)
						if vErr == nil && res.Passed {
							qs.Passed = true
							qs.Earned = sc.Weight
							qs.Attempted = true
							session.Scores[qID] = qs
						}
					}
				}
			}

			_, _ = exam.GenerateScoreReport(session)
			_ = exam.ClearSession()
			return nil
		},
	}

	examCmd.AddCommand(startExamCmd, statusExamCmd, nextExamCmd, prevExamCmd, gotoExamCmd, finishExamCmd)

	tuiCmd := &cobra.Command{
		Use:   "tui",
		Short: "Launch the interactive CKA scenario browser and lab dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}

	progressCmd := &cobra.Command{
		Use:   "progress",
		Short: "Display domain breakdown and overall CKA practice progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			scenarios, err := scenario.LoadAll(scenariosDir)
			if err != nil {
				return err
			}

			ui.PrintBanner("CKA-PRACTICE", "Certified Kubernetes Administrator Interactive Exam & Practice Labs")
			ui.PrintHeader("CKA CURRICULUM PROGRESS BREAKDOWN")

			domainStats := make(map[string]struct{ total, completed, totalPts, earnedPts int })
			for _, sc := range scenarios {
				ds := domainStats[sc.Domain]
				ds.total++
				ds.totalPts += sc.Weight
				if appProgress.IsCompleted(sc.ID) {
					ds.completed++
					ds.earnedPts += sc.Weight
				}
				domainStats[sc.Domain] = ds
			}

			columns := []ui.TableColumn{
				{Title: "DOMAIN", Width: 28},
				{Title: "COMPLETED", Width: 12},
				{Title: "POINTS EARNED", Width: 16},
				{Title: "PERCENTAGE", Width: 12},
			}

			var rows [][]string
			for domain, ds := range domainStats {
				pct := 0.0
				if ds.totalPts > 0 {
					pct = float64(ds.earnedPts) / float64(ds.totalPts) * 100
				}
				rows = append(rows, []string{
					domain,
					fmt.Sprintf("%d/%d", ds.completed, ds.total),
					fmt.Sprintf("%d/%d pts", ds.earnedPts, ds.totalPts),
					fmt.Sprintf("%.0f%%", pct),
				})
			}

			ui.RenderTable(columns, rows)
			pct := float64(appProgress.EarnedScore())
			fmt.Printf("\nOverall Score: %s / %s (%s) | Passing Threshold: %s\n\n",
				ui.GreenBold(fmt.Sprintf("%d pts", appProgress.EarnedScore())),
				ui.YellowBold("100 pts"),
				ui.YellowBold(fmt.Sprintf("%.1f%%", pct)),
				ui.GreenBold("66 pts (66%)"),
			)
			return nil
		},
	}

	rootCmd.AddCommand(clusterCmd, listCmd, practiceCmd, verifyCmd, hintCmd, solveCmd, resetCmd, examCmd, progressCmd, tuiCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runTUI() error {
	scenarios, err := scenario.LoadAll(scenariosDir)
	if err != nil {
		return fmt.Errorf("loading scenarios from %s: %w", scenariosDir, err)
	}

	model := tui.NewModel(scenarios, scenariosDir, flagContext, appProgress)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func displayScenario(sc *scenario.Scenario, targetCtx, headerExtra string) {
	header := fmt.Sprintf("TASK: %s (%s)", sc.Title, sc.ID)
	if headerExtra != "" {
		header = fmt.Sprintf("%s - %s", headerExtra, header)
	}
	ui.PrintHeader(header)

	fmt.Printf("• Domain:          %s (%s)\n", sc.Domain, ui.FormatDomain(sc.Domain))
	fmt.Printf("• Weight:          %s\n", ui.YellowBold(fmt.Sprintf("%d points", sc.Weight)))
	fmt.Printf("• Target Context:  %s\n", ui.CyanBold(targetCtx))
	if sc.TargetNamespace != "" {
		fmt.Printf("• Namespace:       %s\n", ui.MagentaBold(sc.TargetNamespace))
	}

	ui.PrintBox("INSTRUCTIONS", sc.Description)

	fmt.Println("\nCommands:")
	fmt.Printf("  • Verify solution:   %s\n", ui.GreenBold("cka verify"))
	fmt.Printf("  • View hints:        %s\n", ui.YellowBold("cka hint"))
	fmt.Printf("  • View solution:     %s\n", ui.CyanBold("cka solve"))
	fmt.Printf("  • Reset problem:     %s\n", ui.RedBold("cka reset"))
}

func resolveContext(explicitContext, scenarioContext string) string {
	if explicitContext != "" {
		return explicitContext
	}
	// Check if scenarioContext exists in kubeconfig
	if scenarioContext != "" {
		cmd := exec.Command("kubectl", "config", "get-contexts", scenarioContext, "--no-headers")
		if err := cmd.Run(); err == nil {
			return scenarioContext
		}
	}
	// Fallback to current-context
	current, err := cluster.GetCurrentContext(context.Background())
	if err == nil && current != "" {
		return current
	}
	return scenarioContext
}

func findScenariosDir() string {
	// 1. Check local scenarios/
	if _, err := os.Stat("scenarios"); err == nil {
		return "scenarios"
	}
	// 2. Check k8s-practice/scenarios/ relative to current working directory
	if _, err := os.Stat("k8s-practice/scenarios"); err == nil {
		return "k8s-practice/scenarios"
	}
	// 3. Check relative to executable directory
	if exe, err := os.Executable(); err == nil {
		// e.g. exe is k8s-practice/cka -> k8s-practice/scenarios
		dir := filepath.Join(filepath.Dir(exe), "scenarios")
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
		// e.g. exe is k8s-practice/bin/cka -> k8s-practice/scenarios
		parentDir := filepath.Join(filepath.Dir(filepath.Dir(exe)), "scenarios")
		if _, err := os.Stat(parentDir); err == nil {
			return parentDir
		}
		// e.g. exe is any-practice/bin/cka -> any-practice/k8s-practice/scenarios
		repoDir := filepath.Join(filepath.Dir(filepath.Dir(exe)), "k8s-practice", "scenarios")
		if _, err := os.Stat(repoDir); err == nil {
			return repoDir
		}
	}
	return "scenarios"
}
