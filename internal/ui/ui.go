package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
)

const dumpScreenKey = tea.KeyF12

type Model struct {
	program *tea.Program
	sendMsg func(tea.Msg)

	main *mainform.Model

	styles config.Styles

	dumper dumper
}

func New(styles config.Styles) *Model {
	model := &Model{
		main:   mainform.New(styles),
		styles: styles,
		dumper: &fileDumper{},
	}

	model.program = tea.NewProgram(model, tea.WithMouseAllMotion(), tea.WithAltScreen())
	model.sendMsg = model.program.Send

	return model
}

func (m *Model) Run() error {
	_, err := m.program.Run()
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
