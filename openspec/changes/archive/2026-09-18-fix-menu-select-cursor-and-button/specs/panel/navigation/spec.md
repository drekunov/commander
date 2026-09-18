## MODIFIED Requirements

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
