## 1. Theming dependency injection (config/theming)

- [x] 1.1 Rename `Config` to `Styles` in `internal/config/config.go` and add `LoadStyles() (Styles, error)` that resolves the cwd `styles.css` override against the embedded default, removing the `init()` side effect
- [x] 1.2 Remove the `Values` global and update `internal/config` package so importing it has no filesystem read at package init
- [x] 1.3 Thread `Styles` through `button.New`, `dialogs.NewInfo`, `dialogs.NewInput`, and `panel.NewPanel` constructors; delete all `config.Values` reads in `internal/ui/widgets/`
- [x] 1.4 Update `wm.New` (and the renderer) to take `Styles` instead of reading `config.Values` in `wm.go:337-436`
- [x] 1.5 Update the composition root (`cmd/commander/main.go`, `internal/ui/ui.go`, `internal/ui/widgets/mainform`) to call `LoadStyles` once and pass styles down
- [x] 1.6 Verify `make build`, `go test -race ./...`, and `go vet ./...` pass; `make run` and the cwd override (`styles.css` in working dir) still work

## 2. Panel data contract and race-free delivery (panel/data-contract)

- [x] 2.1 Prepend a header row (attribute names) in `internal/connectors/filesystem/connector.go` `ReadDir`; empty directories still return no rows
- [x] 2.2 Move column/row construction into `panel.Model.SetData([]app.AttributeList)` (header → columns, remainder → rows); remove `TableView()` and `SetTableView()` from `panel.Model`
- [x] 2.3 Delete the feature-envy conversion code in `internal/ui/panels.go` (`rowFromAttrList`, column/row building); `ui.Model.SetData` now only forwards data
- [x] 2.4 Deliver listing data as a tea message (`panelDataMsg`) via `program.Send` from `ui.Model.SetData`; handle it in `mainform.Model.Update` by calling `panel.SetData` on the event loop
- [x] 2.5 Keep the `Panel.SetData` interface signature unchanged so `App` callers compile untouched
- [x] 2.6 Verify with `go test -race ./...` that the listing update no longer races with rendering; `make run` shows the full directory including the first entry

## 3. Dialog layer (ui/dialogs)

- [x] 3.1 Extract the shared title-bar/footer/border shell from `info.go` and `input.go` into an embedded `dialogBase` widget (removes duplicated `headerView`/`footerView`/`View` scaffolding)
- [x] 3.2 Delete the dead `Free()` methods from `info.go` and `input.go` (double-close hazard)
- [x] 3.3 Implement `Password` with a masked `textinput` (`EchoMode: EchoPassword`) instead of delegating to the plain `Input`
- [x] 3.4 Add `dialogs/select.go` with single-select (one of N) and multi-select (subset with toggle) models using bubbles primitives; wire `Select` and `SelectMultiple` in `internal/ui/dialogs.go` to them
- [x] 3.5 Replace `Confirm`'s `Input` + `== "Y"` string comparison with a dedicated confirm model returning true/false (removes the magic string from logic)
- [x] 3.6 Add `context.Context` to every blocking dialog method in `app.UI.Dialog` (`Confirm`, `Select`, `SelectMultiple`, `Password`) and select on `ctx.Done()` in each block so F10-quit can never hang the caller
- [x] 3.7 Update the `app` interface (`internal/app/ui.go`) and any internal callers to the new signatures; verify `make build` and `make lint`

## 4. Window manager decoupling (behavior-preserving)

- [x] 4.1 Move `renderWindow`, `buildTitleBar`, `buildResizeGrip` into a `renderer` type in `wm/renderer.go` (pure state → string)
- [x] 4.2 Split mouse handling (`handleMouse`, `handleMousePress`) into `wm/mouse.go`, keeping `Manager`'s public API unchanged
- [x] 4.3 Introduce named constants for the border frame math (`-2`/`-3` in `wm.go:316-317`), resize-grip offsets (`window.go:57`), and dialog default sizes, sharing them between renderer and hit-testing
- [x] 4.4 Verify `go run ./cmd/wm-example` (drag, resize, focus, Tab-cycle, F10) and `make run` behave identically to before the split

## 5. Tests, lint, and verification

- [x] 5.1 Unit-test `config.LoadStyles` (override present/absent) and `ParseCSS`/`ApplyStyle` without cwd side effects
- [x] 5.2 Unit-test the filesystem connector header contract against `t.TempDir()` fixtures and the panel `SetData` column/row mapping (first entry included)
- [x] 5.3 Unit-test the dialog cancel path (context cancelled while a dialog is open returns without hang) using a fake program seam
- [x] 5.4 Add `wm` renderer tests (title bar composition, canvas stamp clipping) for the extracted pure functions
- [x] 5.5 Add an `app`-level test with fake `UI`/`Connector` covering the `ReadDir` error dialog path and successful `SetData` path
- [x] 5.6 Run `go test -race ./...`, `make lint`, and `go vet ./...`; confirm `openspec validate refactor-arch-review` still passes
