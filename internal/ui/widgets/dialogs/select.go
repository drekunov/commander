package dialogs

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
)

// Select presents a list of options. In single-select mode Enter picks the
// focused option; in multi-select mode Space toggles options and Enter
// confirms the whole selection.
type Select struct {
	dialogBase

	text     string
	options  []string
	cursor   int
	selected map[int]bool
	multi    bool
	canceled bool

	done chan struct{}
}

func NewSelect(styles config.Styles, options []string, multi bool) *Select {
	return &Select{
		dialogBase: newDialogBase(styles),
		options:    options,
		multi:      multi,
		selected:   make(map[int]bool),
		done:       make(chan struct{}, 1),
	}
}

func (m *Select) Init() tea.Cmd {
	return nil
}

func (m *Select) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		m.handleKey(key)
	}

	return m, nil
}

func (m *Select) View() string {
	if !m.visible {
		return ""
	}

	lines := make([]string, 0, len(m.options))

	for index, option := range m.options {
		line := option

		if m.multi {
			marker := "[ ]"
			if m.selected[index] {
				marker = "[x]"
			}

			line = marker + " " + option
		}

		line = m.optionStyle(index == m.cursor).Render(line)

		lines = append(lines, line)
	}

	body := strings.Join(lines, "\n")

	if m.text != "" {
		body = lipgloss.JoinVertical(lipgloss.Center, m.textStyle().Render(m.text), body)
	}

	return m.render(body)
}

// Canceled reports whether the dialog was closed with Escape instead of a
// choice.
func (m *Select) Canceled() bool {
	return m.canceled
}

func (m *Select) SetText(text string) {
	m.text = text
}

func (m *Select) Selection() string {
	if len(m.options) == 0 {
		return ""
	}

	return m.options[m.cursor]
}

func (m *Select) Selections() []string {
	selected := make([]string, 0, len(m.options))

	for i, option := range m.options {
		if m.selected[i] {
			selected = append(selected, option)
		}
	}

	return selected
}

func (m *Select) Done() chan struct{} {
	return m.done
}

func (m *Select) handleKey(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < len(m.options)-1 {
			m.cursor++
		}
	case tea.KeySpace:
		if m.multi {
			m.selected[m.cursor] = !m.selected[m.cursor]
		}
	case tea.KeyEnter:
		m.close(false)
	case tea.KeyEsc:
		m.close(true)
	}
}

// close hides the dialog and signals Done; canceled records an Escape.
func (m *Select) close(canceled bool) {
	if !m.visible {
		return
	}

	m.canceled = canceled
	m.visible = false

	m.done <- struct{}{}
}
