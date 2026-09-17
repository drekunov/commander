## Why

The panel shows entries in the connector's read order and offers no way to sort them. Let the Menu button (F2) cycle a per-panel sort mode across the columns, with directories always grouped above files, so the user can order a listing by name, size, date, or time.

## What Changes

- The `2Menu` action becomes real: it cycles the focused panel's sort mode through Name, Size, Date, and Time, each ascending then descending (Name↑, Name↓, Size↑, Size↓, Date↑, Date↓, Time↑, Time↓, then repeat). It no longer opens a mock dialog.
- Each panel keeps its own sort mode and direction; only the focused panel changes.
- Directories and files sort within their own groups, and directories always appear above files; the synthetic `..` row stays at the top.
- Size sorts numerically while still displaying human-readable text.
- The active column's header shows a direction marker (for example `Name▲` / `Name▼`).
- The listing is re-sorted in place, and the selected entry stays selected across a re-sort.
- No change to navigation, theming, the columns, or the other button actions.

## Capabilities

### New Capabilities
- `panel/sorting`: the per-panel sort mode, its cycle, the directory/file grouping, the numeric size ordering, and the header indicator.

### Modified Capabilities
- `ui/function-button-bar`: the Menu button selects the panel sort mode instead of running a mock.
- `connectors/filesystem`: the Size value is a sortable size that displays human-readable text.

## Impact

- `internal/app/connector.go`: a sortable `Size` value type.
- `internal/connectors/filesystem/connector.go`: return the sortable size for the Size column.
- `internal/ui/widgets/panel/panel.go`: sort state, value comparison, directory/file grouping, the header marker, and the in-place re-sort.
- `internal/ui/widgets/mainform/mainwindow.go`: expose the focused panel so the Menu action can reach it.
- `internal/ui/ui.go`: route the Menu action to the focused panel's sort cycle.
- Tests in the panel, connector, and ui packages.
- Assumptions: the default sort is Name ascending; name comparison is case-insensitive; sorting is local (the listing is not re-read); the cursor follows the selected entry by name.
