## 1. Config styles and classes

- [x] 1.1 Add `DialogCursorStyle`, `DialogOptionStyle`, `WindowTitleStyle`, `WindowTitleUnfocusedStyle`, and `WindowGripStyle` to `config.Styles` (`internal/config/config.go`).
- [x] 1.2 Load them from `.dialog-cursor`, `.dialog-option`, `.window-title`, `.window-title-unfocused`, and `.window-grip` in `LoadStyles`.
- [x] 1.3 Add the new classes with Norton Commander Classic values to the embedded `internal/config/styles.css`, keeping the existing `.dialog-box` frame class.
- [x] 1.4 Provide the fallbacks from design D2 so an override that omits the new classes keeps the current appearance.

## 2. Dialog set styling

- [x] 2.1 Render the Select dialog's focused option row with `DialogCursorStyle` and the other rows with `DialogOptionStyle` (`internal/ui/widgets/dialogs/select.go`).
- [x] 2.2 Render Info, Input, and Confirm body text/input/buttons over the `DialogBoxStyle` background so each body line carries the dialog background (`info.go`, `input.go`, `confirm.go`, and the shared `dialogBase.render`).
- [x] 2.3 Remove the `>` cursor prefix from Select option rows and render them flush-left (multi-select keeps its checkbox at the left edge), so focus is shown only by the dialog cursor style (`select.go`).
- [x] 2.4 Give the select option rows left and right padding through the theme's `.dialog-option`/`.dialog-cursor` `padding`, so the menu text is inset and the highlight row gets inner padding (`styles.css`).
- [x] 2.5 Title the placeholder mock dialog `Unimplemented` while its body keeps the requested action name (`internal/ui/ui.go`).
- [x] 2.6 Remove the event-loop-blocking `sendMsg` from `showMock` so the dialog renders on the normal post-Update render and the application stays responsive (`internal/ui/ui.go`).

## 3. Window frame chrome

- [x] 3.1 Replace the hardcoded focused/unfocused title colors in `buildTitleBar` with `WindowTitleStyle`/`WindowTitleUnfocusedStyle` composed over the frame background (`internal/ui/widgets/wm/renderer.go`).
- [x] 3.2 Render the resize grip with `WindowGripStyle` composed over the frame background instead of the hardcoded frame border color.

## 4. Tests

- [x] 4.1 Config test: the new selectors load into their `Styles` fields and a working-directory override changes them.
- [x] 4.2 Dialogs test: the Select focused row uses `DialogCursorStyle` and unfocused rows use `DialogOptionStyle`; dialog body lines carry the frame background.
- [x] 4.3 Renderer test: focused and unfocused titles use their injected styles and the grip uses the grip style.
- [x] 4.4 Dialogs test: option rows carry no `>` marker and start at the left edge in both single- and multi-select modes.
- [x] 4.5 Dialogs/config test: a themed option padding insets the rendered option row and loads from the embedded theme.
- [x] 4.6 ui test: the mock dialog opened by an unimplemented button is titled `Unimplemented`.
- [x] 4.7 ui test: `showMock` posts no message to the program from the event loop.

## 5. Verification

- [x] 5.1 Run `go build ./...`, `go vet ./...`, and `go test -race ./...` and keep them green.
- [x] 5.2 Run `make lint` and confirm golangci-lint passes.
- [x] 5.3 Manual smoke via `make run`: open the Menu sort window, an Info mock dialog, and an Input prompt, and confirm they render in the Norton Commander Classic scheme; drop a `styles.css` override and confirm the dialog, menu options, and frame title/grip follow it.
