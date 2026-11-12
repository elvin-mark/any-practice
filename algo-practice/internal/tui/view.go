package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	switch m.viewState {
	case ViewSolutionConfirm:
		return m.viewSolutionConfirm()
	case ViewSolution:
		return m.viewSolution()
	case ViewMenu:
		return m.viewMenu()
	default:
		return m.viewNormal()
	}
}

func (m *Model) viewNormal() string {
	var b strings.Builder
	ex := m.CurrentExercise()
	if ex == nil {
		return "No exercises found. Run 'algo init' to unpack exercises."
	}

	// 1. Header Bar
	completedCount := m.progress.Stats()
	totalCount := len(m.repo.Exercises)
	progressText := fmt.Sprintf("[%d/%d Solved]", completedCount, totalCount)

	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		titleStyle.Render("⚡ ALGO-PRACTICE"),
		trackBadgeStyle.Render(fmt.Sprintf("%s", ex.Track)),
		progressBadgeStyle.Render(progressText),
	)
	b.WriteString(header + "\n\n")

	// 2. Exercise Title & Difficulty
	diffBadge := diffEasyStyle.Render("Easy")
	if strings.EqualFold(ex.Difficulty, "medium") {
		diffBadge = diffMedStyle.Render("Medium")
	} else if strings.EqualFold(ex.Difficulty, "hard") {
		diffBadge = diffHardStyle.Render("Hard")
	}

	statusIcon := "○"
	if m.progress.IsCompleted(ex.ID) {
		statusIcon = lipgloss.NewStyle().Foreground(colorSuccess).Render("● (Solved)")
	}

	exHeader := fmt.Sprintf("%s  %s  %s",
		boldText.Render(fmt.Sprintf("#%d %s", ex.Order, ex.Name)),
		diffBadge,
		statusIcon,
	)
	b.WriteString(exHeader + "\n")

	// Target Complexity & Editable file path
	complexityText := subtleText.Render(fmt.Sprintf("Complexity: Time %s | Space %s", ex.Complexity.Time, ex.Complexity.Space))
	fileText := mutedText.Render(fmt.Sprintf("File: %s", ex.EntryFilePath()))
	b.WriteString(complexityText + "\n" + fileText + "\n\n")

	// 3. Problem Description Card
	descContent := fmt.Sprintf("%s\n\n%s", boldText.Render("Description:"), ex.Description)
	b.WriteString(cardStyle.Width(m.width-4).Render(descContent) + "\n\n")

	// 4. Status / Result Box
	if m.isChecking {
		checkingStyle := lipgloss.NewStyle().Foreground(colorWarning).Bold(true)
		b.WriteString(cardStyle.Width(m.width-4).Render(checkingStyle.Render("⏳ Running AST validation, unit tests, and complexity benchmarks...")) + "\n\n")
	} else if m.result != nil {
		if !m.result.ASTValid {
			// AST Failed
			var violations strings.Builder
			violations.WriteString(lipgloss.NewStyle().Foreground(colorError).Bold(true).Render("❌ AST Validation Failed: Forbidden Standard Library Usage\n\n"))
			if m.result.ParseErr != nil {
				violations.WriteString(fmt.Sprintf("Syntax error: %v\n", m.result.ParseErr))
			}
			for _, v := range m.result.ASTViolations {
				violations.WriteString(fmt.Sprintf("  • %s\n", v))
			}
			violations.WriteString("\n" + subtleText.Render("Please implement the algorithm from scratch without standard library cheat functions."))
			b.WriteString(errorCardStyle.Width(m.width-4).Render(violations.String()) + "\n\n")
		} else if !m.result.TestsPassed {
			// Unit Tests Failed
			var testFail strings.Builder
			testFail.WriteString(lipgloss.NewStyle().Foreground(colorError).Bold(true).Render("❌ Unit Tests Failed\n\n"))
			output := truncateLines(m.result.TestOutput, 12)
			testFail.WriteString(output)
			b.WriteString(errorCardStyle.Width(m.width-4).Render(testFail.String()) + "\n\n")
		} else if m.result.Passed {
			// Passed!
			var passBox strings.Builder
			passBox.WriteString(lipgloss.NewStyle().Foreground(colorSuccess).Bold(true).Render("✓ ALL CHECKS PASSED!\n\n"))
			passBox.WriteString(fmt.Sprintf("  • AST Verification:  Clean (no forbidden packages)\n"))
			passBox.WriteString(fmt.Sprintf("  • Unit Tests:        Passed\n"))
			if m.result.NsPerOp > 0 {
				passBox.WriteString(fmt.Sprintf("  • Benchmark Runtime: %.2f ns/op\n", m.result.NsPerOp))
				passBox.WriteString(fmt.Sprintf("  • Allocations:       %d B/op  (%d allocs/op)\n", m.result.AllocBytes, m.result.AllocsPerOp))
			}
			passBox.WriteString("\n" + lipgloss.NewStyle().Foreground(colorHighlight).Bold(true).Render("Press [n] to advance to the next exercise!"))
			b.WriteString(successCardStyle.Width(m.width-4).Render(passBox.String()) + "\n\n")
		}
	}

	// 5. Hint Box (if active)
	if m.hintLevel > 0 && len(ex.Hints) >= m.hintLevel {
		hintTitle := fmt.Sprintf("💡 Hint (%d/%d):", m.hintLevel, len(ex.Hints))
		hintContent := fmt.Sprintf("%s\n\n%s",
			lipgloss.NewStyle().Foreground(colorWarning).Bold(true).Render(hintTitle),
			ex.Hints[m.hintLevel-1],
		)
		b.WriteString(cardStyle.BorderForeground(colorWarning).Width(m.width-4).Render(hintContent) + "\n\n")
	}

	// 6. Footer Hotkeys
	hintsLabel := "[h] Hint"
	if len(ex.Hints) > 0 {
		hintsLabel = fmt.Sprintf("[h] Hint (%d/%d)", m.hintLevel, len(ex.Hints))
	}

	keys := []string{
		keyBadgeStyle.Render(hintsLabel),
		keyBadgeStyle.Render("[s] Solution"),
		keyBadgeStyle.Render("[n] Next"),
		keyBadgeStyle.Render("[p] Prev"),
		keyBadgeStyle.Render("[m] Menu"),
		keyBadgeStyle.Render("[r] Re-run"),
		keyBadgeStyle.Render("[q] Quit"),
	}
	footer := footerStyle.Width(m.width - 4).Render(strings.Join(keys, "  "))
	b.WriteString(footer)

	return b.String()
}

func (m *Model) viewSolutionConfirm() string {
	boxContent := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWarning).
		Render("⚠️  Reference Solution Confirmation\n\n") +
		"You have not completed this exercise yet.\n" +
		"Are you sure you want to view the reference solution?\n\n" +
		lipgloss.NewStyle().Bold(true).Foreground(colorTextLight).Render("[y] Yes, reveal solution    [n/esc] No, keep trying")

	return "\n\n" + cardStyle.
		BorderForeground(colorWarning).
		Width(min(65, m.width-4)).
		Padding(1, 2).
		Render(boxContent)
}

func (m *Model) viewSolution() string {
	ex := m.CurrentExercise()
	title := "Reference Solution"
	if ex != nil {
		title = fmt.Sprintf("Reference Solution: %s", ex.Name)
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		titleStyle.Render(title),
		subtleText.Render(" (Press [esc] or [q] to close)"),
	)

	return fmt.Sprintf("%s\n\n%s\n\n%s",
		header,
		cardStyle.Width(m.width-4).Render(m.viewport.View()),
		subtleText.Render("[↑/↓] Scroll  [esc/q] Return to Exercise"),
	)
}

func (m *Model) viewMenu() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("📋 EXERCISE CURRICULUM BROWSER") + "\n\n")
	b.WriteString(subtleText.Render("Use [↑/↓] or [j/k] to navigate, [Enter] to select, [esc] to close.") + "\n\n")

	currentTrack := ""
	for i, ex := range m.repo.Exercises {
		if ex.Track != currentTrack {
			currentTrack = ex.Track
			b.WriteString("\n" + boldText.Foreground(colorHighlight).Render(fmt.Sprintf("▸ Track: %s", currentTrack)) + "\n")
		}

		cursor := "  "
		if i == m.menuCursor {
			cursor = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render("❯ ")
		}

		check := "[ ]"
		if m.progress.IsCompleted(ex.ID) {
			check = lipgloss.NewStyle().Foreground(colorSuccess).Render("[✓]")
		}

		line := fmt.Sprintf("%s%s %d. %-32s (%s)", cursor, check, ex.Order, ex.Name, ex.Difficulty)
		if i == m.menuCursor {
			line = lipgloss.NewStyle().Bold(true).Foreground(colorTextLight).Render(line)
		} else {
			line = subtleText.Render(line)
		}
		b.WriteString(line + "\n")
	}

	return b.String()
}

func truncateLines(s string, maxLines int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s
	}
	return strings.Join(lines[:maxLines], "\n") + fmt.Sprintf("\n... (%d more lines truncated)", len(lines)-maxLines)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
