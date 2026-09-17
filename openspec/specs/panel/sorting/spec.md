# sorting Specification

## Purpose

Lets the user order a panel listing by name, size, date, or time, keeping directories grouped above files, so a directory can be browsed in the order that suits the task.

## Requirements

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

### Requirement: Directories sort above files
A sorted listing SHALL place every directory entry above every file entry. Within each group, entries SHALL be ordered by the active sort column and direction. The synthetic `..` row SHALL remain the first row.

#### Scenario: Directories grouped first
- **WHEN** a listing is sorted by any column
- **THEN** all directory entries appear above all file entries

#### Scenario: Ordering within a group
- **WHEN** a listing is sorted by a column ascending or descending
- **THEN** each group's entries are ordered by that column in that direction

### Requirement: Size sorts numerically
The Size column SHALL order entries by their size as a number while the cell continues to display a human-readable size.

#### Scenario: Numeric size order
- **WHEN** a listing is sorted by Size ascending
- **THEN** a 208-byte file appears before a 3.9K file

### Requirement: Active sort column is indicated
The header SHALL mark the active sort column with a direction indicator, and the other columns SHALL be unmarked.

#### Scenario: Indicator follows the mode
- **WHEN** the sort is Name descending
- **THEN** the Name header shows a descending marker and no other column shows one

### Requirement: Selection survives a re-sort
When the sort mode changes, the panel SHALL keep the selected entry selected, placing the cursor on that entry in the new order.

#### Scenario: Cursor follows the selected entry
- **WHEN** the sort mode changes while an entry is selected
- **THEN** the cursor is on the same entry in the new order
