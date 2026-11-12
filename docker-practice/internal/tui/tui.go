package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/elvin-mark/any-practice/docker-practice/internal/progress"
	"github.com/elvin-mark/any-practice/docker-practice/internal/scenario"
	"github.com/elvin-mark/any-practice/docker-practice/internal/ui"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Verify   key.Binding
	Hint     key.Binding
	Solution key.Binding
	Reset    key.Binding
	Quit     key.Binding
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
	Solution: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "solution"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Model represents the interactive Docker Practice TUI state
type Model struct {
	scenarios        []*scenario.Scenario
	scenariosDir     string
	prog             *progress.Progress
	cursor           int
	statusMsg        string
	activeScenarioID string
	hintsRevealed    map[string]int
	showSolution     bool
	keys             keyMap
}

// NewModel initializes the Docker Practice interactive TUI
func NewModel(scenarios []*scenario.Scenario, scenariosDir string, prog *progress.Progress) Model {
	activeID := ""
	if prog != nil {
		activeID = prog.GetActive()
	}

	return Model{
		scenarios:        scenarios,
		scenariosDir:     scenariosDir,
		prog:             prog,
		cursor:           0,
		statusMsg:        "Ready. Select a scenario and press 'Enter' or 'p' to start practicing.",
		activeScenarioID: activeID,
		hintsRevealed:    make(map[string]int),
		showSolution:     false,
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
				m.showSolution = false
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.scenarios)-1 {
				m.cursor++
				m.showSolution = false
			}

		case key.Matches(msg, m.keys.Enter):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			ctx := context.Background()
			m.statusMsg = fmt.Sprintf("🚀 Setting up lab environment for %s...", sc.ID)
			if err := scenario.Setup(ctx, sc); err != nil {
				m.statusMsg = fmt.Sprintf("Setup error: %v", err)
			} else {
				m.activeScenarioID = sc.ID
				m.statusMsg = fmt.Sprintf("Lab active: %s. Follow instructions in problem description.", sc.ID)
			}
			if m.prog != nil {
				m.prog.RecordStart(sc.ID)
				_ = m.prog.Save()
			}

		case key.Matches(msg, m.keys.Verify):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			m.statusMsg = fmt.Sprintf("⏳ Grading scenario %s against Docker engine...", sc.ID)

			return m, func() tea.Msg {
				ctx := context.Background()
				res, err := scenario.Verify(ctx, sc)
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
			if m.prog != nil {
				m.prog.RecordHint(sc.ID)
				_ = m.prog.Save()
			}

		case key.Matches(msg, m.keys.Solution):
			m.showSolution = !m.showSolution
			if m.showSolution {
				m.statusMsg = "Displaying official solution guide. Press 's' to toggle back."
			} else {
				m.statusMsg = "Displaying problem instructions."
			}

		case key.Matches(msg, m.keys.Reset):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			ctx := context.Background()
			m.statusMsg = fmt.Sprintf("Resetting scenario %s...", sc.ID)
			if err := scenario.Setup(ctx, sc); err != nil {
				m.statusMsg = fmt.Sprintf("Reset error: %v", err)
			} else {
				m.statusMsg = fmt.Sprintf("Scenario %s lab has been reset to starting state.", sc.ID)
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
		ui.TitleBadgeStyle.Render("🐳 DOCKER-PRACTICE"),
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
		} else if sc.ID == m.activeScenarioID {
			statusMark = ui.StatusActive
		}

		line := fmt.Sprintf("  %s %-24s %s", statusMark, sc.ID, ui.FormatDifficulty(sc.Difficulty))
		if i == m.cursor {
			line = lipgloss.NewStyle().
				Background(ui.ColorActiveBg).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true).
				Render(line)
		}
		leftCol.WriteString(line + "\n")
	}

	// Right Details Panel
	sc := m.scenarios[m.cursor]
	var rightCol strings.Builder

	statusBadge := ui.StatusPending + " Unsolved"
	if m.prog != nil && m.prog.IsCompleted(sc.ID) {
		statusBadge = ui.StatusCompleted + " Solved"
	} else if sc.ID == m.activeScenarioID {
		statusBadge = ui.StatusActive + " Active"
	}

	detailsHeader := fmt.Sprintf("%s\n%s • %s • %d pts • %s",
		ui.CyanText.Render(sc.Title),
		ui.FormatDomain(sc.Domain),
		ui.FormatDifficulty(sc.Difficulty),
		sc.Weight,
		statusBadge,
	)
	rightCol.WriteString(detailsHeader + "\n")
	rightCol.WriteString(ui.SubtleText.Render(strings.Repeat("─", 54)) + "\n\n")

	contentToRender := sc.Description
	if m.showSolution {
		contentToRender = sc.Solution
		if strings.TrimSpace(contentToRender) == "" {
			contentToRender = "No solution guide provided for this scenario."
		}
		rightCol.WriteString(ui.GreenText.Render("📖 OFFICIAL SOLUTION:\n\n"))
	}

	// Truncate/wrap description
	lines := strings.Split(contentToRender, "\n")
	if len(lines) > 18 {
		lines = lines[:18]
		lines = append(lines, ui.SubtleText.Render("... (run CLI 'dockp solution "+sc.ID+"' to read full document)"))
	}
	rightCol.WriteString(strings.Join(lines, "\n") + "\n")

	// Render Split Layout
	leftBox := lipgloss.NewStyle().
		Width(48).
		Render(leftCol.String())

	rightBox := ui.CardStyle.
		Width(58).
		Height(22).
		Render(rightCol.String())

	splitLayout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
	b.WriteString(splitLayout + "\n\n")

	// Status Message
	statusBox := lipgloss.NewStyle().
		Foreground(ui.ColorText).
		Background(lipgloss.Color("#252233")).
		Padding(0, 1).
		Render("ℹ " + m.statusMsg)
	b.WriteString(statusBox + "\n\n")

	// Footer Keybindings
	keysHelp := [][2]string{
		{"↑/↓", "navigate"},
		{"enter/p", "practice"},
		{"v", "verify"},
		{"h", "hint"},
		{"s", "solution"},
		{"r", "reset"},
		{"q", "quit"},
	}
	b.WriteString(ui.RenderKeyHelp(keysHelp) + "\n")

	return b.String()
}
