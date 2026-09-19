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

// SetRefresher registers the sink that re-reads a batch of panels (wired to the
// app's batch reader by the composition root).
func (m *Model) SetRefresher(refresher func([]app.NavRequest)) {
	m.refresher = refresher
}

// refreshPanels asks the app to re-read both panels' current directories. It is
// a no-op when no refresher is wired or neither panel has a directory yet.
func (m *Model) refreshPanels() {
	if m.refresher == nil {
		return
	}

	requests := make([]app.NavRequest, 0, len(m.main.Panels()))

	for _, panelModel := range m.main.Panels() {
		if panelModel == nil || panelModel.Dir() == "" {
			continue
		}

		panelModel.PreserveSelection()

		requests = append(requests, app.NavRequest{Panel: panelModel.ID(), Dir: panelModel.Dir()})
	}

	if len(requests) == 0 {
		return
	}

	m.refresher(requests)
}

// SetData forwards a listing for one panel and its directory through the tea
// program so the table state is only mutated on the event loop.
func (m *Model) SetData(panelID app.PanelID, dir string, data []app.AttributeList) {
	if m.sendMsg == nil {
		return
	}

	m.sendMsg(panel.DataMsg{Panel: panelID, Path: dir, Data: data})
}
