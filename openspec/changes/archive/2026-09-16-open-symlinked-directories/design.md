## Context

See proposal.md — Why. The panel decides whether Enter navigates by reading the selected entry's `IsDir` attribute (`entryIsDir` in `internal/ui/widgets/panel/panel.go`); the filesystem connector sets that attribute from `os.ReadDir`'s `DirEntry.IsDir()`, which is false for every symbolic link. The connector interface (`app.Connector`) and the panel data contract need no signature changes.

## Goals / Non-Goals

**Goals:**
- Enter on a symlink that resolves to a directory shows that directory's listing, using the existing `IsDir`-driven navigation path.
- Links that do not resolve to a directory stay non-navigable and never surface an error dialog.

**Non-Goals:**
- No visual symlink badge or new column; the existing `Type` cell already shows `L---------`.
- No change to how Backspace, the `..` row, or read failures behave.
- No symlink loop detection beyond what the OS `stat` call provides.

## Decisions

**Resolve links in the connector, not the panel or app.**
- Chosen: `filesystem.ReadDir` classifies each entry by `entry.IsDir()` first, and for symlink entries (`entry.Type()&os.ModeSymlink != 0`) calls `os.Stat(filepath.Join(path, entry.Name()))`, setting `IsDir` from the resolved `FileInfo.IsDir()`. A failed `stat` (broken link, permission) yields `IsDir=false`.
- Why: the panel already routes Enter purely off `IsDir`, and the app already follows the path when reading, so one classification change makes both work. Keeping filesystem knowledge in the connector preserves the `Connector` abstraction — the panel stays data-source agnostic.
- Alternatives considered:
  - *Panel attempts navigation for every entry and lets the app probe*: would require a new connector capability (e.g. `Stat`/`IsDir`) and risks an error dialog on Enter over a regular file; rejected as more interface churn and worse failure UX.
  - *Add a separate resolved-directory attribute* (keep `IsDir` as the link's own type): preserves the current `IsDir=false` display but adds a second directory concept the panel must consult; rejected as unnecessary complexity for the stated requirement.

**Display follows classification.** A symlinked directory now shows `IsDir=true`; the `Type` column keeps showing the link, so the entry is still identifiable as a symlink.

## Risks / Trade-offs

- [A broken or looping symlink makes `os.Stat` fail or block] → `stat` failure maps to `IsDir=false`, so Enter is a no-op; loops are resolved by the kernel's own symlink limit.
- [One extra `stat` per symlink entry on large listings] → Only symlink entries are stat'd, not every entry; acceptable for an interactive TUI.
- [Display change: symlinked dirs now read as directories in the `IsDir` column] → Intended; the `Type` column still distinguishes links.

## Migration Plan

None — behavior-only change, no data or interface migration. Rollback is reverting the connector classification.

## Open Questions

None.
