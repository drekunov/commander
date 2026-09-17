## ADDED Requirements

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

## MODIFIED Requirements

### Requirement: Select dialog returns a chosen option
The select dialog SHALL present a list of options, let the user choose one, and return the chosen option. The dialog SHALL be cancelable with Escape, in which case it returns no option.

#### Scenario: Single selection made
- **WHEN** the user chooses one of the presented options and confirms
- **THEN** the call returns exactly that option

#### Scenario: Selection canceled
- **WHEN** the user presses Escape
- **THEN** the dialog closes and returns no option
