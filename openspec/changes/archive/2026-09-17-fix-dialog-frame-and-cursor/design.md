## Context

See proposal.md — Why. `dialogs.dialogBase.render` (`base.go`) draws the dialog's own box: a `DialogBoxStyle` (double border, cyan) frame with a title header and a footer, then `lipgloss.Place`s it. The window manager's renderer (`renderer.go`) also frames every window — a `[ Title ]` title bar, side borders, and a bottom resize grip, using the same `DialogBoxStyle`. `addDialogWindow` (`dialogs.go`) sizes dialog windows to half the canvas. `dialogs.Select` styles its focused option with `ActiveButtonStyle`.

## Goals / Non-Goals

**Goals:**
- One frame per dialog (the window frame), a window sized to its content, and a Select cursor using the injected cursor style.
- Apply to every dialog.

**Non-Goals:**
- Changing modal behavior, the sort logic, the theme classes, or the panel frames.

## Decisions

**1. The dialog renders only its body.**
- `dialogBase.render` returns the body, appending the footer as a line when set; it no longer draws a border or title, and no longer places the content. The window title (set on the window) is the single title. The now-unused header/footer border helpers are removed.
- Alternative: drop the window frame for dialogs — rejected because the window chrome, title, and grip are shared.

**2. The dialog window is sized to its content and title.**
- `addDialogWindow` measures `content.View()` with `lipgloss.Width`/`Height` and the decorated title with `wm.TitleWidth`, takes the wider of the two for the inner width, adds `wm.FrameCols`/`wm.FrameRows`, clamps to the canvas, and centers. Sizing to the title keeps the caption from being truncated by the window frame.
- Risk: a dialog whose content changes size after opening keeps its initial window size; the current dialogs are fixed-size, so this is acceptable.

**3. The Select cursor uses the cursor style.**
- `Select.View` renders the focused row with `styles.CursorStyle` instead of `ActiveButtonStyle`.

**4. The sort window's prompt is its caption.**
- The sort window is opened with the caption `Sort: Select a column` and no body text, so the instruction sits in the window frame (the dialog's caption is the window title) instead of as a body line.
- Alternative: keep the caption `Sort` and the body prompt — rejected per the request.

## Risks / Trade-offs

- [Content-sized windows] → a very large dialog could exceed the canvas; the window is clamped to the canvas and the content clips.
- [Caption truncation] → the window is sized to the wider of the content and the decorated title, so the caption stays fully visible.
- [Footer placement] → the footer becomes a body line rather than a border caption; dialogs that use a footer (Info/Error) show it below the content.
- [Removed helpers] → `headerView`/`footerView` and `dialogFrameCols` become unused and are removed.
- [Shared base] → Info, Input, and Confirm are re-laid-out too; their tests and the smoke confirm no regression.

## Migration Plan

None. Rollback restores the dialog box and the half-screen sizing.

## Open Questions

None.
