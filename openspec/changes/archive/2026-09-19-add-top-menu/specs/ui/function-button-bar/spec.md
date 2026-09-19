## MODIFIED Requirements

### Requirement: Non-quit actions are mocks that report their name
Activating any button from 1Help through 8Delete, except the Menu button, SHALL run a placeholder mock action that reports which action was requested without performing a real file operation, and SHALL NOT modify any file, directory, panel data, or the application state. The Menu button SHALL instead open the panel sort window. The 9PullDn button SHALL instead activate the top menu bar.

#### Scenario: Mock action reports its name
- **WHEN** the user activates a non-quit button such as 5Copy
- **THEN** the application reports a placeholder result naming the requested action, and no file or directory is created, changed, or removed

#### Scenario: Mock activation is safe on an empty selection
- **WHEN** the user activates a non-quit button while no file entry is selected
- **THEN** the mock action still reports its placeholder result and the application does not crash or error

#### Scenario: Menu selects the sort mode
- **WHEN** the user activates the Menu button
- **THEN** the focused panel's sort window opens and no mock dialog is shown

#### Scenario: Pull down opens the top menu
- **WHEN** the user activates the 9PullDn button
- **THEN** the top menu bar activates and no mock dialog is shown
