# Rekunov Commander

[![Go](https://github.com/drekunov/gc/actions/workflows/go.yml/badge.svg)](https://github.com/drekunov/gc/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/drekunov/gc)](https://goreportcard.com/report/github.com/drekunov/gc)

A TUI file browser built with Go using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework. Norton Commander inspired.

## Features

- Interactive table view for directory contents
- Dialog windows (info, input) with blocking API
- CSS-based theming (Norton Commander classic theme included)
- Window manager with z-index, drag-to-move, drag-to-resize, and focus management

## Quick Start

```bash
# Build and run
go run ./cmd/commander

# Run the window manager example
go run ./cmd/wm-example
```

## CSS Theming

Styles are defined in `internal/config/styles.css` and embedded at build time. To override, place a `styles.css` file in the working directory.

Supported CSS properties:

| Property | Maps to |
|---|---|
| `color` | Foreground color |
| `background-color` | Background color |
| `padding` | Padding (top right bottom left) |
| `margin` | Margin (top right bottom left) |
| `border-style` | `rounded`, `double`, `thick`, `hidden`, `normal` |
| `border-color` | Border foreground color |
| `text-decoration` | `underline` |

Selectors: `.button`, `.active-button`, `.pressed-button`, `.dialog-box`, `.text`, `.input`, `.table`

## Window Manager

The `wm` package (`internal/ui/widgets/wm`) provides a compositing window manager for Bubble Tea applications.

```go
manager := wm.New()
manager.SetSize(width, height)

// Add a window with content, title, position, and size.
id := manager.Add(myTeaModel, "My Window", x, y, w, h)

// Programmatic control.
manager.Move(id, newX, newY)
manager.Resize(id, newW, newH)
manager.Focus(id)
manager.SetZIndex(id, 10)
manager.Remove(id)
```

Mouse interactions:
- Click title bar + drag to move
- Click bottom-right grip (◢) + drag to resize
- Click any window to bring to front and focus

Keyboard input is routed only to the focused window.

## Architecture

```
cmd/
  commander/          Entry point
  wm-example/         Window manager demo
internal/
  app/                Core application logic + interfaces
  config/             CSS-based Lipgloss theme system
  connectors/         Data source implementations (filesystem)
  ui/                 Bubble Tea UI layer
    widgets/
      wm/             Window manager (z-index, move, resize)
      mainform/       Main form container
      panel/          Table view widget
      dialogs/        Info and Input dialogs
      button/         Button widget
```

## Controls

| Key | Action |
|---|---|
| F10 | Quit |
| Tab | Cycle window focus (wm-example) |
| Arrow keys | Navigate tables/lists |
| Enter | Confirm / submit |

## SAST Tools

[PVS-Studio](https://pvs-studio.com/pvs-studio/?utm_source=website&utm_medium=github&utm_campaign=open_source) - static analyzer for C, C++, C#, and Java code.
