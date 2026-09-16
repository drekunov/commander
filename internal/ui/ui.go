package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/buttonbar"
	"github.com/drekunov/gc/internal/ui/widgets/dialogs"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
	"github.com/drekunov/gc/internal/ui/widgets/panel"
)

const dumpScreenKey = tea.KeyF12

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

// handleBarActivation resolves a function-button activation. Quit is the one
// real action; every other button reports a mock placeholder through the Info
// dialog. This switch is the seam where real actions later replace the mocks.
func (m *Model) handleBarActivation(action buttonbar.Action) tea.Cmd {
	if action == buttonbar.ActionQuit {
		return tea.Quit
	}

	go m.mockInfo(action)

	return nil
}

// mockInfo shows the placeholder Info dialog for a not-yet-implemented bar
// action. The dialog flow blocks, so it runs on a goroutine and unblocks as
// soon as the lifecycle channel closes (quit never leaks it).
func (m *Model) mockInfo(action buttonbar.Action) {
	info := dialogs.NewInfo(m.styles)
	info.SetText(action.Name() + " is not implemented yet")
	info.SetTitle("Mock")
	info.SetVisible(true)

	id := m.addDialogWindow(info, "Mock")
	defer m.main.WM().Remove(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-m.quit:
	case <-info.Done():
	}
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
