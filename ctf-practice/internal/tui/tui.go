package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/elvin-mark/any-practice/ctf-practice/internal/challenge"
	"github.com/elvin-mark/any-practice/ctf-practice/internal/progress"
	"github.com/elvin-mark/any-practice/ctf-practice/internal/ui"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Submit key.Binding
	Hint   key.Binding
	Cancel key.Binding
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
	Submit: key.NewBinding(
		key.WithKeys("s", "enter"),
		key.WithHelp("s/enter", "submit flag"),
	),
	Hint: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "hint"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Model represents the interactive CTF TUI state.
type Model struct {
	registry      *challenge.Registry
	challenges    []challenge.Challenge
	prog          *progress.Progress
	cursor        int
	statusMsg     string
	hintsRevealed map[string]int
	keys          keyMap
	textInput     textinput.Model
	isSubmitting  bool
	width         int
	height        int
}

// NewModel initializes the CTF interactive TUI.
func NewModel(reg *challenge.Registry, prog *progress.Progress) Model {
	challenges := reg.List()

	ti := textinput.New()
	ti.Placeholder = "FLAG{...}"
	ti.CharLimit = 120
	ti.Width = 40

	m := Model{
		registry:      reg,
		challenges:    challenges,
		prog:          prog,
		cursor:        0,
		hintsRevealed: make(map[string]int),
		keys:          defaultKeyMap,
		textInput:     ti,
		isSubmitting:  false,
		width:         80,
		height:        24,
	}

	if prog != nil && prog.GetActive() != "" {
		for i, ch := range challenges {
			if ch.ID == prog.GetActive() {
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

	case tea.KeyMsg:
		if m.isSubmitting {
			switch msg.Type {
			case tea.KeyEsc:
				m.isSubmitting = false
				m.textInput.Blur()
				m.textInput.Reset()
				return m, nil

			case tea.KeyEnter:
				flagValue := strings.TrimSpace(m.textInput.Value())
				if flagValue == "" {
					m.isSubmitting = false
					m.textInput.Blur()
					return m, nil
				}

				if m.cursor >= 0 && m.cursor < len(m.challenges) {
					ch := m.challenges[m.cursor]
					valid, err := m.registry.ValidateFlag(ch.ID, flagValue)
					if err != nil {
						m.statusMsg = ui.RedText.Render(fmt.Sprintf("✖ Error: %v", err))
					} else if valid {
						if m.prog != nil {
							m.prog.RecordCompletion(ch.ID, ch.Points)
							_ = m.prog.Save()
						}
						m.statusMsg = ui.GreenText.Render(fmt.Sprintf("✔ PASS: Flag correct! Awarded %d points.", ch.Points))
					} else {
						if m.prog != nil {
							m.prog.RecordAttempt(ch.ID)
							_ = m.prog.Save()
						}
						m.statusMsg = ui.RedText.Render("✖ FAIL: Incorrect flag. Try again!")
					}
				}

				m.isSubmitting = false
				m.textInput.Blur()
				m.textInput.Reset()
				return m, nil

			case tea.KeyCtrlC:
				return m, tea.Quit

			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		}

		// Normal navigation mode
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.challenges) - 1
			}
			m.updateActiveID()

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.challenges)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
			m.updateActiveID()

		case key.Matches(msg, m.keys.Submit):
			if len(m.challenges) == 0 {
				return m, nil
			}
			m.isSubmitting = true
			m.textInput.Reset()
			m.textInput.Focus()
			return m, textinput.Blink

		case key.Matches(msg, m.keys.Hint):
			if len(m.challenges) == 0 {
				return m, nil
			}
			ch := m.challenges[m.cursor]
			if len(ch.Hints) == 0 {
				m.statusMsg = ui.SubtleText.Render("No hints available for this challenge.")
				return m, nil
			}
			curHint := m.hintsRevealed[ch.ID]
			if curHint < len(ch.Hints) {
				m.statusMsg = fmt.Sprintf("%s %s", ui.YellowText.Render(fmt.Sprintf("💡 Hint %d/%d:", curHint+1, len(ch.Hints))), ch.Hints[curHint])
				m.hintsRevealed[ch.ID] = curHint + 1
				if m.prog != nil {
					m.prog.RecordHint(ch.ID, curHint)
					_ = m.prog.Save()
				}
			} else {
				m.statusMsg = fmt.Sprintf("%s %s", ui.YellowText.Render("💡 All hints revealed:"), ch.Hints[len(ch.Hints)-1])
			}
		}
	}

	return m, nil
}

func (m *Model) updateActiveID() {
	if m.prog != nil && m.cursor >= 0 && m.cursor < len(m.challenges) {
		ch := m.challenges[m.cursor]
		m.prog.SetActive(ch.ID)
		m.prog.RecordStart(ch.ID)
		_ = m.prog.Save()
	}
}

// View renders the terminal user interface.
func (m Model) View() string {
	if len(m.challenges) == 0 {
		return ui.CardStyle.Render("No challenges found in challenges directory.")
	}

	// 1. Header Banner
	bannerBadge := ui.TitleBadgeStyle.Render(" CTF-PRACTICE ")
	bannerSubtitle := ui.SubtleText.Render("Capture The Flag Training & Security Labs")

	completedCount := 0
	totalScore := 0
	maxScore := 0
	for _, ch := range m.challenges {
		maxScore += ch.Points
		if m.prog != nil && m.prog.IsCompleted(ch.ID) {
			completedCount++
			totalScore += ch.Points
		}
	}

	progressStr := fmt.Sprintf("Solved: %s/%d  Score: %s/%d pts",
		ui.GreenText.Render(fmt.Sprintf("%d", completedCount)),
		len(m.challenges),
		ui.CyanText.Render(fmt.Sprintf("%d", totalScore)),
		maxScore,
	)

	header := fmt.Sprintf(" %s  %s   %s", bannerBadge, bannerSubtitle, progressStr)

	// 2. Left Pane: Challenge List
	var leftLines []string
	leftLines = append(leftLines, ui.BoldText.Render("AVAILABLE CHALLENGES"))
	leftLines = append(leftLines, ui.SubtleText.Render("──────────────────────────────────────"))

	for i, ch := range m.challenges {
		isDone := m.prog != nil && m.prog.IsCompleted(ch.ID)
		statusBadge := ui.FormatStatus(isDone)

		catBadge := lipgloss.NewStyle().
			Foreground(ui.ColorSecondary).
			Bold(true).
			Render(fmt.Sprintf("[%s]", ch.Category))

		ptsBadge := ui.FormatPoints(ch.Points)

		line := fmt.Sprintf("%s %-8s %s %s", statusBadge, catBadge, ptsBadge, ch.ID)

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
		Width(46).
		Height(15).
		Render(strings.Join(leftLines, "\n"))

	// 3. Right Pane: Selected Challenge Detail
	cur := m.challenges[m.cursor]
	isDone := m.prog != nil && m.prog.IsCompleted(cur.ID)

	var rightLines []string
	rightLines = append(rightLines, fmt.Sprintf("%s %s", ui.FormatStatus(isDone), ui.BoldText.Render(cur.Title)))
	rightLines = append(rightLines, fmt.Sprintf("%s  %s  %s  %s",
		ui.CategoryBadgeStyle.Render(cur.Category),
		ui.FormatDifficulty(cur.Difficulty),
		ui.FormatPoints(cur.Points),
		ui.SubtleText.Render("ID: "+cur.ID),
	))
	rightLines = append(rightLines, ui.SubtleText.Render("─────────────────────────────────────────────"))

	descLines := strings.Split(strings.TrimSpace(cur.Description), "\n")
	for _, dl := range descLines {
		rightLines = append(rightLines, ui.BoldText.Render(dl))
	}

	rightLines = append(rightLines, "")
	hintsCount := len(cur.Hints)
	revealed := m.hintsRevealed[cur.ID]
	rightLines = append(rightLines, ui.SubtleText.Render(fmt.Sprintf("💡 Hints: %d/%d revealed (press 'h' to reveal)", revealed, hintsCount)))

	if m.isSubmitting {
		rightLines = append(rightLines, "")
		rightLines = append(rightLines, ui.CyanText.Render("Enter Flag: ")+m.textInput.View())
		rightLines = append(rightLines, ui.SubtleText.Render("(Enter: verify flag, Esc: cancel)"))
	}

	rightPane := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorCardBorder).
		Padding(0, 1).
		Width(54).
		Height(15).
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
		{"s/enter", "submit flag"},
		{"h", "hint"},
		{"q", "quit"},
	}
	help := ui.RenderKeyHelp(keys)

	return fmt.Sprintf("\n%s\n\n%s\n\n%s%s\n", header, body, status, help)
}
