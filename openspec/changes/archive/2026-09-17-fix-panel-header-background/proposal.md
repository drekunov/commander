## Why

The panel table's column header row renders with a different background than the listing body. The panel paints the themed `.table` background by wrapping the whole table, but bubbles' default header style is bold with no background, and its SGR reset clears the inherited background. The header strip therefore falls back to the terminal default background while the rows below stay on the themed panel background, so each panel looks visibly broken.

## What Changes

- The panel's column header row SHALL paint the same background as the panel's table body, so the header strip is continuous with the listing.
- The header background SHALL be derived from the injected table style, so a `styles.css` override changes the header and body together.
- This SHALL hold for both the focused and the unfocused panel.
- No change to column titles, column widths, header text weight/padding, row rendering, the cursor bar, navigation, the connector, or the data contract.

## Capabilities

### New Capabilities
- `panel/rendering`: visual presentation of the table panel; introduced with the requirement that the column header row shares the panel's table background.

### Modified Capabilities
<!-- None: this fixes the rendering of an existing panel region without changing navigation or data behavior. -->

## Impact

- `internal/ui/widgets/panel/panel.go`: the table styles built for rendering give the header the injected `.table` background; no other widget changes.
- Assumption: "different background" is a defect and "fixed" means the header uses the panel's table background (currently `#000080`). The header keeps its bold text, padding, and column widths.
