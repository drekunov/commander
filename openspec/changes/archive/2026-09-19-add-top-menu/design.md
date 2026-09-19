## Context

See proposal.md - Why. Existing constraints that shape the approach:

- The bottom bar (`internal/ui/widgets/buttonbar`) is a lock-free widget that owns its action identity, renders one full-width row, and emits a single `ActivateMsg{Action}`. The ui (`internal/ui/ui.go`) is the only consumer that resolves actions (`handleBarActivation`), which is the seam where real actions replace mocks.
- `mainform.Model` (`internal/ui/widgets/mainform/mainwindow.go`) owns layout: it sizes the `wm` canvas to `height - functionBarHeight` and renders `canvas + "\n" + bar`. Key routing there sends function keys to the bar before windows, and sends every other key to the window manager; a modal dialog suppresses F1-F9.
- Styles reach widgets by injection only (`config.Styles`); `LoadStyles` maps CSS selectors to fields and the embedded `styles.css` is overridable from the working directory.
- The `wm` canvas composites styled windows with z-index. A menu is not a window: it must overlay the top of the canvas without owning a frame.

## Goals / Non-Goals

**Goals:**

- A top menu bar widget that is independent of the button bar and reusable/testable in isolation.
- One action seam: top-menu item activation reuses the existing `buttonbar` action dispatch rather than a parallel mechanism.
- The pull-down overlays window content; opening a menu never resizes or refocuses windows.
- Every visible menu part is a themable CSS class with graceful fallback when a class is absent.

**Non-Goals:**

- Turning the top menu into a `wm` window or giving it a resize/move frame.
- Replacing the bottom function bar or changing its rendering.
- Wiring new real file operations; non-Files items stay mocks.
- Changing dialog behavior or the window manager's focus/drag/resize rules.

## Decisions

### D1. New `internal/ui/widgets/topmenu` package, mirroring `buttonbar`

The widget owns the menu tree, the active/opened state, key and mouse handling, and rendering. It takes `config.Styles` by injection and exposes `Activate()`, `Deactivate()`, `Toggle()`, `Active()`, `Update(msg)`, `View(width)`, and mouse hit-testing. It does not import `ui`, so it stays testable without the app model.

*Alternative considered*: extend `buttonbar` to hold both bars. Rejected: the two bars have different state machines (single row of actions vs. hierarchical menus) and combining them would bloat the bar and its tests.

### D2. Menu model: classic captions, Files maps to real actions, everything else mocks

Top-level menus and pull-down items (label, optional action):

| Menu (hotkey) | Items |
|---|---|
| Left (L) | Brief, Full, Info, Tree, Quick view, Sort by..., Filter... |
| Files (F) | View→`ActionView`, Edit→`ActionEdit`, Copy→`ActionCopy`, RenMov→`ActionRenMov`, Mkdir→`ActionMkdir`, Delete→`ActionDelete` |
| Commands (C) | Find file, History, Swap panels, Compare directories |
| Options (O) | Configuration, Save editor, Save setup |
| Right (R) | Brief, Full, Info, Tree, Quick view, Sort by..., Filter... |

An item with an action dispatches it through the bar's action contract; an item without an action is a placeholder that reports its own label through the existing mock dialog. The item list is data in the package, so adding real actions later is a table edit.

*Alternative considered*: derive menus from `buttonbar.actions`. Rejected: the user chose classic Norton Commander labels, which are wider than the ten bar actions.

### D3. Activation goes through the existing action seam; the bar still owns F9

F9 already maps to `ActionPullDn` in `buttonbar.ActionForKey`, so the bar emits `ActivateMsg{ActionPullDn}` with its pressed-feedback behavior intact. `ui.handleBarActivation` gains an `ActionPullDn` case that calls `m.main.TopMenu().Activate()`.

Once the bar is active, `mainform.handleKey` routes keys to the top menu before the window manager, and the top menu handles F9 itself as a toggle so the active-state key path is single-owner. While a modal dialog is open, `mainform` keeps suppressing F1-F9, so the top menu cannot activate.

*Alternative considered*: bind a new `tea.KeyF9` handler in `ui.handleGlobalKey` and bypass the bar. Rejected: it would duplicate the bar's key mapping and lose the pressed-feedback contract the bar spec already defines.

### D4. The pull-down is a compact box overlaid on the canvas

`mainform.View` becomes `topMenu + "\n" + canvas + "\n" + bar`. The widget's `View(width)` returns the top row, and when a menu is open also returns the pull-down as a compact multi-line box plus the column where its left edge starts. `mainform` composites that box onto the canvas with the window manager's ANSI-aware overlay (`wm.OverlayAt`), so the panels stay visible on either side of the box and no full-width blank rows are produced.

This keeps the menu above the panels while avoiding empty rows: the box is only as wide as its widest item, and the columns around it are the panels' own content.

*Alternative considered*: render the pull-down as full-width rows. Rejected: the columns beside the box become empty lines that visually blank the panels.

*Alternative considered*: render the pull-down as a borderless `wm` window. Rejected: it would add a borderless mode to the frame renderer and drag the window manager's focus/z-order rules into a transient overlay.

### D5. New style fields with fallbacks

Add to `config.Styles`, each from a selector:

| Field | Selector | Purpose |
|---|---|---|
| `TopBarStyle` | `.menu-top-bar` | Full-width top menu row background |
| `MenuCaptionStyle` | `.menu-caption` | Inactive caption |
| `MenuCaptionActiveStyle` | `.menu-caption-active` | Active caption |
| `MenuHotkeyStyle` | `.menu-hotkey` | Highlighted caption letter |
| `PulldownStyle` | `.menu-pulldown` | Pull-down background and item text |
| `PulldownCursorStyle` | `.menu-pulldown-cursor` | Focused pull-down item |

The widget composes each over `TopBarStyle`/`PulldownStyle`; an unset class falls back to `ButtonBarStyle`/`MenuBackgroundStyle`, so an override written before this change still renders a usable menu. Caption padding comes from the `.menu-caption`/`.menu-caption-active` horizontal padding and is rendered with the caption style, so the active highlight covers the padding as well as the text; the width uses the larger of the two states' padding so the row does not shift on activation. An active caption applies the active-caption background to its hotkey letter as well, so the selection highlight spans the whole caption instead of stopping before the first letter, and a pull-down item applies the cursor style to its text only, leaving the box borders on the pull-down frame style.

*Alternative considered*: reuse `.button-bar`/`.active-button`. Rejected: the top menu needs its own caption, hotkey, and pull-down rows, and the request asks for a self-contained CSS style.

### D6. Placeholder items reuse the mock dialog

`ui.showMock` is generalized from `buttonbar.Action` to a text argument (`showMockText`), keeping the single-dialog guard and the off-event-loop wait. Top-menu placeholder items call it with their label; Files items call `handleBarActivation`.

*Alternative considered*: a second mock dialog for menu items. Rejected: it would break the "mock activation does not stack dialogs" guarantee for the same Info dialog.

### D7. Layout and hit-testing coordinates

`mainform` reserves `topMenuHeight = 1` above and `functionBarHeight = 1` below; the window manager canvas height becomes `height - 2`, and panel windows keep `y = 0` within that canvas. Mouse routing checks the top row (`Y == 0`) first, then, when a pull-down is open, whether the press falls inside the box (`1 <= Y <= itemCount + 1` and `X` within the box), before handing clicks to windows. Presses outside the box fall through to the panels. Both the widget and `mainform` share the same row and column math so a click lands on the item the user sees.

## Risks / Trade-offs

- **Full-row pull-down hides panel rows while open** → Acceptable for a transient overlay; the menu closes on activation or Escape, and the spec only requires that panels yield the top row, not that the pull-down never covers content.
- **`mainform` tests assert `canvas + "\n" + bar`** → update `mainwindow_test.go` for the new three-part layout and reduced canvas height.
- **Existing test asserting `9PullDn` shows a mock** → update to assert the top menu activates instead.
- **Narrow terminals and long classic labels** → truncate caption and item text with `ansi.Truncate` and pad with `lipgloss.Width`, as `buttonbar` already does; add a narrow-width test.
- **ANSI styling leaking across the pull-down boundary** → composite the box with the window manager's ANSI-aware `wm.OverlayAt`, which keeps the canvas styling to the left and right of the box intact.
- **Menu item opens a dialog while the menu is still active** → deactivate the menu and close the pull-down before dispatching the item's command.

## Migration Plan

No migration. Apply in order: `config` (fields, selectors, embedded stylesheet) → `topmenu` widget (model, keys, mouse, render) → `mainform` (layout, key/mouse routing, a `TopMenu()` accessor) → `ui` (the `ActionPullDn` case and generalized mock) → tests. Keep the tree buildable at each step, then run `go build ./...`, `go vet ./...`, `go test -race ./...`, and `make lint`.

## Open Questions

None.
