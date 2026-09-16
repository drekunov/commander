## Why

The commander's two panels are static: the app reads `/` once, fills the left panel, and the user can only move a cursor, not browse. A file manager needs keyboard directory navigation — jump to the top, step down, descend into a directory, and go back up — so the panels become actually usable.

## What Changes

- Both panels become live, independent file browsers: each panel keeps its own current directory and listing, and keyboard navigation acts only on the currently focused panel. Initially both show the filesystem root.
- **Left arrow** jumps the cursor to the first entry of the focused panel's list; **Right arrow** steps the cursor down one entry (clamped at the bottom). Up/Down keep their existing one-row behavior.
- **Enter** on a directory descends one level: the focused panel is reloaded with that directory's listing. Enter on a regular file does nothing.
- **Backspace** ascends one level to the focused panel's parent directory; at the filesystem root it does nothing.
- After a listing change the panel cursor is placed back on the first entry.
- **BREAKING (internal)**: the `app.Panel.SetData` signature changes to carry the target panel and its directory path, because listings are now delivered to either panel independently. No external callers exist; only `internal/app`, `internal/ui`, and their tests are affected.
- Failed directory reads (permission denied, missing directory) do not change the panel's listing and are reported through the existing error dialog; the application keeps running.
- The data contract of `Connector.ReadDir` and table rendering is unchanged — the existing `ReadDir(path)` already reads any directory.

## Capabilities

### New Capabilities
- `panel/navigation`: keyboard browsing of the two independent panels — cursor jump/step keys, Enter descend / Backspace ascend against the real filesystem, per-panel state, and safe handling of unreadable directories.

### Modified Capabilities
<!-- None: panel/data-contract's connector contract and rendering requirements are unchanged. -->

## Impact

- `internal/app`: `App.Run` starts by loading both panels at the filesystem root; a new navigation request path (UI → app) lets the app read the requested directory and re-deliver it to the correct panel. `app.Panel.SetData` gains panel + directory parameters. Errors surface via the existing dialog flow.
- `internal/ui`: `ui.Model` owns the navigation sink and forwards panel navigation requests to the app; `SetData` delivers a per-panel, per-directory message to the correct panel.
- `internal/ui/widgets/panel`: each panel tracks its current directory, exposes its selected entry, emits Enter/Backspace navigation requests for the focused panel, resets the cursor after a reload, and extends the table keymap so Left jumps to the top and Right steps down.
- `internal/ui/widgets/mainform`: applies listings to the left or right panel per the delivered message and updates the window title with the panel's directory.
- No connector, theming, dialog, or window-manager spec/behavior changes.
