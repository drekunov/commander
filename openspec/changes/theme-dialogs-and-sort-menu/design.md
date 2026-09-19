## Context

See proposal.md - Why. Existing constraints that shape the approach:

- `config.LoadStyles` maps CSS selectors to fields on `config.Styles`; widgets receive the struct by injection and must not read global state.
- The window frame is shared: the `wm` renderer draws the frame for both panel windows and dialog windows, and already takes its border/background from `DialogBoxStyle` (`.dialog-box`). Title colors (focused `#FFFFFF`, unfocused border color) and the grip color are hardcoded in `renderer.go`.
- The Select dialog (used as the sort window) styles its focused row with `CursorStyle` (`.cursor`), which is also the file panel's cursor, so the two cannot be themed apart.
- Dialog body text uses `.text`/`.input` with a foreground but no background, so the frame's blue background is only painted on the renderer's padding, not behind the text itself.
- `select_test.go` asserts the Select cursor uses the panel cursor style, and `renderer_test.go` asserts the title bar contains the title; both will need adjustment.

## Goals / Non-Goals

**Goals:**

- Every visible part of a dialog (body text, input, options, buttons, frame, title, grip) is driven by an injected CSS class.
- The sort window's option rows are themeable independently of the file panel cursor.
- The embedded `styles.css` expresses the Norton Commander Classic scheme for all of it, and a working-directory override can retheme any class.

**Non-Goals:**

- No selectable/multiple named themes; one embedded theme plus the working-dir override.
- No change to the window manager's geometry, focus, or drag/resize behavior.
- No change to the file panel's own cursor class or its selection behavior.
- No new CSS properties beyond what `ApplyStyle` already supports.

## Decisions

### D1. Add explicit style fields and CSS classes

Add to `config.Styles`, each loaded from a selector:

| Field | Selector | Purpose |
|---|---|---|
| `DialogCursorStyle` | `.dialog-cursor` | Sort-window / Select focused option row |
| `DialogOptionStyle` | `.dialog-option` | Sort-window / Select unfocused option rows |
| `WindowTitleStyle` | `.window-title` | Focused window title |
| `WindowTitleUnfocusedStyle` | `.window-title-unfocused` | Unfocused window title |
| `WindowGripStyle` | `.window-grip` | Bottom-border resize grip |

`DialogBoxStyle` (`.dialog-box`) keeps drawing the frame border and background for every window, so existing overrides that only set `.dialog-box` keep working.

- *Alternative considered*: reuse `.cursor`/`.text` and rely on the theme author to differentiate. Rejected: it cannot express a sort menu distinct from the panel cursor, which is the point of the request.
- *Alternative considered*: a generic `.window` prefix for the frame style. Rejected: renaming `.dialog-box` would break existing `styles.css` overrides and the current tests.

### D2. Fallback to the current look when a new class is absent

`config.LoadStyles` populates the new fields with `ApplyStyle(stylesheet[selector])`. Widgets compose each new style over an existing base so an unset selector degrades to today's appearance:

- `DialogCursorStyle` falls back to `CursorStyle`, `DialogOptionStyle` falls back to the dialog body style.
- `WindowTitleStyle`/`WindowTitleUnfocusedStyle`/`WindowGripStyle` fall back to the frame border foreground (and the existing focus distinction when the focused title class is unset).

This keeps third-party overrides written before this change visually unchanged.

### D3. Dialog body inherits the frame background

Dialog widgets render each body line through the injected styles composed over the `DialogBoxStyle` background, so text, input, options, and buttons carry the frame background instead of the terminal default. The renderer continues padding the remaining cells with the same background. The background is inherited rather than required in each class, so a theme can still override any single class.

- *Alternative considered*: require `background-color` on every dialog class in `styles.css`. Rejected: a partial override would leave gaps, and the requirement is that content is continuous with the frame regardless of theme.

### D4. Frame chrome reads the injected styles

`wm.renderer.buildTitleBar` and `buildResizeGrip` use `WindowTitleStyle` (focused), `WindowTitleUnfocusedStyle` (unfocused), and `WindowGripStyle`, each composed over the frame background, replacing the hardcoded colors. The frame border and background keep using `DialogBoxStyle`.

- *Alternative considered*: keep hardcoded title/grip and expose only colors. Rejected: the request is that the dialog set be set by CSS, and a theme must be able to change weight/decoration, not just color.

### D5. Norton Commander Classic values in the embedded stylesheet

Update `internal/config/styles.css` with the new classes in the NC Classic palette: dialog background blue, focused option cyan-on-black, focused title white bold, unfocused title cyan, grip cyan. Exact values live in the stylesheet, not in the spec.

## Risks / Trade-offs

- **`select_test.go` asserts the old `.cursor` behavior** → update it to assert `.dialog-cursor` for the focused row and `.dialog-option` for the rest.
- **`renderer_test.go` may assume hardcoded title colors** → update to assert the injected focused/unfocused title styles differ and come from the theme.
- **Background inheritance on already-rendered ANSI content** → compose the style before rendering each span (as `buttonbar` already does with `Inherit`), not by wrapping rendered strings.
- **New fields left unset in a user override** → covered by D2 fallbacks.

## Migration Plan

No migration. Apply in order: `config` (fields + stylesheet) → `dialogs` (option/body styles) → `wm` (frame chrome), keeping the tree buildable; then update tests and run `go build ./...`, `go vet ./...`, `go test -race ./...`, `make lint`.

## Open Questions

None.
