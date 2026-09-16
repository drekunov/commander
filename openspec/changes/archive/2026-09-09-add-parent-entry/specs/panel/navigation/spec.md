## MODIFIED Requirements

### Requirement: Enter descends into the selected directory
Pressing Enter while a directory entry is selected SHALL make that directory the focused panel's current directory and reload the panel with its listing, placing the cursor on the first row. Pressing Enter on a regular file SHALL do nothing and SHALL NOT change the listing. Opening an empty directory SHALL show only the ".." row and SHALL NOT report an error.

#### Scenario: Enter descends into a directory
- **WHEN** the user presses Enter on a selected directory entry
- **THEN** the focused panel shows that directory's listing and the cursor is on its first row

#### Scenario: Enter on a file does nothing
- **WHEN** the user presses Enter on a regular file
- **THEN** the panel's directory and listing are unchanged and no error is reported

#### Scenario: Enter into an empty directory
- **WHEN** the user presses Enter on an empty directory
- **THEN** the panel shows only the ".." row, with no file entries, and does not report an error

## ADDED Requirements

### Requirement: Listings lead with a parent (..) entry
When the focused panel's current directory is not the filesystem root, the panel SHALL show a ".." entry as the first row of its listing. Pressing Enter on the ".." entry SHALL change the panel's current directory to the parent and reload its listing, behaving like Backspace. When the current directory is the filesystem root, the panel SHALL NOT show a ".." row. Empty directories SHALL still show the ".." row.

#### Scenario: Parent entry leads the listing
- **WHEN** a panel shows a non-root directory
- **THEN** the first row of its listing is the ".." entry, followed by the directory's entries

#### Scenario: Enter on the parent entry ascends
- **WHEN** the user presses Enter while the ".." row is selected
- **THEN** the focused panel shows the parent directory's listing

#### Scenario: No parent entry at the filesystem root
- **WHEN** a panel's current directory is the filesystem root
- **THEN** no ".." row is shown in its listing

#### Scenario: Empty directory keeps its parent entry
- **WHEN** a panel shows an empty non-root directory
- **THEN** its listing shows the ".." row and no file entries
