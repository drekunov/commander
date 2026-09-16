## Purpose

Provides the full set of modal dialog interactions so that every method on the dialog contract behaves as documented and no dialog can hang the application.

## ADDED Requirements

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
The select dialog SHALL present a list of options, let the user choose one, and return the chosen option.

#### Scenario: Single selection made
- **WHEN** the user chooses one of the presented options and confirms
- **THEN** the call returns exactly that option

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
Every dialog call SHALL return promptly when the application is quitting, independent of user input.

#### Scenario: Application quits while a dialog is open
- **WHEN** the user quits the application while any dialog is open
- **THEN** the dialog call returns and the application exits without deadlock
