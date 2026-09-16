## MODIFIED Requirements

### Requirement: Dialogs never hang the application
Every dialog call SHALL return promptly when the application is quitting, independent of user input and independent of whether the caller supplied a cancellable context.

#### Scenario: Application quits while a dialog is open
- **WHEN** the user quits the application while any dialog is open
- **THEN** the dialog call returns and the application exits without deadlock

#### Scenario: Error dialog opened with a non-cancellable context returns on quit
- **WHEN** the application quits while an Error or Warning dialog is open and the caller did not supply a cancellable context
- **THEN** the dialog call returns and the application exits without deadlock
