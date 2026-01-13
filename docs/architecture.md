# GC File Manager - Architecture

## Class Diagram

```mermaid
classDiagram
    direction TB

    %% ========================================
    %% LAYER 1: ENTRY POINT
    %% ========================================
    namespace EntryPoint {
        class main {
            +main()
            +Run() error
        }
    }

    %% ========================================
    %% LAYER 2: APPLICATION LOGIC
    %% ========================================
    namespace ApplicationLayer {
        class UI {
            <<interface>>
        }

        class Dialog {
            <<interface>>
            +Info(ctx, title, footer, message)
            +Error(title, message)
            +Warning(title, message)
            +Confirm(title, message) bool
            +Input(ctx, title, footer, message) string
        }

        class Panel {
            <<interface>>
            +SetData(data AttributeList[])
        }

        class Connector {
            <<interface>>
            +Name() string
            +ReadDir(path) AttributeList[]
            +ReadFile(path) byte[]
        }

        class App {
            -ui UI
            -connector Connector
            +New(ui, connector) App
            +Run(ctx) error
        }

        class AttributeList {
            <<alias>>
            Attribute[]
        }

        class Attribute {
            +Name string
            +Value any
        }
    }

    UI --|> Dialog : extends
    UI --|> Panel : extends
    App --> UI : uses
    App --> Connector : uses
    Connector ..> AttributeList : returns

    %% ========================================
    %% LAYER 3: TERMINAL UI
    %% ========================================
    namespace PresentationLayer {
        class Model {
            <<tea.Model>>
            -program *tea.Program
            -mainform *mainform.Model
            +New() Model
            +Run(ctx) error
            +SetData(data)
            +Info(ctx, title, footer, msg)
            +Init() tea.Cmd
            +Update(msg) tea.Model
            +View() string
        }
    }

    Model ..|> UI : implements

    %% ========================================
    %% WIDGETS
    %% ========================================
    namespace Widgets {
        class MainformModel {
            <<tea.Model>>
            -panel *panel.Model
            -windows windowsList
            +AddWindow(w) WindowID
            +RemoveWindow(id)
            +View() string
        }

        class windowsList {
            <<map>>
            map~WindowID, tea.Model~
        }

        class PanelModel {
            <<tea.Model>>
            -table table.Model
            -header string
            -footer string
            +SetTableView(cols, rows)
            +View() string
        }

        class InfoDialog {
            <<tea.Model>>
            -title string
            -button *button.Model
            -done chan
            +Done() chan
            +View() string
        }

        class InputDialog {
            <<tea.Model>>
            -input textinput.Model
            -button *button.Model
            -done chan
            +Value() string
            +View() string
        }

        class ButtonModel {
            <<tea.Model>>
            -label string
            -focused bool
            +View() string
        }
    }

    Model *-- MainformModel : delegates
    MainformModel *-- PanelModel : contains
    MainformModel *-- windowsList : manages
    windowsList o-- InfoDialog : stores
    windowsList o-- InputDialog : stores
    InfoDialog *-- ButtonModel : contains
    InputDialog *-- ButtonModel : contains

    %% ========================================
    %% CONNECTORS
    %% ========================================
    namespace DataLayer {
        class FileSystem {
            +New() FileSystem
            +Name() string
            +ReadDir(path) AttributeList[]
            +ReadFile(path) byte[]
        }
    }

    FileSystem ..|> Connector : implements

    %% ========================================
    %% CONFIG
    %% ========================================
    namespace Configuration {
        class Config {
            +Buttons ButtonStyles
            +Dialogs DialogStyles
            +Table TableStyles
            +Text TextStyles
        }

        class Values {
            <<global>>
            Config instance
        }
    }

    %% ========================================
    %% CREATION RELATIONSHIPS
    %% ========================================
    main ..> Model : creates
    main ..> FileSystem : creates
    main ..> App : creates
```

## Component Diagram

```mermaid
flowchart TB
    subgraph Entry["cmd/gc (Entry Point)"]
        main[main.go]
    end

    subgraph App["internal/app (Application Layer)"]
        AppStruct[App]
        UIInterface[UI Interface]
        ConnInterface[Connector Interface]
        AttrList[AttributeList]
    end

    subgraph UI["internal/ui (Presentation Layer)"]
        UIModel[Model]
        subgraph Widgets["widgets"]
            Mainform[mainform.Model]
            Panel[panel.Model]
            Windows[windowsList]
            Dialogs[Info / Input]
            Button[button.Model]
        end
    end

    subgraph Connectors["internal/connectors (Data Layer)"]
        FS[filesystem.FileSystem]
    end

    subgraph Config["internal/config"]
        ConfigStruct[Config / Values]
    end

    main --> AppStruct
    main --> UIModel
    main --> FS

    AppStruct --> UIInterface
    AppStruct --> ConnInterface

    UIModel -.->|implements| UIInterface
    FS -.->|implements| ConnInterface

    UIModel --> Mainform
    Mainform --> Panel
    Mainform --> Windows
    Windows --> Dialogs
    Dialogs --> Button

    ConnInterface -->|returns| AttrList
    AttrList -->|SetData| UIInterface

    UIModel -.-> ConfigStruct
    Mainform -.-> ConfigStruct
    Panel -.-> ConfigStruct
    Dialogs -.-> ConfigStruct
```

## Sequence Diagram - Startup Flow

```mermaid
sequenceDiagram
    participant main
    participant App
    participant FileSystem
    participant UI as ui.Model
    participant Mainform as mainform.Model
    participant Panel as panel.Model

    main->>FileSystem: New()
    main->>UI: New()
    UI->>Mainform: New()
    Mainform->>Panel: New()
    main->>App: New(ui, connector)

    par Run concurrently
        main->>UI: Run(ctx)
        UI->>UI: tea.NewProgram().Run()
    and
        main->>App: Run(ctx)
        App->>FileSystem: ReadDir("/")
        FileSystem-->>App: []AttributeList
        App->>UI: SetData(data)
        UI->>Panel: SetTableView(columns, rows)
    end
```

## Sequence Diagram - Dialog Flow

```mermaid
sequenceDiagram
    participant App
    participant UI as ui.Model
    participant Mainform as mainform.Model
    participant Windows as windowsList
    participant Dialog as Info Dialog
    participant Button

    App->>UI: Info(ctx, title, footer, message)
    UI->>Dialog: NewInfo()
    UI->>Dialog: SetTitle(), SetText(), SetVisible(true)
    UI->>Mainform: AddWindow(dialog)
    Mainform->>Windows: store dialog with WindowID

    loop Wait for completion
        UI->>UI: select ctx.Done() or dialog.Done()
    end

    Note over Dialog: User presses Enter
    Dialog->>Dialog: visible = false
    Dialog->>Dialog: done <- struct{}{}

    UI->>Mainform: RemoveWindow(id)
    Mainform->>Windows: delete(windows, id)
```

## Layer Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      cmd/gc/main.go                             │
│                      (Entry Point)                              │
│  Creates instances, runs UI and App concurrently via errgroup  │
└───────────────────────────┬─────────────────────────────────────┘
                            │
            ┌───────────────┼───────────────┐
            │               │               │
            ▼               ▼               ▼
┌───────────────────┐ ┌───────────┐ ┌──────────────────┐
│   internal/ui     │ │internal/  │ │ internal/        │
│   (Presentation)  │ │app (App)  │ │ connectors       │
│                   │ │           │ │ (Data Layer)     │
│ Implements:       │ │ Defines:  │ │                  │
│ - UI interface    │ │ - UI      │ │ Implements:      │
│ - Dialog          │ │ - Panel   │ │ - Connector      │
│ - Panel           │ │ - Dialog  │ │                  │
│                   │ │ - Conn.   │ │ filesystem/      │
│ widgets/          │ │           │ │ FileSystem       │
│ - mainform        │ │ App       │ │                  │
│ - panel           │ │ orchestr. │ │                  │
│ - dialogs         │ │ both      │ │                  │
│ - button          │ │           │ │                  │
└───────────────────┘ └───────────┘ └──────────────────┘
            │               │               │
            └───────────────┼───────────────┘
                            │
                            ▼
              ┌─────────────────────────────┐
              │     internal/config         │
              │  (Global Lipgloss Styles)   │
              │                             │
              │  Config.Buttons             │
              │  Config.Dialogs             │
              │  Config.Table               │
              │  Config.Text                │
              └─────────────────────────────┘
```

## Data Flow

```
┌────────────────┐     ReadDir()     ┌─────────────────┐
│   Connector    │ ◄──────────────── │       App       │
│  (FileSystem)  │                   │                 │
└───────┬────────┘                   └────────┬────────┘
        │                                     │
        │ []AttributeList                     │ SetData()
        │                                     │
        ▼                                     ▼
┌────────────────────────────────────────────────────────┐
│                        UI                              │
│  ┌──────────────────────────────────────────────────┐  │
│  │                  mainform.Model                  │  │
│  │  ┌────────────────┐  ┌───────────────────────┐   │  │
│  │  │  panel.Model   │  │     windowsList       │   │  │
│  │  │                │  │  ┌─────┐  ┌─────────┐ │   │  │
│  │  │  ┌──────────┐  │  │  │Info │  │  Input  │ │   │  │
│  │  │  │  table   │  │  │  │     │  │         │ │   │  │
│  │  │  └──────────┘  │  │  └─────┘  └─────────┘ │   │  │
│  │  └────────────────┘  └───────────────────────┘   │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
```

## Key Design Patterns

| Pattern | Usage |
|---------|-------|
| **Interface Segregation** | UI = Dialog + Panel (separate concerns) |
| **Composition** | Model wraps mainform, mainform contains panel + windows |
| **Strategy** | Connector interface for pluggable data sources |
| **Observer (Channel)** | Dialog.Done() channel signals completion |
| **MVU (Model-View-Update)** | Bubble Tea pattern throughout UI |
| **Singleton** | config.Values global instance |
