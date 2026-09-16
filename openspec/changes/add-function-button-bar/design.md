## Context

The commander renders two file panels through a `wm.Manager` canvas that owns all windows (panels + dialogs) and routes keys only to the focused window. `mainform.Model` (widgets/mainform) owns the WM, lays the two panels out to the full screen height, and forwards every message to the WM. The app's stylesheet is already the "Norton Commander classic theme" and styles are injected as `config.Styles` (see proposal.md - Why). Current input surface: F10 quits (handled at `ui.Model` and `mainform`), F12 dumps the screen (`ui.Model`), Tab cycles WM focus, mouse drives WM drag/resize/focus. F1-F9 are unbound.

Constraints shaping the approach:

- Widgets never know about each other; composition happens in `mainform` / `ui.Model`. Message flow is a single Bubble Tea program rooted at `ui.Model` — a widget cannot push messages sideways, only produce `tea.Cmd`s whose results land back in the root `ui.Model`.
- Blocking dialogs (`ui.Info`, `ui.Confirm`, ...) are only ever invoked from non-event-loop goroutines today (`app.Run`); they add a WM window, `Send(tea.ResumeMsg{})`, and block on a channel.
- Theming: `Styles` fields map 1:1 to CSS classes in `internal/config/styles.css`; style values are cosmetic, so extending them is an implementation change, not a spec change.
- `make build`, `make lint` (gofumpt + gci + golangci-lint v2), `go test -race ./...`, and `make run` are the gates.

## Goals / Non-Goals

**Goals:**

- Ship the visible, always-present bottom function-button bar (F1-F10 Norton-Commander labels) as a widget with real key and mouse activation.
- Give every non-quit button a mock handler that reports its name through the existing Info dialog, proving the activation path end-to-end.
- Keep 10Quit/F10 quitting the app (unchanged real action, only in-app exit path).
- Leave a typed, injectable seam (ActionID + activation callback) so real file actions can replace mocks later without touching rendering, focus, or input handling.

**Non-Goals:**

- No real file operations (copy/move/delete/view/edit), no panel command menu, no changes to `Connector`/`SetData`/panel data contract, no changes to window drag/resize/focus or the F12 dump.
- No top menu bar; only the bottom function-button bar (see proposal).
- No changes to the app-layer (`internal/app`) interfaces.

## Decisions

### D1. New `buttonbar` widget; `mainform` owns it and renders it below the WM canvas

Add `internal/ui/widgets/buttonbar`. It is a pure widget: it knows its own state (focused index, pressed flash, ten buttons), renders exactly one row via the existing themed `button`-style rendering, and emits an activation as a typed message (`buttonbar.ActivateMsg{Action: ...}`). `mainform.Model` composes it and `ui.Model` owns the action handlers.

- *Why `mainform`, not `ui.Model`*: `mainform` already owns screen layout (panel sizing). Rendering is `mainform.View() = wm.View() + "\n" + bar.View()`.
- *Why a dedicated widget rather than reusing `wm` windows*: the bar is global UI chrome, not a framed, focusable, draggable window; putting it in the window stack would fight the z-order/drag logic.

### D2. Reserve the bottom row by shrinking the WM to `height - 1`

`mainform` keeps full-screen `width`/`height` from `tea.WindowSizeMsg`, but sizes the WM to `(width, height-1)` via `wm.SetSize` (do **not** forward the unmodified `WindowSizeMsg` into `wm.Update`, which would restore full height) and lays panels out to `height-1`. The bar row is outside the WM canvas, so dialogs (WM windows) can never overlap it and it is visible on every frame — including while dialogs are open.

- *Alternative considered*: overlay the bar by replacing the WM's last canvas row. Rejected: canvas rows are drawn bottom-up by windows; a stable chrome row needs its own reserved area.
- Consequence: `addDialogWindow` centering uses `wm.Height()` (now `height-1`) — dialogs sit one row higher; acceptable.

### D3. Activate-by-action messages: bar produces `ActivateMsg`, `ui.Model` owns handlers

`buttonbar.Model` never calls dialogs or quits. On activation (key, click, Enter) it flips its pressed state and returns a `tea.Cmd` resolving to `buttonbar.ActivateMsg{Action}`. `mainform.Update` relays that cmd; Bubble Tea delivers it to the root `ui.Model.Update`, which switches on the action:

- `ActionQuit` → `tea.Quit` (this is the only code path that turns an on-bar 10Quit activation into an exit; F10 already quits higher up).
- Any other action → spawn `go m.mockInfo(action)` (below).

- *Why a message hop instead of a callback into `ui`*: it keeps `buttonbar` and `mainform` dependency-free (no import of `ui`/dialogs) and uses Bubble Tea's normal message flow; the widget just declares what happened.
- This is the replaceable-action seam from the spec: swapping a mock means changing `ui.Model`'s `Action`→handler table, nothing in the widget.

### D4. Function keys are global; arrow/Enter navigation needs bar focus

- **F1-F9**: `mainform.Update` intercepts `tea.KeyMsg` F1-F9 *before* forwarding to the WM and maps each to its button. This matches NC muscle memory and needs no focus. (F10 already quits in `ui.Model`; F12 dump stays in `ui.Model`. Both take precedence over the bar because they never reach `mainform`'s bar path — the bar's 10Quit is still reachable by mouse/Enter and quits via D3.)
- **Arrow/Enter**: only while the bar holds keyboard focus. A mouse press on the bar row enters bar-focus mode; Left/Right move the highlighted button; Enter activates it; Tab exits back to the WM focus cycle; any mouse press outside the bar row clears bar focus. When the bar is not focused, Left/Right/Enter keep flowing to the focused window (panel) exactly as today.
- `mainform.Update` inspects `tea.MouseMsg`: on press at row `height-1` it forwards to the bar and clears bar focus when the press is above the bar; all other mouse messages (notably motion/release during window drags) are still forwarded to the WM unchanged so drag behavior is untouched.

### D5. Mock handlers reuse the existing blocking Info dialog, from a short-lived goroutine

A non-quit activation calls `ui.Model.mockInfo(action)`, which runs the existing `Info`-style blocking flow (`addDialogWindow` + `Send(ResumeMsg)` + select on dialog `Done()`/ctx) on a goroutine, exactly the pattern `app.Run` already uses for the read-error dialog. The dialog reports the mock, e.g. title `Mock` / message `"<action> is not implemented yet"`. `ui.Model` keeps a context cancelled when the tea program exits, so the goroutine unblocks instead of leaking.

- *Alternative considered*: a non-blocking dialog add from the event loop. Rejected: reimplements the blocking wait/removal logic; the goroutine approach reuses the code path and concurrency guarantees that already ship dialogs today.
- Consequence: while a mock dialog is open the WM focus moves to it; closing it refocuses the top window. Acceptable for placeholder behavior.

### D6. Ten buttons with shared layout math and new themed styles

`buttonbar` carries a fixed table of the ten actions with NC labels (`1Help 2Menu 3View 4Edit 5Copy 6RenMov 7Mkdir 8Delete 9PullDn 10Quit`) in order. The row divides the terminal width into ten equal segments; each label is truncated to fit with its trailing name kept visible; the row is padded/filled to end exactly at the right edge (single row, never wraps).

Theme: reuse the existing `.button`/`.active-button`/`.pressed-button` styles (the CSS is already NC-classic) and add a `.button-bar` fill style so the strip has a background. New `config.Styles` field `ButtonBarStyle` wired in `LoadStyles` + `styles.css`. No change to `config/theming` spec behavior — only themed cosmetics.

### D7. Verification strategy

- Unit tests for the new seams: `buttonbar` rendering (ten segments, right-edge fit, pressed flash) and activation mapping (F1..F10, click col→button, Left/Right/Enter focus) with injected styles — no terminal, no WM.
- A `mainform`-level test that F-keys produce an activation message and that bar-row presses do not reach the WM.
- Manual smoke: `make run` — bar visible across the full width, F5 shows the "Copy not implemented" dialog and closes without hanging, F10 quits, F12 still dumps; `go test -race ./...` and `make lint` green.

## Risks / Trade-offs

- **Stealing F1-F9 from focused windows** → they are unbound today, but a future dialog could want one (e.g. F1 help). Mitigation: keep the action dispatch table in one place (`ui.Model`) so a later policy (intercept only when no modal is open) is a one-line change; for the mock phase the bar wins.
- **Bar-focus mode intercepts arrows** → could confuse panel users. Mitigation: arrows are only intercepted while a click put the bar in focus mode; the default state (F-keys global, arrows to panels) is unchanged.
- **Mock dialog from a goroutine mutates WM state concurrently with the event loop** → pre-existing, proven pattern (error dialog from `app.Run`); WM list access is mutex-guarded. Mitigation: reuse it rather than inventing a new concurrent surface.
- **Panels lose one row of height** → cosmetic; NC reserves a bottom strip too.
- **Spamming F5 can stack several mock dialogs** → each mock Info blocks until confirmed; a user can only open one at a time this way unless they spam rapidly while one is open. Acceptable placeholder behavior; real actions later will gate on modality.

## Migration Plan

New feature behind existing composition root (`ui.New(styles)` / `mainform.New`), no external API change, no data migration. Rollback is a revert of the change's commits; the WM canvas and panel sizing revert with it.

## Open Questions

- Whether the eventual "real" actions (Copy/Move/...) are triggered from `ui.Model` through an `app`-owned callback or a new ui→app message channel — deliberately deferred until a real action exists; the `Action` dispatch table is the seam and its shape (func table today) can be wrapped without touching the widget.
