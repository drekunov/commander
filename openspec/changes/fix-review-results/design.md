## Context

The current working tree implements the function-button bar and panel navigation (see `add-function-button-bar`). A code review of that work produced a prioritized list of defects. This change fixes them. Relevant current-state constraints:

- `internal/ui/dialogs.go` implements the `app.Dialog` methods by adding a dialog window to the WM, sending `tea.ResumeMsg{}`, and blocking on `select { <-ctx.Done(); <-<dialog>.Done() }`. `Error`/`Warning` pass `context.Background()`.
- `ui.Model` owns a `quit chan struct{}` that `stop()` closes on exit, but only `mockInfo` selects on it.
- `mainform.New` calls `left.Focus()` and `right.Focus()`; `panel.Focus/Blur` only toggle the bubbles table cursor, and nothing wires WM window focus to panel focus.
- `config.Styles.ButtonBarStyle` (`.button-bar`) is loaded but never rendered.
- `buttonbar.View()` clears its `pressed` state, and `colorize()` hand-copies a subset of lipgloss attributes.
- `App.Navigate` uses a `select`/`default` that silently drops requests when the buffered channel is full.

See proposal.md - Why for motivation.

## Goals / Non-Goals

**Goals:**

- Every dialog call returns on quit, regardless of the context passed, with no `Dialog` interface change.
- Only the focused panel renders a cursor; cursor follows WM focus.
- The bar's themed background is actually visible; the pressed flash is side-effect-free in `View()`.
- Navigation is reliable under rapid input and mock dialogs do not accumulate.

**Non-Goals:**

- No change to the `app.Dialog`/`app.Panel` interface signatures (except where an internal-only caller change is unavoidable).
- No new connectors, no new themes, no behavior changes to `wm-example`.
- No rework of the overall message-driven architecture.

## Decisions

### D1. Quit safety via the existing `m.quit` channel (not a context change)

Add `case <-m.quit:` as an additional branch in every blocking `select` in `dialogs.go` (`Info`, `runInput`, `Select`, `SelectMultiple`, `Confirm`), returning the method's zero value. This makes every dialog return when `stop()` closes `m.quit`, independent of whether the caller supplied a cancellable context, and keeps `Error`/`Warning` signatureless.

- *Alternative considered*: thread `context.Context` through `Error`/`Warning` (breaking interface change) and require callers to pass the app ctx. Rejected: `m.quit` already exists for exactly this purpose; the interface stays stable and no caller discipline is required.

### D2. Mirror WM focus to panel cursor via an explicit sync

Remove `right.Focus()` from `mainform.New`. Add `mainform.syncPanelFocus()` that sets `panel.Focus()` for the focused panel window and `panel.Blur()` for the other. Call it after construction and after `m.wm.Update(msg)` (which can change focus via mouse clicks and Tab). `panel.Focus/Blur` are idempotent and cheap, so per-Update sync is acceptable.

- *Alternative considered*: push focus into the WM (`Content` implementing a focus interface). Rejected: `tea.Model` has no focus contract; keeps the change localized to mainform.

### D3. Paint the bar background without restructuring the renderer

Render the strip's background from `ButtonBarStyle` while keeping per-button foreground/bold from `Button/Active/Pressed` styles. The segment loop already guarantees exact total width, so the background is painted uniformly across the full row by making the segment background come from `ButtonBarStyle`. Normal buttons keep the strip's own foreground (cyan) rather than `.button`'s black, which would be unreadable on the strip; only the pressed/active foregrounds are overlaid.

- *Alternative considered*: build a full-width `ButtonBarStyle` base string and overlay button labels via `lipgloss.Place`. Rejected: ANSI overlay is fiddly; the existing single-pass segment loop already handles width, so overriding the segment background is simpler and preserves current behavior.

### D4. Replace hand-rolled `colorize` with lipgloss inheritance

Use lipgloss's `Style.Copy()`/`Inherit` to derive the per-segment style (button style over `ButtonBarStyle` background) instead of re-selecting Foreground/Background/Bold/Underline by hand. This preserves attributes the theme may add (padding, italic) without new switch arms.

- *Alternative considered*: keep `colorize` and extend the copied field list. Rejected: silently drops future attributes and duplicates lipgloss.

### D5. Single mock dialog, updated in place

Store the open mock dialog's `*dialogs.Info` and window id on `ui.Model`. On activation, if one is open, update its text in place (send a re-render); otherwise create it. Clear the reference when the mock goroutine finishes (quit or `Done()`).

- *Alternative considered*: guard with a `sync.Mutex`/flag and drop duplicate activations. Rejected: updating in place is friendlier (the latest action is always reported) and matches the "reports its name" intent.

### D6. Latest-wins navigation (no silent drop)

Replace the `select`/`default` drop in `App.Navigate` with a mutex-protected single-slot "latest request" plus a signal channel. `Navigate` stores the request and signals non-blockingly; `Run` drains the slot under the mutex. The newest request is never lost.

- *Alternative considered*: on full channel, drain the oldest buffered request and enqueue the new one. Rejected: still a bounded drop under extreme load; the single-slot latest-wins pattern is simpler to reason about and matches "the last request wins".

### D7. Move the pressed-flash reset out of `View()`

Clear `buttonbar.pressed` at the start of `Update` (before message handling) instead of at the end of `View()`. The flash persists from the activating frame until the next message, which is the intended "brief" feedback, and `View()` becomes side-effect-free (so F12 dumps and re-renders cannot eat the flash).

### D8. Consistent, case-insensitive attribute matching

Add a shared `attrValue(entry, name)` helper doing case-insensitive lookup, and use it in `entryName`, `entryIsDir`, and `parentRow`. This removes the exact-case vs case-insensitive mismatch between `panel.go` helpers.

### D9. Harden `Action` accessors and trim panel complexity

`Action.Label()`/`Name()` return "" for out-of-range values instead of panicking. In `panel.go`, extract listing→columns/rows construction into small helpers (`columnsFor`, `displayRowsFor`) and keep the `colTitles` fallback with a comment explaining the empty-directory-with-no-prior-listing case (not speculative — it is reachable via `SetData(dir, nil)`).

### D10. Stop double-delivering `DataMsg`

After `mainform.Update` applies `panel.DataMsg`, return early instead of falling through to `m.wm.Update(msg)`, which currently broadcasts the message to every window (ignored, but wasteful and confusing).

## Risks / Trade-offs

- **`m.quit` is closed only by `stop()` after `Run()` returns** → a dialog opened and then left open while the program is *not* quitting still blocks until closed, exactly as today. No change.
- **`syncPanelFocus` runs on every `Update`** → a tiny per-frame cost for idempotent `Focus`/`Blur` calls. Acceptable for two panels.
- **Segment background override changes bar appearance** → the strip now shows `.button-bar` background instead of per-button background. This is the intended fix; verify with `make run` and the `styles.css` override path.
- **Latest-wins navigation coalesces intermediate directories** → rapid Enter presses may skip intermediate listings. This is the desired "last request wins" behavior and matches the spec.

## Migration Plan

No migration. All changes are internal; apply in dependency order (config → buttonbar → panel → mainform → ui/dialogs → app), keeping the tree buildable, then add regression tests and run `go test -race ./...`, `go vet ./...`, and `make lint`.

## Open Questions

None.
