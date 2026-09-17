## Context

See proposal.md — Why. The panel emits a `NavigateMsg` and the app asynchronously reads the directory and calls `panel.SetData(dir, data)`, which today always does `SetCursor(0)` (`internal/ui/widgets/panel/panel.go:158`). The panel's `m.dir` still names the child directory at the moment the ascent request is emitted, but `SetData` overwrites it with the parent before the cursor is chosen. The connector and app need no change.

## Goals / Non-Goals

**Goals:**
- Restore the cursor to the child entry on both ascent paths (Backspace and Enter on `..`) and at every level.
- Keep the fallback simple and deterministic when the child is not present.

**Non-Goals:**
- No cursor memory for descending, panel focus changes, or refreshes (the existing "reset to first row" behavior stays for those).
- No app- or connector-level state; the selection stays per panel.

**Decisions**

**Record the child name at ascent time, consume it on the next delivery.**
- Chosen: `Model` gains a `pendingSelect string`. The ascent path (`ascendCmd`, reached by both Backspace and Enter on the `..` row) sets `pendingSelect = filepath.Base(m.dir)` before returning the `NavigateMsg`. `SetData` computes the cursor row by matching `pendingSelect` against the delivered entries' `Name` attribute (offsetting by one when the synthetic `..` row is shown), then clears `pendingSelect`.
- Why: the child name is available at ascent time and survives until the listing arrives, so no app-layer or message-shape change is needed. Matching by name is robust to listing reordering and to the `..` row's position.
- Alternatives considered:
  - *Remember the cursor index and restore it in the parent*: the index refers to the child's listing, not the parent's, so it is meaningless after ascent; rejected.
  - *Have the app return the previous directory with the listing*: changes the `SetData`/`Panel` contract and the `DataMsg` shape for a purely presentational concern; rejected.
  - *Store the full child path*: `filepath.Base` is the entry name in the parent and avoids carrying a path the panel must re-split; rejected as unnecessary.

**Clear the pending name on a descend request.** Entering a directory emits a descend `NavigateMsg`; clearing `pendingSelect` there prevents a stale ascent name from leaking into the new listing if requests interleave.

**Fallback.** If no delivered entry matches `pendingSelect`, the cursor stays on row 0, preserving the current behavior for refreshes and for parents that no longer contain the child.

## Risks / Trade-offs

- [Stale pending name consumed by an unrelated listing if requests interleave] → cleared on every `SetData` and on descend; a non-match falls back to row 0.
- [Name collisions or renamed entries between the two listings] → a missing match is a harmless fallback to the first row.

## Migration Plan

None — behavior-only change, no data or interface migration. Rollback is reverting the panel change.

## Open Questions

None.
