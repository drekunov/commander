# dialogs Specification

## Purpose

Provides the full set of modal dialog interactions so that every method on the dialog contract behaves as documented and no dialog can hang the application.

## Requirements

### Requirement: Info dialog
The dialog system SHALL display a message with an OK control and remain open until the user confirms it or the provided context is cancelled.

#### Scenario: User confirms info dialog
- **WHEN** the user presses the confirm key on an open info dialog
- **THEN** the dialog closes and the call returns

#### Scenario: Context cancelled
- **WHEN** the surrounding context is cancelled while the info dialog is open
- **THEN** the call returns without waiting for user input

### Requirement: Input dialog returns entered text
The dialog system SHALL display a text field and return the text the user entered when confirmed.

#### Scenario: User submits input
- **WHEN** the user types text and confirms the input dialog
- **THEN** the call returns the typed text

### Requirement: Password dialog masks input
The password dialog SHALL behave like the input dialog but MUST NOT echo the entered characters.

#### Scenario: Masked password entry
- **WHEN** the user types into a password dialog
- **THEN** the entered characters are not displayed and the typed value is returned on confirm

### Requirement: Select dialog returns a chosen option
The select dialog SHALL present a list of options, let the user choose one, and return the chosen option. The dialog SHALL be cancelable with Escape, in which case it returns no option.

#### Scenario: Single selection made
- **WHEN** the user chooses one of the presented options and confirms
- **THEN** the call returns exactly that option

#### Scenario: Selection canceled
- **WHEN** the user presses Escape
- **THEN** the dialog closes and returns no option

### Requirement: Dialogs are modal
While a dialog window is open, it SHALL capture input: function keys F1 through F9 and mouse clicks outside the dialog SHALL be ignored, so neither the function-button bar nor the panels react. The quit key (F10) SHALL still quit the application.

#### Scenario: Function keys are ignored
- **WHEN** a dialog is open and the user presses a function key F1 through F9
- **THEN** no button is activated and the dialog stays open

#### Scenario: Mouse clicks outside the dialog are ignored
- **WHEN** a dialog is open and the user clicks outside it
- **THEN** no other window is focused and the dialog stays open

#### Scenario: Quit still works
- **WHEN** a dialog is open and the user presses F10
- **THEN** the application quits

### Requirement: Multi-select dialog returns chosen options
The multi-select dialog SHALL present a list of options, let the user choose any subset, and return all chosen options.

#### Scenario: Multiple options chosen
- **WHEN** the user selects several of the presented options and confirms
- **THEN** the call returns exactly the chosen options

#### Scenario: Nothing chosen
- **WHEN** the user confirms without selecting anything
- **THEN** the call returns an empty list

### Requirement: Confirm dialog returns a boolean decision
The confirm dialog SHALL present a yes/no choice and return whether the user accepted.

#### Scenario: User accepts
- **WHEN** the user accepts the prompt
- **THEN** the call returns true

#### Scenario: User declines
- **WHEN** the user declines the prompt
- **THEN** the call returns false

### Requirement: Dialogs never hang the application
Every dialog call SHALL return promptly when the application is quitting, independent of user input and independent of whether the caller supplied a cancellable context.

#### Scenario: Application quits while a dialog is open
- **WHEN** the user quits the application while any dialog is open
- **THEN** the dialog call returns and the application exits without deadlock

#### Scenario: Error dialog opened with a non-cancellable context returns on quit
- **WHEN** the application quits while an Error or Warning dialog is open and the caller did not supply a cancellable context
- **THEN** the dialog call returns and the application exits without deadlock

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
