# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

gc is a terminal-based file manager written in Go, built with the Charm stack (Bubble Tea, Lip Gloss, Bubbles). It provides a TUI for navigating and managing files.

## Commands

**Build:**
```bash
go build ./...
```

**Run:**
```bash
go run ./cmd/gc
```

**Lint (formats and checks code):**
```bash
make lint
```
This runs gofumpt, gci (import sorting), and golangci-lint.

**Update dependencies:**
```bash
make modup
```

**Run tests:**
```bash
go test ./...
```

**Run a single test:**
```bash
go test -run TestName ./path/to/package
```

## Architecture

The application follows a three-layer architecture with clear separation of concerns:

### Layer 1: Entry Point (`cmd/gc/main.go`)
Bootstraps the application by creating UI and App instances, then runs them concurrently using an errgroup. The UI and App communicate through interfaces.

### Layer 2: Application Logic (`internal/app/`)
- `App` struct orchestrates between UI and data connectors
- Defines interfaces (`UI`, `Connector`) that abstract dependencies
- `Connector` interface allows pluggable data sources (filesystem, future: remote, archives)
- `AttributeList` is the universal data format passed between connectors and UI

### Layer 3: Terminal UI (`internal/ui/`)
- Built on Bubble Tea's Model-View-Update pattern
- `Model` wraps the tea.Program and delegates to `mainform.Model`
- `mainform` manages the main window, panels, and floating windows (dialogs)
- `windowsList` maintains a map of modal windows with unique IDs

### Connectors (`internal/connectors/`)
Data source adapters implementing the `Connector` interface. Currently only `filesystem` exists, which reads local directories.

### UI Widgets (`internal/ui/widgets/`)
Reusable Bubble Tea components: `panel` (table display with header/footer), `button`, and `dialogs` (info, input modals).

### Styling (`internal/config/`)
Global lipgloss styles initialized at startup via `config.Values`.

## Key Patterns

- All UI components follow Bubble Tea's `Init()`, `Update()`, `View()` pattern
- Dialogs use channels (`Done()`) to signal completion back to callers
- Data flows: Connector → AttributeList → UI.SetData() → table rows
- F10 quits the application
