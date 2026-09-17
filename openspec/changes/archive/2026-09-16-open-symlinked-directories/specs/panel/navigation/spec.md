## MODIFIED Requirements

### Requirement: Enter descends into the selected directory
Pressing Enter while a directory entry is selected SHALL make that directory the focused panel's current directory and reload the panel with its listing, placing the cursor on the first row. A symbolic link whose target is a directory SHALL count as a directory entry. Pressing Enter on a regular file, or on a symbolic link that does not resolve to a directory, SHALL do nothing and SHALL NOT change the listing. Opening an empty directory SHALL show only the ".." row and SHALL NOT report an error.

#### Scenario: Enter descends into a directory
- **WHEN** the user presses Enter on a selected directory entry
- **THEN** the focused panel shows that directory's listing and the cursor is on its first row

#### Scenario: Enter descends into a symlinked directory
- **WHEN** the user presses Enter on a symbolic link whose target is a directory
- **THEN** the focused panel shows the target directory's listing and the cursor is on its first row

#### Scenario: Enter on a symlink to a file does nothing
- **WHEN** the user presses Enter on a symbolic link whose target is a regular file
- **THEN** the panel's directory and listing are unchanged and no error is reported

#### Scenario: Enter on a broken symlink does nothing
- **WHEN** the user presses Enter on a symbolic link whose target does not exist
- **THEN** the panel's directory and listing are unchanged and no error is reported

#### Scenario: Enter on a file does nothing
- **WHEN** the user presses Enter on a regular file
- **THEN** the panel's directory and listing are unchanged and no error is reported

#### Scenario: Enter into an empty directory
- **WHEN** the user presses Enter on an empty directory
- **THEN** the panel shows only the ".." row, with no file entries, and does not report an error
