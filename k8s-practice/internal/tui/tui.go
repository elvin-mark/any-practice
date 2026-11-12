package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/elvin-mark/any-practice/k8s-practice/internal/cluster"
	"github.com/elvin-mark/any-practice/k8s-practice/internal/exam"
	"github.com/elvin-mark/any-practice/k8s-practice/internal/progress"
	"github.com/elvin-mark/any-practice/k8s-practice/internal/scenario"
	"github.com/elvin-mark/any-practice/k8s-practice/internal/ui"
)

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Enter   key.Binding
	Verify  key.Binding
	Hint    key.Binding
	Reset   key.Binding
	Cluster key.Binding
	Quit    key.Binding
}

var defaultKeyMap = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", "p"),
		key.WithHelp("enter/p", "practice"),
	),
	Verify: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "verify"),
	),
	Hint: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "hint"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
	),
	Cluster: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "cluster"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Model represents the interactive CKA TUI state
type Model struct {
	scenarios        []*scenario.Scenario
	scenariosDir     string
	targetContext    string
	prog             *progress.Progress
	cursor           int
	statusMsg        string
	activeScenarioID string
	hintsRevealed    map[string]int
	keys             keyMap
}

// NewModel initializes the CKA interactive TUI
func NewModel(scenarios []*scenario.Scenario, scenariosDir string, targetContext string, prog *progress.Progress) Model {
	activeID := ""
	if prog != nil {
		activeID = prog.GetActive()
	}

	return Model{
		scenarios:        scenarios,
		scenariosDir:     scenariosDir,
		targetContext:    targetContext,
		prog:             prog,
		cursor:           0,
		statusMsg:        "Ready. Select a scenario and press 'Enter' to start practicing.",
		activeScenarioID: activeID,
		hintsRevealed:    make(map[string]int),
		keys:             defaultKeyMap,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type verifyResultMsg struct {
	scenarioID string
	passed     bool
	message    string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.scenarios) - 1
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.scenarios)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}

		case key.Matches(msg, m.keys.Enter):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			ctx := context.Background()

			tCtx := sc.TargetContext
			if m.targetContext != "" {
				tCtx = m.targetContext
			}
			if tCtx != "" {
				_ = cluster.SwitchContext(ctx, tCtx)
			}

			if err := scenario.Setup(ctx, sc, tCtx); err != nil {
				m.statusMsg = fmt.Sprintf("Setup warning: %v", err)
			} else {
				m.statusMsg = fmt.Sprintf("Scenario %s environment active! Solve challenge and press 'v' to verify.", sc.ID)
			}
			m.activeScenarioID = sc.ID
			if m.prog != nil {
				m.prog.RecordStart(sc.ID)
				_ = m.prog.Save()
			}
			_, _ = exam.StartPracticeSession(sc.ID)

		case key.Matches(msg, m.keys.Verify):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			m.statusMsg = fmt.Sprintf("⏳ Grading scenario %s against cluster...", sc.ID)

			tCtx := sc.TargetContext
			if m.targetContext != "" {
				tCtx = m.targetContext
			}

			return m, func() tea.Msg {
				ctx := context.Background()
				res, err := scenario.Verify(ctx, sc, tCtx)
				if err != nil {
					return verifyResultMsg{
						scenarioID: sc.ID,
						passed:     false,
						message:    fmt.Sprintf("Verification error: %v", err),
					}
				}
				return verifyResultMsg{
					scenarioID: sc.ID,
					passed:     res.Passed,
					message:    res.Message,
				}
			}

		case key.Matches(msg, m.keys.Hint):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			if len(sc.Hints) == 0 {
				m.statusMsg = "No hints available for this scenario."
				return m, nil
			}
			curHint := m.hintsRevealed[sc.ID]
			if curHint < len(sc.Hints) {
				m.statusMsg = fmt.Sprintf("💡 Hint %d/%d: %s", curHint+1, len(sc.Hints), sc.Hints[curHint])
				m.hintsRevealed[sc.ID] = curHint + 1
			} else {
				m.statusMsg = fmt.Sprintf("💡 All hints revealed. Latest: %s", sc.Hints[len(sc.Hints)-1])
			}

		case key.Matches(msg, m.keys.Reset):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			ctx := context.Background()
			tCtx := sc.TargetContext
			if m.targetContext != "" {
				tCtx = m.targetContext
			}
			m.statusMsg = fmt.Sprintf("Resetting scenario %s...", sc.ID)
			if err := scenario.Setup(ctx, sc, tCtx); err != nil {
				m.statusMsg = fmt.Sprintf("Reset warning: %v", err)
			} else {
				m.statusMsg = fmt.Sprintf("Scenario %s has been reset to starting state.", sc.ID)
			}

		case key.Matches(msg, m.keys.Cluster):
			ctx := context.Background()
			clusters, err := cluster.ListClusters(ctx)
			if err != nil {
				m.statusMsg = fmt.Sprintf("Cluster error: %v", err)
			} else if len(clusters) == 0 {
				m.statusMsg = "No Kind clusters found. Run 'cka cluster init' to create them."
			} else {
				curCtx, _ := cluster.GetCurrentContext(ctx)
				m.statusMsg = fmt.Sprintf("Clusters: %s | Active Context: %s", strings.Join(clusters, ", "), curCtx)
			}
		}

	case verifyResultMsg:
		if msg.passed {
			if m.prog != nil {
				weight := 5
				for _, sc := range m.scenarios {
					if sc.ID == msg.scenarioID {
						weight = sc.Weight
						break
					}
				}
				m.prog.RecordCompletion(msg.scenarioID, weight)
				_ = m.prog.Save()
			}
			m.statusMsg = fmt.Sprintf("🎉 SUCCESS! %s: %s", msg.scenarioID, msg.message)
		} else {
			m.statusMsg = fmt.Sprintf("❌ FAILED: %s: %s", msg.scenarioID, msg.message)
		}
	}

	return m, nil
}

func (m Model) View() string {
	if len(m.scenarios) == 0 {
		return "No scenarios found. Check scenarios directory."
	}

	var b strings.Builder

	// Calculate statistics
	completedCount := 0
	totalScore := 0
	earnedScore := 0
	for _, sc := range m.scenarios {
		totalScore += sc.Weight
		if m.prog != nil && m.prog.IsCompleted(sc.ID) {
			completedCount++
			earnedScore += sc.Weight
		}
	}
	progressPct := 0.0
	if len(m.scenarios) > 0 {
		progressPct = float64(completedCount) / float64(len(m.scenarios)) * 100
	}

	// 1. Header Bar
	headerBar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		ui.TitleBadgeStyle.Render("☸ CKA-PRACTICE"),
		lipgloss.NewStyle().MarginLeft(1).Render(
			ui.ProgressBadgeStyle.Render(fmt.Sprintf("[%d/%d Solved • %d/%d pts • %.0f%%]",
				completedCount, len(m.scenarios), earnedScore, totalScore, progressPct)),
		),
	)
	b.WriteString(headerBar + "\n\n")

	// Split view: Left list, Right details
	var leftCol strings.Builder
	currentDomain := ""

	for i, sc := range m.scenarios {
		if sc.Domain != currentDomain {
			currentDomain = sc.Domain
			domainHeader := fmt.Sprintf("\n  %s", ui.FormatDomain(currentDomain))
			leftCol.WriteString(domainHeader + "\n")
		}

		statusMark := ui.StatusPending
		if m.prog != nil && m.prog.IsCompleted(sc.ID) {
			statusMark = ui.StatusCompleted
		}

		activeIndicator := " "
		if sc.ID == m.activeScenarioID {
			activeIndicator = ui.CyanText.Render("▶")
		}

		itemText := fmt.Sprintf("%s%s %-24s (%2d pts)", activeIndicator, statusMark, sc.ID, sc.Weight)
		if i == m.cursor {
			leftCol.WriteString(lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(ui.ColorActiveBg).
				Padding(0, 1).
				Render("> "+itemText) + "\n")
		} else {
			leftCol.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#C0C0C0")).
				Padding(0, 1).
				Render("  "+itemText) + "\n")
		}
	}

	// Details Pane
	curScenario := m.scenarios[m.cursor]
	var rightCol strings.Builder

	domainBadge := ui.FormatDomain(curScenario.Domain)
	diffBadge := ui.FormatDifficulty(curScenario.Difficulty)
	weightBadge := ui.YellowText.Render(fmt.Sprintf("[%d Points]", curScenario.Weight))

	rightCol.WriteString(fmt.Sprintf("%s  %s  %s\n\n", domainBadge, diffBadge, weightBadge))
	rightCol.WriteString(ui.BoldText.Render(curScenario.Title) + "\n")
	rightCol.WriteString(ui.SubtleText.Render(fmt.Sprintf("ID: %s | Target Context: %s | Namespace: %s",
		curScenario.ID, curScenario.TargetContext, curScenario.TargetNamespace)) + "\n\n")

	rightCol.WriteString(ui.BoldText.Render("Scenario Instructions:") + "\n")
	rightCol.WriteString(curScenario.Description + "\n\n")

	// Hints status
	hintsRevealed := m.hintsRevealed[curScenario.ID]
	rightCol.WriteString(fmt.Sprintf("%s %d (Revealed: %d)\n",
		ui.SubtleText.Render("Hints Available:"), len(curScenario.Hints), hintsRevealed))

	// Layout side by side
	listPanel := lipgloss.NewStyle().Width(44).Render(leftCol.String())
	detailsPanel := ui.CardStyle.Width(54).Padding(1, 2).Render(rightCol.String())

	content := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, detailsPanel)
	b.WriteString(content + "\n\n")

	// Status Message
	if m.statusMsg != "" {
		b.WriteString(ui.CyanText.Render("Status: "+m.statusMsg) + "\n\n")
	}

	// Keymap Help
	helpLine := ui.RenderKeyHelp([][2]string{
		{"↑/↓", "navigate"},
		{"enter/p", "practice"},
		{"v", "verify"},
		{"h", "hint"},
		{"c", "cluster"},
		{"r", "reset"},
		{"q", "quit"},
	})
	b.WriteString(helpLine + "\n")

	return b.String()
}
