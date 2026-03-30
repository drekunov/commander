package mainform

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/panel"
)

type Model struct {
	width, height int

	panel panel.Model
}

func New() *Model {
	return &Model{
		panel: panel.NewPanel(),
	}
}

func (m *Model) Init() tea.Cmd {
	return m.panel.Init()
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.panel, cmd = m.panel.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		if msg.Type == tea.KeyF10 {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width - 2
		m.height = msg.Height - 2
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	m.panel.SetWidth(m.width)
	m.panel.SetHeight(m.height)

	panelView := m.panel.View()

	out := config.Values.DialogBoxStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(panelView)

	return out
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

func (m *Model) SetPanel(panel panel.Model) {
	m.panel = panel
}

func (m *Model) Panel() panel.Model {
	return m.panel
}

func (m *Model) Width() int {
	return m.width
}

func (m *Model) Height() int {
	return m.height
}
