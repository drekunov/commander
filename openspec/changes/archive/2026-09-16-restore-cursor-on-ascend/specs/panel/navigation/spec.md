## ADDED Requirements

### Requirement: Ascending restores the cursor to the directory just left
When the focused panel ascends to a parent directory — whether by Backspace or by Enter on the ".." row — it SHALL place the cursor on the listing entry for the directory it just left. Each level of ascent SHALL restore the cursor for its own directory, so ascending several levels in succession returns to each directory in turn. When the parent listing has no entry for the directory just left, the panel SHALL place the cursor on the first row.

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
