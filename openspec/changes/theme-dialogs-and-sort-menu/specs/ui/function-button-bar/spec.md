## ADDED Requirements

### Requirement: Mock dialog is titled Unimplemented
When activating a button that is not yet implemented, the system SHALL show the placeholder mock dialog with the title `Unimplemented` and a body that names the requested action.

#### Scenario: Unimplemented action shows the Unimplemented dialog
- **WHEN** the user activates a button such as 3View
- **THEN** the dialog title reads `Unimplemented` and its body names the requested action

#### Scenario: Repeated activation keeps the title
- **WHEN** the user activates another unimplemented button while the mock dialog is open
- **THEN** the dialog stays titled `Unimplemented` and updates its body to the new action

### Requirement: Opening the mock dialog does not block the application
Opening or updating the placeholder mock dialog SHALL NOT block the application: the dialog SHALL be shown and the application SHALL keep rendering and responding to input while it is open.

#### Scenario: Application stays responsive
- **WHEN** an unimplemented button opens the mock dialog
- **THEN** the dialog is shown and the application still responds to other keys such as the quit key

#### Scenario: Re-activation updates in place
- **WHEN** the mock dialog is already tracked and another unimplemented action is reported
- **THEN** the dialog text updates and the application stays responsive
