## Context

See proposal.md — Why. `panel.SetData` positions the cursor with `m.tableView.SetCursor(idx)` (`internal/ui/widgets/panel/panel.go`). The bundled `charmbracelet/bubbles` table's `SetCursor` calls `UpdateViewport`, which re-anchors the rendered content window around the cursor, but it never changes `viewport.YOffset`. The viewport renders from `YOffset`, so when the selected row lies past the visible window the highlight is off-screen. Reproduced with a 60-entry listing: after ascending, no row is highlighted.

## Goals / Non-Goals

**Goals:**
- Make the panel scroll so the restored/selected row is visible after any `SetData`, including tall listings.
- Keep the selected row unchanged and avoid touching arrow-key navigation (which already scrolls).

**Non-Goals:**
- No change to the bubbles dependency or vendoring it.
- No change to which entry is selected, the `..` row, or the connector/app.

## Decisions

**Position the cursor with `GotoTop()` then `MoveDown(idx)`.**
- Chosen: in `SetData`, replace `SetCursor(idx)` with `m.tableView.GotoTop()` followed by `m.tableView.MoveDown(idx)`. `GotoTop` resets both the cursor and the viewport offset to the top; `MoveDown` then moves the cursor and updates `YOffset` so the cursor stays visible. This is verified against several pre-states (scrolled listing replaced by another, scrolled listing then a small child then the tall parent) and indices (0, 2, mid, last), where `SetCursor` alone leaves the row hidden for far indices.
- Why: it uses only the table's public API, needs no dependency change, and works for the reset-to-first-row case (`MoveDown(0)`) as well as far restores.
- Alternatives considered:
  - *`SetCursor(idx)` then `MoveDown(0)`*: fails when a stale non-zero `YOffset` from a previous listing leaves a near-top cursor hidden.
  - *`SetCursor(0)` + `MoveUp(0)` + `MoveDown(idx)`*: also works but is more obscure about intent.
  - *`GotoTop()` twice then `MoveDown(idx)`*: works, but the second call is only needed to work around `MoveUp`'s offset math; a single `GotoTop` already resets when the cursor is at the top.
  - *Vendor/patch the table to expose `SetYOffset`*: rejected as dependency churn for a one-line positioning concern.

## Risks / Trade-offs

- [The fix depends on bubbles' `MoveUp`/`MoveDown` viewport heuristics] → covered by a panel test that renders a tall listing and asserts the restored row is visible, plus a pty smoke.
- [Empty listings] → `GotoTop`/`MoveDown(0)` on zero rows leave the cursor at 0 with nothing to scroll, matching the current behavior.
- [Two table calls instead of one] → negligible cost on a listing delivery.

## Migration Plan

None — behavior-only fix. Rollback is restoring the single `SetCursor` call.

## Open Questions

None.
