## Why

The bottom function-key bar lays its ten buttons out by splitting the terminal width into ten segments, so a width that is not divisible by ten makes the leftmost buttons one cell wider than the rest. Within each segment the number and name chips have different widths per label, so the visible buttons are uneven. Give every button the same width: a two-cell number field plus a six-cell name field, with the leftover columns split into equal gaps.

## What Changes

- Lay the ten buttons out with equal width: each button is a two-cell function-key number field followed by a six-cell name field (eight cells total).
- Split the columns left over after the ten buttons into equal gaps between adjacent buttons, so the row still spans the full terminal width.
- When the terminal is too narrow for ten eight-cell buttons, shrink the buttons equally and truncate labels so the bar still fits on one row.
- Mouse hit regions follow the new layout: a click maps to the button whose span (button plus the gap after it) contains the column.
- No change to the number, name, or background styles; the label split; action handling; or focus behavior.

## Capabilities

### New Capabilities
<!-- None: this refines the layout of an existing bar. -->

### Modified Capabilities
- `ui/function-button-bar`: adds a requirement that the ten buttons have equal width, and refines the existing button and terminal-width requirements accordingly.

## Impact

- `internal/ui/widgets/buttonbar/buttonbar.go`: replace the per-segment width with the fixed two-cell number plus six-cell name layout and equal inter-button gaps; update `actionIndexAtColumn` for the new hit regions.
- Tests in the buttonbar package.
- Assumption: the gaps are the nine gaps between adjacent buttons (no leading or trailing gap); a leftover that does not divide evenly is spread over the leftmost gaps. The number field is right-aligned (left-padded) so names line up.
