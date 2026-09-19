## MODIFIED Requirements

### Requirement: Listing columns are Name, Size, Date, Time
The filesystem connector SHALL return a header row declaring the columns Name, Size, Date, and Time in that order, and one row per entry with that entry's values for those columns. The Name column SHALL be marked flexible so the panel sizes it to fill the remaining width; Size, Date, and Time SHALL keep their fixed widths. The connector SHALL also return a hidden IsDir attribute so the panel can detect directories without showing a directory column.

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

#### Scenario: Name column is flexible
- **WHEN** the header row is inspected
- **THEN** the Name attribute is marked flexible and the Size, Date, and Time attributes are not
