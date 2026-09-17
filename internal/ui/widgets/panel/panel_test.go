//nolint:testpackage // tests read the unexported table state
package panel

import (
	"fmt"
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
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
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

func TestBackspaceRestoresCursorToChild(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetPanelID(app.PanelLeft)
	model.Focus()

	model.SetData("/a", testListing("d:b", "f1"))

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = mustPanelModel(t, updated)

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = mustPanelModel(t, updated)

	model.SetData("/a/b", testListing("d:c", "f2"))

	if model.tableView.Cursor() != 0 {
		t.Fatalf("cursor after descending = %d, want 0", model.tableView.Cursor())
	}

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if cmd == nil {
		t.Fatal("Backspace produced no command")
	}

	model.SetData("/a", testListing("d:b", "f1"))

	if model.tableView.Cursor() != 1 {
		t.Errorf("cursor after ascent = %d, want 1 (entry b)", model.tableView.Cursor())
	}
}

func TestEnterOnParentRowRestoresCursorToChild(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetPanelID(app.PanelLeft)
	model.Focus()

	model.SetData("/a/b", testListing("d:c"))

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter on the parent row produced no command")
	}

	model.SetData("/a", testListing("d:b", "f1"))

	if model.tableView.Cursor() != 1 {
		t.Errorf("cursor after ascent = %d, want 1 (entry b)", model.tableView.Cursor())
	}
}

func TestAscendRestoresCursorAtEachLevel(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetPanelID(app.PanelLeft)
	model.Focus()

	model.SetData("/a", testListing("d:b"))
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = mustPanelModel(t, updated)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = mustPanelModel(t, updated)

	model.SetData("/a/b", testListing("d:c"))
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = mustPanelModel(t, updated)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = mustPanelModel(t, updated)

	model.SetData("/a/b/c", testListing("f1"))

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if cmd == nil {
		t.Fatal("Backspace produced no command")
	}

	model.SetData("/a/b", testListing("d:c"))

	if model.tableView.Cursor() != 1 {
		t.Fatalf("cursor in /a/b = %d, want 1 (entry c)", model.tableView.Cursor())
	}

	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if cmd == nil {
		t.Fatal("Backspace produced no command")
	}

	model.SetData("/a", testListing("d:b"))

	if model.tableView.Cursor() != 1 {
		t.Errorf("cursor in /a = %d, want 1 (entry b)", model.tableView.Cursor())
	}
}

func TestAscendFallsBackWhenChildMissing(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetPanelID(app.PanelLeft)
	model.Focus()

	model.SetData("/a/b", testListing("f1"))

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if cmd == nil {
		t.Fatal("Backspace produced no command")
	}

	model.SetData("/a", testListing("f1", "f2"))

	if model.tableView.Cursor() != 0 {
		t.Errorf("cursor after ascent = %d, want 0 (fallback)", model.tableView.Cursor())
	}
}

func TestReloadResetsCursorToTop(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("fa", "fb", "fc"))
	model.Focus()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = mustPanelModel(t, updated)

	if model.tableView.Cursor() != 1 {
		t.Fatalf("cursor = %d, want 1 after one Down", model.tableView.Cursor())
	}

	model.SetData("/a", testListing("dx", "dy"))

	if model.tableView.Cursor() != 0 {
		t.Errorf("cursor after reload = %d, want 0", model.tableView.Cursor())
	}
}

func TestRightJumpsToBottom(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("fa", "fb", "fc"))
	model.Focus()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = mustPanelModel(t, updated)

	if model.tableView.Cursor() != 2 {
		t.Errorf("cursor after Right = %d, want last row 2", model.tableView.Cursor())
	}
}

func TestRightOnLastRowStaysPut(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("fa", "fb", "fc"))
	model.Focus()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = mustPanelModel(t, updated)

	if model.tableView.Cursor() != 2 {
		t.Fatalf("cursor = %d, want last row 2 before extra Right", model.tableView.Cursor())
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = mustPanelModel(t, updated)

	if model.tableView.Cursor() != 2 {
		t.Errorf("cursor after Right on last row = %d, want 2", model.tableView.Cursor())
	}
}

func TestLeftJumpsToTop(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("fa", "fb", "fc"))
	model.Focus()

	for range 2 {
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
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

func TestUpAndDownStepOneRow(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetData("/", testListing("fa", "fb", "fc"))
	model.Focus()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = mustPanelModel(t, updated)

	if model.tableView.Cursor() != 1 {
		t.Fatalf("cursor after one Down = %d, want 1", model.tableView.Cursor())
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = mustPanelModel(t, updated)

	if model.tableView.Cursor() != 0 {
		t.Errorf("cursor after one Up = %d, want 0", model.tableView.Cursor())
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

	model := NewPanel(config.Styles{CursorStyle: lipgloss.NewStyle().Background(lipgloss.Color("#C0C0C0"))})
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

// tallListing returns a listing of n files followed by a directory named
// target, so the listing is taller than a small panel viewport.
func tallListing(n int) []app.AttributeList {
	names := make([]string, 0, n+1)
	for i := range n {
		names = append(names, fmt.Sprintf("file%02d", i))
	}

	return testListing(append(names, "d:target")...)
}

func TestAscendScrollsRestoredEntryIntoView(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetPanelID(app.PanelLeft)
	model.SetWidth(40)
	model.SetHeight(5)
	model.Focus()

	listing := tallListing(19)

	model.SetData("/a", listing)

	// Render once so the panel applies its viewport size, as the app does
	// every frame, before navigating.
	model.View()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = mustPanelModel(t, updated)

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = mustPanelModel(t, updated)

	model.SetData("/a/target", nil)
	model.View()

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if cmd == nil {
		t.Fatal("Backspace produced no command")
	}

	model.SetData("/a", listing)

	if model.tableView.Cursor() != 20 {
		t.Fatalf("cursor after ascent = %d, want 20 (entry target)", model.tableView.Cursor())
	}

	if view := model.View(); !strings.Contains(view, "target") {
		t.Errorf("restored entry not visible in rendered view:\n%s", view)
	}
}

func TestReloadShowsFirstRow(t *testing.T) {
	t.Parallel()

	model := NewPanel(config.Styles{})
	model.SetWidth(40)
	model.SetHeight(5)
	model.Focus()

	listing := tallListing(19)

	model.SetData("/", listing)

	// Render once so the panel applies its viewport size, as the app does
	// every frame, before scrolling.
	model.View()

	for range 19 {
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = mustPanelModel(t, updated)
	}

	model.SetData("/", listing)

	if view := model.View(); !strings.Contains(view, "file00") {
		t.Errorf("first row not visible after reload:\n%s", view)
	}
}

func TestFocusedHeaderSharesTableBackground(t *testing.T) {
	t.Parallel()

	tableBackground := lipgloss.Color("#000080")
	model := NewPanel(config.Styles{TableStyle: lipgloss.NewStyle().Background(tableBackground)})
	model.Focus()

	if got := model.tableStyles().Header.GetBackground(); got != tableBackground {
		t.Errorf("focused header background = %v, want %v", got, tableBackground)
	}
}

func TestBlurredHeaderSharesTableBackground(t *testing.T) {
	t.Parallel()

	tableBackground := lipgloss.Color("#000080")
	model := NewPanel(config.Styles{TableStyle: lipgloss.NewStyle().Background(tableBackground)})
	model.Focus()
	model.Blur()

	if got := model.tableStyles().Header.GetBackground(); got != tableBackground {
		t.Errorf("blurred header background = %v, want %v", got, tableBackground)
	}
}

func TestHeaderFollowsCustomTableBackground(t *testing.T) {
	t.Parallel()

	tableBackground := lipgloss.Color("#123456")
	model := NewPanel(config.Styles{TableStyle: lipgloss.NewStyle().Background(tableBackground)})

	if got := model.tableStyles().Header.GetBackground(); got != tableBackground {
		t.Errorf("header background = %v, want injected %v", got, tableBackground)
	}
}

func TestCursorUsesInjectedStyle(t *testing.T) {
	t.Parallel()

	cursor := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#C0C0C0"))

	model := NewPanel(config.Styles{CursorStyle: cursor})
	model.Focus()

	got := model.tableStyles().Selected

	if got.GetBackground() != cursor.GetBackground() {
		t.Errorf("cursor background = %v, want %v", got.GetBackground(), cursor.GetBackground())
	}

	if got.GetForeground() != cursor.GetForeground() {
		t.Errorf("cursor foreground = %v, want %v", got.GetForeground(), cursor.GetForeground())
	}

	if got.GetBold() != cursor.GetBold() {
		t.Errorf("cursor bold = %v, want %v", got.GetBold(), cursor.GetBold())
	}

	model.Blur()

	if blurred := model.tableStyles().Selected; blurred.GetBackground() != lipgloss.NewStyle().GetBackground() {
		t.Errorf("blurred cursor background = %v, want no highlight", blurred.GetBackground())
	}
}
