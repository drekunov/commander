# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Build the application
go build ./cmd/commander

# Run the application
go run ./cmd/commander

# Lint (requires gofumpt, gci, golangci-lint)
make lint

# Update dependencies
make modup

# Run tests
go test ./...

# Run a single test
go test -run TestName ./path/to/package
```

## Architecture Overview

GC is a TUI (Terminal User Interface) file browser built with Go using the Bubble Tea framework. It displays directory contents in an interactive table with dialog windows for user interaction.

### Key Dependencies
- `github.com/charmbracelet/bubbletea` - TUI framework (Elm architecture)
- `github.com/charmbracelet/bubbles` - Pre-built TUI components
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/evertras/bubble-table` - Table widget

### Package Structure

```
cmd/gc/main.go          # Entry point - wires up components
internal/
  app/                  # Core application logic
    app.go              # App struct, orchestrates UI and connectors
    ui.go               # UI interface (Dialog, Panel methods)
    connector.go        # Connector interface and data types
  config/config.go      # Global Lipgloss styles
  connectors/           # Data source implementations
    filesystem/         # Filesystem connector (ReadDir, ReadFile)
  ui/                   # User interface layer
    ui.go               # Main Bubble Tea Model
    dialogs.go          # Dialog implementations
    panels.go           # Panel data handling
    widgets/            # Reusable components
      mainform/         # Main window container + windows manager
      panel/            # Table view widget
      dialogs/          # Input and Info dialogs
      button/           # Button widget
```

### Data Flow

1. `main.go` creates a Connector (filesystem) and UI instance, passes both to App
2. App and UI run concurrently via `errgroup`
3. App reads data via Connector and calls `UI.SetData()` to populate the panel
4. UI renders mainform containing the panel; dialogs overlay as windows
5. F10 triggers shutdown via context cancellation

### Core Interfaces

**Connector** (`internal/app/connector.go`):
```go
type Connector interface {
    Name() string
    ReadDir(path string) ([]AttributeList, error)
    ReadFile(path string) (string, error)
}
```

**UI** (`internal/app/ui.go`): Defines Dialog and Panel methods that the App uses to interact with the UI layer.

### Widget Hierarchy

```
mainform (mainwindow.go)
├── panel (table view)
└── windowsList (overlay dialogs)
    ├── InfoDialog
    ├── InputDialog
    └── (future: Password, Select, Confirm, etc.)
```

### Styling

All styles are defined in `internal/config/config.go` as global Lipgloss styles. Modify this file to change the theme.

### Extensibility

- **New data sources**: Implement the `Connector` interface in `internal/connectors/`
- **New dialogs**: Add to `internal/ui/widgets/dialogs/` following the InfoDialog/InputDialog pattern
- **New widgets**: Create in `internal/ui/widgets/` implementing the Bubble Tea Model interface