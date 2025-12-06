package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
	"github.com/evertras/bubble-table/table"
)

func rowFromAttrList(attrList app.AttributeList) table.Row {
	rawData := make(table.RowData, len(attrList))
	for key, attr := range attrList {
		rawData[key.String()] = attr
	}

	return table.NewRow(rawData)
}

func (m *Model) SetData(data []app.AttributeList) {
	panel := m.main.Panel()
	panel.SetVisible(true)
	view := panel.TableView()

	if len(data) == 0 {
		return
	}

	header := data[0]

	cols := make([]table.Column, 0, len(header))
	for key := range header {
		cols = append(cols, table.NewColumn(key.String(), key.String(), 1))
	}

	view = view.WithColumns(cols)

	rows := make([]table.Row, 0, len(data))
	for _, row := range data {
		rows = append(rows, rowFromAttrList(row))
	}

	view = view.WithRows(rows).BorderRounded().
		WithBaseStyle(config.Values.DialogBoxStyle).
		WithPageSize(10).
		Focused(true)

	panel.SetTableView(view)
	m.main.SetPanel(panel)
	m.program.Send(tea.ResumeMsg{})
}
