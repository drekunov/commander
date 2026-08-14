//nolint:testpackage // tests read the unexported table state
package panel

import (
	"testing"

	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
)

func TestSetData(t *testing.T) {
	t.Parallel()

	panelModel := NewPanel(config.Styles{})

	data := []app.AttributeList{
		{
			{AttrName: "Name", AttrValue: "Name"},
			{AttrName: "IsDir", AttrValue: "IsDir"},
		},
		{
			{AttrName: "Name", AttrValue: "a.txt"},
			{AttrName: "IsDir", AttrValue: "false"},
		},
		{
			{AttrName: "Name", AttrValue: "sub"},
			{AttrName: "IsDir", AttrValue: "true"},
		},
	}

	panelModel.SetData(data)

	cols := panelModel.tableView.Columns()
	if len(cols) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(cols))
	}

	if cols[0].Title != "Name" {
		t.Errorf("cols[0].Title = %q, want Name", cols[0].Title)
	}

	rows := panelModel.tableView.Rows()
	if len(rows) != 2 {
		t.Fatalf("expected 2 entry rows, got %d", len(rows))
	}

	if rows[0][0] != "a.txt" {
		t.Errorf("first entry not rendered: rows[0][0] = %q", rows[0][0])
	}

	if rows[1][0] != "sub" {
		t.Errorf("rows[1][0] = %q, want sub", rows[1][0])
	}

	if rows[1][1] != "true" {
		t.Errorf("rows[1][1] = %q, want true", rows[1][1])
	}
}

func TestSetDataNoRows(t *testing.T) {
	t.Parallel()

	panelModel := NewPanel(config.Styles{})

	panelModel.SetData(nil)

	if len(panelModel.tableView.Rows()) != 0 {
		t.Error("table should be unchanged when no data is provided")
	}
}
