## Why

An architectural review of the current codebase found a data-corruption contract mismatch (the first directory entry is silently consumed as a table header), a data race between the app goroutine and the Bubble Tea event loop, four `Dialog` interface methods that return stubbed garbage, a latent deadlock in `Confirm`, and a global mutable `config.Values` singleton that blocks testability and theming extensibility. These defects surface as wrong UI data, intermittent races under `-race`, and a dialog API that lies to callers.

## What Changes

- **Fix the panel data contract**: `Connector.ReadDir` now returns `[header, rows...]`; the filesystem connector emits a proper header row, and `Panel` consumes the data through a new `SetData` method that owns column/row construction. The full listing renders with the first entry intact.
- **Eliminate the render data race**: panel data is delivered through the Bubble Tea event loop (tea message) instead of mutating `table.Model` from the app goroutine.
- **Make all `Dialog` methods functional**: `Select`, `SelectMultiple`, and `Password` get real implementations; `Confirm` is context-aware and cannot hang the app; `Info`/`Input`/`Confirm` share one dialog-shell implementation (removes copy-pasted `headerView`/`footerView`).
- **Replace the `config.Values` global with dependency injection**: widgets receive a `config.Styles` value at construction; the cwd `styles.css` override is resolved once at startup by the composition root. **BREAKING**: `config.Values` is removed.
- **Decouple the window manager**: extract rendering and mouse handling out of `wm.Manager` into focused types; all window-manager behavior (z-order, focus, drag, resize, key routing) is preserved.
- **Name magic layout constants** (window border math, dialog default sizes, table column width) and make the confirm sentinel configurable/decoupled from presenter text.

## Capabilities

### New Capabilities
- `panel/data-contract`: Defines the `ReadDir` → `SetData` row/header contract and race-free data delivery to the table.
- `ui/dialogs`: All `Dialog` interface methods behave as documented, with a shared, cancellable dialog shell.
- `config/theming`: Theme styles are injected per widget instead of read from a global mutable singleton.

### Modified Capabilities
- None. `openspec/specs/` has no main specs yet; all capabilities above are new.

## Impact

- **Code**: `internal/app/connector.go`, `internal/app/ui.go`, `internal/ui/ui.go`, `internal/ui/dialogs.go`, `internal/ui/panels.go`, `internal/ui/widgets/panel/`, `internal/ui/widgets/dialogs/`, `internal/ui/widgets/button/`, `internal/config/`, `internal/ui/widgets/wm/`, `internal/ui/widgets/mainform/`.
- **Connectors**: `filesystem` (and any future connector) must return a header row first. **BREAKING** contract change for `Connector.ReadDir`.
- **UI API**: `Panel` interface gains `SetData`; `panel.TableView()`/`SetTableView()` removed from public surface. **BREAKING** internal API change.
- **Dependencies**: none added. Uses existing bubbletea/lipgloss machinery.
- **Tests**: no tests exist today; this change introduces seams (injected styles, fake `UI`, message-driven dialogs) that make the affected packages unit-testable.
