## Context

See proposal.md - Why. Existing constraints that shape the approach:

- `App` keeps a single latest-wins navigation slot (`pending *NavRequest`), so calling `Navigate` once per panel would drop the first request. Refreshing both panels needs a batch path.
- `App.Run` seeds both panels from `rootDir` by calling `loadPanel` for `PanelLeft` then `PanelRight`, so delivering two `SetData` calls from one loop iteration is already an established shape.
- `panel.Model.SetData` currently repositions the cursor to the first row (or to `pendingSelect` when ascending) on every delivery, so a naive re-delivery would lose the selection.
- `mainform` owns both panels and holds the modal `dialogOpen()` state; global keys such as F10 are already handled above the window manager.
- `ui.Model` owns the app-facing sinks (`navigator`) and handles global keys (F10, F12) before forwarding to `mainform`.

## Goals / Non-Goals

**Goals:**

- `Ctrl+R` re-reads both panels' current directories in one operation.
- The operation is a no-op while a modal dialog is open and does not change focus.
- The selected entry survives the refresh when it still exists.
- Keep the existing latest-wins behavior of ordinary navigation unchanged.

**Non-Goals:**

- No periodic/automatic refresh, file watching, or polling.
- No change to the `app.UI` interface.
- No new dependency and no upgrade of Bubble Tea.
- No refresh of a panel that has no current directory yet.

## Decisions

### D1. Handle `Ctrl+R` in `ui.Model.Update`, guarded by modal state

Add `Ctrl+R` to the `tea.KeyMsg` branch of `ui.Model.Update`, next to the existing F10/F12 handling. It calls `m.main.DialogOpen()` (a new exported wrapper over `mainform.dialogOpen`) and, when no dialog is open, invokes a new `refresher` sink with a batch of `app.NavRequest`s built from both panels' current directories.

- *Alternative considered*: handle the key in `mainform` and return a message for `ui` to resolve. Rejected: `ui` already hosts global keys and owns the app sink, so handling it there (with an exported modal-state check) is fewer moving parts.
- *Alternative considered*: bind a bar button instead. Rejected: the request is for a keyboard shortcut; the bar's buttons are mock actions and F5 is already Copy.

### D2. Batch read via `App.Refresh`, keeping latest-wins for `Navigate`

Generalize `App.pending` from `*NavRequest` to `[]NavRequest`. `takePending` returns the slice and `Run` handles every request in it. `Navigate` stores a single-element slice (unchanged latest-wins semantics); add `Refresh(requests []NavRequest)` that stores the whole batch and wakes the loop. Replacements remain atomic under the existing mutex.

- *Alternative considered*: merge requests per panel inside `Navigate`. Rejected: it silently changes single-key navigation semantics and the documented "latest request wins" behavior; a separate batch entry point is clearer.
- *Alternative considered*: add a second signal channel for refresh. Rejected: one pending slot with a slice is simpler and reuses the existing wake-up path.

### D3. Preserve the selection with an explicit panel method

`SetData` keeps its current "reset to the first row (or to `pendingSelect`)" behavior, which existing tests (`TestReloadResetsCursorToTop`, `TestReloadShowsFirstRow`) lock in; a same-directory delivery is not by itself a refresh. Add `panel.Model.PreserveSelection()`, which stores the currently selected entry name in the existing `pendingSelect` field. `ui.refreshPanels` calls it on each panel before requesting the re-read, so the subsequent `SetData` restores and scrolls to that entry; if the entry is gone, `cursorForPendingSelect` returns 0 and the first row is selected. `SetData` clears `pendingSelect`, so the mark is one-shot.

- *Alternative considered*: auto-preserve whenever `SetData` receives the current directory. Rejected: it broke the documented "reload resets to top" behavior and its tests; production same-directory deliveries are not always refreshes.
- *Alternative considered*: add a separate `panel.Refresh` right. Rejected: the panel does not read directories itself; `PreserveSelection` reuses the existing `pendingSelect` restoration path with no new state.

### D4. Wiring

- Add `ui.Model.SetRefresher(func([]app.NavRequest))` and call it from `cmd/commander/main.go` with `appInstance.Refresh`.
- Add `mainform.Model.Panels() []*panel.Model` (left, right) and `panel.Model.ID() app.PanelID` so `ui` can build the request batch.

## Risks / Trade-offs

- **`Ctrl+R` may be sent by some terminals as a rune sequence** → Bubble Tea maps `ctrl+r` to `KeyCtrlR`; verify in the manual smoke.
- **Two panel reads per refresh** → one extra `ReadDir` per panel per keypress, performed synchronously in the app loop as today; acceptable for a user-initiated action.
- **Selection preservation depends on entry names** → renamed entries fall back to the first row by design; document in the spec scenario.
- **Generalizing `pending` to a slice touches the navigation hot path** → existing `Navigate` stores one element, so `TestNavigateKeepsLatestRequest` and the other app tests must stay green as the regression guard.

## Migration Plan

No migration. Apply core first (app batch → panel preservation → mainform accessors → ui key + sink → wiring), keep the tree buildable, add tests, then run `go build ./...`, `go vet ./...`, `go test -race ./...`, and `make lint`.

## Open Questions

None.
