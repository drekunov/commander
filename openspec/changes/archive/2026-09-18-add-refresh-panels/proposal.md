## Why

The panels show a listing captured when the directory was read; if files change on disk (a build output appears, a log rotates, a download completes), the only way to see the change today is to navigate away and back. A refresh key gives the user a direct way to re-read the listings.

## What Changes

- Add a `Ctrl+R` shortcut that re-reads the current directory of **both** panels and redraws their listings.
- Keep the selected entry selected when it still exists after the refresh; fall back to the first row when it does not.
- Ignore the shortcut while a modal dialog is open, so dialogs stay modal.
- Add an app-level batch read so both panel listings can be requested together without the existing latest-wins navigation slot dropping one of them.

## Capabilities

### New Capabilities

- `panel/refresh`: re-reading both panels' current directories on demand and preserving the selection across the refresh.

### Modified Capabilities

None.

## Impact

- **Code**: `internal/ui/ui.go`, `internal/ui/panels.go` or a new refresh helper, `internal/ui/widgets/mainform/mainwindow.go`, `internal/ui/widgets/panel/panel.go`, `internal/app/app.go`, `cmd/commander/main.go`.
- **Tests**: new app test for the batch read; panel test for selection preservation on re-delivery; mainform/ui tests for the key being routed only when no dialog is open.
- **UI API**: internal-only additions (`SetRefresher`, `mainform.Panels`, `panel.ID`); no `app.UI` interface change.
- **Dependencies**: none added; `Ctrl+R` is already decoded by the pinned Bubble Tea.
