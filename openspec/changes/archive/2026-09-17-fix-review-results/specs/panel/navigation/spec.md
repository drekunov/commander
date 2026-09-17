## ADDED Requirements

### Requirement: Only the focused panel shows the selection cursor
The selection cursor SHALL be rendered only on the panel that currently holds focus; the other panel SHALL NOT render a cursor. When window focus changes, cursor visibility SHALL follow the focus.

#### Scenario: Cursor follows window focus
- **WHEN** the left panel holds focus
- **THEN** only the left panel renders its selection cursor and the right panel does not

#### Scenario: Focus change moves the cursor
- **WHEN** the user changes focus from one panel to the other
- **THEN** the cursor appears only on the newly focused panel

### Requirement: Navigation is reliable under rapid input
Repeated navigation requests (Enter on a directory, Enter on the parent row, or Backspace) SHALL be applied so the focused panel ultimately displays the directory the user most recently requested; no valid request SHALL be dropped merely because the navigation queue is busy.

#### Scenario: Rapid navigation honors the latest request
- **WHEN** the user issues several navigation requests in quick succession
- **THEN** the focused panel ends on the directory of the last request
