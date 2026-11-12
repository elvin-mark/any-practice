package tui

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"

	"github.com/elvin-mark/algo-practice/internal/exercise"
	"github.com/elvin-mark/algo-practice/internal/progress"
	"github.com/elvin-mark/algo-practice/internal/runner"
)

type ViewState int

const (
	ViewNormal ViewState = iota
	ViewMenu
	ViewSolutionConfirm
	ViewSolution
)

// Messages
type CheckStartedMsg struct{}
type CheckFinishedMsg struct {
	Result *runner.Result
}
type FileChangedMsg struct{}

type Model struct {
	repo        *exercise.Repository
	progress    *progress.Progress
	activeIdx   int
	result      *runner.Result
	isChecking  bool
	hintLevel   int
	viewState   ViewState
	menuCursor  int
	width       int
	height      int
	viewport    viewport.Model
	watcher     *fsnotify.Watcher
	watchChan   chan struct{}
	solutionTxt string
	statusMsg   string
}

func NewModel(repo *exercise.Repository, prog *progress.Progress) *Model {
	activeIdx := 0
	currID := prog.GetCurrent()
	if currID != "" {
		if idx := repo.IndexOf(currID); idx >= 0 {
			activeIdx = idx
		}
	} else if len(repo.Exercises) > 0 {
		prog.SetCurrent(repo.Exercises[0].ID)
		_ = prog.Save()
	}

	vp := viewport.New(80, 20)

	m := &Model{
		repo:       repo,
		progress:   prog,
		activeIdx:  activeIdx,
		hintLevel:  0,
		viewState:  ViewNormal,
		menuCursor: activeIdx,
		viewport:   vp,
		watchChan:  make(chan struct{}, 1),
	}

	m.initWatcher()
	return m
}

func (m *Model) CurrentExercise() *exercise.Exercise {
	if len(m.repo.Exercises) == 0 {
		return nil
	}
	return m.repo.Exercises[m.activeIdx]
}

func (m *Model) initWatcher() {
	ex := m.CurrentExercise()
	if ex == nil {
		return
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	m.watcher = w
	_ = w.Add(ex.Dir)

	go func() {
		var lastTrigger time.Time
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					if filepath.Ext(event.Name) == ".go" {
						now := time.Now()
						// Debounce rapid events within 150ms
						if now.Sub(lastTrigger) > 150*time.Millisecond {
							lastTrigger = now
							select {
							case m.watchChan <- struct{}{}:
							default:
							}
						}
					}
				}
			case _, ok := <-w.Errors:
				if !ok {
					return
				}
			}
		}
	}()
}

func (m *Model) updateWatchDir() {
	if m.watcher == nil {
		return
	}
	ex := m.CurrentExercise()
	if ex == nil {
		return
	}
	for _, dir := range m.watcher.WatchList() {
		_ = m.watcher.Remove(dir)
	}
	_ = m.watcher.Add(ex.Dir)
}

func (m *Model) loadSolution() {
	ex := m.CurrentExercise()
	if ex == nil {
		m.solutionTxt = "No exercise selected."
		return
	}
	solPath := ex.SolutionFilePath()
	data, err := os.ReadFile(solPath)
	if err != nil {
		m.solutionTxt = "// Solution file not found for this exercise."
		return
	}
	m.solutionTxt = string(data)
}

func (m *Model) runCheck() tea.Cmd {
	ex := m.CurrentExercise()
	if ex == nil {
		return nil
	}

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		res := runner.Run(ctx, ex)
		return CheckFinishedMsg{Result: res}
	}
}

func waitForFileChange(ch <-chan struct{}) tea.Cmd {
	return func() tea.Msg {
		<-ch
		return FileChangedMsg{}
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		waitForFileChange(m.watchChan),
		m.runCheck(),
	)
}
