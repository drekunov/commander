## Context

See proposal.md — Why. Current state that shapes the approach:

- `ncTableStyles()` (`internal/ui/widgets/panel/panel.go`) hardcodes the cursor as `Bold(true)` with black foreground on `#C0C0C0`; it is not themeable.
- `config.Styles` carries `ButtonStyle`, `ActiveButtonStyle`, `PressedButtonStyle`, `DialogBoxStyle`, `TextStyle`, `InputStyle`, `TableStyle`, and `ButtonBarStyle`, but has no cursor or menu-label style.
- `ApplyStyle` (`internal/config/css.go`) supports `color`, `background-color`, `text-decoration: underline`, `padding`, `margin`, `border-style`, and `border-color` — but not `font-weight`.
- `buttonbar.styleFor` renders each full-width segment with `ButtonBarStyle`, or an overlay of the strip with the active/pressed button style; labels and strip are the same pixels.
- `.button` is shared by the generic button widget (`internal/ui/widgets/button/button.go`) and the confirm dialog (`internal/ui/widgets/dialogs/confirm.go`), so it must not change.

## Goals / Non-Goals

**Goals:**
- The cursor's colors and weight come entirely from an injected `.cursor` style.
- Default bottom-menu labels use an injected `.menu-label` style equal to `.cursor`, drawn on the existing bar strip.
- Keep the cursor a full-width bar, keep the blurred-panel no-cursor behavior, and keep focused/pressed buttons distinct.
- Keep the CSS engine's supported-property set coherent by adding `font-weight: bold`.

**Non-Goals:**
- Changing the shared `.button` class, the confirm dialog, or the generic button widget.
- Changing bar layout/width math, action handling, focus behavior, or the strip style.
- Supporting CSS properties beyond `font-weight`.

## Decisions

**1. Add `font-weight: bold` support to `ApplyStyle`.**
- Chosen: recognize `font-weight` with value `bold` and apply `Bold(true)`.
- Why: the cursor is currently bold; expressing bold in CSS keeps `.cursor` and `.menu-label` truly equal and fully themeable. Alternatives: dropping bold (changes the cursor's look) — rejected.

**2. New `.cursor` class + `CursorStyle`; the panel uses it.**
- Chosen: add `CursorStyle: ApplyStyle(stylesheet[".cursor"])` and set the table's `Selected` from `m.styles.CursorStyle` in `tableStyles()`, keeping the existing focus branch (blur clears `Selected`) and the full-width right padding.
- Why: makes the cursor themeable and matches the equal-style requirement. Alternative: keep it hardcoded — rejected.

**3. New `.menu-label` class + `MenuLabelStyle`, set equal to `.cursor` in the default stylesheet; the bar draws default label text with it on the strip.**
- Chosen: render the label text with the state's text style and fill the rest of the segment with `ButtonBarStyle`. For the default (non-focused, non-pressed) state the text style is `MenuLabelStyle`; for focused/pressed it is the existing `overlay(ButtonBarStyle, Active/Pressed)`.
- Why: keeps the strip and focused/pressed states intact while letting the default labels match the cursor. Alternatives: changing the shared `.button` (restyles dialogs/generic buttons) — rejected; styling the whole segment (changes the strip) — rejected.

**4. Width safety in the bar.**
- The label is truncated so the rendered label width plus strip padding fills exactly `segWidth` cells. The rendered label width is measured with `lipgloss.Width` so a theme override that adds padding/margin cannot overflow the segment.

## Risks / Trade-offs

- [Width math with themed padding/margin on the label style] → measure with `lipgloss.Width`, truncate the label to fit, and fill the remainder with the strip style.
- [No label class present in a custom theme] → `MenuLabelStyle` is empty, the label renders unstyled and the segment is filled by the strip, matching today's strip look.
- [Bold support changing other classes] → only classes that declare `font-weight: bold` are affected; no existing class declares it, so current rendering is unchanged.
- [Cursor blur behavior] → the existing blurred branch clears `Selected`, so an unfocused panel still shows no cursor.

## Migration Plan

None — theme/rendering change. Rollback is reverting the stylesheet classes and their use.

## Open Questions

None.
