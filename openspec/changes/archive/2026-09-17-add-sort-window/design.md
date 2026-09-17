## Context

See proposal.md — Why. The Menu action is handled in `ui.handleBarActivation`, which currently calls the focused panel's `CycleSort` on the event loop. The dialog infrastructure (`ui.dialogs.go`) opens windows via `addDialogWindow` (the window manager is mutex-guarded) and blocks on the dialog's `Done()` channel, so callers run on a non-event-loop goroutine. `dialogs.Select` is a single-select list that returns the focused option on Enter and has no cancel.

The mainform's `handleKey` intercepts F10 (quit) and the function keys before the window manager sees them, and its `handleMouse` intercepts the bar row. A dialog window therefore does not currently block function keys or clicks on other windows, so it is only partially modal.

## Goals / Non-Goals

**Goals:**
- The Menu button opens a modal column list; choosing a column sorts by it (inverting the current column); Escape cancels.
- Reuse the existing Select dialog.

**Non-Goals:**
- Changing the grouping, numeric size, header marker, or selection survival.
- Changing any other button.

## Decisions

**1. Add Escape to `dialogs.Select`.**
- On Escape, mark the dialog canceled, hide it, and signal `Done()`; expose `Canceled()`. `ui.Select` returns no option (an empty string) when canceled.
- Alternative: a dedicated sort dialog — rejected in favor of reusing Select.

**2. Panel API.**
- Add `ColumnTitles()` (the visible titles without the header marker) and `SortBy(title)` (find the column; invert when it is the current column, otherwise select it ascending; then re-sort and reposition). Remove `CycleSort`, since the cycle is gone. `SortBy` reuses the existing `sortRows` and cursor-positioning logic.

**3. Asynchronous window.**
- The Menu action, on the event loop, spawns a goroutine that opens the Select dialog through the existing `ui.Select` (which blocks that goroutine, not the event loop) and, on a non-empty choice, sends a `sortColumnMsg` to the event loop. The ui handles the message by calling the focused panel's `SortBy`, so the panel is mutated on the event loop and delivery stays race-free.

**4. Options.**
- The window lists the panel's column titles; the header marker already shows the active column, so the window does not repeat it.

**5. Modal dialogs.**
- While any dialog window is open, the mainform SHALL stop routing F1–F9 to the button bar and SHALL ignore mouse clicks outside the dialog; F10 still quits. The mainform detects an open dialog through the window manager (the focused window's content is not a panel) or a dedicated modal flag.
- Alternative: leave function keys live — rejected because it lets a second dialog open and breaks modality.

## Risks / Trade-offs

- [Shared Select behavior] → Escape now cancels every Select, returning no option; existing callers already treat an empty result as no choice (the context/quit path).
- [Goroutine and windows] → the window manager's Add/Remove are mutex-guarded, matching the existing dialog pattern; the panel is mutated on the event loop via the message.
- [Removed cycle] → `CycleSort` is removed and the panel tests move to `SortBy`.
- [Empty columns] → the Menu action does nothing when the focused panel has no columns.
- [Modality vs. quit] → F10 is exempt so the user can always quit; dialogs already return on quit.
- [Detecting a dialog] → the mainform identifies an open dialog by the focused window's content type; a focused non-panel window is treated as modal.

## Migration Plan

None. Rollback restores `CycleSort` and the direct cycle.

## Open Questions

None.
