## Why

The focused panel already lets the left arrow jump straight to the top of the list, but the right arrow only steps down a single row. Reaching the bottom of a long directory listing therefore requires many repeated key presses. Making the right arrow jump to the bottom gives the two horizontal arrows a fast, symmetric top/bottom navigation pair.

## What Changes

- **Right arrow** on the focused panel now moves the cursor directly to the **last entry** of the list, instead of stepping down one entry.
- **Left arrow** keeps its existing behavior: jump to the first entry.
- Up/Down arrows keep their existing one-row stepping behavior; the cursor is still clamped at the last entry.
- No change to Enter/Backspace, panel focus, listing contents, or the connector contract.

## Capabilities

### New Capabilities
<!-- None: this modifies the existing panel navigation behavior. -->

### Modified Capabilities
- `panel/navigation`: the "Arrow keys jump and step through the list" requirement changes — Right arrow jumps to the last entry (bottom) rather than stepping down one row, making Left/Right a top/bottom jump pair while Up/Down retain single-row stepping.

## Impact

- `internal/ui/widgets/panel`: the table keymap binds `right` to the bottom-jump action (`GotoBottom`) instead of line-down, and the corresponding widget tests are updated.
- No changes to `internal/app`, `internal/ui`, `internal/config`, connectors, dialogs, or the window manager.
