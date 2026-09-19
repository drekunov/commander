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
	"github.com/drekunov/gc/internal/ui/widgets/topmenu"
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

	// refresher resolves a batch refresh of both panels (see SetRefresher).
	refresher func([]app.NavRequest)

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
		if cmd, handled := m.handleGlobalKey(msg); handled {
			return m, cmd
		}

	case buttonbar.ActivateMsg:
		m.main.Bar().ClearPressed()

		return m, m.handleBarActivation(msg.Action)

	case topmenu.ActivateMsg:
		return m, m.handleTopMenuActivation(msg)

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

// handleGlobalKey processes keys the ui owns before windows see them: F10
// quits, F12 dumps the screen, and Ctrl+R refreshes both panels unless a modal
// dialog is open. It reports whether the key was handled.
func (m *Model) handleGlobalKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.Type {
	case tea.KeyF10:
		return tea.Quit, true
	case dumpScreenKey:
		return m.dumpCmd(m.View()), true
	case tea.KeyCtrlR:
		if !m.main.DialogOpen() {
			m.refreshPanels()
		}

		return nil, true
	}

	return nil, false
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

	case buttonbar.ActionPullDn:
		m.main.TopMenu().Toggle()

		return nil
	}

	m.showMockText(action.Name() + " is not implemented yet")

	return nil
}

// handleTopMenuActivation resolves an item chosen in the top menu. Files items
// carry a function-button action and route through the shared dispatch; every
// other item is a placeholder reported through the mock dialog. The menu is
// deactivated before dispatch so a dialog never opens under an active menu.
func (m *Model) handleTopMenuActivation(msg topmenu.ActivateMsg) tea.Cmd {
	m.main.TopMenu().Deactivate()

	if msg.Action != 0 {
		return m.handleBarActivation(msg.Action)
	}

	m.showMockText(msg.Label + " is not implemented yet")

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

// showMockText reports a not-yet-implemented action. Only one mock dialog is
// kept open: a later activation updates the existing dialog's text in place
// instead of stacking another window and goroutine. It runs on the event loop,
// so it must not send a message to the program (the message channel is
// unbuffered and single-reader): the window is rendered by the normal
// post-Update render.
func (m *Model) showMockText(text string) {
	m.mockMu.Lock()

	if m.mockDialog != nil {
		m.mockDialog.SetText(text)
		m.mockMu.Unlock()

		return
	}

	info := dialogs.NewInfo(m.styles)
	info.SetText(text)
	info.SetVisible(true)

	m.mockDialog = info
	m.mockWinID = m.addDialogWindow(info, "Unimplemented")
	m.mockMu.Unlock()

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

	m.closeDialogWindow(id)
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
