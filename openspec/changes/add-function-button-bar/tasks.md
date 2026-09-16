## 1. Button-bar theme styles

- [x] 1.1 Add a `.button-bar` class to `internal/config/styles.css` (NC-classic fill/background for the bottom strip, distinct from `.button`)
- [x] 1.2 Add a `ButtonBarStyle lipgloss.Style` field to `config.Styles` (`internal/config/config.go:13`) and populate it in `LoadStyles` via `ApplyStyle(stylesheet[".button-bar"])`

## 2. `buttonbar` widget

- [x] 2.1 Create `internal/ui/widgets/buttonbar/buttonbar.go`: an `Action` type with the ten ordered actions (`1Help 2Menu 3View 4Edit 5Copy 6RenMov 7Mkdir 8Delete 9PullDn 10Quit`) plus the label table and an F-key (F1-F10) → action map
- [x] 2.2 Add an `ActivateMsg{Action}` message type and a `tea.Cmd` factory that resolves to it
- [x] 2.3 Implement `Model` state (focused index, pressed action for one-frame flash, focused bool) and `New(styles config.Styles)` using the existing `button`-style rendering paths
- [x] 2.4 Render exactly one full-width row: ten equal segments, labels truncated to fit each segment, row padded/filled so it ends at the right edge and never wraps
- [x] 2.5 Handle input in `Update`: Left/Right move focus, Enter activates the focused button (only meaningful while the bar is focused), F1-F10 keys activate directly, mouse press hit-tests the segment under the column and activates it
- [x] 2.6 Pressed feedback: the activated button renders in its pressed styling on the frame of activation

## 3. `mainform` integration

- [x] 3.1 Reserve the bottom row: in `mainform.Update` handle `tea.WindowSizeMsg` by calling `m.wm.SetSize(width, height-1)` instead of letting the raw message reach `wm.Update`; keep panels laid out to `height-1`; guard against degenerate heights (≤ 1 row → render the bar only)
- [x] 3.2 Compose the view: `mainform.View()` renders the WM canvas joined with the one-row bar on the bottom line (skip the WM part when its canvas is empty)
- [x] 3.3 Intercept F1-F9 key messages in `mainform.Update` *before* forwarding to the WM, map them to bar activations, and return the resulting `ActivateMsg` cmd so the focused window never receives them (F10 quit and F12 dump handling stays untouched)
- [x] 3.4 Intercept a mouse press on the bottom bar row in `mainform.Update` → bar focus on + activate the clicked button; a press above the bar row clears bar focus; all other mouse messages (motion/release, drags) keep flowing to the WM unchanged
- [x] 3.5 Bar-focus mode: when the bar is focused, route Left/Right/Enter to the bar; Tab exits bar focus back to the WM focus cycle; otherwise those keys reach the focused window as today

## 4. Action dispatch and mocks in `ui.Model`

- [x] 4.1 Add a `buttonbar.ActivateMsg` case to `ui.Model.Update` (`internal/ui/ui.go`) switching on the action: `Quit` → `tea.Quit`; every other action → start the mock report
- [x] 4.2 Implement the mock report as a goroutine that reuses the existing blocking `Info` flow (add dialog window, `Send(tea.ResumeMsg{})`, wait on dialog `Done()`/ctx) and shows a placeholder naming the action (e.g. title `Mock`, message `"<action> is not implemented yet"`)
- [x] 4.3 Give `ui.Model` a context cancelled when the tea program exits (create in `New`, cancel at the end of `Run`) and pass it to the mock goroutine so quit never leaks a waiting dialog goroutine

## 5. Tests and verification

- [x] 5.1 Unit-test `buttonbar`: ten labels render left-to-right in order; the row ends at the right edge at narrow widths without wrapping; F-key/click/Enter map to the right action; pressed flash appears and clears
- [x] 5.2 Unit-test that `mainform` returns an `ActivateMsg` cmd for F-keys and that bar-row mouse presses are consumed while above-bar presses still reach the WM path
- [x] 5.3 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 5.4 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 5.5 Manual smoke with `make run`: the bar spans the full width below both panels; F5 (and mouse click on 5Copy) opens the mock dialog naming Copy and closes on confirm; pressing a panel key still works; F10 quits; F12 still writes a screen dump; panels are not obscured by the bar
