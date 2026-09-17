## Why

The bottom function-key bar renders each label (for example `1Help`) as a single styled span, so the function-key number is not distinguished from the name. Split the number from the name and give it its own CSS-driven style: the number defaults to a black background with light gray foreground, inverting the name's light gray background with black foreground.

## What Changes

- Split each button label into its numeric prefix and its name.
- Add a `.menu-number` CSS class and a `MenuNumberStyle`; the numeric prefix renders with it (black background, light gray `#C0C0C0` foreground, bold).
- The name keeps the injected `.menu-label` style; a focused or pressed name keeps its distinct active or pressed style.
- The numeric prefix keeps the menu-number style on focused and pressed buttons.
- No change to the shared `.button` class, the bar strip, the segment width, action handling, or the cursor.

## Capabilities

### New Capabilities
<!-- None: this refines the styling of an existing bar region. -->

### Modified Capabilities
- `ui/function-button-bar`: the "Button labels use the injected menu-label style" requirement is extended so each label is split into a numeric prefix (menu-number style) and a name (menu-label style), with the number keeping its style on focused and pressed buttons.

## Impact

- `internal/config/styles.css`, `internal/config/config.go`: a new `.menu-number` class and `MenuNumberStyle`.
- `internal/ui/widgets/buttonbar/buttonbar.go`: render the number and name as separate spans while keeping each segment's exact width.
- Tests in the config and buttonbar packages.
- Assumption: `.menu-number` is bold, matching the name's weight.
