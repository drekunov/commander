## Why

The panel's selection cursor is hardcoded in Go (`internal/ui/widgets/panel/panel.go`) and cannot be themed, while the bottom function-key menu draws its labels with the bar strip style. The cursor and the menu labels therefore look unrelated, and neither can be controlled from `styles.css`. Give both a CSS-driven style and set them equal so the cursor and the menu labels share one look.

## What Changes

- Add a `.cursor` CSS class (black text on `#C0C0C0`, bold) and a `CursorStyle` in the injected styles; the panel's selection cursor uses it instead of the built-in style.
- Add a dedicated bottom-menu label class `.menu-label`, set equal to `.cursor`, and a `MenuLabelStyle`; the function-button bar renders each non-focused, non-pressed label's text with it on the existing `.button-bar` strip.
- Extend the CSS engine to support `font-weight: bold` so a theme can express bold, which is needed to keep the cursor's current bold look and keep the cursor and the menu labels equal.
- Keep the focused (`.active-button`) and pressed (`.pressed-button`) button states distinct.
- Leave the shared `.button` class unchanged, since it also styles the confirm dialog and the generic button widget.

## Capabilities

### New Capabilities
<!-- None: this themes existing regions and extends the existing theming engine. -->

### Modified Capabilities
- `config/theming`: the CSS engine gains `font-weight: bold` support, so a theme (including a working-directory override) can express bold text.
- `panel/rendering`: the panel cursor takes its appearance from the injected `.cursor` style rather than a built-in style.
- `ui/function-button-bar`: non-focused, non-pressed button labels are drawn with the injected `.menu-label` style on the bar strip.

## Impact

- `internal/config/styles.css`, `internal/config/css.go`, `internal/config/config.go`: new `.cursor` and `.menu-label` classes, `font-weight` support, and `CursorStyle`/`MenuLabelStyle` fields.
- `internal/ui/widgets/panel/panel.go`: the selection cursor uses the injected `CursorStyle`.
- `internal/ui/widgets/buttonbar/buttonbar.go`: default label text uses the injected `MenuLabelStyle` on the bar strip.
- Tests in the config, panel, and buttonbar packages.
- Assumption: the bottom-menu label is styled by a new dedicated class (not the shared `.button`) so dialogs and the generic button widget are unaffected.
