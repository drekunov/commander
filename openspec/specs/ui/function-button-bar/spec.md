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

### Requirement: Buttons are reachable and activate by mouse click
The system SHALL accept a mouse click on a rendered button as an activation of that button, and SHALL support keyboard focus navigation across the bar so the focused button can be activated with Enter.

#### Scenario: Click on a button activates it
- **WHEN** the user clicks on a rendered button
- **THEN** that button is activated and the click is not treated as a window or panel action

#### Scenario: Keyboard navigation between buttons
- **WHEN** the bar has keyboard focus
- **THEN** left and right arrow keys move the focus between adjacent buttons and pressing Enter activates the focused button

### Requirement: Activated buttons give visible feedback
When a button is activated, the system SHALL render it in its pressed state at the moment of activation and SHALL clear that pressed state when the activation is dispatched, so the button returns to its normal style before or as the consequence of the action appears. A button whose action opens a dialog SHALL NOT remain pressed while the dialog is open.

#### Scenario: Pressed state shown on activation
- **WHEN** a button is activated by key or mouse
- **THEN** the button is shown in its pressed styling at the moment of activation

#### Scenario: Pressed state ends when the action is dispatched
- **WHEN** a button is activated and its action has been dispatched
- **THEN** the button is rendered in its normal style and is no longer shown as pressed

#### Scenario: Button does not stay pressed while its dialog is open
- **WHEN** a button such as Menu opens a dialog
- **THEN** the button is rendered unpressed for as long as the dialog remains open

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

### Requirement: Mock actions expose a replaceable action contract
The system SHALL route every button activation through a single action-dispatch mechanism so that replacing a mock behavior with a real implementation does not change how buttons are rendered, focused, or activated.

#### Scenario: Mock handler can be swapped for a real one
- **WHEN** a mock action's behavior is replaced by a real action implementation
- **THEN** the button's label, focus handling, and key and mouse activation continue to work unchanged

### Requirement: Button labels use the injected menu-label style
The bar SHALL render each button's label split into its numeric prefix and its name. For a button that is neither focused nor pressed, the numeric prefix SHALL use the injected menu-number style and the name SHALL use the injected menu-label style, both drawn on the bar's menu-background base. When the button is focused or pressed, the numeric prefix SHALL keep the menu-number style while the name SHALL use the active or pressed style. The focused and pressed buttons SHALL keep their distinct active and pressed styles for the name.

#### Scenario: Default labels use the menu-label style on the strip
- **WHEN** the bar renders a button that is neither focused nor pressed
- **THEN** its name uses the menu-label style and its numeric prefix uses the menu-number style, while the surrounding bar background keeps the menu-background style

#### Scenario: Focused and pressed labels stay distinct
- **WHEN** a button is focused or pressed
- **THEN** its numeric prefix still uses the menu-number style and its name uses the active or pressed style rather than the menu-label style

#### Scenario: Labels follow a theme override
- **WHEN** the injected menu-number or menu-label style defines colors different from the embedded default
- **THEN** the number and name render with those colors

### Requirement: Bar renders a full-width themed background
The bar SHALL paint its full-width background using the injected menu-background style across the bottom row, falling back to the button-bar style for any attribute the menu-background style leaves unset. The ten buttons SHALL be drawn on top of it.

#### Scenario: Background strip spans the bar row
- **WHEN** the bar is rendered
- **THEN** the full width of the bottom row is painted with the menu-background style and the ten buttons overlay it

### Requirement: Mock activation does not stack dialogs
Repeatedly activating non-quit buttons while a mock result dialog is already open SHALL NOT open additional dialogs, and SHALL NOT accumulate idle goroutines.

#### Scenario: Repeated activation shows a single mock dialog
- **WHEN** the user activates a non-quit button several times in quick succession
- **THEN** at most one mock result dialog is open at a time and no dialog or goroutine accumulates

### Requirement: Buttons have equal width
The bar SHALL lay its ten buttons out with equal width. Each button SHALL render its function-key number in a two-cell field followed by its name in a six-cell field, and any columns left over after the ten buttons SHALL be split into equal gaps between adjacent buttons so the row spans the full terminal width. When the terminal is too narrow for ten eight-cell buttons, the buttons SHALL shrink to equal widths and labels SHALL be truncated to fit.

#### Scenario: All buttons share one width
- **WHEN** the bar renders at a width that fits ten eight-cell buttons
- **THEN** every button occupies the same number of columns, with its number in a two-cell field and its name in a six-cell field

#### Scenario: Leftover columns become equal gaps
- **WHEN** the terminal width leaves columns unused after the ten buttons
- **THEN** those columns are split into equal gaps between adjacent buttons and the row still spans the full width

#### Scenario: Narrow terminal shrinks buttons equally
- **WHEN** the terminal is too narrow for ten eight-cell buttons
- **THEN** the buttons shrink to equal widths and each label is truncated to fit
