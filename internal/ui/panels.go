package ui

import (
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/ui/widgets/panel"
)

func (m *Model) SetData(data []app.AttributeList) {
	if m.sendMsg == nil {
		return
	}

	m.sendMsg(panel.DataMsg{Data: data})
}
