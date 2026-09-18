# refresh Specification

## Purpose

Lets the user re-read the current directory of both file panels on demand, so listings reflect the filesystem without navigating away and back.

## Requirements

### Requirement: Refresh re-reads both panels
Pressing `Ctrl+R` SHALL make both panels re-read their current directories and redraw their listings. The shortcut SHALL be ignored while a modal dialog is open, and refreshing SHALL NOT change which window or panel holds focus.

#### Scenario: Both panels are refreshed
- **WHEN** the user presses `Ctrl+R` while neither panel is under a modal dialog
- **THEN** each panel re-reads its current directory and renders the refreshed listing, including changes made on disk since the directory was last read

#### Scenario: A modal dialog blocks refresh
- **WHEN** a modal dialog holds focus and the user presses `Ctrl+R`
- **THEN** no panel is re-read and the dialog stays open

#### Scenario: Focus is unchanged by refresh
- **WHEN** the user presses `Ctrl+R`
- **THEN** the focused panel remains the focused panel and the other panel stays unfocused

### Requirement: Refresh preserves the selected entry
When a panel re-reads the same directory, it SHALL keep the entry that was selected before the refresh selected after it, and SHALL scroll the listing so that entry stays visible. When the previously selected entry is no longer present, the panel SHALL select its first row.

#### Scenario: Selection survives a refresh
- **WHEN** the user presses `Ctrl+R` while an entry is selected and that entry still exists
- **THEN** the cursor remains on that entry after the listing is redrawn

#### Scenario: Missing entry falls back to the first row
- **WHEN** the previously selected entry has been removed from the directory
- **THEN** the cursor moves to the first row of the refreshed listing

#### Scenario: Restored selection is scrolled into view
- **WHEN** the selected entry lies outside the visible rows after the refresh
- **THEN** the listing scrolls so that entry is visible and highlighted
