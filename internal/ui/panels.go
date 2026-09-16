package ui

import (
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/ui/widgets/panel"
)

// SetNavigator registers the sink that resolves panel navigation requests
// (wired to the app's directory reader by the composition root).
func (m *Model) SetNavigator(navigator func(app.PanelID, string)) {
	m.navigator = navigator
}

// SetData forwards a listing for one panel and its directory through the tea
// program so the table state is only mutated on the event loop.
func (m *Model) SetData(panelID app.PanelID, dir string, data []app.AttributeList) {
	if m.sendMsg == nil {
		return
	}

	m.sendMsg(panel.DataMsg{Panel: panelID, Path: dir, Data: data})
}
