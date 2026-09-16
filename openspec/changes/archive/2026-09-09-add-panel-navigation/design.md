## Context

The commander shows two `panel` windows in a `wm.Manager`. Today the app reads `/` once in `app.Run`, delivers it to the left panel via `ui.SetData`, and blocks on `ctx.Done()`; the right panel is never loaded. Panels render generic attribute rows (header row + entries, per `panel/data-contract`) through bubbles `table`; the left panel's table is focused so Up/Down already move its cursor. Keys reach a panel only when the WM routes them to the focused window, and the function-button bar already claims F1-F10 plus arrows while it holds focus (see `add-function-button-bar`). See proposal.md - Why.

Constraints that shape the approach:

- `Connector.ReadDir(path string)` already reads any directory; the filesystem connector emits one row per entry with `Name` and `IsDir` attributes plus a per-entry `path` attribute holding the containing directory.
- `app` and `ui` run on separate goroutines; `app` already calls `ui.SetData`/`ui.Info` cross-goroutine and those flow through `program.Send`/WM locks.
- Bubbles `table` moves its cursor only from key messages when focused; `LineUp`/`LineDown`/`GotoTop` are key bindings, and `Model.SetCursor(n)` can position it directly.
- There is no existing ui→app channel: gestures cannot reach the connector today.
- Gates: `make build`, `make lint` (gofumpt + gci + golangci-lint v2, `default: all`), `go test -race ./...`, plus pty-driven `make run`.

## Goals / Non-Goals

**Goals:**

- Both panels become independent browsers with their own directory and cursor; navigation targets only the focused panel.
- Left = cursor to top, Right = cursor one down (Up/Down unchanged); Enter descends into the selected directory, Backspace ascends to the parent; both reload the focused panel for real via the connector.
- Unreadable directories never change the listing and surface through the existing error dialog; the app keeps running.
- Keep everything else (theming, dialogs, WM, function-button bar, screen dump, connector contract) behaviorally unchanged.

**Non-Goals:**

- No file actions (open/copy/move/delete), no history/forward-back, no path autocompletion, no sorting/filtering, no hidden-file toggles.
- No changes to the WM's focus/drag logic or to how the function-button bar reserves keys.
- No connector or data-contract changes (`ReadDir`/header/rows contract is untouched).

## Decisions

### D1. Panel browsing state lives in the UI; the app is a stateless directory reader

The app keeps no per-panel directory state. The focused panel already knows its current directory and its selected entry, so the UI computes the exact target path and asks the app to read it:

- **Enter** target = `filepath.Join(currentDir, selectedName)` when the selected entry `IsDir`; otherwise no request.
- **Backspace** target = `filepath.Dir(currentDir)`; when the parent equals the current directory (filesystem root) no request is sent.

The app just performs `connector.ReadDir(target)` and delivers the result back to the requesting panel. This keeps one source of truth (the listing the user sees) and avoids duplicating path bookkeeping in the app.

### D2. ui→app navigation requests travel over a channel; `ui` gets a navigator sink

`app.App` owns a buffered channel `navCh`. A new method `func (a *App) Navigate(panel app.PanelID, dir string)` enqueues a request (non-blocking; buffer 8). `app.Run` becomes a select loop: load both panels at `/`, then serve `ctx.Done()` and incoming navigation requests until cancelled — reading the directory and calling `ui.SetData(panel, dir, data)` on success, or `ui.Error` (existing blocking dialog, called from the app goroutine) on failure, leaving the panel's previous listing intact.

The composition root (`cmd/commander/main.go`) wires the sink: `uiInstance.SetNavigator(appInstance.Navigate)`. `ui.Model` stores a `func(app.PanelID, string)` and fires it when it receives a panel navigation request message.

- *Alternative considered*: ui calls a `Connector` directly. Rejected: it would leak the data source into the view layer and bypass the `app` orchestration/error contract.
- *Why non-blocking*: the request originates in the Bubble Tea event loop; blocking there would stall rendering. The app goroutine is the only reader, so buffered sends never drop in practice.

### D3. Active-panel detection falls out of WM focus routing

The WM already routes key messages only to the focused window's content. A navigation request therefore always originates from the focused panel, and no panel can move while a dialog owns focus (its window is the focused one and handles Enter itself). The panel widget emits the request; `mainform` and `ui.Model` only relay it. Both panel tables are focused from construction so the one receiving keys reacts immediately — the WM's routing, not table focus, decides which panel acts.

- Each panel is assigned an id (`app.PanelLeft` / `app.PanelRight`) by `mainform`, and every request carries that id so the app can deliver the reply to the right panel.

### D4. Left/Right are implemented as table keymap extensions; Enter/Backspace as widget messages

In the panel widget's `NewPanel`, extend the bubbles table keymap per panel instance:

- `GotoTop` gains the `left` key; `LineDown` gains the `right` key (defaults `up`/`down`/`k`/`j` preserved). Pressing Left jumps the cursor to row 0; Right steps one row and clamps at the bottom — matching `panel/navigation`.
- `Enter`/`Backspace` are not table keys. The panel's `Update` intercepts them (it only sees keys when it is the focused window) and returns a `tea.Cmd` that resolves to `panel.NavigateMsg{Panel, Dir}`. A dialog-focused Enter never reaches a panel; a bar-focused Enter/arrow is already claimed by the button bar before the WM.

- *Why keymap extension over message interception*: cursor moves must still honor table viewport scrolling and clamping, which the built-in handlers already do; reimplementing them in the panel would duplicate bubbles internals.

### D5. Delivered listings address a panel and carry its directory

`panel.DataMsg` gains `Panel app.PanelID` and `Path string`; `ui.SetData` and `panel.SetData` gain the same parameters (**BREAKING** internal API — `app.Panel` interface and callers/tests updated). `mainform` routes the message to the left or right window, updates the window title to the directory, and calls `panel.SetData(dir, data)`, which stores the current directory, rebuilds the table, and resets the cursor with `SetCursor(0)` so a reload lands on the first entry.

### D6. Errors never mutate the listing

Only a successful `ReadDir` calls `SetData`. On error the app calls the existing `Error` dialog (footer `[Error]`) with the connector message; the panel keeps its previous directory/listing/cursor and the app loop continues to the next request. Because the error dialog is a WM window, the underlying panel keeps focus behavior per the WM rules and remains usable once the dialog closes.

### D7. Verification strategy

- App-level tests with a fake `UI` (records `SetData(panel, dir, data)`, exposes `Error`) and a fake `Connector` over an in-memory dir map: startup loads both panels at `/`; a navigate request reloads the requested panel; a failing directory calls `Error` and issues no `SetData`.
- Panel-widget tests: Enter on a directory entry yields `NavigateMsg` with the joined child path; Enter on a file yields nothing; Backspace from a nested dir yields the parent path; Backspace at `/` yields nothing; cursor jumps to top on Left and steps down on Right; a reload resets the cursor.
- `mainform`/`ui` tests: `DataMsg` with a panel id updates that panel and its title; `NavigateMsg` invokes the configured navigator with the right id and dir.
- pty smoke (`make run`): both panels show `/`, Tab/click focus moves between them, Right steps down, Left jumps to top, Enter descends (window title changes), Backspace ascends, an unreadable directory (e.g. `/root`) reports the error and keeps the listing, F10 quits, F12 still dumps.

## Risks / Trade-offs

- **Enter/Backspace steal focus from panels that later grow file actions** → file actions will bind to other keys or an explicit mode; the navigation contract (`panel/navigation`) keeps Enter/Backspace reserved for directory traversal.
- **Left = top vs classic panel-switch** → the user chose Left=top/Right=down for this change; switching panels stays on Tab/click (unchanged). If panel-switch-on-arrow is wanted later it is a keymap-only change.
- **Reload resets the cursor to the top** → acceptable; re-selecting the previous entry on ascends is a later nicety and does not change the spec.
- **Rapid navigation enqueues reads the app serializes** → reads are sequential on the app goroutine; an older queued request could apply after a newer one. Buffer stays small (8) and each request carries an absolute path, so at worst the panel shows a stale-but-valid directory; no data race (all delivery via `program.Send`).
- **Long directory titles truncate** → the WM title bar already clips plain-text titles at the window width; full path shown, clipped. Cosmetic.

## Migration Plan

Internal API break confined to `internal/app`, `internal/ui`, `internal/ui/widgets/{panel,mainform}` and their tests; no external callers. Rollback is reverting the change's commits; `app.Panel.SetData` gains parameters atomically with the new message shape. The app loop replaces the current load-once body, so the first commit of this change must keep `go build`/`go test` green at the seam.

## Open Questions

- Whether to show the full current path in the window title or just its base — cosmetic; full path with existing truncation is chosen and can be revisited without spec impact.
