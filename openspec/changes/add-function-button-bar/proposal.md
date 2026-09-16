## Why

The commander is a dual-panel file manager, but it has no command affordances: all actions live behind undocumented keys (F10 quits, F12 dumps the screen). A Norton-Commander-style bottom function-button bar gives users a visible, discoverable row of actions and establishes the UI scaffolding and action-wiring seam that real file operations will attach to later.

## What Changes

- Add a bottom function-button bar, rendered as a fixed row under the two file panels, mirroring the classic Norton Commander F1-F10 labels: `1Help 2Menu 3View 4Edit 5Copy 6RenMov 7Mkdir 8Delete 9PullDn 10Quit`.
- Each button is rendered from the existing themed `button` widget so it shows focus/pressed state consistent with the rest of the UI.
- A button is activated either by its matching function key (F1-F10) or by a mouse click on the button; the bar supports keyboard focus navigation (left/right) so `Enter` activates the focused button too.
- **Nine buttons are mocks for now**: activating any of them (1Help-9PullDn) invokes a stub action handler that reports a placeholder result (e.g. an info dialog naming the action) instead of performing a real file operation. The bar exposes a typed action-handler contract so real commands can replace the mocks without reworking the widget.
- **10Quit stays real**: activating it (F10) quits the application exactly as F10 does today — it is the app's only in-app exit path and must keep working; it is the one pre-wired real action.
- The layout reserves a one-row-high bar across the full terminal width, below the two panels, and keeps it visible regardless of window/focus state.

## Capabilities

### New Capabilities
- `ui/function-button-bar`: The bottom function-key button bar — its rendered layout, keyboard/mouse activation, and the mocked action handlers behind each button.

### Modified Capabilities
<!-- None: existing capabilities keep their behavior. -->

## Impact

- New widget package under `internal/ui/widgets/` (e.g. `buttonbar`) composed of the existing `button` widget.
- `internal/ui/widgets/mainform` reserves the bottom bar row in `layoutPanels` and forwards F1-F10 / mouse messages to the bar; panel heights shrink by one row.
- `internal/ui/widgets/wm` hit-testing extended so clicks on the bar are delivered to the bar rather than a window.
- `internal/config` styles extended with a button-bar style (bottom bar, active/inactive keys) driven by the CSS theme engine; no global style state.
- Mock action handlers live behind a small handler interface so the `app`/`ui` layer can inject real actions later. No `app.Connector` or panel data-contract changes.
