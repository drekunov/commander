## Why

The commander has a bottom function-key bar but no top menu line, so the classic Norton Commander interaction of pulling down menus with F9 is missing. Adding a top menu bar gives the application its familiar menu-driven navigation and a single visible home for grouped commands.

## What Changes

- Add a top menu bar rendered as the first screen row, showing the classic Norton Commander top-level menus: Left, Files, Commands, Options, Right.
- Each top-level menu caption highlights one hotkey letter (L, F, C, O, R); pressing that letter while the menu is active selects that menu, and pressing the caption hotkey opens its pull-down.
- F9 activates the top menu bar; while active, Left/Right move between top-level menus, Up/Down move within an open pull-down, Enter opens or activates, and Escape closes the pull-down and then deactivates the bar.
- Each top-level menu opens a pull-down list of items; selecting an item reuses the existing function-button action dispatch (Menu opens the sort window, Quit quits, and the remaining actions show the existing mock report) so no second action mechanism is introduced.
- Give the top menu bar its own CSS classes for the bar background, an inactive caption, an active/highlighted caption, the hotkey letter, the pull-down background, and a focused pull-down item, all overridable from a working-directory `styles.css`.
- Keep the bottom function-button bar working; the 9PullDn button and F9 now activate the top menu instead of reporting a mock.

## Capabilities

### New Capabilities

- `ui/top-menu`: the top menu bar — classic NC top-level menus, F9 activation, hotkey-letter selection, pull-down navigation, and reuse of the function-button action dispatch.

### Modified Capabilities

- `config/theming`: the theme SHALL expose style classes for the top menu bar, its captions (inactive and active), hotkey letters, pull-down background, and focused pull-down item, so a working-directory override can restyle the menu.
- `ui/function-button-bar`: the 9PullDn button and the F9 key SHALL activate the top menu instead of running a mock action.

## Impact

- **Code**: new `internal/ui/widgets/topmenu` widget; `internal/ui/widgets/mainform/mainwindow.go` (top row layout and key routing); `internal/ui/ui.go` (route the Menu action to the top menu, keep action dispatch); `internal/config/config.go` and `internal/config/styles.css` (new classes); `internal/ui/widgets/buttonbar/buttonbar.go` (9PullDn semantics unchanged as an action, behavior now owned by the ui).
- **Tests**: new `internal/ui/widgets/topmenu/*_test.go`; updated `internal/ui/widgets/mainform/mainwindow_test.go`, the `internal/ui` test files, and `internal/config/config_test.go`.
- **UI API**: internal-only; `config.Styles` gains fields, the main form gains a top row, no `app` interface change.
- **Dependencies**: none added.
