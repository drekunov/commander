## ADDED Requirements

### Requirement: Dialog windows have a single frame
A dialog SHALL render its content inside the single window-manager frame: the dialog widget SHALL NOT draw its own border or title, and the window SHALL be sized to fit the dialog content plus the frame. A dialog's footer SHALL be rendered as a body line.

#### Scenario: Only one frame is drawn
- **WHEN** a dialog is shown
- **THEN** exactly one border surrounds its content (the window frame) and it has a single title

#### Scenario: The window fits the content
- **WHEN** a dialog is shown
- **THEN** the window is sized to the dialog content plus the frame rather than a fixed fraction of the screen

#### Scenario: A footer is a body line
- **WHEN** a dialog has a footer
- **THEN** the footer text appears inside the frame as a body line, not as a second border

### Requirement: Select dialog cursor uses the cursor style
The Select dialog SHALL render its focused option row with the injected cursor style.

#### Scenario: Cursor row uses the cursor style
- **WHEN** the Select dialog renders its focused option
- **THEN** that row uses the injected cursor style, distinct from the unfocused options

### Requirement: Sort window caption carries its prompt
The sort window SHALL use the caption `Sort: Select a column` and SHALL NOT render the instruction as a body line.

#### Scenario: Caption shows the instruction
- **WHEN** the sort window is opened
- **THEN** its caption reads `Sort: Select a column` and no body line repeats the instruction
