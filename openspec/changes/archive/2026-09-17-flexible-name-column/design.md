## Context

See proposal.md — Why. `panel.columnsFor` turns each non-hidden header attribute into a column with a fixed width (the attribute's `Width`, or `defaultColumnWidth`). `SetData` stores the columns and `contentWidth`; `View` sets the table width each frame, and `mainform` calls `panel.SetWidth` on layout. Column widths are never recomputed after `SetData`, so a wider panel leaves empty space to the right of the last column.

## Goals / Non-Goals

**Goals:**
- The flexible column fills the width left after the fixed columns and the per-cell padding.
- Fixed columns keep their widths; the flexible column shrinks to a minimum of eight cells on narrow panels.
- Recompute the widths when the panel is resized.

**Non-Goals:**
- Changing the columns themselves, the hidden flag, navigation, theming, or the function-button bar.
- Horizontal scrolling or per-column flex weights.

## Decisions

**1. `app.Attribute` gains `Flex bool`.**
- The header marks the flexible column, keeping the panel connector-agnostic. Alternative: the panel hardcodes the Name column — rejected.

**2. The panel records column specs.**
- `columnsFor` records each visible column's title, declared width, and flex flag. A helper builds the `table.Column`s from the specs and the current panel width.

**3. Flex width formula.**
- `flexWidth = panelWidth - 2*numColumns - sum(fixed widths)`, clamped to a minimum of eight cells. When no column is flexible, every column uses its declared width.

**4. Recompute on resize.**
- The width-building helper is called from `SetData` (specs are known) and from `SetWidth` (the panel was resized), and updates `contentWidth` to the resulting total so the selected-row padding still fills the remaining cells.

**5. Multiple flexible columns.**
- The first flexible column flexes; any others keep their declared widths.

## Risks / Trade-offs

- [Width before the first layout] → `panelWidth` may be 0 before the first layout; the flex width clamps to the eight-cell minimum and is corrected on the next `SetWidth`.
- [Very narrow panel] → the fixed columns may exceed the panel; the flexible column stays at eight cells and the table clips the remainder.
- [contentWidth] → recomputed from the actual widths, so the selected-row padding does not double-count the empty space.

## Migration Plan

None. Rollback restores fixed widths and ignores the flex flag.

## Open Questions

None.
