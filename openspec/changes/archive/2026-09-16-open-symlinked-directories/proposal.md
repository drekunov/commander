## Why

Pressing Enter on a directory shows its contents, but a symbolic link that points to a directory is reported as a non-directory (`IsDir=false`), so Enter does nothing on it. On a typical Linux root this hides `/bin`, `/lib`, `/sbin`, and similar links, which classic file managers happily enter.

## What Changes

- An entry whose path is a **symbolic link to a directory** is treated as a directory: pressing Enter on it shows that directory's listing, exactly like a real directory.
- Symbolic links that point to a file, or that are broken, remain non-directories: Enter still does nothing and reports no error.
- Because the entry is now classified as a directory, its `IsDir` cell shows `true`; the `Type` cell still shows the link (`L---------`), so the symlink is still visible.
- No change to Backspace, the `..` parent row, unreadable-directory handling, the connector interface, or the panel data contract.

## Capabilities

### New Capabilities
<!-- None: this refines an existing navigation behavior. -->

### Modified Capabilities
- `panel/navigation`: the "Enter descends into the selected directory" requirement changes — a symbolic link whose target is a directory now counts as a directory entry for Enter, with new scenarios for symlinked directories and for links that are not directories.

## Impact

- `internal/connectors/filesystem`: `ReadDir` determines `IsDir` by resolving symlinks (follow the link's target) instead of relying only on `DirEntry.IsDir()`. The `Connector` interface and its signatures are unchanged.
- `internal/ui/widgets/panel` and `internal/app`: no code change — they already navigate whenever the selected entry's `IsDir` attribute is true.
- Tests: connector tests gain symlink-to-directory, symlink-to-file, and broken-symlink cases.
- Assumption: resolving links adds one `stat` per symlink entry; acceptable for a TUI listing.
