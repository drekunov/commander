## Why

The bottom function-key bar paints its full-width background with the `.button-bar` strip style (blue). The bar background should be a separate, independently themeable concern, and its default should match the menu-number treatment: a black background with a light gray foreground.

## What Changes

- Add a `.menu-background` CSS class and a `MenuBackgroundStyle`; it defaults to the menu-number look (black `#000000` background, light gray `#C0C0C0` foreground, bold).
- The bar paints its full-width background with the injected menu-background style, falling back to the button-bar style for any attribute the menu-background style leaves unset.
- The number, the name, and the focused and pressed labels are all drawn on the menu-background base, so the new background applies in every state.
- `.button-bar` stays as the fallback base; no other class changes.
- No change to the label split, the number and name styles, the segment width, action handling, or focus behavior.

## Capabilities

### New Capabilities
<!-- None: this refines the styling of an existing bar region. -->

### Modified Capabilities
- `ui/function-button-bar`: the "Bar renders a full-width themed background" requirement paints the background with the injected menu-background style; the "Button labels use the injected menu-label style" requirement draws the number and name on the menu-background base.

## Impact

- `internal/config/styles.css`, `internal/config/config.go`: a `.menu-background` class and a `MenuBackgroundStyle`.
- `internal/ui/widgets/buttonbar/buttonbar.go`: use the menu-background base for the padding between labels and as the inheritance base for the number, name, and active/pressed styles.
- Tests in the config and buttonbar packages.
- Assumption: `.menu-background` falls back to `.button-bar` for unset attributes, so a theme that omits it keeps today's strip.
