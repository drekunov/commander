## Context

Panels now navigate for real (see `panel/navigation`): each listing is delivered as connector rows (header row + entries, `panel/data-contract`), Enter descends into the selected directory, and Backspace ascends to the parent. The panel stores the connector entries (`m.rows`) and maps the table cursor directly onto them (cursor 0 == first entry). Empty directories currently render zero rows, and the filesystem root has no `..`. This change makes "one level up" visible and Enter-activatable via a leading `..` row while keeping Backspace. See proposal.md - Why.

Constraints that shape the approach:

- `..` must be a presentation affordance, not a data row: connectors and `panel/data-contract` stay unchanged, so no app/connector/mainform changes.
- Column layout is data-driven: the panel builds table columns from whatever header row a connector returns (filesystem: `path`, `Name`, `IsDir`, `Type`), so a synthetic row's cells must be derived from those column titles generically.
- All cursor movements already go through the bubbles `table` (keymap `left`→top, `right`→down) and `SetCursor(0)` resets on reload; the panel's own selection logic reads the table cursor and maps it to stored entries.
- Gates: `make build`, `make lint`, `go test -race ./...`, pty `make run`.

## Goals / Non-Goals

**Goals:**

- A synthetic `..` entry leads every non-root listing, including empty ones, and is absent at the filesystem root.
- Enter on `..` ascends exactly like Backspace; Backspace behavior is unchanged.
- Both panels behave the same; connector/data contract and app navigation request path are untouched.

**Non-Goals:**

- No connector/data-contract change; no `..` returned by data sources; no file operations on `..`.
- No changes to cursor keys, reload semantics, descend/ascend mechanics, error handling, or the WM.

## Decisions

### D1. `..` is synthesized by the panel widget

`panel.SetData(dir, data)` keeps the connector rows as the source of truth for real entries and, when the current directory has a parent (`filepath.Dir(dir) != dir`), prepends one synthetic display row labeled `..`. The synthetic row's cells are derived from the stored column titles: the `Name` column (matched case-insensitively; first column when absent) shows `..`, a `path` column shows the current directory, an `IsDir` column shows `true`, and any other column is left blank. If no column titles are known yet, a minimal single `Name` column is used so the row can render.

- *Why not in the connector*: `..` is navigation chrome, not a filesystem entry; injecting it at the data layer would leak UI concerns into every connector consumer and contradict `panel/data-contract`'s generic rows contract.

### D2. Display rows are offset from stored entries by the parent row

The table renders `[dotdotRow] + entryRows` when a parent exists and `entryRows` otherwise. Real entries stay in `m.rows`; the panel maps a table cursor to either the parent row (cursor 0 with a parent present) or a real entry (`cursor - 1`). All existing helpers (`selectedEntry`, `navigationCmd`) route through this single mapping so Right/Left stepping, viewport scrolling, and row selection keep working unchanged.

### D3. Enter on the parent row ascends; the request path is unchanged

`navigationCmd` gains one branch: when Enter is pressed on the parent row it emits the same `NavigateMsg` for `filepath.Dir(dir)` that Backspace already emits. Directory-descend and file-no-op branches are untouched, and Backspace keeps its existing behavior. Empty directories render only the parent row, so the user can always leave through the list; after a reload the cursor rests on row 0 (the `..` row when present), consistent with the existing "cursor on its first row" contract.

- *Why reuse the message*: the app side already reads any directory and delivers it back to the requesting panel; ascending is just a navigate request with the parent path.

### D4. No parent row at the filesystem root

The parent row is gated on `filepath.Dir(dir) != dir`, so the root listing never shows `..` and Enter/Backspace at the root keep their existing no-op semantics.

### D5. Verification strategy

- Panel unit tests: `..` leads a non-root listing and is absent at the root; Enter on `..` yields a `NavigateMsg` to the parent; the real directory below `..` still descends on Enter (index mapping); an empty non-root directory renders only the `..` row; Right/Left stepping still work across the offset.
- Existing panel navigation tests updated for the shifted indexing (cursor 0 is now the parent row when a parent exists).
- Gates: `go build`, `go test -race ./...`, `go vet ./...`, `make lint`, and a pty smoke checking a non-root panel shows a `..` row at the top and Enter on it returns to the parent while Backspace still ascends.

## Risks / Trade-offs

- **Cursor 0 now means "parent", not the first entry** → all index mapping funnels through one place (D2); existing tests catch drift. Users can still reach the first real entry with a single Right.
- **Synthetic cell contents are cosmetic** → only the `Name` cell (`..`) and navigational identity matter; attribute cells for `path`/`IsDir` are display-only and may not match a future connector's semantics. Accepted: `..` is an affordance, not an entry.
- **Empty-dir display changes from "no rows" to "one `..` row"** → a listing the user is inside always offers a way out; the existing "show an empty listing without error" contract is updated accordingly in the spec.

## Migration Plan

Pure panel-widget change; no API, app, or data-contract changes. Rollback is reverting the change's commits. Existing panel tests that assert cursor 0 == first entry need the offset update at the same time the display logic changes, so the change lands as one buildable unit.

## Open Questions

- Whether a fresh reload should rest the cursor on `..` or on the first real entry — cosmetic; the existing "first row" contract (row 0) is kept, and either choice would not change the specs beyond the current wording.
