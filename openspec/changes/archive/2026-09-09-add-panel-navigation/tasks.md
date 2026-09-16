## 1. App layer: panel-aware delivery and a navigation request channel

- [x] 1.1 Add `type PanelID int` with `PanelLeft`/`PanelRight` constants in `internal/app`; change the `Panel` interface method to `SetData(panel PanelID, dir string, data []AttributeList)` (`internal/app/ui.go`)
- [x] 1.2 Define a `NavRequest{Panel PanelID; Dir string}` type and add `navCh chan NavRequest` plus `func (a *App) Navigate(panel PanelID, dir string)` that enqueues a request (buffered channel)
- [x] 1.3 Rewrite `app.Run` (`internal/app/app.go`): load `/` into both panels at startup (a read error surfaces via the existing Info dialog and the app still starts), then serve a select loop over `ctx.Done()` and `navCh`, reading the requested directory on each request
- [x] 1.4 On a successful read call `ui.SetData(panel, dir, data)`; on a failed read call the existing error dialog with the connector message and issue no `SetData`, leaving the panel's previous listing untouched

## 2. Panel widget: directory tracking, key semantics, navigation requests

- [x] 2.1 Extend `panel.DataMsg` with `Panel app.PanelID` and `Path string`; change `panel.SetData` to `SetData(dir string, data []app.AttributeList)` storing the current directory and entry rows and resetting the cursor with `SetCursor(0)`
- [x] 2.2 Add a `NavigateMsg{Panel app.PanelID; Dir string}` type plus a `tea.Cmd` factory that resolves to it
- [x] 2.3 Give each panel an id (field + `SetPanelID(app.PanelID)` or constructor change) and extend the table keymap in `NewPanel`: add `left` to `GotoTop` and `right` to `LineDown`, preserving the default keys
- [x] 2.4 Intercept Enter/Backspace in `panel.Update` (keys arrive only while the panel is the focused window): on Enter over a directory entry return a `NavigateMsg` for `filepath.Join(dir, selectedName)`; Enter over a regular file returns nothing; Backspace returns `NavigateMsg` for `filepath.Dir(dir)` unless the parent equals the current directory (root) in which case nothing is sent
- [x] 2.5 Add a helper that maps the table cursor to the selected entry's `Name` and `IsDir` attributes from the stored attribute rows

## 3. UI routing and composition

- [x] 3.1 `mainform` assigns each panel its `app.PanelID` and focuses both table cursors at construction; handle the extended `panel.DataMsg` by applying the listing to the correct window and setting that window's title to the delivered path
- [x] 3.2 `ui.Model` (`internal/ui/panels.go`) implements the new `SetData(panel, dir, data)` by sending a `panel.DataMsg{Panel, Path, Data}` through the tea program
- [x] 3.3 `ui.Model` stores a `navigator func(app.PanelID, string)` set via `SetNavigator`, and a `panel.NavigateMsg` case in `Update` invokes it (fire-and-forget, returning nil)
- [x] 3.4 Wire the composition root (`cmd/commander/main.go`): after constructing the app, call `uiInstance.SetNavigator(appInstance.Navigate)`
- [x] 3.5 Update all remaining `SetData` call sites and fakes (internal/app tests, internal/ui tests, any widget tests) to the new signature

## 4. Tests and verification

- [x] 4.1 App tests with fake `UI`/`Connector`: startup delivers `/` to both panels; a `Navigate` request reads and delivers the requested directory to the requested panel; an unreadable directory calls `Error` and produces no `SetData`
- [x] 4.2 Panel widget tests: Enter on a directory yields `NavigateMsg` with the joined child path; Enter on a file yields no message; Backspace from a nested directory yields the parent path; Backspace at `/` yields no message; a delivered listing resets the cursor to the top
- [x] 4.3 Panel widget tests for the keymap: with the panel focused, Right moves the cursor one row down (clamping at the bottom) and Left jumps to the first row
- [x] 4.4 `mainform`/`ui` tests: `DataMsg` addressed to the right panel updates the right window (left untouched) and sets the window title; `NavigateMsg` invokes the configured navigator with the correct panel id and directory
- [x] 4.5 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 4.6 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 4.7 Manual pty smoke (`make run`): both panels show `/`; Tab/click moves focus between them; Right steps down, Left jumps to the top, Enter descends (window title changes), Backspace ascends; opening an unreadable directory (e.g. `/root`) reports the error and keeps the current listing; the other panel still navigates; F10 quits; F12 still writes a screen dump
