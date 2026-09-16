## Why

Users can already ascend with the Backspace key, but the UI gives no visible cue that a parent directory exists or how to reach it. A leading `..` entry — the classic Norton Commander / Midnight Commander convention — makes "one level up" discoverable and mouse/Enter-friendly, while Backspace stays as the keyboard shortcut.

## What Changes

- The panel (presentation layer only) prepends a synthetic `..` row as the first row of the focused panel's listing whenever the panel is not at the filesystem root. Connectors and the `panel/data-contract` are unchanged — `..` is a navigation affordance, not a data row.
- Pressing Enter on the `..` row ascends one level, exactly like Backspace; Backspace keeps its existing ascend behavior.
- At the filesystem root no `..` row is rendered (there is no parent to go to).
- Empty directories still render their `..` row so a user can always leave via the list.
- Both panels behave identically.

## Capabilities

### New Capabilities
<!-- None. -->

### Modified Capabilities
- `panel/navigation`: the list now leads with a `..` parent entry (hidden at the root) whose Enter activation ascends one level; empty-directory wording updated to reflect that the parent entry remains visible.

## Impact

- `internal/ui/widgets/panel`: synthesizes and renders the `..` row per listing, remaps cursor-to-entry indexing by the presence of the row, routes Enter on `..` to an ascend request, and hides the row at the root.
- No `app`, `ui`, `connector`, `mainform`, WM, or data-contract changes; the existing navigation request path (Enter → `NavigateMsg` → app read → reload) is reused unchanged.
- Tests in the panel package updated/extended for the `..` row behavior.
