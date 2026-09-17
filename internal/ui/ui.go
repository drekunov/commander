package ui

import (
	"context"
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/buttonbar"
	"github.com/drekunov/gc/internal/ui/widgets/dialogs"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
	"github.com/drekunov/gc/internal/ui/widgets/panel"
)

const (
	dumpScreenKey = tea.KeyF12

	// sortWindowCaption is the sort window's caption; the instruction lives in
	// the caption rather than as a body line.
	sortWindowCaption = "Sort by"
)

type Model struct {
	program *tea.Program
	sendMsg func(tea.Msg)

	main *mainform.Model

	styles config.Styles

	dumper dumper

	// navigator resolves panel navigation requests (see SetNavigator).
	navigator func(app.PanelID, string)

	// quit is closed when the tea program exits so goroutines that show
	// dialogs never outlive the application.
	quit chan struct{}

	// mockMu guards the single mock dialog shared by repeated bar activations.
	mockMu     sync.Mutex
	mockDialog *dialogs.Info
	mockWinID  int
}

func New(styles config.Styles) *Model {
	model := &Model{
		main:   mainform.New(styles),
		styles: styles,
		dumper: &fileDumper{},
		quit:   make(chan struct{}),
	}

	model.program = tea.NewProgram(model, tea.WithMouseAllMotion(), tea.WithAltScreen())
	model.sendMsg = model.program.Send

	return model
}

func (m *Model) Run() error {
	_, err := m.program.Run()
	m.stop()

	if err != nil {
		return fmt.Errorf("program run failed: %w", err)
	}

	return nil
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.main.Init(), tea.HideCursor)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		if msg.Type == tea.KeyF10 {
			return m, tea.Quit
		}

		if msg.Type == dumpScreenKey {
			return m, m.dumpCmd(m.View())
		}

	case buttonbar.ActivateMsg:
		return m, m.handleBarActivation(msg.Action)

	case sortColumnMsg:
		if msg.column != "" {
			if p := m.main.FocusedPanel(); p != nil {
				p.SortBy(msg.column)
			}
		}

		return m, nil

	case panel.NavigateMsg:
		m.handleNavigate(msg)

		return m, nil
	}

	var cmd tea.Cmd

	m.main, cmd = m.main.Update(msg)

	return m, cmd
}

func (m *Model) View() string {
	return m.main.View()
}

// dumpCmd returns a command that writes the captured frame to a dump file,
// ignoring write errors so a failed dump never disturbs the running UI.
func (m *Model) dumpCmd(frame string) tea.Cmd {
	return func() tea.Msg {
		_ = m.dumper.Dump(frame)

		return nil
	}
}

// handleBarActivation resolves a function-button activation. Quit and Menu are
// real actions; every other button reports a mock placeholder through the Info
// dialog. This switch is the seam where real actions later replace the mocks.
func (m *Model) handleBarActivation(action buttonbar.Action) tea.Cmd {
	switch action {
	case buttonbar.ActionQuit:
		return tea.Quit

	case buttonbar.ActionMenu:
		return m.sortWindowCmd()
	}

	m.showMock(action)

	return nil
}

// sortColumnMsg carries the column chosen in the sort window back to the event
// loop, where the focused panel is mutated.
type sortColumnMsg struct {
	column string
}

// sortWindowCmd opens the modal column list for the focused panel and resolves
// to the chosen column. The dialog blocks, so it runs as a command off the
// event loop; an empty column means the window was canceled.
func (m *Model) sortWindowCmd() tea.Cmd {
	focused := m.main.FocusedPanel()
	if focused == nil {
		return nil
	}

	titles := focused.ColumnTitles()
	if len(titles) == 0 {
		return nil
	}

	return func() tea.Msg {
		return sortColumnMsg{column: m.Select(context.Background(), sortWindowCaption, "", titles)}
	}
}

// showMock reports a not-yet-implemented bar action. Only one mock dialog is
// kept open: a later activation updates the existing dialog's text in place
// instead of stacking another window and goroutine. Runs on the event loop.
func (m *Model) showMock(action buttonbar.Action) {
	text := action.Name() + " is not implemented yet"

	m.mockMu.Lock()

	if m.mockDialog != nil {
		m.mockDialog.SetText(text)
		m.mockMu.Unlock()

		m.sendMsg(tea.ResumeMsg{})

		return
	}

	info := dialogs.NewInfo(m.styles)
	info.SetText(text)
	info.SetVisible(true)

	m.mockDialog = info
	m.mockWinID = m.addDialogWindow(info, "Mock")
	m.mockMu.Unlock()

	m.sendMsg(tea.ResumeMsg{})

	go m.waitMock(info)
}

// waitMock releases the mock dialog when it is dismissed or the application
// quits, clearing the tracked reference and removing its window.
func (m *Model) waitMock(info *dialogs.Info) {
	select {
	case <-m.quit:
	case <-info.Done():
	}

	m.mockMu.Lock()

	if m.mockDialog != info {
		m.mockMu.Unlock()

		return
	}

	id := m.mockWinID
	m.mockDialog = nil
	m.mockWinID = 0
	m.mockMu.Unlock()

	m.main.WM().Remove(id)
}

// handleNavigate forwards a panel navigation request to the app's directory
// reader. It never blocks the event loop: the app replies asynchronously.
func (m *Model) handleNavigate(msg panel.NavigateMsg) {
	if m.navigator != nil {
		m.navigator(msg.Panel, msg.Dir)
	}
}

// stop closes the lifecycle channel, releasing any goroutine still waiting on
// a dialog. It is safe to call more than once.
func (m *Model) stop() {
	select {
	case <-m.quit:
	default:
		close(m.quit)
	}
}
