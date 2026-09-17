## Context

See proposal.md — Why. `panel.Model` builds its table styles in two places: `ncTableStyles()` overrides only `Selected`, and `tableStyles()` adds the focus-aware `Selected` variant. `View()` then wraps `tableView.View()` in `m.styles.TableStyle` (the injected `.table` style, background `#000080`). Bubbles' default `Header` style is `Bold(true).Padding(0, 1)` with no background; its bold render emits an SGR reset that clears the background inherited from the outer `.table` wrap, so header cells fall back to the terminal default while body cells keep `.table`'s background.

## Goals / Non-Goals

**Goals:**
- The column header row paints the same background as the table body, focused or unfocused.
- The header background comes from the injected table style, so a `styles.css` override stays consistent across header and body.
- Minimal, localized change with no dependency changes.

**Non-Goals:**
- Changing header text, weight, padding, column widths, or introducing a distinct header color.
- Changing the cursor bar, row rendering, navigation, the connector, or the data contract.
- Reworking the theming system or adding new CSS classes.

## Decisions

**Decision: give the header style the injected table background where render styles are built.**
- Chosen: in the panel's render-style builder, set the table styles' `Header` background from `m.styles.TableStyle.GetBackground()`, leaving the header's bold weight and padding intact. Apply it in the single builder that `View()` uses, so focused and unfocused panels both get it.
- Why: it targets the exact cause — the header's own reset clears the inherited background — and keeps the injected theme as the single source of truth without new config surface.
- Alternatives considered:
  - *Add a `.table-header` CSS class and a `TableHeaderStyle` field.* More flexible, but adds config surface for one region and lets a theme author desynchronize header and body backgrounds.
  - *Suppress the header's SGR reset (render header cells unstyled).* Relies on the outer `.table` background bleeding through; fragile, and any future header styling reintroduces the gap.
  - *Set the header background once in `NewPanel`.* Insufficient: `View()` rebuilds styles from `ncTableStyles()` every frame, so the assignment must live where render styles are produced.

## Risks / Trade-offs

- [Header foreground contrast] → the header keeps its existing default foreground and bold weight; only the background changes, and the body already uses that foreground on that background.
- [Theme with no table background] → `GetBackground()` returns no color, so the assignment is a no-op and header and body both remain unbackgrounded and consistent.
- [Terminal color approximation] → lipgloss maps the theme color to the terminal palette for header and body identically, so the two stay matched.

## Migration Plan

None — behavior-only rendering fix. Rollback is removing the header background assignment.

## Open Questions

None.
