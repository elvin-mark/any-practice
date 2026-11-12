package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/elvin-mark/any-practice/docker-practice/internal/progress"
	"github.com/elvin-mark/any-practice/docker-practice/internal/scenario"
	"github.com/elvin-mark/any-practice/docker-practice/internal/tui"
	"github.com/elvin-mark/any-practice/docker-practice/internal/ui"
)

var (
	scenariosDir string
	flagProgress string
	appProgress  *progress.Progress
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "dockp [command]",
		Short: "Docker Practice - Master Docker containerization, architecture, and security",
		Long: `dockp is an interactive CLI and TUI training suite for mastering Docker.
Practice building ultra-lean multi-stage images, advanced container lifecycle management,
custom bridge networks, persistent volume data pipelines, Docker Compose stacks, and
production container security hardening.`,
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

	rootCmd.PersistentFlags().StringVarP(&scenariosDir, "scenarios-dir", "s", scenario.FindScenariosDir(), "Path to scenarios directory")
	rootCmd.PersistentFlags().StringVarP(&flagProgress, "progress", "p", progress.DefaultPath(), "Path to progress state file")

	rootCmd.AddCommand(
		newTuiCmd(),
		newListCmd(),
		newStartCmd(),
		newVerifyCmd(),
		newHintCmd(),
		newSolutionCmd(),
		newResetCmd(),
		newProgressCmd(),
		newCleanCmd(),
		newCheckCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runTUI() error {
	scenarios, err := scenario.LoadAll(scenariosDir)
	if err != nil {
		return fmt.Errorf("loading scenarios from %s: %w", scenariosDir, err)
	}

	m := tui.NewModel(scenarios, scenariosDir, appProgress)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func newTuiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch the interactive split-pane TUI browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}
}

func newListCmd() *cobra.Command {
	var filterDomain string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all Docker practice scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			scenarios, err := scenario.LoadAll(scenariosDir)
			if err != nil {
				return err
			}

			ui.PrintBanner("DOCKER-PRACTICE", "Hands-on Docker Containerization Labs")
			ui.PrintHeader("AVAILABLE DOCKER PRACTICE SCENARIOS")

			columns := []ui.TableColumn{
				{Title: "Status", Width: 6},
				{Title: "Scenario ID", Width: 26},
				{Title: "Domain", Width: 12},
				{Title: "Difficulty", Width: 10},
				{Title: "Pts", Width: 4},
				{Title: "Title", Width: 38},
			}

			var rows [][]string
			activeID := appProgress.GetActive()

			for _, sc := range scenarios {
				if filterDomain != "" && !strings.Contains(strings.ToLower(sc.Domain), strings.ToLower(filterDomain)) {
					continue
				}

				status := ui.StatusPending
				if appProgress.IsCompleted(sc.ID) {
					status = ui.StatusCompleted
				} else if sc.ID == activeID {
					status = ui.StatusActive
				}

				rows = append(rows, []string{
					status,
					sc.ID,
					ui.FormatDomain(sc.Domain),
					ui.FormatDifficulty(sc.Difficulty),
					fmt.Sprintf("%2d", sc.Weight),
					sc.Title,
				})
			}

			fmt.Print(ui.RenderTable(columns, rows))
			fmt.Printf("\nLegend: %s Solved  %s Active  %s Unsolved\n", ui.StatusCompleted, ui.StatusActive, ui.StatusPending)
			fmt.Println("Run 'dockp start <scenario-id>' to begin practicing a scenario.")
			return nil
		},
	}
	cmd.Flags().StringVarP(&filterDomain, "domain", "d", "", "Filter scenarios by domain name")
	return cmd
}

func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "start <scenario-id>",
		Aliases: []string{"practice"},
		Short:   "Start a Docker practice lab scenario",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			sc, err := scenario.FindByID(scenariosDir, id)
			if err != nil {
				return err
			}

			ui.PrintBanner()
			ui.PrintHeader(fmt.Sprintf("STARTING SCENARIO: %s", sc.ID))

			ctx := context.Background()
			if err := scenario.Setup(ctx, sc); err != nil {
				return err
			}

			appProgress.RecordStart(sc.ID)
			_ = appProgress.Save()

			ui.PrintSubHeader("Scenario Overview")
			fmt.Printf("Title:      %s\n", ui.BoldText.Render(sc.Title))
			fmt.Printf("Domain:     %s\n", sc.Domain)
			fmt.Printf("Difficulty: %s\n", ui.FormatDifficulty(sc.Difficulty))
			fmt.Printf("Points:     %d pts\n", sc.Weight)
			fmt.Printf("Lab Path:   %s\n", ui.CyanText.Render(sc.Path))

			ui.PrintSubHeader("Task Instructions")
			fmt.Println(sc.Description)

			ui.PrintSubHeader("Next Steps")
			fmt.Println("1. Perform the tasks described in the problem above.")
			fmt.Println("2. Run 'dockp verify' to validate your solution.")
			fmt.Println("3. Run 'dockp hint' if you get stuck.")
			return nil
		},
	}
}

func newVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify [scenario-id]",
		Short: "Verify your solution for the active or specified scenario",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var id string
			if len(args) == 1 {
				id = args[0]
			} else {
				id = appProgress.GetActive()
				if id == "" {
					return fmt.Errorf("no active scenario. Provide an ID or run 'dockp start <id>' first")
				}
			}

			sc, err := scenario.FindByID(scenariosDir, id)
			if err != nil {
				return err
			}

			ui.PrintBanner()
			ui.PrintHeader(fmt.Sprintf("VERIFYING SCENARIO: %s", sc.ID))

			ctx := context.Background()
			res, err := scenario.Verify(ctx, sc)
			if err != nil {
				ui.PrintFailure("Verification runner error: %v", err)
				return nil
			}

			if len(res.Checks) > 0 {
				ui.PrintSubHeader("Verification Checks")
				for _, check := range res.Checks {
					if check.Passed {
						ui.PrintSuccess("%s", check.Name)
					} else {
						ui.PrintFailure("%s", check.Name)
					}
				}
			}

			if res.Passed {
				appProgress.RecordCompletion(sc.ID, sc.Weight)
				_ = appProgress.Save()
				ui.PrintHeader("RESULT: PASSED")
				ui.PrintSuccess("Scenario %s completed successfully! +%d points awarded.", sc.ID, sc.Weight)
			} else {
				ui.PrintHeader("RESULT: FAILED")
				ui.PrintFailure("Scenario %s did not pass all assertions.", sc.ID)
				fmt.Printf("\nDetailed Output:\n%s\n", res.Output)
				fmt.Println("Tip: Run 'dockp hint' to view clues or 'dockp solution " + sc.ID + "' if needed.")
			}

			return nil
		},
	}
}

func newHintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hint [scenario-id]",
		Short: "Reveal hints progressively for the active or specified scenario",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var id string
			if len(args) == 1 {
				id = args[0]
			} else {
				id = appProgress.GetActive()
				if id == "" {
					return fmt.Errorf("no active scenario. Specify scenario ID or run 'dockp start <id>' first")
				}
			}

			sc, err := scenario.FindByID(scenariosDir, id)
			if err != nil {
				return err
			}

			if len(sc.Hints) == 0 {
				ui.PrintInfo("No hints available for scenario %s.", sc.ID)
				return nil
			}

			hintsViewed := appProgress.RecordHint(sc.ID)
			_ = appProgress.Save()

			ui.PrintBanner()
			ui.PrintHeader(fmt.Sprintf("HINTS FOR %s", sc.ID))

			maxShow := hintsViewed
			if maxShow > len(sc.Hints) {
				maxShow = len(sc.Hints)
			}

			for i := 0; i < maxShow; i++ {
				ui.PrintHint("Hint %d/%d: %s", i+1, len(sc.Hints), sc.Hints[i])
			}

			if maxShow < len(sc.Hints) {
				fmt.Printf("\n%s additional hint(s) remaining. Run 'dockp hint' again to reveal the next hint.\n",
					strconv.Itoa(len(sc.Hints)-maxShow))
			} else {
				fmt.Println("\nAll hints revealed.")
			}

			return nil
		},
	}
}

func newSolutionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "solution [scenario-id]",
		Short: "Display the official reference solution for a scenario",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var id string
			if len(args) == 1 {
				id = args[0]
			} else {
				id = appProgress.GetActive()
				if id == "" {
					return fmt.Errorf("no active scenario. Specify scenario ID or run 'dockp start <id>' first")
				}
			}

			sc, err := scenario.FindByID(scenariosDir, id)
			if err != nil {
				return err
			}

			ui.PrintBanner()
			ui.PrintHeader(fmt.Sprintf("OFFICIAL SOLUTION GUIDE: %s", sc.ID))

			if strings.TrimSpace(sc.Solution) == "" {
				ui.PrintInfo("No solution guide found for scenario %s.", sc.ID)
				return nil
			}

			fmt.Println(sc.Solution)
			return nil
		},
	}
}

func newResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset [scenario-id]",
		Short: "Reset a scenario back to its initial lab state",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var id string
			if len(args) == 1 {
				id = args[0]
			} else {
				id = appProgress.GetActive()
				if id == "" {
					return fmt.Errorf("no active scenario. Specify scenario ID or run 'dockp start <id>' first")
				}
			}

			sc, err := scenario.FindByID(scenariosDir, id)
			if err != nil {
				return err
			}

			ui.PrintBanner()
			ui.PrintHeader(fmt.Sprintf("RESETTING SCENARIO: %s", sc.ID))

			ctx := context.Background()
			if err := scenario.Reset(ctx, sc); err != nil {
				return fmt.Errorf("reset failed: %w", err)
			}

			ui.PrintSuccess("Scenario %s lab has been reset to starting state.", sc.ID)
			return nil
		},
	}
}

func newProgressCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "progress",
		Short: "Display overall Docker learning progress and domain mastery",
		RunE: func(cmd *cobra.Command, args []string) error {
			scenarios, err := scenario.LoadAll(scenariosDir)
			if err != nil {
				return err
			}

			ui.PrintBanner()
			ui.PrintHeader("DOCKER PRACTICE MASTERY PROGRESS")

			domainTotal := make(map[string]int)
			domainSolved := make(map[string]int)
			domainPoints := make(map[string]int)
			domainEarned := make(map[string]int)

			totalScore := 0
			earnedScore := 0
			solvedCount := 0

			for _, sc := range scenarios {
				domainTotal[sc.Domain]++
				domainPoints[sc.Domain] += sc.Weight
				totalScore += sc.Weight

				if appProgress.IsCompleted(sc.ID) {
					solvedCount++
					domainSolved[sc.Domain]++
					domainEarned[sc.Domain] += sc.Weight
					earnedScore += sc.Weight
				}
			}

			overallPct := 0.0
			if len(scenarios) > 0 {
				overallPct = float64(solvedCount) / float64(len(scenarios)) * 100
			}

			fmt.Printf("Overall Completion: %s (%d / %d scenarios solved • %.1f%%)\n",
				ui.ProgressBadgeStyle.Render(fmt.Sprintf("%.0f%%", overallPct)),
				solvedCount, len(scenarios), overallPct)
			fmt.Printf("Total Score:        %s (%d / %d points)\n\n",
				ui.YellowText.Render(fmt.Sprintf("%d pts", earnedScore)),
				earnedScore, totalScore)

			columns := []ui.TableColumn{
				{Title: "Domain", Width: 38},
				{Title: "Solved", Width: 10},
				{Title: "Score", Width: 12},
				{Title: "Mastery", Width: 10},
			}

			var rows [][]string
			// Deterministic domain ordering
			orderedDomains := []string{
				scenario.DomainDockerfile,
				scenario.DomainLifecycle,
				scenario.DomainStorage,
				scenario.DomainNetworking,
				scenario.DomainCompose,
				scenario.DomainSecurity,
				scenario.DomainTroubleshooting,
			}

			for _, dom := range orderedDomains {
				tot := domainTotal[dom]
				if tot == 0 {
					continue
				}
				sol := domainSolved[dom]
				pts := domainPoints[dom]
				erd := domainEarned[dom]
				pct := float64(sol) / float64(tot) * 100

				rows = append(rows, []string{
					dom,
					fmt.Sprintf("%d / %d", sol, tot),
					fmt.Sprintf("%d / %d pts", erd, pts),
					fmt.Sprintf("%5.1f%%", pct),
				})
			}

			fmt.Print(ui.RenderTable(columns, rows))
			return nil
		},
	}
}

func newCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean up all containers, networks, and volumes created during practice",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.PrintBanner()
			ui.PrintHeader("CLEANING PRACTICE LAB RESOURCES")

			ctx := context.Background()

			// Remove practice labeled containers
			rmContainersCmd := exec.CommandContext(ctx, "bash", "-c", "docker rm -f $(docker ps -aq --filter label=practice=docker) 2>/dev/null || true")
			_ = rmContainersCmd.Run()

			// Prune practice networks
			rmNetsCmd := exec.CommandContext(ctx, "bash", "-c", "docker network prune -f --filter label=practice=docker 2>/dev/null || true")
			_ = rmNetsCmd.Run()

			// Prune practice volumes
			rmVolsCmd := exec.CommandContext(ctx, "bash", "-c", "docker volume prune -f --filter label=practice=docker 2>/dev/null || true")
			_ = rmVolsCmd.Run()

			ui.PrintSuccess("Cleaned up Docker practice containers, networks, and test volumes.")
			return nil
		},
	}
}

func newCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check local Docker engine & CLI prerequisites",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.PrintBanner()
			ui.PrintHeader("CHECKING DOCKER PREREQUISITES")

			// 1. Docker CLI
			dockerVersionCmd := exec.Command("docker", "version", "--format", "{{.Client.Version}}")
			out, err := dockerVersionCmd.Output()
			if err != nil {
				ui.PrintFailure("Docker CLI is not installed or not in PATH.")
				return nil
			}
			ui.PrintSuccess("Docker Client Version: %s", strings.TrimSpace(string(out)))

			// 2. Docker Daemon Connectivity
			daemonPingCmd := exec.Command("docker", "info", "--format", "{{.ServerVersion}}")
			out, err = daemonPingCmd.Output()
			if err != nil {
				ui.PrintFailure("Cannot connect to Docker daemon. Is docker service running?")
				return nil
			}
			ui.PrintSuccess("Docker Engine Server Version: %s", strings.TrimSpace(string(out)))

			// 3. Docker Compose
			composeVersionCmd := exec.Command("docker", "compose", "version")
			out, err = composeVersionCmd.Output()
			if err != nil {
				ui.PrintWarning("Docker Compose v2 plugin not detected (some multi-container labs require it).")
			} else {
				ui.PrintSuccess("Docker Compose: %s", strings.TrimSpace(string(out)))
			}

			ui.PrintSuccess("All Docker prerequisites met! Ready for practice.")
			return nil
		},
	}
}
