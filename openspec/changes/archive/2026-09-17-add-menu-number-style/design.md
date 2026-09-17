## Context

See proposal.md — Why. `buttonbar.View` currently renders each button's whole label (`def.label`, e.g. `1Help`) with one style from `styleFor`, then fills the rest of the segment with the strip style. `config.Styles` already carries `MenuLabelStyle` (from `.menu-label`) but has no number style, and each `actionDef` already stores both `label` (`1Help`) and `name` (`Help`).

## Goals / Non-Goals

**Goals:**
- Split each label into its numeric prefix and its name.
- The number uses an injected `.menu-number` style; the name keeps the injected `.menu-label` / active / pressed styles.
- Keep each segment's exact width and existing truncation behavior.

**Non-Goals:**
- Changing the shared `.button` class, the bar strip, the segment layout, action handling, or focus behavior.
- Adding CSS properties beyond those the engine already supports (`color`, `background-color`, `font-weight`, and the rest).

## Decisions

**1. Derive the numeric prefix from the label, not the action.**
- Chosen: `number := strings.TrimSuffix(def.label, def.name)`, so the prefix is exactly the part of the rendered label before the name.
- Why: it needs no new import and cannot drift from the label. Alternative: format the action number with `strconv` — equivalent but adds an import.

**2. Add `.menu-number` + `MenuNumberStyle`; the number always uses it.**
- Chosen: render the number with `m.styles.MenuNumberStyle.Inherit(m.styles.ButtonBarStyle)`, so a theme that omits `.menu-number` falls back to the strip style (today's look).
- Why: the number keeps its default style on focused and pressed buttons, while the name changes.

**3. Keep `styleFor` as the name style.**
- `styleFor` continues to return the name's style: `MenuLabelStyle.Inherit(ButtonBarStyle)` by default, or the active/pressed overlay. Only `View` changes to render two spans.

**4. Width safety.**
- Render the number, truncate it to `segWidth` if needed, then render the name into the remaining width (truncated to fit), and fill any remainder with the strip style. Widths are measured with `lipgloss.Width` and truncation uses `ansi.Truncate`, so themed padding/margin cannot overflow the segment and the row stays exactly the bar width.

## Risks / Trade-offs

- [Narrow segments] → the number is truncated first, then the name; the segment still fills exactly `segWidth`.
- [Theme omits `.menu-number`] → the number inherits the strip style, matching the previous whole-label strip look.
- [Number on focused/pressed buttons] → stays default by design; only the name switches to active/pressed.

## Migration Plan

None — rendering/theme change. Rollback is reverting the split rendering and the class.

## Open Questions

None.
