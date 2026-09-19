## ADDED Requirements

### Requirement: Menu button opens the sort window
Activating the Menu button SHALL open a modal sort window for the focused panel that lists its columns. Choosing a column SHALL sort the panel by that column ascending, unless it is the current sort column, in which case its direction SHALL invert. Escape SHALL close the window without changing the sort. Each panel SHALL keep its own sort mode and direction, and only the focused panel SHALL change.

#### Scenario: Choosing a different column
- **WHEN** the focused panel is sorted by Name ascending and the user chooses Size
- **THEN** the panel sorts by Size ascending

#### Scenario: Choosing the current column inverts
- **WHEN** the focused panel is sorted by Name ascending and the user chooses Name
- **THEN** the panel sorts by Name descending

#### Scenario: Canceling leaves the sort unchanged
- **WHEN** the sort window is open and the user presses Escape
- **THEN** the window closes and the sort is unchanged

#### Scenario: Each panel sorts independently
- **WHEN** the sort window is used
- **THEN** only the focused panel's sort mode changes and the other panel keeps its own

## REMOVED Requirements

### Requirement: Menu button cycles the sort mode
**Reason**: The Menu button now opens a column-selection window instead of cycling the sort directly.
**Migration**: Activate the Menu button and choose a column in the window; choosing the current column inverts its direction.
