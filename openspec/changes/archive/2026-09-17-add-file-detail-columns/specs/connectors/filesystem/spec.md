## Purpose

Describes the local filesystem directory listing the connector supplies to the panel, so the panel can show each entry's name, size, and modification date and time while still telling directories from files.

## ADDED Requirements

### Requirement: Listing columns are Name, Size, Date, Time
The filesystem connector SHALL return a header row declaring the columns Name, Size, Date, and Time in that order, and one row per entry with that entry's values for those columns. It SHALL also return a hidden IsDir attribute so the panel can detect directories without showing a directory column.

#### Scenario: Header declares the four columns
- **WHEN** the connector reads a directory containing entries
- **THEN** the header row declares Name, Size, Date, and Time in that order

#### Scenario: File entries
- **WHEN** the connector reads a regular file
- **THEN** its Name is the file name, its Size is a human-readable size, and its Date and Time are the file's modification date (YYYY-MM-DD) and time (HH:MM:SS)

#### Scenario: Directory entries
- **WHEN** the connector reads a directory
- **THEN** its Size cell is `<DIR>` and its Date and Time are the directory's modification date and time

#### Scenario: Directory flag is hidden
- **WHEN** the connector returns an entry
- **THEN** it carries an IsDir attribute marked hidden, so the panel can detect directories without rendering a directory column
