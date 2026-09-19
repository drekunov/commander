## Context

See proposal.md — Why. The filesystem connector's `ReadDir` (`internal/connectors/filesystem/connector.go`) returns a header row `Name, IsDir, Type` and one row per entry. `panel.columnsFor` turns every header attribute into a column, `rowFromAttrList` renders every attribute value, and `parentRow` synthesizes the `..` row from the column titles. `entryIsDir` reads the `IsDir` attribute from the entry data. `app.Attribute` is `{AttrName string; AttrValue any}`.

## Goals / Non-Goals

**Goals:**
- Columns Name, Size, Date, Time with compact widths (11, 6, 10, 8) and the agreed formats.
- Directory detection via a hidden `IsDir` attribute, with no IsDir column.
- Keep the panel connector-agnostic: the hidden flag and widths come from the header.

**Non-Goals:**
- Changing navigation, the cursor, theming, or the function-button bar.
- Dynamic column sizing or horizontal scrolling.

## Decisions

**1. `app.Attribute` gains `Hidden bool` and `Width int`.**
- `Hidden` marks an attribute that supplies data but is not a column; `Width` on a header attribute sets the column width (0 means the panel's default).
- Alternatives: hardcode the filesystem schema in the panel (rejected — couples the panel to one connector); return a richer header type (rejected — more churn for the same effect).

**2. Filesystem header and entries.**
- Header: `Name{Width:11}`, `Size{Width:6}`, `Date{Width:10}`, `Time{Width:8}`, `IsDir{Hidden:true}`.
- Each entry mirrors those attributes with values; `IsDir` carries the directory flag.

**3. Size format.**
- Base 1024 with one decimal and a `K`/`M`/`G` suffix (for example `4.0K`, `1.2M`); directories show `<DIR>`. A helper formats the value.

**4. Date and time from the modification time.**
- `info.ModTime().Format("2006-01-02")` and `Format("15:04:05")`, from `entry.Info()` — for a symlink this is the link itself.

**5. Panel excludes hidden attributes.**
- `columnsFor` skips hidden attributes and uses each header attribute's `Width`; `rowFromAttrList` skips hidden attributes; `setDefaultColumns` keeps the display titles and widths. `entryIsDir` is unchanged because it reads the entry data, which still contains `IsDir`.

**6. Parent `..` row.**
- `parentRow` builds a row for the display columns: Name `..`, Size `<DIR>` when a Size column exists, and blank Date and Time.

## Risks / Trade-offs

- [Widths too narrow] → Name is 11 cells, so long names are truncated by the table. This is the chosen compact layout; flexible sizing can be a later change.
- [Hidden-attribute mismatch] → the panel derives hidden columns from the header; an entry that omits the hidden attribute simply loses that data (directory detection falls back to false). The connector keeps header and entries consistent.
- [Symlink size/time] → `entry.Info()` reports the link, not its target; a symlinked directory is still detected through `isDirEntry`, which stats the target.
- [Existing tests] → panel tests build their own schemas, so the hidden/width behavior needs new cases and tests asserting a `Type`/`IsDir` column must be updated.

## Migration Plan

None — data/schema change. Rollback restores the Name/IsDir/Type header and the uniform `defaultColumnWidth`.

## Open Questions

None.
