# navigation Specification

## Purpose

Lets the user browse the filesystem from the keyboard in either of the two independent file panels, so both panels become live navigable browsers instead of a single static listing.

## Requirements

### Requirement: Each panel browses an independent directory
Each panel SHALL display its own independent directory listing, and navigating in one panel SHALL NOT change the other panel's listing or cursor position. When the application starts, both panels SHALL display the filesystem root.

#### Scenario: Navigating one panel leaves the other untouched
- **WHEN** the user changes the directory of one panel
- **THEN** the other panel keeps its previous directory, listing, and cursor position

#### Scenario: Both panels start at the filesystem root
- **WHEN** the application starts
- **THEN** each panel displays the root directory's listing

### Requirement: Navigation acts only on the focused panel
Keyboard navigation SHALL affect only the panel that currently holds focus. While any other window, such as a dialog, holds focus, navigation keys SHALL NOT change any panel.

#### Scenario: The focused panel responds to navigation
- **WHEN** the right panel is focused and the user navigates
- **THEN** only the right panel's directory or cursor changes and the left panel is unaffected

#### Scenario: A dialog keeps navigation away from panels
- **WHEN** a modal dialog holds focus
- **THEN** Enter and arrow keys operate on the dialog and do not move or reload the panel beneath it

### Requirement: Arrow keys jump and step through the list
Pressing the left arrow on the focused panel SHALL move the cursor to the first entry of the list. Pressing the right arrow SHALL move the cursor one entry down, stopping at the last entry. Up and down arrows SHALL continue stepping one row at a time.

#### Scenario: Left arrow jumps to the top
- **WHEN** the cursor is on a middle entry and the user presses the left arrow
- **THEN** the cursor moves to the first entry of the list

#### Scenario: Right arrow steps down one entry
- **WHEN** the user presses the right arrow
- **THEN** the cursor moves to the next entry, unless the cursor is already on the last entry, in which case it stays there

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

### Requirement: Backspace ascends to the parent directory
Pressing Backspace SHALL make the focused panel's parent directory its current directory and reload its listing. At the filesystem root, Backspace SHALL do nothing.

#### Scenario: Backspace moves up one level
- **WHEN** the focused panel is inside a nested directory and the user presses Backspace
- **THEN** the panel shows the parent directory's listing

#### Scenario: Backspace at the root is a no-op
- **WHEN** the focused panel is at the filesystem root and the user presses Backspace
- **THEN** the panel's directory and listing are unchanged

### Requirement: Unreadable directories are handled safely
When a directory cannot be read, the focused panel SHALL keep its current directory and listing, the failure SHALL be reported to the user, the application SHALL continue running, and the other panel SHALL remain unaffected.

#### Scenario: Failed read keeps the current listing
- **WHEN** the user tries to open a directory that cannot be read
- **THEN** the focused panel keeps showing its previous directory and listing and reports the error

#### Scenario: Failure does not disturb the application
- **WHEN** a directory read fails in one panel
- **THEN** the application keeps running and the other panel continues to work normally

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
