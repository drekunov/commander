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

### Requirement: Arrow keys jump to the list edges and step through the list
Pressing the left arrow on the focused panel SHALL move the cursor to the first entry of the list. Pressing the right arrow SHALL move the cursor to the last entry of the list. Up and down arrows SHALL continue stepping one row at a time. The cursor SHALL remain clamped at the first and last entries.

#### Scenario: Left arrow jumps to the top
- **WHEN** the cursor is on a middle entry and the user presses the left arrow
- **THEN** the cursor moves to the first entry of the list

#### Scenario: Right arrow jumps to the bottom
- **WHEN** the cursor is on an entry above the last one and the user presses the right arrow
- **THEN** the cursor moves to the last entry of the list

#### Scenario: Right arrow on the last entry stays put
- **WHEN** the cursor is already on the last entry and the user presses the right arrow
- **THEN** the cursor stays on the last entry

#### Scenario: Up and down arrows still step one row
- **WHEN** the user presses the down arrow on a non-last entry or the up arrow on a non-first entry
- **THEN** the cursor moves exactly one row in that direction

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

### Requirement: Ascending restores the cursor to the directory just left
When the focused panel ascends to a parent directory — whether by Backspace or by Enter on the ".." row — it SHALL place the cursor on the listing entry for the directory it just left, and SHALL scroll the listing so that entry is visible. Each level of ascent SHALL restore the cursor for its own directory, so ascending several levels in succession returns to each directory in turn. When the parent listing has no entry for the directory just left, the panel SHALL place the cursor on the first row.

#### Scenario: Backspace restores the cursor
- **WHEN** the panel is inside /a/b and the user presses Backspace
- **THEN** the panel shows /a and the cursor is on the entry for b

#### Scenario: Enter on the parent row restores the cursor
- **WHEN** the panel is inside /a/b and the user presses Enter on the ".." row
- **THEN** the panel shows /a and the cursor is on the entry for b

#### Scenario: Each level remembers its own child
- **WHEN** the user descends from /a into b, then into c, and then ascends one level at a time
- **THEN** the cursor is on c in /a/b and then on b in /a

#### Scenario: Missing entry falls back to the first row
- **WHEN** the panel ascends to a parent listing that has no entry for the directory just left
- **THEN** the cursor is on the first row of that listing

#### Scenario: Descending starts at the first row
- **WHEN** the user enters a directory
- **THEN** the cursor is on the first row of that directory's listing

#### Scenario: Restored entry is scrolled into view
- **WHEN** the restored directory lies outside the visible rows of a parent listing that is taller than the panel
- **THEN** the listing scrolls so the restored entry is visible and highlighted

### Requirement: Only the focused panel shows the selection cursor
The selection cursor SHALL be rendered only on the panel that currently holds focus; the other panel SHALL NOT render a cursor. Cursor visibility SHALL follow window focus, including when focus returns to a panel after a dialog window closes, not only while the application is processing an input message.

#### Scenario: Cursor follows window focus
- **WHEN** the left panel holds focus
- **THEN** only the left panel renders its selection cursor and the right panel does not

#### Scenario: Focus change moves the cursor
- **WHEN** the user changes focus from one panel to the other
- **THEN** the cursor appears only on the newly focused panel

#### Scenario: Cursor returns after a dialog closes
- **WHEN** a dialog window closes and focus returns to a panel
- **THEN** that panel renders its selection cursor again without requiring further input

### Requirement: Navigation is reliable under rapid input
Repeated navigation requests (Enter on a directory, Enter on the parent row, or Backspace) SHALL be applied so the focused panel ultimately displays the directory the user most recently requested; no valid request SHALL be dropped merely because the navigation queue is busy.

#### Scenario: Rapid navigation honors the latest request
- **WHEN** the user issues several navigation requests in quick succession
- **THEN** the focused panel ends on the directory of the last request
