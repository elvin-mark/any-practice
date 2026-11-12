package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/elvin-mark/any-practice/net-practice/internal/progress"
	"github.com/elvin-mark/any-practice/net-practice/internal/scenario"
	"github.com/elvin-mark/any-practice/net-practice/internal/ui"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Start  key.Binding
	Verify key.Binding
	Reset  key.Binding
	Hint   key.Binding
	Quit   key.Binding
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
	Start: key.NewBinding(
		key.WithKeys("enter", "s"),
		key.WithHelp("enter/s", "start lab"),
	),
	Verify: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "verify"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
	),
	Hint: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "hint"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Model represents the interactive Net Practice TUI state.
type Model struct {
	registry      *scenario.Registry
	scenarios     []*scenario.Scenario
	prog          *progress.Progress
	cursor        int
	statusMsg     string
	hintsRevealed map[string]int
	keys          keyMap
	width         int
	height        int
}

type verifyResultMsg struct {
	scenarioID string
	passed     bool
	message    string
}

// NewModel initializes the Net Practice TUI.
func NewModel(reg *scenario.Registry, prog *progress.Progress) Model {
	scenarios := reg.List()

	m := Model{
		registry:      reg,
		scenarios:     scenarios,
		prog:          prog,
		cursor:        0,
		hintsRevealed: make(map[string]int),
		keys:          defaultKeyMap,
		width:         80,
		height:        24,
	}

	if prog != nil && prog.GetActive() != "" {
		for i, sc := range scenarios {
			if sc.ID == prog.GetActive() {
				m.cursor = i
				break
			}
		}
	}

	return m
}

// Init starts the Bubbletea program.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles UI events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case verifyResultMsg:
		if msg.passed {
			if m.prog != nil {
				m.prog.RecordCompletion(msg.scenarioID)
				_ = m.prog.Save()
			}
			m.statusMsg = ui.GreenText.Render("✔ PASS: Network assertions passed! Well done.")
		} else {
			if m.prog != nil {
				m.prog.RecordAttempt(msg.scenarioID)
				_ = m.prog.Save()
			}
			m.statusMsg = ui.RedText.Render(fmt.Sprintf("✖ FAIL: Verification failed:\n%s", msg.message))
		}
		return m, nil

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
			m.updateActiveID()

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.scenarios)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
			m.updateActiveID()

		case key.Matches(msg, m.keys.Start):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			labDir, err := scenario.SetupLab(sc)
			if err != nil {
				m.statusMsg = ui.RedText.Render(fmt.Sprintf("✖ Setup Error: %v", err))
				return m, nil
			}
			if m.prog != nil {
				m.prog.RecordStart(sc.ID)
				m.prog.SetActive(sc.ID)
				_ = m.prog.Save()
			}
			m.statusMsg = fmt.Sprintf("%s\n%s %s",
				ui.GreenText.Render("✔ Network lab initialized!"),
				ui.YellowText.Render("Run:"),
				ui.BoldText.Render("cd "+labDir),
			)

		case key.Matches(msg, m.keys.Verify):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			m.statusMsg = ui.CyanText.Render(fmt.Sprintf("⏳ Verifying scenario %s...", sc.ID))

			return m, func() tea.Msg {
				passed, outMsg, err := scenario.VerifyLab(sc)
				if err != nil {
					return verifyResultMsg{
						scenarioID: sc.ID,
						passed:     false,
						message:    fmt.Sprintf("Internal verify error: %v", err),
					}
				}
				return verifyResultMsg{
					scenarioID: sc.ID,
					passed:     passed,
					message:    outMsg,
				}
			}

		case key.Matches(msg, m.keys.Reset):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			labDir, err := scenario.SetupLab(sc)
			if err != nil {
				m.statusMsg = ui.RedText.Render(fmt.Sprintf("✖ Reset Error: %v", err))
				return m, nil
			}
			m.statusMsg = fmt.Sprintf("%s\n%s %s",
				ui.YellowText.Render("⚠ Lab workspace reset to starting state!"),
				ui.YellowText.Render("Run:"),
				ui.BoldText.Render("cd "+labDir),
			)

		case key.Matches(msg, m.keys.Hint):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			if len(sc.Hints) == 0 {
				m.statusMsg = ui.SubtleText.Render("No hints available for this scenario.")
				return m, nil
			}
			curHint := m.hintsRevealed[sc.ID]
			if curHint < len(sc.Hints) {
				m.statusMsg = fmt.Sprintf("%s %s", ui.YellowText.Render(fmt.Sprintf("💡 Hint %d/%d:", curHint+1, len(sc.Hints))), sc.Hints[curHint])
				m.hintsRevealed[sc.ID] = curHint + 1
				if m.prog != nil {
					m.prog.RecordHint(sc.ID, curHint)
					_ = m.prog.Save()
				}
			} else {
				m.statusMsg = fmt.Sprintf("%s %s", ui.YellowText.Render("💡 All hints revealed:"), sc.Hints[len(sc.Hints)-1])
			}
		}
	}

	return m, nil
}

func (m *Model) updateActiveID() {
	if m.prog != nil && m.cursor >= 0 && m.cursor < len(m.scenarios) {
		sc := m.scenarios[m.cursor]
		m.prog.SetActive(sc.ID)
		_ = m.prog.Save()
	}
}

// View renders the terminal user interface.
func (m Model) View() string {
	if len(m.scenarios) == 0 {
		return ui.CardStyle.Render("No scenarios found in scenarios directory.")
	}

	// 1. Header Banner
	bannerBadge := ui.TitleBadgeStyle.Render(" NET-PRACTICE ")
	bannerSubtitle := ui.SubtleText.Render("Computer Networking, Protocols & Packet Labs")

	completedCount := 0
	for _, sc := range m.scenarios {
		if m.prog != nil && m.prog.IsCompleted(sc.ID) {
			completedCount++
		}
	}

	progressStr := fmt.Sprintf("Progress: %s/%d completed (%.1f%%)",
		ui.GreenText.Render(fmt.Sprintf("%d", completedCount)),
		len(m.scenarios),
		float64(completedCount)/float64(len(m.scenarios))*100,
	)

	header := fmt.Sprintf(" %s  %s   %s", bannerBadge, bannerSubtitle, progressStr)

	// 2. Left Pane: Scenarios List
	var leftLines []string
	leftLines = append(leftLines, ui.BoldText.Render("AVAILABLE NETWORK LABS"))
	leftLines = append(leftLines, ui.SubtleText.Render("──────────────────────────────────────"))

	for i, sc := range m.scenarios {
		isDone := m.prog != nil && m.prog.IsCompleted(sc.ID)
		statusBadge := ui.FormatStatus(isDone)

		catBadge := lipgloss.NewStyle().
			Foreground(ui.ColorSecondary).
			Bold(true).
			Render(fmt.Sprintf("[%s]", sc.Category))

		diffBadge := ui.FormatDifficulty(sc.Difficulty)

		line := fmt.Sprintf("%s %-12s %s %s", statusBadge, catBadge, diffBadge, sc.ID)

		if i == m.cursor {
			line = lipgloss.NewStyle().
				Background(ui.ColorActiveBg).
				Foreground(ui.ColorText).
				Bold(true).
				Render("▶ " + line)
		} else {
			line = "  " + line
		}
		leftLines = append(leftLines, line)
	}

	leftPane := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorCardBorder).
		Padding(0, 1).
		Width(48).
		Height(16).
		Render(strings.Join(leftLines, "\n"))

	// 3. Right Pane: Selected Scenario Detail
	cur := m.scenarios[m.cursor]
	isDone := m.prog != nil && m.prog.IsCompleted(cur.ID)
	labPath := scenario.LabDir(cur.ID)

	var rightLines []string
	rightLines = append(rightLines, fmt.Sprintf("%s %s", ui.FormatStatus(isDone), ui.BoldText.Render(cur.Title)))
	rightLines = append(rightLines, fmt.Sprintf("%s  %s  %s",
		ui.CategoryBadgeStyle.Render(cur.Category),
		ui.FormatDifficulty(cur.Difficulty),
		ui.SubtleText.Render("ID: "+cur.ID),
	))
	rightLines = append(rightLines, ui.SubtleText.Render("─────────────────────────────────────────────"))

	descLines := strings.Split(strings.TrimSpace(cur.Description), "\n")
	for _, dl := range descLines {
		rightLines = append(rightLines, dl)
	}

	rightLines = append(rightLines, "")
	rightLines = append(rightLines, ui.CyanText.Render("🎯 GOAL:"))
	goalLines := strings.Split(strings.TrimSpace(cur.Goal), "\n")
	for _, gl := range goalLines {
		rightLines = append(rightLines, "  "+gl)
	}

	rightLines = append(rightLines, "")
	rightLines = append(rightLines, fmt.Sprintf("%s %s",
		ui.YellowText.Render("📁 Workspace:"),
		ui.BoldText.Render(labPath),
	))

	hintsCount := len(cur.Hints)
	revealed := m.hintsRevealed[cur.ID]
	rightLines = append(rightLines, ui.SubtleText.Render(fmt.Sprintf("💡 Hints: %d/%d revealed (press 'h' to reveal)", revealed, hintsCount)))

	rightPane := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorCardBorder).
		Padding(0, 1).
		Width(56).
		Height(16).
		Render(strings.Join(rightLines, "\n"))

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, " ", rightPane)

	// 4. Status Bar / Notifications
	status := ""
	if m.statusMsg != "" {
		status = m.statusMsg + "\n"
	}

	// 5. Help Footer
	keys := [][2]string{
		{"↑/↓", "navigate"},
		{"enter/s", "start lab"},
		{"v", "verify"},
		{"r", "reset"},
		{"h", "hint"},
		{"q", "quit"},
	}
	help := ui.RenderKeyHelp(keys)

	return fmt.Sprintf("\n%s\n\n%s\n\n%s%s\n", header, body, status, help)
}
