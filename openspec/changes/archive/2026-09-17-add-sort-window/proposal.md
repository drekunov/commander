## Why

The Menu button currently cycles the panel's sort mode blindly, so reaching a particular column takes several presses. Replace the cycle with a sort window: the Menu button opens a modal list of the panel's columns, and the user chooses the sort column directly.

## What Changes

- The Menu button opens a modal sort window listing the focused panel's columns (Name, Size, Date, Time).
- Choosing a column sorts the panel by it ascending; choosing the current sort column inverts its direction.
- Escape closes the window without changing the sort.
- The direct cycle behavior is removed.
- The Select dialog gains Escape-to-cancel, returning no option.
- Every dialog window becomes modal: while a dialog is open, function keys F1 through F9 and mouse clicks outside it are ignored, and F10 still quits.
- No change to the directory/file grouping, the numeric size ordering, the header marker, or selection survival.

## Capabilities

### New Capabilities
<!-- None: this replaces how the sort column is chosen. -->

### Modified Capabilities
- `panel/sorting`: the Menu button opens the sort window, replacing the direct cycle.
- `ui/function-button-bar`: the Menu button opens the panel sort window, and function keys F1 through F9 are inactive while a dialog is open.
- `ui/dialogs`: the Select dialog can be canceled with Escape, and every dialog window is modal.

## Impact

- `internal/ui/widgets/dialogs/select.go`: Escape handling and a canceled state.
- `internal/ui/dialogs.go`: `Select` returns no option when canceled.
- `internal/ui/widgets/panel/panel.go`: expose the column titles and a `SortBy` operation; remove the cycle.
- `internal/ui/ui.go`: the Menu action opens the sort window and applies the chosen column.
- Tests in the panel, dialogs, and ui packages.
- Assumptions: the window lists the visible column titles without a current-sort indicator; the sort applies to the focused panel; the window is modal.
