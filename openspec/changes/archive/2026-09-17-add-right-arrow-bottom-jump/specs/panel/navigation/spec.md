## REMOVED Requirements

### Requirement: Arrow keys jump and step through the list
**Reason**: The right arrow no longer steps down one entry; it now jumps to the last entry, so this requirement's behavior contract changes rather than being extended. A MODIFIED delta cannot drop the superseded "Right arrow steps down one entry" scenario, so the requirement is replaced.
**Migration**: Use the Down arrow to step one row at a time. The Left arrow still jumps to the first entry; the Right arrow now jumps to the last entry.

## ADDED Requirements

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
