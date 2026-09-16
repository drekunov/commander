## Context

The app has three layers: `internal/app` (interfaces + orchestration), `internal/ui` (Bubble Tea), and `internal/ui/widgets` (WM, panels, dialogs, buttons). `main.go` runs the tea program and `app.Run` concurrently in an errgroup. See proposal.md - Why for the defect list.

Constraints that shape the approach:

- The `app.Connector` and `app.UI` interfaces are the seams between layers; `ui.Model` is the only `UI` implementation.
- Widgets currently read styles from the process-global `config.Values` (a mutable package var set in `init()` with a cwd `styles.css` read).
- `Connector.ReadDir` returns `[]app.AttributeList`; `ui/panels.go` currently treats `data[0]` as a header and `data[1:]` as rows, but `filesystem.ReadDir` emits no header — the first entry is lost.
- `ui.Model` blocking-dialog methods add a window to the WM, `Send(tea.ResumeMsg{})`, and block on a channel. `Confirm` blocks without a context.
- Bubbles' `table.Model` is not goroutine-safe; `SetData` mutates it from the app goroutine while the event loop renders.
- No tests exist. `make build`, `make lint` (gofumpt + gci + golangci-lint v2, `default: all`) are the gates.

## Goals / Non-Goals

**Goals:**

- Fix the two blockers (header/rows contract, render data race) with minimal behavioral surface change.
- Make every `Dialog` method functional, cancellable, and hang-free; remove the copy-pasted dialog shell.
- Eliminate the `config.Values` global in favor of injected styles.
- Extract rendering/mouse handling from `wm.Manager` without changing its public behavior.
- Land each step with `go build`, `go vet`, and `go test -race ./...` green, and add unit tests at the new seams.

**Non-Goals:**

- No new connectors, no directory navigation/selection features, no new themes.
- No changes to the `wm-example` behaviors beyond what the refactor forces.
- No package/dependency changes; no CI changes.
- No conversion of the whole app to a message-driven architecture — only the data-delivery path becomes message-driven.

## Decisions

### D1. Keep `AttributeList` as the contract shape; fix the header contract, don't introduce a typed schema

`filesystem.ReadDir` prepends a header row (attribute names) before entry rows, matching what `SetData` already assumes. A generic connector thus stays a "generic spreadsheet of rows" and future connectors (FTP, archive) need no schema coupling.

- *Alternative considered*: a typed `DirEntry{Name, IsDir, Type, Size}` return. Rejected: it hard-codes directory semantics into the `app` layer and would force the table to always render the same columns, contradicting the current generic `AttributeList` design.
- *Contract documentation*: the header-row requirement is written into the `panel/data-contract` spec so connector authors can't guess.

### D2. Deliver listing data as a tea message; keep the `Panel.SetData` signature

`App` keeps calling `ui.SetData(data)` (interface unchanged), but the implementation forwards it as `program.Send(panelDataMsg{data})`. `mainform.Update` applies the message by calling `panel.SetData`, which owns column/row construction (replacing `ui/panels.go`'s feature-envy code). `panel.TableView()/SetTableView()` are removed from the public surface.

- *Why*: Bubble Tea owns the `table.Model` state in the event loop; the only correct fix for the race is to make the mutation happen inside the loop. A mutex in `panel` would make the memory access safe but still tear the bubbles table's internal state and could deadlock with `View()`.
- *Alternative considered*: locking in `panel`. Rejected for the reason above and because it adds concurrency complexity to a widget that should be single-threaded.

### D3. Replace the `config.Values` global with injected `config.Styles`

The `Config` struct is renamed `Styles` and becomes a plain value. `config.LoadStyles()` resolves the cwd `styles.css` override (preserving today's runtime override behavior) and is called once from the composition root. Widget constructors take styles: `panel.NewPanel(styles)`, `dialogs.NewInfo(styles)`, `button.New(text, focused, styles)`. `wm.Manager` takes styles via `wm.New(styles)` and the renderer uses them.

- *Why*: the global is the single biggest blocker for unit tests and for instance-scoped theming (spec `config/theming`). Removing it turns every widget's `View()` into a pure function of (state, styles).
- *Alternative considered*: keep a global but make it immutable post-init. Rejected: it still prevents two widgets/instances with different styles and still reads the filesystem at package import.

### D4. Shared dialog shell via an embedded base widget

Extract the copy-pasted title-bar/footer/border composition (`info.go` and `input.go` both have identical `headerView`/`footerView`/`View` scaffolding) into a `dialogBase` embedded by `Info`, `Input`, and the new selection dialogs. `Free()` (dead, double-close hazard) is removed.

- *Why*: the shell is the most theme-sensitive part of the dialogs; one implementation means one place to fix theming.
- *Alternative considered*: one parameterized `Dialog` widget holding a content `tea.Model`. Rejected: Bubble Tea's `Update`/`View` idiom composes better through embedding, and each dialog keeps its own state machine.

### D5. New selection dialogs; password via echo-mask; confirm without string comparison

- `dialogs/select.go`: single-select (one of N, returns chosen option) and multi-select (subset with toggle, returns all chosen). Implemented with bubbles primitives; no new dependency. Multi-select uses a simple checkbox list model (checkmark toggle + Enter) since bubbles has no built-in multi-select.
- `Password` uses `textinput` with `EchoMode: textinput.EchoPassword`.
- `Confirm` gets a dedicated model that reports true/false; the `ui` layer stops reusing `Input` and comparing against `"Y"` (removes the magic string from logic).
- All blocking dialog calls gain a `context.Context` parameter (`Info`/`Input` already have one); `Confirm`, `Select`, `SelectMultiple`, `Password` select on `ctx.Done()` so F10-quit can never leave `app.Run` goroutine blocked (spec `ui/dialogs` — "Dialogs never hang").

### D6. Decouple `wm.Manager` without behavior change

Move `renderWindow`, `buildTitleBar`, `buildResizeGrip` into a `renderer` type (`wm/renderer.go`) and the mouse handlers into `wm/mouse.go` (still bound to `Manager`, but split for readability/testability). `Manager`'s public API (`Add`/`Remove`/`Focus`/`Move`/`Resize`/`Update`/`View`) is unchanged; the canvas stays as-is. Border math (`-2`/`-3`) and resize-grip offsets become named constants shared between `renderer` and `window` hit-testing.

- *Why*: the manager is 573 lines mixing state, input, and rendering; the rendering half is pure (state → string) and unit-testable once split. This is behavior-preserving, so it lives only in design/tasks — no spec delta.

### D7. Verification strategy

- `go test -race ./...` on every step (new unit tests at seams: connector header contract, panel `SetData`/message handling, dialog cancel paths, config parse + override, wm renderer snapshots).
- Manual smoke: `make run` and `go run ./cmd/wm-example` — the listing must show the full directory including the first entry; dialogs (Info via a triggered read error, plus a temporary debug binding if needed) must open, close, and not block on quit.
- `make lint` after the final step (formatters rewrite in place).

## Risks / Trade-offs

- **Message-delivered listing arrives after the program starts** → the first frames may briefly show an empty table. Mitigation: acceptable for a TUI; no ordering dependency exists today (`SetData` runs concurrently anyway).
- **Changing `Dialog` signatures (adding `ctx`) is a breaking internal API change** → no external callers exist; only `internal/ui` and the (currently unused) interface methods are affected. Mitigation: updated `app/ui.go` interface; keep signatures consistent across all dialog methods.
- **wm extraction may shift rendering pixels** → mitigation: renderer unit tests (canvas stamp + title bar composition) and manual wm-example check; keep public API byte-compatible.
- **Reworking `Select`/`SelectMultiple` from stubs to real UI is new surface** → mitigation: scope to minimal bubbles-based list; ship behind existing interface (callers today are none).

## Migration Plan

Apply in dependency order, keeping the tree buildable at each step:

1. `config/theming`: `Styles` type + `LoadStyles`, update all widget constructors and `wm`, remove `Values`.
2. `panel/data-contract`: header row in filesystem connector; `panel.SetData`; message-based delivery; remove `TableView/SetTableView`; delete old `panels.go` conversion code.
3. `ui/dialogs`: shared `dialogBase`; real `Select`/`SelectMultiple`/`Password`; dedicated `Confirm`; `ctx` on all blocking methods; remove `Free()`.
4. `ui/window-manager`: renderer/mouse extraction, named constants.
5. Tests for the seams; `make build`, `go test -race ./...`, `make lint` green.

Rollback: each step is a revertible commit; the message-driven delivery can fall back to the direct-call path without interface changes.

## Open Questions

- Whether single-select should reuse bubbles' `list` widget or a shared minimal list model with multi-select — defer to implementation; both satisfy the spec.
- Whether the `Panel` interface stays in `app` or moves closer to `ui` — cosmetic; no spec impact.
