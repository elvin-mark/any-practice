package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/elvin-mark/any-practice/algo-practice/internal/exercise"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 8
		return m, nil

	case FileChangedMsg:
		m.isChecking = true
		if ex := m.CurrentExercise(); ex != nil {
			m.progress.RecordAttempt(ex.ID)
		}
		cmds = append(cmds, m.runCheck(), waitForFileChange(m.watchChan))
		return m, tea.Batch(cmds...)

	case CheckFinishedMsg:
		m.isChecking = false
		m.result = msg.Result
		if msg.Result != nil && msg.Result.Passed {
			if ex := m.CurrentExercise(); ex != nil {
				m.progress.MarkCompleted(ex.ID, msg.Result.NsPerOp, msg.Result.AllocBytes, msg.Result.AllocsPerOp)
				_ = m.progress.Save()
			}
		}
		return m, nil

	case tea.KeyMsg:
		switch m.viewState {
		case ViewSolutionConfirm:
			switch msg.String() {
			case "y", "Y":
				m.loadSolution()
				m.viewState = ViewSolution
				m.viewport.SetContent(m.solutionTxt)
			case "n", "N", "esc", "q":
				m.viewState = ViewNormal
			}
			return m, nil

		case ViewSolution:
			switch msg.String() {
			case "esc", "q", "s":
				m.viewState = ViewNormal
				return m, nil
			default:
				var vpCmd tea.Cmd
				m.viewport, vpCmd = m.viewport.Update(msg)
				return m, vpCmd
			}

		case ViewMenu:
			switch msg.String() {
			case "esc", "q", "m":
				m.viewState = ViewNormal
				return m, nil
			case "up", "k":
				if m.menuCursor > 0 {
					m.menuCursor--
				}
			case "down", "j":
				if m.menuCursor < len(m.repo.Exercises)-1 {
					m.menuCursor++
				}
			case "enter":
				m.activeIdx = m.menuCursor
				m.hintLevel = 0
				m.viewState = ViewNormal
				m.updateWatchDir()
				if ex := m.CurrentExercise(); ex != nil {
					m.progress.SetCurrent(ex.ID)
					_ = m.progress.Save()
				}
				m.isChecking = true
				return m, m.runCheck()
			}
			return m, nil

		case ViewNormal:
			switch msg.String() {
			case "ctrl+c", "q":
				if m.watcher != nil {
					_ = m.watcher.Close()
				}
				return m, tea.Quit

			case "n":
				if m.activeIdx < len(m.repo.Exercises)-1 {
					m.activeIdx++
					m.hintLevel = 0
					m.menuCursor = m.activeIdx
					m.updateWatchDir()
					if ex := m.CurrentExercise(); ex != nil {
						m.progress.SetCurrent(ex.ID)
						_ = m.progress.Save()
					}
					m.isChecking = true
					return m, m.runCheck()
				}

			case "p":
				if m.activeIdx > 0 {
					m.activeIdx--
					m.hintLevel = 0
					m.menuCursor = m.activeIdx
					m.updateWatchDir()
					if ex := m.CurrentExercise(); ex != nil {
						m.progress.SetCurrent(ex.ID)
						_ = m.progress.Save()
					}
					m.isChecking = true
					return m, m.runCheck()
				}

			case "h":
				ex := m.CurrentExercise()
				if ex != nil && len(ex.Hints) > 0 {
					if m.hintLevel < len(ex.Hints) {
						m.hintLevel++
					} else {
						m.hintLevel = 0
					}
				}

			case "s":
				ex := m.CurrentExercise()
				if ex != nil {
					if m.progress.IsCompleted(ex.ID) {
						m.loadSolution()
						m.viewState = ViewSolution
						m.viewport.SetContent(m.solutionTxt)
					} else {
						m.viewState = ViewSolutionConfirm
					}
				}

			case "m":
				m.viewState = ViewMenu
				m.menuCursor = m.activeIdx

			case "r":
				m.isChecking = true
				return m, m.runCheck()
			}
		}
	}

	return m, nil
}

// SwitchToExercise selects a specific exercise by ID.
func (m *Model) SwitchToExercise(id string) tea.Cmd {
	idx := m.repo.IndexOf(id)
	if idx >= 0 {
		m.activeIdx = idx
		m.menuCursor = idx
		m.hintLevel = 0
		m.updateWatchDir()
		m.progress.SetCurrent(id)
		_ = m.progress.Save()
		m.isChecking = true
		return m.runCheck()
	}
	return nil
}

func (m *Model) NextExercise() *exercise.Exercise {
	return m.repo.Next(m.CurrentExercise().ID)
}
