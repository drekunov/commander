## Why

Descending into a directory and then pressing Backspace (or Enter on the `..` row) resets the cursor to the first row, so the directory you just came from is no longer highlighted and you lose your place. Classic file managers put the cursor back on the directory you left, which keeps tree navigation oriented.

## What Changes

- Ascending to a parent — by **Backspace** or by **Enter on the `..` row** — places the cursor on the listing entry for the directory that was just left, instead of the first row.
- Each level remembers its own child, so ascending several levels in a row returns the cursor to each directory in turn.
- If the parent listing has no entry for the directory just left, the cursor falls back to the first row.
- Entering a directory still places the cursor on the first row.
- No change to Backspace's directory semantics, the `..` row, unreadable-directory handling, the connector interface, or the app layer.

## Capabilities

### New Capabilities
<!-- None: this refines an existing navigation behavior. -->

### Modified Capabilities
- `panel/navigation`: adds a requirement that ascending to a parent restores the cursor to the directory just left (with fallback to the first row), covering both the Backspace and `..` ascent paths and multi-level ascent.

## Impact

- `internal/ui/widgets/panel`: `Model` gains a per-panel pending-selection field; the ascent path records the base name of the directory being left, and the next delivered listing places the cursor on the matching entry before clearing the pending name. The app, connector, and data contract are unchanged.
- Assumption: the child is matched by its entry name; a name that is absent from the parent listing falls back to the first row.
