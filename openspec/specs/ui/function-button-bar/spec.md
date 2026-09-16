# function-button-bar Specification

## Purpose

Provides a Norton-Commander-style bottom function-key button bar so the commander shows its available actions and captures each one behind a single activation path, with real actions wired in later.

## Requirements

### Requirement: Function-button bar is always visible at the bottom
The system SHALL render a function-button bar as a fixed row at the bottom of the screen, below the file panels, spanning the full terminal width, on every frame while the application runs.

#### Scenario: Bar spans the full width
- **WHEN** the application renders at any supported terminal size
- **THEN** the bottom screen row is occupied by the function-button bar across the full width and no panel or window content overlaps it

#### Scenario: Bar persists across focus changes and dialogs
- **WHEN** the user changes the active window or opens a dialog
- **THEN** the function-button bar remains visible at the bottom of the screen

#### Scenario: Panels yield the bottom row
- **WHEN** the function-button bar is shown
- **THEN** the file panels are sized so their content ends above the bar and the last panel row is not obscured

### Requirement: Bar shows ten numbered function buttons
The bar SHALL display ten buttons, ordered left to right, each labeled with its function-key number and a Norton-Commander-style name: 1Help, 2Menu, 3View, 4Edit, 5Copy, 6RenMov, 7Mkdir, 8Delete, 9PullDn, 10Quit.

#### Scenario: All ten buttons rendered in order
- **WHEN** the application starts and the bar is rendered
- **THEN** all ten buttons appear left to right in the order 1Help through 10Quit

#### Scenario: Bar fits the terminal width
- **WHEN** the terminal is narrower than the buttons' natural combined width
- **THEN** the bar still renders as a single row that ends at the right edge of the screen without wrapping or corrupting the display

### Requirement: Function keys activate their matching button
Pressing a function key F1 through F10 SHALL activate the button with the same number (F1 activates 1Help, ..., F10 activates 10Quit), regardless of which window or panel currently has focus.

#### Scenario: Function key activates its button
- **WHEN** the user presses a function key while any panel is focused
- **THEN** the correspondingly numbered button is activated and the focused window does not also react to that key

#### Scenario: Focused button key map
- **WHEN** the user presses F10
- **THEN** the 10Quit button is activated and the application quits

### Requirement: Buttons are reachable and activate by mouse click
The system SHALL accept a mouse click on a rendered button as an activation of that button, and SHALL support keyboard focus navigation across the bar so the focused button can be activated with Enter.

#### Scenario: Click on a button activates it
- **WHEN** the user clicks on a rendered button
- **THEN** that button is activated and the click is not treated as a window or panel action

#### Scenario: Keyboard navigation between buttons
- **WHEN** the bar has keyboard focus
- **THEN** left and right arrow keys move the focus between adjacent buttons and pressing Enter activates the focused button

### Requirement: Activated buttons give visible feedback
When a button is activated, the system SHALL briefly render it in its pressed state and then render the consequence of the action, so the user sees which button was triggered.

#### Scenario: Pressed state shown on activation
- **WHEN** a button is activated by key or mouse
- **THEN** the button is shown in its pressed styling at the moment of activation

### Requirement: Non-quit actions are mocks that report their name
Activating any button from 1Help through 9PullDn SHALL run a placeholder mock action that reports which action was requested without performing a real file operation, and SHALL NOT modify any file, directory, panel data, or the application state.

#### Scenario: Mock action reports its name
- **WHEN** the user activates a non-quit button such as 5Copy
- **THEN** the application reports a placeholder result naming the requested action, and no file or directory is created, changed, or removed

#### Scenario: Mock activation is safe on an empty selection
- **WHEN** the user activates a non-quit button while no file entry is selected
- **THEN** the mock action still reports its placeholder result and the application does not crash or error

### Requirement: Mock actions expose a replaceable action contract
The system SHALL route every button activation through a single action-dispatch mechanism so that replacing a mock behavior with a real implementation does not change how buttons are rendered, focused, or activated.

#### Scenario: Mock handler can be swapped for a real one
- **WHEN** a mock action's behavior is replaced by a real action implementation
- **THEN** the button's label, focus handling, and key and mouse activation continue to work unchanged
