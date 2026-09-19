## Context

See proposal.md — Why. The panel holds the delivered `AttributeList` rows in `m.rows` and builds display rows from them; the connector returns entries in read order. F2 is intercepted by the button bar (`mainform.handleKey`) and delivered as an `ActivateMsg` to the ui, where `handleBarActivation` routes Quit to quit and every other action to `showMock`. The panel's columns come from the header (`colSpecs`), and hidden attributes (IsDir) are excluded from display.

## Goals / Non-Goals

**Goals:**
- The Menu button cycles the focused panel's sort mode; per-panel state.
- Directories above files, within-group ordering, and the `..` row first.
- Numeric size ordering, a header indicator, and the selection surviving a re-sort.

**Non-Goals:**
- Changing navigation, theming, the columns, or the other button actions.
- Re-reading the directory on a sort (the panel re-sorts the rows it holds).

## Decisions

**1. Sort state on the panel.**
- `panel.Model` gains `sortCol int` and `sortAsc bool` (default Name ascending). `CycleSort` advances Name↑, Name↓, Size↑, Size↓, Date↑, Date↓, Time↑, Time↓ and wraps.

**2. Routing the Menu action.**
- The `2Menu` action reaches the ui as an `ActivateMsg`; `handleBarActivation` calls the focused panel's `CycleSort`. The mainform gains `FocusedPanel()`, returning the focused window whose content is a `*panel.Model`.

**3. Value comparison.**
- The panel compares the active column's values with a helper: two `app.Size` values compare by bytes; otherwise values compare as case-insensitive strings. The Date and Time strings (`YYYY-MM-DD`, `HH:MM:SS`) sort correctly as text.

**4. Sortable size.**
- Add `app.Size` (an `int64` whose `String()` renders the human-readable form, moved from the connector's `humanSize`). The connector returns `app.Size` for the Size column; directories keep the `<DIR>` string, so the directories group compares equal by size and files compare numerically.

**5. Grouping.**
- Sort with a stable sort that first partitions by the hidden IsDir attribute (directories before files) and then orders each group by the active column and direction. The synthetic `..` row is separate and stays first.

**6. Indicator.**
- `applyColumnWidths` appends `▲` or `▼` to the active column's header title, so the marker refreshes whenever the sort mode or the columns change.

**7. Selection across a re-sort.**
- Before re-sorting, record the selected entry's name; after sorting, reuse the pending-select mechanism to place the cursor on that entry.

## Risks / Trade-offs

- [app.Size coupling] → the panel's comparator knows `app.Size`; acceptable because the panel already imports `app`.
- [Directories under a size sort] → every directory shows `<DIR>`, so a stable sort keeps the directories group in its previous relative order; this matches the grouped design.
- [Indicator width] → the marker is part of the title and is truncated like any header cell; the Name column is wide enough in practice.
- [Empty or unsized panel] → `CycleSort` does nothing when there are no columns.

## Migration Plan

None. Rollback removes the sort state and the Menu handling.

## Open Questions

None.
