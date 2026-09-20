## Why

The file panels render every entry name in the same color, so directories, executables, symlinks, and common file categories are visually indistinguishable — a regression from the classic commander look the theme imitates. The colors should come from the existing CSS theming system so a `styles.css` override can restyle them like every other widget.

## What Changes

- Classify each directory entry into a file kind: `directory`, `symlink`, `executable`, and the extension-based categories `image`, `archive`, `source`, and `config`.
- Have the filesystem connector carry the kind on each entry as a hidden attribute, alongside the existing hidden `IsDir` flag.
- Render the entry's **Name cell** in the injected style for its kind, so only the name is colored and Size/Date/Time keep the table styling.
- Add theme classes `.file-directory`, `.file-executable`, `.file-symlink`, `.file-image`, `.file-archive`, `.file-source`, and `.file-config` to the CSS engine and the embedded Norton Commander Classic theme, each loaded into its own `config.Styles` field.
- Keep the focused cursor row uniform: the cursor style wins over the kind color, and a blurred panel shows the kind colors again.
- Leave entries whose kind has no matching class, and any unsupported kind, in the default table style.

## Capabilities

### New Capabilities
<!-- None: the behavior extends existing capabilities. -->

### Modified Capabilities

- `connectors/filesystem`: each entry SHALL carry a hidden file-kind classification derived from its type, mode, and name.
- `panel/rendering`: the panel SHALL color the Name cell of each entry by its kind using injected styles, except on the focused cursor row.
- `config/theming`: the theme SHALL expose a style class per file kind, loaded into its own `Styles` field and overridable from a working-directory `styles.css`.

## Impact

- **Code**: `internal/app/connector.go` (kind constants), `internal/connectors/filesystem/connector.go` (classification), `internal/config/config.go` and `internal/config/styles.css` (classes and fields), `internal/ui/widgets/panel/panel.go` (Name-cell coloring).
- **Tests**: `internal/connectors/filesystem/connector_test.go`, `internal/config/config_test.go`, `internal/ui/widgets/panel/panel_test.go`.
- **UI API**: internal-only; `config.Styles` gains fields, no `app.Connector` interface change.
- **Dependencies**: none added; reuses `charmbracelet/x/ansi` for ANSI-aware splicing.
