## Context

See proposal.md — Why. `buttonbar.View` renders the numeric prefix, the name, and the strip padding that fills the rest of each segment. The padding uses `ButtonBarStyle`; `numberStyle()` and `styleFor()` inherit from `ButtonBarStyle`; the focused and pressed labels use `overlay(ButtonBarStyle, Active/Pressed)`. `config.Styles` carries `ButtonBarStyle` but has no dedicated bar-background style.

## Goals / Non-Goals

**Goals:**
- A dedicated menu-background style for the bar background, defaulting to the menu-number look (black background, light gray foreground, bold).
- Use it as the base for the padding and for the number, name, and focused/pressed labels, so the background applies in every state.

**Non-Goals:**
- Changing the label split, the number and name styles, the segment width, action handling, or focus behavior.
- Removing `.button-bar`.

## Decisions

**1. Add `.menu-background` + `MenuBackgroundStyle`, defaulting to the menu-number look.**
- Chosen: `MenuBackgroundStyle: ApplyStyle(stylesheet[".menu-background"])` with `color: #C0C0C0`, `background-color: #000000`, `font-weight: bold`.
- Why: gives the bar background its own themeable style whose default matches the number.

**2. Base the bar background on the menu-background style with a button-bar fallback.**
- Chosen: `barBackground() = MenuBackgroundStyle.Inherit(ButtonBarStyle)`.
- Why: a theme that omits `.menu-background` keeps today's strip instead of losing the background. Alternative: use `MenuBackgroundStyle` directly — rejected as less robust.

**3. Apply the base everywhere.**
- Use `barBackground()` for the padding, as the base in `numberStyle()`, as the default base in `styleFor()`, and as the base for the active/pressed overlays. This makes the new background apply to all states.

## Risks / Trade-offs

- [Theme omits `.menu-background`] → it inherits the button-bar strip, matching today's look.
- [Number background equals the base] → the number's own black background matches the new base, so the number reads as light-gray text on the bar while the name stays a light-gray chip — the intended split look.
- [Focused/pressed] → their foreground and underline stay distinct on the new base.

## Migration Plan

None — theme/rendering change. Rollback is reverting the class and its use.

## Open Questions

None.
