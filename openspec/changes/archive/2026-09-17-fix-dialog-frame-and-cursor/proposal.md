## Why

A dialog window shows two nested frames: the window manager draws a window frame (title bar, border, resize grip) and the dialog widget draws its own bordered box with a second title. The inner box wastes space and looks broken, and the dialog's focused option is styled with the button style rather than the cursor style. Remove the inner box so the window frame is the single frame, size the window to the dialog content, and style the Select cursor with `.cursor`.

## What Changes

- Dialog widgets stop drawing their own border and title; the window-manager frame is the only frame, and a dialog's footer becomes a body line.
- A dialog window is sized to fit its content plus the frame, instead of half the screen.
- The Select dialog's focused option row is rendered with the injected `.cursor` style.
- The sort window's caption reads `Sort: Select a column`; the instruction is no longer repeated as a body line.
- This applies to every dialog (Info, Error, Input, Select, Confirm).
- No change to the modal behavior, the sort logic, the theme classes, or the panel frames.

## Capabilities

### New Capabilities
<!-- None: this fixes how existing dialogs render. -->

### Modified Capabilities
- `ui/dialogs`: a dialog renders inside a single content-sized window frame with no inner box, the Select dialog's cursor uses the injected cursor style, and the sort window's caption carries its instruction.

## Impact

- `internal/ui/widgets/dialogs/base.go`: render only the body plus the footer line, with no border, title, or placement.
- `internal/ui/dialogs.go`: size a dialog window to its content plus the frame.
- `internal/ui/widgets/dialogs/select.go`: use the cursor style for the focused option row.
- Tests in the dialogs and ui packages.
- Assumptions: the window title is the only title (the dialog's own title is no longer drawn); a dialog's footer is rendered as a body line; the window frame keeps its title bar and resize grip.
