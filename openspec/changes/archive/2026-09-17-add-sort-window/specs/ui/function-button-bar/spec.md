## MODIFIED Requirements

### Requirement: Non-quit actions are mocks that report their name
Activating any button from 1Help through 9PullDn, except the Menu button, SHALL run a placeholder mock action that reports which action was requested without performing a real file operation, and SHALL NOT modify any file, directory, panel data, or the application state. The Menu button SHALL instead open the panel sort window.

#### Scenario: Mock action reports its name
- **WHEN** the user activates a non-quit button such as 5Copy
- **THEN** the application reports a placeholder result naming the requested action, and no file or directory is created, changed, or removed

#### Scenario: Mock activation is safe on an empty selection
- **WHEN** the user activates a non-quit button while no file entry is selected
- **THEN** the mock action still reports its placeholder result and the application does not crash or error

#### Scenario: Menu selects the sort mode
- **WHEN** the user activates the Menu button
- **THEN** the focused panel's sort window opens and no mock dialog is shown

### Requirement: Function keys activate their matching button
Pressing a function key F1 through F10 SHALL activate the button with the same number (F1 activates 1Help, ..., F10 activates 10Quit), regardless of which window or panel currently has focus, except while a dialog is open: then F1 through F9 SHALL NOT activate a button. F10 SHALL always quit.

#### Scenario: Function key activates its button
- **WHEN** the user presses a function key while any panel is focused
- **THEN** the correspondingly numbered button is activated and the focused window does not also react to that key

#### Scenario: Focused button key map
- **WHEN** the user presses F10
- **THEN** the 10Quit button is activated and the application quits

#### Scenario: Dialog blocks function keys
- **WHEN** a dialog is open and the user presses F1 through F9
- **THEN** no button is activated and the dialog stays open
