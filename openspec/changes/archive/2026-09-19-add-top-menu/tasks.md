## 1. Config styles and classes

- [x] 1.1 Add `TopBarStyle`, `MenuCaptionStyle`, `MenuCaptionActiveStyle`, `MenuHotkeyStyle`, `PulldownStyle`, and `PulldownCursorStyle` to `config.Styles` (`internal/config/config.go`).
- [x] 1.2 Load them from `.menu-top-bar`, `.menu-caption`, `.menu-caption-active`, `.menu-hotkey`, `.menu-pulldown`, and `.menu-pulldown-cursor` in `LoadStyles`.
- [x] 1.3 Add the new classes with Norton Commander Classic values to the embedded `internal/config/styles.css`.
- [x] 1.4 Provide fallbacks in the top-menu widget so an override that omits the new classes still renders a usable menu (design D5).

## 2. Top-menu widget

- [x] 2.1 Create `internal/ui/widgets/topmenu/topmenu.go` with the menu tree (Left, Files, Commands, Options, Right) and item tables, including the Files item-to-`buttonbar.Action` mapping and the placeholder labels (design D2).
- [x] 2.2 Implement active state: `Activate`, `Deactivate`, `Toggle`, `Active`, and the active-menu and focused-item cursors.
- [x] 2.3 Implement key handling for F9 toggle, Left/Right menu movement, Down to open/move, Up to move, Enter to open/activate, Escape to close then deactivate, and per-caption hotkey letters (case-insensitive).
- [x] 2.4 Implement `View(width)` for the top row: captions in order with each hotkey letter drawn in the hotkey style, the active caption in the active style across the whole caption (the hotkey letter keeps the active background), themeable horizontal caption padding included in the highlight, padded to the full width on the bar background.
- [x] 2.5 Implement the pull-down render as a compact box (border plus items) with the focused item's cursor style applied to the item text only (never the frame), truncated to the box width, plus the box's start column, so `mainform` can overlay it without blanking the panels (design D4, D7).
- [x] 2.6 Implement mouse handling: caption and item hit-testing by row and column, reporting whether the press was consumed and returning the same activation result as the keyboard path; presses outside the open box fall through to the panels.
- [x] 2.7 Return an activation message that carries either a `buttonbar.Action` (Files items) or a placeholder label (all other items).

## 3. Main-form integration

- [x] 3.1 Add `topMenuHeight = 1`, own a `topmenu.Model`, and expose it with a `TopMenu()` accessor (`internal/ui/widgets/mainform/mainwindow.go`).
- [x] 3.2 Render the new three-part layout: top menu row, window-manager canvas, function bar; reduce the window-manager height by the top row and keep panel windows at `y = 0` in the canvas.
- [x] 3.3 Route keys to the top menu before the window manager while it is active, and handle F9 there as a toggle; keep F1-F9 suppressed while a modal dialog is open.
- [x] 3.4 Route mouse presses on the top row, and on open pull-down rows, to the top menu before windows; consume such clicks so panels do not handle them.

## 4. Action dispatch and mocks

- [x] 4.1 Add an `ActionPullDn` case to `ui.handleBarActivation` that activates the top menu instead of showing a mock (`internal/ui/ui.go`).
- [x] 4.2 Handle top-menu activation messages: dispatch Files item actions through `handleBarActivation`, and report placeholder item labels through the existing single mock dialog.
- [x] 4.3 Generalize `showMock` to `showMockText(text string)` so placeholder menu items reuse the same single-dialog guard (`internal/ui/ui.go`).
- [x] 4.4 Deactivate the menu and close the pull-down before dispatching an item's action so a dialog opens with the menu closed.

## 5. Tests

- [x] 5.1 Config test: the new top-menu selectors load into their `Styles` fields and a working-directory override changes them.
- [x] 5.2 Top-menu test: captions render in order with the correct hotkey letters, the active caption uses the active style, and its hotkey letter keeps the active background so the highlight spans the whole caption.
- [x] 5.3 Top-menu test: F9 activates and toggles, Left/Right move menus, Down opens and moves, Enter activates, and Escape closes then deactivates.
- [x] 5.4 Top-menu test: a caption hotkey (both cases) opens the matching menu, and letters do not open a menu while the widget is inactive.
- [x] 5.5 Top-menu test: a Files item yields its `buttonbar.Action` while a placeholder item yields its label; a mouse click matches the keyboard activation for the same cell.
- [x] 5.6 Top-menu test: the pull-down renders as a compact box (not full-width) with the focused item's cursor covering the item text only, not the box frame, and truncates on a narrow terminal.
- [x] 5.9 Top-menu test: a click outside the pull-down box is not consumed, so it falls through to the panels.
- [x] 5.10 Top-menu test: a pull-down row keeps its borders on the pull-down frame style and applies the cursor style only to the item text.
- [x] 5.11 Top-menu test: themed caption padding widens each caption and the active segment spans the text plus its padding.
- [x] 5.7 Main-form test: the view is top row plus canvas plus bar, the canvas height is reduced by one, and top-row/dialog key and mouse routing behaves as specified.
- [x] 5.8 ui test: `ActionPullDn` activates the top menu instead of opening a mock, and a placeholder item reports its label through the single mock dialog.

## 6. Verification

- [x] 6.1 Run `go build ./...`, `go vet ./...`, and `go test -race ./...` and keep them green.
- [x] 6.2 Run `make lint` and confirm golangci-lint passes.
- [x] 6.3 Manual smoke via `make run`: press F9, navigate with arrows and the L/F/C/O/R hotkeys, open each pull-down, activate a Files item and a placeholder item, and confirm the bar renders in the NC Classic scheme; drop a `styles.css` override and confirm the top menu, captions, hotkeys, and pull-down follow it.
