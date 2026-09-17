## Context

See proposal.md — Why. `buttonbar.View` walks the ten actions; `segmentWidth(idx)` returns `width / 10` plus one for the leftmost `width % 10` segments, and each segment renders the number chip, the name chip, then background padding. `actionIndexAtColumn` accumulates `segmentWidth` to map a click column to a button. The widest name is six cells (`RenMov`, `Delete`, `PullDn`) and the widest number is two (`10`).

## Goals / Non-Goals

**Goals:**
- Every button has the same width: a two-cell number field followed by a six-cell name field.
- Leftover columns become equal gaps between adjacent buttons; the row spans the full width.
- Mouse mapping follows the new layout.

**Non-Goals:**
- Changing the number, name, or background styles, the label split, action handling, or focus behavior.
- Adding new theme classes.

## Decisions

**1. Fixed button width of eight cells (number 2 + name 6).**
- Chosen because the widest name is six cells and the widest number is two.
- Alternative: derive from the longest full label (seven, e.g. `6RenMov`) — rejected in favor of the longest-name width.

**2. Right-align the number field (left-pad to two cells).**
- Keeps the names aligned and keeps the original label text (`1Help`) as a substring.
- Alternative: right-pad (`1 Help`) — rejected; it breaks the substring and looks no better.

**3. Layout algorithm.**
- `buttonWidth = min(8, width / buttonCount)`; `leftover = width - buttonCount*buttonWidth`; split `leftover` across the `buttonCount - 1 = 9` internal gaps (floor per gap, remainder over the leftmost gaps).
- When `width / 10 >= 8` the buttons are eight cells and the gaps absorb the rest; when narrower the buttons shrink to `width / 10` and the gaps absorb `width % 10`.
- Alternative: distribute over eleven gaps including the edges — rejected because the first and last buttons would not touch the edges.

**4. Hit regions.**
- A button's hit region is `[start, nextStart)`: its span plus the gap after it, so a click anywhere maps to a button and the last button's region reaches the end of the row.

**5. Truncation.**
- Each field is truncated to its allotted cells (number to two, name to six, or fewer when the button is narrower).

## Risks / Trade-offs

- [Narrow terminals] → buttons shrink to `width / 10` and labels truncate, preserving the single-row guarantee.
- [Leftover not divisible by nine] → gaps differ by at most one cell and the row still spans the full width.
- [Very wide terminals] → buttons stay eight cells and the gaps grow, which is the intended fixed-button look.
- [Existing tests] → `TestRowFillsExactWidth` still holds (the row is exactly the width) and `TestRenderTenButtonsInOrder` still holds because each label remains a substring.

## Migration Plan

None — rendering change. Rollback is restoring the segment-width layout.

## Open Questions

None.
