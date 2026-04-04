package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
)

func rowFromAttrList(attrList app.AttributeList) table.Row {
	rawData := make(table.Row, 0, len(attrList))
	for _, attr := range attrList {
		rawData = append(rawData, fmt.Sprintf("%v", attr.AttrValue))
	}

	return rawData
}

func (m *Model) SetData(data []app.AttributeList) {
	panel := m.main.Panel()
	if panel == nil {
		return
	}

	view := panel.TableView()

	if len(data) == 0 {
		return
	}

	header := data[0]

	cols := make([]table.Column, 0, len(header))
	for _, attr := range header {
		cols = append(
			cols,
			table.Column{
				Title: attr.AttrName,
				Width: 12,
			},
		)
	}

	view.SetColumns(cols)

	// data[0] is the header row used for column titles above; skip it here.
	rows := make([]table.Row, 0, len(data)-1)
	for _, row := range data[1:] {
		rows = append(rows, rowFromAttrList(row))
	}

	view.SetRows(rows)

	panel.SetTableView(view)
	m.program.Send(tea.ResumeMsg{})
}
