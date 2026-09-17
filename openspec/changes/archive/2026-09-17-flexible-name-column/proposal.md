## Why

The panel's columns all have fixed widths, so on a wide terminal the row leaves empty space to the right of the Time column instead of using it. Let the Name column flex to fill the space left after the fixed columns, and shrink on narrow panels, so names get the room and the row fills the panel.

## What Changes

- The header may mark one column flexible; the filesystem connector marks Name flexible while Size, Date, and Time keep their fixed widths (6, 10, 8).
- The panel sizes the flexible column to the panel width minus the fixed columns and the per-cell padding, with a minimum of eight cells; fixed columns keep their declared widths.
- Column widths are recomputed when the panel is resized, so the layout follows the terminal.
- No change to the columns themselves, the hidden IsDir flag, navigation, theming, or the function-button bar.

## Capabilities

### New Capabilities
<!-- None: this refines how an existing column is sized. -->

### Modified Capabilities
- `panel/data-contract`: the header may mark a column flexible, and the panel sizes it to fill the remaining width.
- `connectors/filesystem`: the Name column is marked flexible.

## Impact

- `internal/app/connector.go`: `Attribute` gains `Flex bool`.
- `internal/ui/widgets/panel/panel.go`: record each column's spec, compute the flexible width from the panel width, and recompute on resize while keeping fixed widths.
- `internal/connectors/filesystem/connector.go`: mark the Name header attribute flexible.
- Tests in the panel and filesystem packages.
- Assumption: at most one column is flexible; if several are marked, the first flexes and the rest keep their declared widths.
