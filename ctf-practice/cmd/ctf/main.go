package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/elvin-mark/any-practice/ctf-practice/internal/challenge"
	"github.com/elvin-mark/any-practice/ctf-practice/internal/progress"
	"github.com/elvin-mark/any-practice/ctf-practice/internal/tui"
	"github.com/elvin-mark/any-practice/ctf-practice/internal/ui"
)

var (
	flagChallengesDir string
	flagProgress      string
	appProgress       *progress.Progress
	appRegistry       *challenge.Registry
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ctf",
		Short: "CTF (Capture The Flag) Practice & Security Training CLI/TUI",
		Long: `Practice cybersecurity and CTF challenges across Cryptography, Web Security,
Reverse Engineering, Binary Exploitation (Pwn), Digital Forensics, and Misc.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Resolve challenges directory
			resolvedDir, err := challenge.ResolveChallengesDir(flagChallengesDir)
			if err != nil {
				return err
			}
			flagChallengesDir = resolvedDir

			// Load registry
			reg, err := challenge.NewRegistry(flagChallengesDir)
			if err != nil {
				return fmt.Errorf("loading challenges from %s: %w", flagChallengesDir, err)
			}
			appRegistry = reg

			// Load progress
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

	rootCmd.PersistentFlags().StringVarP(&flagChallengesDir, "challenges-dir", "c", "", "Path to challenges directory")
	rootCmd.PersistentFlags().StringVarP(&flagProgress, "progress", "p", progress.DefaultPath(), "Path to progress file")

	// Subcommands
	rootCmd.AddCommand(tuiCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(showCmd())
	rootCmd.AddCommand(submitCmd())
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
		Short: "Launch interactive CTF dashboard (default)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}
}

func listCmd() *cobra.Command {
	var filterCategory string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all CTF challenges and completion status",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.PrintBanner("CTF-PRACTICE", "Capture The Flag Training & Security Labs")

			challenges := appRegistry.List()
			if len(challenges) == 0 {
				ui.PrintWarning("No challenges found in %s", flagChallengesDir)
				return nil
			}

			columns := []ui.TableColumn{
				{Title: "#", Width: 3},
				{Title: "CHALLENGE ID", Width: 26},
				{Title: "CATEGORY", Width: 12},
				{Title: "DIFFICULTY", Width: 12},
				{Title: "POINTS", Width: 8},
				{Title: "STATUS", Width: 8},
				{Title: "TITLE", Width: 32},
			}

			var rows [][]string
			solvedCount := 0
			totalPoints := 0
			earnedPoints := 0
			idx := 1

			for _, ch := range challenges {
				if filterCategory != "" && !strings.EqualFold(ch.Category, filterCategory) {
					continue
				}

				totalPoints += ch.Points
				isDone := appProgress.IsCompleted(ch.ID)
				if isDone {
					solvedCount++
					earnedPoints += ch.Points
				}

				rows = append(rows, []string{
					fmt.Sprintf("%2d", idx),
					ui.CyanText.Render(ch.ID),
					ui.CategoryBadgeStyle.Render(ch.Category),
					ui.FormatDifficulty(ch.Difficulty),
					fmt.Sprintf("%d", ch.Points),
					ui.FormatStatus(isDone),
					ch.Title,
				})
				idx++
			}

			ui.RenderTable(columns, rows)

			pct := 0.0
			if len(challenges) > 0 {
				pct = float64(solvedCount) / float64(len(challenges)) * 100
			}
			fmt.Printf("\nProgress: %s challenges solved (%.1f%%) • Total Score: %s / %d pts\n\n",
				ui.GreenText.Render(fmt.Sprintf("%d/%d", solvedCount, len(challenges))),
				pct,
				ui.CyanText.Render(fmt.Sprintf("%d", earnedPoints)),
				totalPoints,
			)

			return nil
		},
	}

	cmd.Flags().StringVar(&filterCategory, "category", "", "Filter challenges by category (Crypto, Web, Reverse, Pwn, Forensics, Misc)")
	return cmd
}

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <challenge-id>",
		Short: "Display details and briefing for a challenge",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			ch, ok := appRegistry.Get(id)
			if !ok {
				ui.PrintFailure("Challenge '%s' not found", id)
				return nil
			}

			appProgress.RecordStart(id)
			_ = appProgress.Save()

			ui.PrintBanner("CTF-PRACTICE", "Challenge Briefing")
			fmt.Printf("\n%s %s\n", ui.BoldText.Render("Title:     "), ch.Title)
			fmt.Printf("%s %s\n", ui.BoldText.Render("ID:        "), ui.CyanText.Render(ch.ID))
			fmt.Printf("%s %s\n", ui.BoldText.Render("Category:  "), ui.CategoryBadgeStyle.Render(ch.Category))
			fmt.Printf("%s %s\n", ui.BoldText.Render("Difficulty:"), ui.FormatDifficulty(ch.Difficulty))
			fmt.Printf("%s %s\n", ui.BoldText.Render("Points:    "), ui.FormatPoints(ch.Points))
			fmt.Printf("%s %s\n", ui.BoldText.Render("Status:    "), ui.FormatStatus(appProgress.IsCompleted(ch.ID)))

			ui.PrintBox("BRIEFING & DESCRIPTION", ch.Description)

			rec := appProgress.GetRecord(id)
			fmt.Printf("Hints Available: %d (viewed: %d) • Use 'ctf hint %s' to reveal\n\n", len(ch.Hints), rec.HintsViewed, id)
			return nil
		},
	}
}

func submitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "submit <challenge-id> <flag>",
		Short: "Submit a flag for validation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			flag := args[1]

			ch, ok := appRegistry.Get(id)
			if !ok {
				ui.PrintFailure("Challenge '%s' not found", id)
				return nil
			}

			valid, err := appRegistry.ValidateFlag(id, flag)
			if err != nil {
				ui.PrintFailure("Validation error: %v", err)
				return nil
			}

			if valid {
				alreadyDone := appProgress.IsCompleted(id)
				appProgress.RecordCompletion(id, ch.Points)
				_ = appProgress.Save()

				if alreadyDone {
					ui.PrintSuccess("Flag verified! (Challenge was already completed) [%d pts]", ch.Points)
				} else {
					ui.PrintSuccess("CORRECT FLAG! Challenge solved! +%d points awarded. Total: %d pts", ch.Points, appProgress.TotalPoints())
				}
			} else {
				appProgress.RecordAttempt(id)
				_ = appProgress.Save()
				ui.PrintFailure("Incorrect flag for '%s'. Try again! (Attempts: %d)", id, appProgress.GetRecord(id).Attempts)
			}

			return nil
		},
	}
}

func hintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hint <challenge-id> [number]",
		Short: "Reveal hints for a challenge",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			ch, ok := appRegistry.Get(id)
			if !ok {
				ui.PrintFailure("Challenge '%s' not found", id)
				return nil
			}

			if len(ch.Hints) == 0 {
				ui.PrintInfo("No hints configured for challenge '%s'", id)
				return nil
			}

			rec := appProgress.GetRecord(id)
			targetIndex := rec.HintsViewed // next unread hint by default
			if len(args) == 2 {
				if parsed, err := strconv.Atoi(args[1]); err == nil && parsed > 0 {
					targetIndex = parsed - 1
				}
			}

			if targetIndex >= len(ch.Hints) {
				ui.PrintWarning("All %d hints have already been revealed for '%s'", len(ch.Hints), id)
				targetIndex = len(ch.Hints) - 1
			}

			appProgress.RecordHint(id, targetIndex)
			_ = appProgress.Save()

			ui.PrintHint("Hint %d/%d for '%s':", targetIndex+1, len(ch.Hints), id)
			fmt.Printf("  %s\n\n", ch.Hints[targetIndex])

			return nil
		},
	}
}

func progressCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "progress",
		Short: "Display detailed CTF learning progress and score breakdown",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.PrintBanner("CTF-PRACTICE", "Scoreboard & Progress Analytics")

			challenges := appRegistry.List()
			if len(challenges) == 0 {
				ui.PrintInfo("No challenges loaded.")
				return nil
			}

			// Group by category
			type catStat struct {
				totalChallenges int
				solvedCount     int
				totalPoints     int
				earnedPoints    int
			}

			catStats := make(map[string]*catStat)
			for _, cat := range appRegistry.Categories() {
				catStats[cat] = &catStat{}
			}

			totalChallenges := len(challenges)
			totalSolved := 0
			totalPoints := 0
			earnedPoints := 0

			for _, ch := range challenges {
				stat := catStats[ch.Category]
				if stat == nil {
					stat = &catStat{}
					catStats[ch.Category] = stat
				}
				stat.totalChallenges++
				stat.totalPoints += ch.Points
				totalPoints += ch.Points

				if appProgress.IsCompleted(ch.ID) {
					stat.solvedCount++
					stat.earnedPoints += ch.Points
					totalSolved++
					earnedPoints += ch.Points
				}
			}

			columns := []ui.TableColumn{
				{Title: "CATEGORY", Width: 16},
				{Title: "SOLVED / TOTAL", Width: 18},
				{Title: "POINTS / MAX", Width: 18},
				{Title: "COMPLETION", Width: 14},
			}

			var rows [][]string
			for _, cat := range appRegistry.Categories() {
				stat := catStats[cat]
				pct := 0.0
				if stat.totalChallenges > 0 {
					pct = float64(stat.solvedCount) / float64(stat.totalChallenges) * 100
				}

				rows = append(rows, []string{
					ui.CategoryBadgeStyle.Render(cat),
					fmt.Sprintf("%d / %d", stat.solvedCount, stat.totalChallenges),
					fmt.Sprintf("%d / %d", stat.earnedPoints, stat.totalPoints),
					fmt.Sprintf("%.1f%%", pct),
				})
			}

			ui.RenderTable(columns, rows)

			overallPct := 0.0
			if totalChallenges > 0 {
				overallPct = float64(totalSolved) / float64(totalChallenges) * 100
			}

			fmt.Println()
			ui.PrintBox("OVERALL PERFORMANCE", fmt.Sprintf(
				"Challenges Solved: %s (%.1f%%)\nTotal Score:       %s / %d pts\nProgress File:     %s",
				ui.GreenText.Render(fmt.Sprintf("%d / %d", totalSolved, totalChallenges)),
				overallPct,
				ui.CyanText.Render(fmt.Sprintf("%d", earnedPoints)),
				totalPoints,
				flagProgress,
			))
			fmt.Println()

			return nil
		},
	}
}
