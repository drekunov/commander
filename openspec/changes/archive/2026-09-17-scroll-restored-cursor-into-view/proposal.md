## Why

Restoring the cursor after ascending works only when the restored row is already on screen. In a listing taller than the panel, the panel sets the table cursor but the table's viewport does not scroll to it, so the highlighted row is off-screen and the cursor appears to vanish — defeating the restore feature exactly where it matters most.

## What Changes

- When a listing is delivered and the selected row lies outside the visible rows, the panel SHALL scroll the listing so that row is visible.
- This applies to the ascent restore (Backspace or Enter on the `..` row) and to the first-row reset used when descending or reloading.
- The row that is selected does not change, and key navigation (arrows) is unchanged.
- No change to the connector, app layer, data contract, or the `..` row.

## Capabilities

### New Capabilities
<!-- None: this fixes the visibility of an existing navigation behavior. -->

### Modified Capabilities
- `panel/navigation`: the "Ascending restores the cursor to the directory just left" requirement gains a clause that the listing SHALL be scrolled so the restored entry is visible, with a scenario for a restored entry outside the visible rows.

## Impact

- `internal/ui/widgets/panel`: the cursor placement in `SetData` changes so the table scrolls the selected row into view. The app, connector, and data contract are unchanged.
- Assumption: the fix targets the panel's programmatic cursor placement only; arrow-key navigation already scrolls correctly.
