## Why

The panel's directory columns are currently Name, IsDir, and Type, which show the entry kind but not the size or timestamps a file manager needs. Replace them with Name, Size, Date, and Time: a human-readable size (`<DIR>` for directories) and the modification date (`YYYY-MM-DD`) and time (`HH:MM:SS`). Directory detection moves to a hidden `IsDir` attribute so navigation keeps working without an IsDir column.

## What Changes

- The filesystem listing header declares Name, Size, Date, and Time, in that order, with compact widths (11, 6, 10, 8), plus a hidden `IsDir` attribute.
- Size is human-readable (for example `4.0K`, `1.2M`); directories show `<DIR>`.
- Date is the modification date (`YYYY-MM-DD`); Time is the modification time (`HH:MM:SS`).
- The data contract gains a hidden-attribute flag and an optional per-column width, so a connector can supply data the panel does not render as a column and can size its columns.
- The panel excludes hidden attributes from its columns and rendered rows, while still reading them for directory detection and the synthetic `..` row.
- No change to navigation, the cursor, theming, or the function-button bar.

## Capabilities

### New Capabilities
- `connectors/filesystem`: the local filesystem directory-listing schema — its columns and the format of each value.

### Modified Capabilities
- `panel/data-contract`: hidden attributes are excluded from columns and rendered rows, and the header may declare each column's width.

## Impact

- `internal/app/connector.go`: `Attribute` gains `Hidden bool` and `Width int`.
- `internal/connectors/filesystem/connector.go`: header and entries produce Name/Size/Date/Time plus a hidden IsDir; helpers for human-readable size and modification date/time.
- `internal/ui/widgets/panel/panel.go`: skip hidden attributes when building columns and rows, honor header widths, and build the `..` row from the display columns.
- Tests in the filesystem and panel packages.
- Assumptions: human-readable size uses base 1024 with one decimal and a `K`/`M`/`G` suffix; the `..` row shows Size `<DIR>` with blank Date and Time; a symlink's size and timestamps come from the link itself.
