## Why

The sort window and the rest of the dialog set are only partially theme-driven: dialog body and option text are rendered without the dialog background, the sort menu reuses the file panel's cursor style, and the window title bar and resize grip use hardcoded colors. A `styles.css` override therefore cannot fully restyle the dialogs, and the embedded Norton Commander Classic scheme looks inconsistent across them.

## What Changes

- Make the whole dialog set — Info, Input, Confirm, and the Select-based sort window — render from injected CSS styles, with the dialog background carried by every body line, not only the padding.
- Give the sort window a dedicated themeable option style for its selected and unselected rows, separate from the file panel's cursor.
- Draw the window frame chrome — title (focused and unfocused) and resize grip — from injected CSS styles instead of hardcoded colors.
- Extend the embedded `styles.css` (Norton Commander Classic) with the new classes and values so the dialog set and frame chrome match the scheme, and remain overridable per class.
- Add regression tests that a working-directory `styles.css` override restyles the dialogs, the sort menu rows, and the frame chrome.

## Capabilities

### New Capabilities

- `ui/window-frame`: theming of the shared window frame chrome (border, background, focused/unfocused title, resize grip) from injected styles.

### Modified Capabilities

- `config/theming`: the theme SHALL expose style classes covering dialog content, dialog menu options, and window frame chrome, so a working-directory override can restyle every part of a dialog.
- `ui/dialogs`: the dialog set and the sort window SHALL render their body, options (selected and unselected), and frame chrome from injected styles, including a menu option style distinct from the file panel cursor.

## Impact

- **Code**: `internal/config/config.go`, `internal/config/styles.css`, `internal/ui/widgets/dialogs/*.go`, `internal/ui/widgets/wm/renderer.go`.
- **Tests**: new/updated tests in `internal/config/config_test.go`, `internal/ui/widgets/dialogs/*_test.go`, and `internal/ui/widgets/wm/renderer_test.go`.
- **UI API**: internal-only; `config.Styles` gains fields, no `app` interface change.
- **Dependencies**: none added.
