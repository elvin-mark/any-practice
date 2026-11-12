package tui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elvin-mark/any-practice/linux-practice/embedded"
	"github.com/elvin-mark/any-practice/linux-practice/internal/config"
	"github.com/elvin-mark/any-practice/linux-practice/internal/container"
	"github.com/elvin-mark/any-practice/linux-practice/internal/progress"
	"github.com/elvin-mark/any-practice/linux-practice/internal/scenario"
	"github.com/elvin-mark/any-practice/linux-practice/internal/ui"
)

// Styling definitions aligned with the Any-Practice design system
var (
	titleStyle = ui.TitleBadgeStyle

	headerStyle = lipgloss.NewStyle().
			Foreground(ui.ColorPrimary).
			Bold(true)

	activeItemStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ui.ColorActiveBg).
			Padding(0, 1)

	inactiveItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#C0C0C0")).
				Padding(0, 1)

	completedStyle = lipgloss.NewStyle().
			Foreground(ui.ColorSuccess).
			Bold(true)

	pendingStyle = lipgloss.NewStyle().
			Foreground(ui.ColorSubtle)

	panelStyle = ui.CardStyle.Padding(1, 2)

	helpStyle = ui.SubtleText

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")).
			Italic(true)
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Verify key.Binding
	Attach key.Binding
	Hint   key.Binding
	Reset  key.Binding
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
	Enter: key.NewBinding(
		key.WithKeys("enter", "s"),
		key.WithHelp("enter/s", "start"),
	),
	Verify: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "verify"),
	),
	Attach: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "attach"),
	),
	Hint: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "hint"),
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

// Model represents the Bubble Tea TUI state.
type Model struct {
	registry     *scenario.Registry
	prog         *progress.Progress
	manager      *container.Manager
	scenarios    []*scenario.Scenario
	cursor       int
	statusMsg    string
	hintRevealed map[string]int
	width        int
	height       int
}

// NewModel creates an initialized TUI model.
func NewModel(reg *scenario.Registry, prog *progress.Progress, mgr *container.Manager) Model {
	scenarios := reg.All()
	return Model{
		registry:     reg,
		prog:         prog,
		manager:      mgr,
		scenarios:    scenarios,
		cursor:       0,
		statusMsg:    "Welcome to Linux Practice! Select a scenario to start practicing.",
		hintRevealed: make(map[string]int),
		width:        80,
		height:       24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultKeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, defaultKeyMap.Up):
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.scenarios) - 1
			}

		case key.Matches(msg, defaultKeyMap.Down):
			if m.cursor < len(m.scenarios)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}

		case key.Matches(msg, defaultKeyMap.Enter):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			ctx := context.Background()

			m.statusMsg = fmt.Sprintf("Preparing and starting scenario %s...", sc.ID)
			if err := m.manager.EnsureBaseImage(ctx, embedded.Dockerfile, config.DefaultImageName); err != nil {
				m.statusMsg = fmt.Sprintf("❌ Error building base image: %v", err)
				return m, nil
			}

			if err := m.manager.StartScenario(ctx, sc); err != nil {
				m.statusMsg = fmt.Sprintf("❌ Error starting scenario: %v", err)
				return m, nil
			}

			m.prog.RecordStart(sc.ID)
			_ = m.prog.Save()

			// Launch interactive shell via tea.ExecProcess
			cName := container.ContainerName(sc.ID)
			cmd := tea.ExecProcess(
				exec.Command("docker", "exec", "-it", "-u", config.DefaultUser, "-w", "/home/"+config.DefaultUser, cName, "/bin/bash"),
				func(err error) tea.Msg {
					return shellExitMsg{err: err}
				},
			)
			return m, cmd

		case key.Matches(msg, defaultKeyMap.Attach):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			ctx := context.Background()
			if !m.manager.IsRunning(ctx, sc.ID) {
				m.statusMsg = fmt.Sprintf("⚠️ Container for %s is not running. Press 'enter' to start it.", sc.ID)
				return m, nil
			}

			cName := container.ContainerName(sc.ID)
			cmd := tea.ExecProcess(
				exec.Command("docker", "exec", "-it", "-u", config.DefaultUser, "-w", "/home/"+config.DefaultUser, cName, "/bin/bash"),
				func(err error) tea.Msg {
					return shellExitMsg{err: err}
				},
			)
			return m, cmd

		case key.Matches(msg, defaultKeyMap.Verify):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			ctx := context.Background()
			passed, out, err := m.manager.Verify(ctx, sc)
			if err != nil {
				m.statusMsg = fmt.Sprintf("❌ Verification error: %v", err)
				return m, nil
			}

			if passed {
				m.prog.RecordCompletion(sc.ID)
				_ = m.prog.Save()
				m.statusMsg = fmt.Sprintf("🎉 SUCCESS! %s passed:\n%s", sc.ID, out)
			} else {
				m.statusMsg = fmt.Sprintf("❌ Failed: %s", out)
			}

		case key.Matches(msg, defaultKeyMap.Hint):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			if len(sc.Hints) == 0 {
				m.statusMsg = "No hints available for this scenario."
				return m, nil
			}
			curHint := m.hintRevealed[sc.ID]
			if curHint < len(sc.Hints) {
				m.statusMsg = fmt.Sprintf("💡 Hint %d/%d: %s", curHint+1, len(sc.Hints), sc.Hints[curHint])
				m.hintRevealed[sc.ID] = curHint + 1
				m.prog.RecordHint(sc.ID, curHint)
				_ = m.prog.Save()
			} else {
				m.statusMsg = fmt.Sprintf("💡 All hints revealed. Latest: %s", sc.Hints[len(sc.Hints)-1])
			}

		case key.Matches(msg, defaultKeyMap.Reset):
			if len(m.scenarios) == 0 {
				return m, nil
			}
			sc := m.scenarios[m.cursor]
			ctx := context.Background()
			m.statusMsg = fmt.Sprintf("Resetting scenario %s...", sc.ID)
			if err := m.manager.StartScenario(ctx, sc); err != nil {
				m.statusMsg = fmt.Sprintf("❌ Reset error: %v", err)
			} else {
				m.statusMsg = fmt.Sprintf("🔄 Scenario %s has been reset to starting state.", sc.ID)
			}
		}

	case shellExitMsg:
		m.statusMsg = "Returned from terminal session. Press 'v' to verify your solution!"
		return m, nil
	}

	return m, nil
}

type shellExitMsg struct {
	err error
}

func (m Model) View() string {
	if len(m.scenarios) == 0 {
		return "No scenarios found. Check your configuration or embedded scenarios."
	}

	var b strings.Builder

	// Calculate statistics
	completedCount := m.prog.CompletedCount()
	progressPct := float64(completedCount) / float64(len(m.scenarios)) * 100

	// 1. Unified Header Bar
	headerBar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		ui.TitleBadgeStyle.Render("🐧 LINUX-PRACTICE"),
		lipgloss.NewStyle().MarginLeft(1).Render(
			ui.ProgressBadgeStyle.Render(fmt.Sprintf("[%d/%d Solved • %.0f%%]", completedCount, len(m.scenarios), progressPct)),
		),
	)
	b.WriteString(headerBar + "\n\n")

	// Split view: Left list, Right details
	var leftCol strings.Builder
	currentTrack := ""

	for i, sc := range m.scenarios {
		if sc.Track != currentTrack {
			currentTrack = sc.Track
			trackHeader := fmt.Sprintf("\n  [%s]", strings.ToUpper(currentTrack))
			leftCol.WriteString(headerStyle.Render(trackHeader) + "\n")
		}

		statusMark := ui.StatusPending
		if m.prog.IsCompleted(sc.ID) {
			statusMark = ui.StatusCompleted
		}

		itemText := fmt.Sprintf("%s %s", statusMark, sc.Title)
		if i == m.cursor {
			leftCol.WriteString(activeItemStyle.Render("> "+itemText) + "\n")
		} else {
			leftCol.WriteString(inactiveItemStyle.Render("  "+itemText) + "\n")
		}
	}

	// Details Pane
	curScenario := m.scenarios[m.cursor]
	var rightCol strings.Builder

	diffBadge := ui.FormatDifficulty(curScenario.Difficulty)
	trackBadge := ui.CategoryBadgeStyle.Render(strings.ToUpper(curScenario.Track))

	rightCol.WriteString(fmt.Sprintf("%s  %s\n\n", trackBadge, diffBadge))
	rightCol.WriteString(ui.BoldText.Render(curScenario.Title) + "\n")
	rightCol.WriteString(ui.SubtleText.Render("ID: "+curScenario.ID) + "\n\n")

	rightCol.WriteString(ui.BoldText.Render("Challenge:") + "\n")
	rightCol.WriteString(curScenario.Description + "\n\n")

	if len(curScenario.Tags) > 0 {
		rightCol.WriteString(fmt.Sprintf("%s %s\n", ui.SubtleText.Render("Tags:"), strings.Join(curScenario.Tags, ", ")))
	}
	rightCol.WriteString(fmt.Sprintf("%s %d mins\n", ui.SubtleText.Render("Est. Time:"), curScenario.EstimatedMinutes))

	// Hints counter
	rec := m.prog.GetRecord(curScenario.ID)
	attempts := rec.Attempts
	hintsUsed := rec.HintsViewed
	rightCol.WriteString(fmt.Sprintf("%s %d | %s %d (%s %d)\n",
		ui.SubtleText.Render("Attempts:"), attempts,
		ui.SubtleText.Render("Hints:"), len(curScenario.Hints),
		ui.SubtleText.Render("Used:"), hintsUsed))

	// Layout side by side
	listPanel := lipgloss.NewStyle().Width(42).Render(leftCol.String())
	detailsPanel := panelStyle.Width(54).Render(rightCol.String())

	content := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, detailsPanel)
	b.WriteString(content + "\n\n")

	// Status Message
	if m.statusMsg != "" {
		b.WriteString(ui.CyanText.Render("Status: "+m.statusMsg) + "\n\n")
	}

	// Keymap Help
	helpLine := ui.RenderKeyHelp([][2]string{
		{"↑/↓", "navigate"},
		{"enter/s", "start/shell"},
		{"a", "attach"},
		{"v", "verify"},
		{"h", "hint"},
		{"r", "reset"},
		{"q", "quit"},
	})
	b.WriteString(helpLine + "\n")

	return b.String()
}
