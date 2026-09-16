//nolint:testpackage // tests read the unexported table state
package panel

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
)

const (
	attrNameLower  = "name"
	attrIsDirLower = "isdir"
	subName        = "sub"
)

// testListing builds a listing whose entries are files unless their name
// carries the "d:" marker, in which case they are directories.
func mustPanelModel(t *testing.T, model tea.Model) *Model {
	t.Helper()

	typed, ok := model.(*Model)
	if !ok {
		t.Fatalf("Update returned %T, want *Model", model)
	}

	return typed
}

func testListing(names ...string) []app.AttributeList {
	data := make([]app.AttributeList, 0, 1+len(names))
	data = append(data, app.AttributeList{
		{AttrName: attrName, AttrValue: attrName},
		{AttrName: attrIsDir, AttrValue: attrIsDir},
	})

	for _, name := range names {
		isDir := strings.HasPrefix(name, "d:")
		if isDir {
			name = strings.TrimPrefix(name, "d:")
		}

		data = append(data, app.AttributeList{
			{AttrName: attrName, AttrValue: name},
			{AttrName: attrIsDir, AttrValue: isDir},
		})
	}

	return data
}

func TestSetData(t *testing.T) {
	t.Parallel()

	panelModel := NewPanel(config.Styles{})

	data := []app.AttributeList{
		{
			{AttrName: attrName, AttrValue: attrName},
			{AttrName: attrIsDir, AttrValue: attrIsDir},
		},
		{
			{AttrName: attrName, AttrValue: "a.txt"},
			{AttrName: attrIsDir, AttrValue: "false"},
		},
		{
			{AttrName: attrName, AttrValue: subName},
			{AttrName: attrIsDir, AttrValue: "true"},
		},
	}

	panelModel.SetData("/", data)

	cols := panelModel.tableView.Columns()
	if len(cols) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(cols))
	}

	if cols[0].Title != attrName {
		t.Errorf("cols[0].Title = %q, want Name", cols[0].Title)
	}

	rows := panelModel.tableView.Rows()
	if len(rows) != 2 {
		t.Fatalf("expected 2 entry rows, got %d", len(rows))
	}

	if rows[0][0] != "a.txt" {
		t.Errorf("first entry not rendered: rows[0][0] = %q", rows[0][0])
	}

	if rows[1][0] != subName {
		t.Errorf("rows[1][0] = %q, want sub", rows[1][0])
	}

	if rows[1][1] != "true" {
		t.Errorf("rows[1][1] = %q, want true", rows[1][1])
	}
}

func TestSetDataNoRows(t *testing.T) {
	t.Parallel()

	panelModel := NewPanel(config.Styles{})

	panelModel.SetData("/", nil)

	if len(panelModel.tableView.Rows()) != 0 {
		t.Error("table should be empty when no data is provided")
	}
}

func TestEnterOnDirectoryEmitsNavigateMsg(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetPanelID(app.PanelRight)
	model.SetData("/", testListing("d:sub"))

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter on a directory produced no command")
	}

	msg, ok := cmd().(NavigateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want NavigateMsg", cmd())
	}

	if msg.Panel != app.PanelRight {
		t.Errorf("NavigateMsg panel = %v, want right", msg.Panel)
	}

	if msg.Dir != "/sub" {
		t.Errorf("NavigateMsg dir = %q, want /sub", msg.Dir)
	}
}

func TestEnterOnFileEmitsNothing(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetPanelID(app.PanelLeft)
	model.SetData("/etc", testListing("fhosts", "fpasswd", "d:sub"))

	model.Focus()
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = mustPanelModel(t, updated)

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Errorf("Enter on a file produced a command %v, want none", cmd())
	}
}

func TestBackspaceEmitsParent(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetPanelID(app.PanelLeft)
	model.SetData("/a/b", testListing("d:c"))

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if cmd == nil {
		t.Fatal("Backspace produced no command")
	}

	msg, ok := cmd().(NavigateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want NavigateMsg", cmd())
	}

	if msg.Dir != "/a" {
		t.Errorf("NavigateMsg dir = %q, want /a", msg.Dir)
	}
}

func TestBackspaceAtRootEmitsNothing(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("d:sub"))

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if cmd != nil {
		t.Errorf("Backspace at the root produced a command %v, want none", cmd())
	}
}

func TestReloadResetsCursorToTop(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("fa", "fb", "fc"))
	model.Focus()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = mustPanelModel(t, updated)

	if model.tableView.Cursor() != 1 {
		t.Fatalf("cursor = %d, want 1 after one Right", model.tableView.Cursor())
	}

	model.SetData("/a", testListing("dx", "dy"))

	if model.tableView.Cursor() != 0 {
		t.Errorf("cursor after reload = %d, want 0", model.tableView.Cursor())
	}
}

func TestRightMovesDownAndClamps(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("fa", "fb", "fc"))
	model.Focus()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = mustPanelModel(t, updated)

	if model.tableView.Cursor() != 1 {
		t.Errorf("cursor after Right = %d, want 1", model.tableView.Cursor())
	}

	for range 2 {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
		model = mustPanelModel(t, updated)
	}

	if model.tableView.Cursor() != 2 {
		t.Errorf("cursor after extra Rights = %d, want last row 2", model.tableView.Cursor())
	}
}

func TestLeftJumpsToTop(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("fa", "fb", "fc"))
	model.Focus()

	for range 2 {
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
		model = mustPanelModel(t, updated)
	}

	if model.tableView.Cursor() != 2 {
		t.Fatalf("cursor = %d, want 2 before Left", model.tableView.Cursor())
	}

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = mustPanelModel(t, updated)

	if model.tableView.Cursor() != 0 {
		t.Errorf("cursor after Left = %d, want 0", model.tableView.Cursor())
	}
}

func TestParentRowLeadsNonRootListing(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/etc", testListing("fhosts", "fpasswd"))

	rows := model.tableView.Rows()
	if len(rows) != 3 {
		t.Fatalf("got %d rows, want 3 (parent + 2 entries)", len(rows))
	}

	if rows[0][0] != parentLabel {
		t.Errorf("first row cell = %q, want %q", rows[0][0], parentLabel)
	}

	if rows[1][0] != "fhosts" {
		t.Errorf("second row cell = %q, want first entry fhosts", rows[1][0])
	}
}

func TestParentRowHiddenAtRoot(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("fa", "fb"))

	rows := model.tableView.Rows()
	if len(rows) != 2 {
		t.Fatalf("got %d rows at the root, want 2", len(rows))
	}

	if rows[0][0] != "fa" {
		t.Errorf("first root row cell = %q, want first entry fa", rows[0][0])
	}
}

func TestEnterOnParentRowAscends(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetPanelID(app.PanelLeft)
	model.SetData("/a/b", testListing("f1"))

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter on the parent row produced no command")
	}

	msg, ok := cmd().(NavigateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want NavigateMsg", cmd())
	}

	if msg.Dir != "/a" {
		t.Errorf("NavigateMsg dir = %q, want /a", msg.Dir)
	}
}

func TestDirectoryBelowParentStillDescends(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/a", testListing("d:sub"))
	model.Focus()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = mustPanelModel(t, updated)

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter on the directory below the parent row produced no command")
	}

	msg, ok := cmd().(NavigateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want NavigateMsg", cmd())
	}

	if msg.Dir != "/a/sub" {
		t.Errorf("NavigateMsg dir = %q, want /a/sub", msg.Dir)
	}
}

func TestAttributeMatchingIsCaseInsensitive(t *testing.T) {
	t.Parallel()

	data := []app.AttributeList{
		{
			{AttrName: attrNameLower, AttrValue: attrNameLower},
			{AttrName: attrIsDirLower, AttrValue: attrIsDirLower},
		},
		{
			{AttrName: attrNameLower, AttrValue: subName},
			{AttrName: attrIsDirLower, AttrValue: true},
		},
	}

	model := NewPanel(config.Styles{})
	model.SetPanelID(app.PanelLeft)
	model.SetData("/a", data)
	model.Focus()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = mustPanelModel(t, updated)

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter on a lowercase-attribute directory produced no command")
	}

	msg, ok := cmd().(NavigateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want NavigateMsg", cmd())
	}

	if msg.Dir != "/a/sub" {
		t.Errorf("NavigateMsg dir = %q, want /a/sub", msg.Dir)
	}
}

func TestEmptyDirectoryShowsParentRow(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/a", nil)

	rows := model.tableView.Rows()
	if len(rows) != 1 {
		t.Fatalf("empty non-root directory rendered %d rows, want 1", len(rows))
	}

	if rows[0][0] != parentLabel {
		t.Errorf("only row cell = %q, want %q", rows[0][0], parentLabel)
	}

	if model.tableView.Cursor() != 0 {
		t.Errorf("cursor = %d, want 0 on the parent row", model.tableView.Cursor())
	}
}

func TestBlurHidesCursor(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("fa", "fb"))
	model.SetWidth(40)
	model.SetHeight(10)

	model.Focus()
	focused := model.tableStyles().Selected.GetBackground()

	model.Blur()
	blurred := model.tableStyles().Selected.GetBackground()

	if focused == blurred {
		t.Error("blurring the panel did not change the cursor highlight")
	}

	if blurred != lipgloss.NewStyle().GetBackground() {
		t.Error("blurred panel still renders a cursor highlight")
	}
}
