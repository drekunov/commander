# top-menu Specification

## Purpose

Provides a Norton-Commander-style top menu bar: a first-row line of pull-down menus that is activated with F9, navigated with the arrow keys or per-caption hotkey letters, and whose items dispatch the same actions as the function-button bar.

## Requirements

### Requirement: Top menu bar is always visible as the first screen row
The system SHALL render a top menu bar as the first row of the screen, above the file panels, on every frame while the application runs.

#### Scenario: Menus rendered in order
- **WHEN** the application renders at any supported terminal size
- **THEN** the first screen row shows the top-level menus Left, Files, Commands, Options, and Right, left to right, and no other menu

#### Scenario: Panels yield the top row
- **WHEN** the top menu bar is shown
- **THEN** the file panels and any window are sized and positioned so their content starts below the top row and none of them overlaps it

#### Scenario: Bar persists across focus changes and dialogs
- **WHEN** the user changes the focused window or opens a dialog
- **THEN** the top menu bar remains visible as the first screen row

### Requirement: Each menu caption carries one highlighted hotkey letter
Every top-level menu caption SHALL contain exactly one hotkey letter, rendered distinctly from the rest of the caption, using the classic Norton Commander letters: Left shows L, Files shows F, Commands shows C, Options shows O, and Right shows R.

#### Scenario: Hotkey letters present and distinct
- **WHEN** the top menu bar is rendered
- **THEN** each caption's hotkey letter is drawn with the injected menu-hotkey style, apart from the caption's other characters

#### Scenario: Hotkey selects and opens a menu
- **WHEN** the bar is active and the user presses the hotkey letter of a caption, in either case
- **THEN** that menu becomes the active menu and its pull-down opens

#### Scenario: Active highlight spans the whole caption
- **WHEN** a menu is active
- **THEN** the active-caption highlight covers the entire caption, including its hotkey letter, so no caption character is left outside the highlight

#### Scenario: Hotkey is inert while the bar is inactive
- **WHEN** the bar is inactive and the user types a letter
- **THEN** the focused panel or window receives the letter and no menu opens

### Requirement: F9 activates the top menu bar
The 9PullDn function button and the F9 key SHALL activate the top menu bar; while the bar is active, activating it again SHALL deactivate it and any open pull-down SHALL close.

#### Scenario: F9 activates the bar
- **WHEN** the user presses F9 or activates the 9PullDn button while the bar is inactive
- **THEN** the bar becomes active, one top-level menu is highlighted, and the focused panel does not also react to the key

#### Scenario: F9 toggles the bar off
- **WHEN** the bar is active and the user presses F9 again
- **THEN** the bar becomes inactive and any open pull-down closes

#### Scenario: Active bar owns navigation keys
- **WHEN** the bar is active and the user presses Left, Right, Up, Down, or Enter
- **THEN** the top menu handles the key and the focused panel does not navigate

#### Scenario: Dialog blocks menu activation
- **WHEN** a modal dialog is open and the user presses F9 or activates the 9PullDn button
- **THEN** the top menu does not activate and the dialog stays open

### Requirement: Keyboard navigation across menus and pull-down items
While the top menu bar is active, Left and Right SHALL move the active menu among the top-level menus, Down SHALL open the active menu's pull-down when closed and move the focus down when open, Up SHALL move the focus up within the open pull-down, Enter SHALL open the pull-down when closed and activate the focused item when open, and Escape SHALL close the pull-down when open and otherwise deactivate the bar.

#### Scenario: Left and right move between menus
- **WHEN** the bar is active and the user presses Left or Right
- **THEN** the highlight moves to the previous or next top-level menu and any open pull-down follows it

#### Scenario: Down opens and moves in a pull-down
- **WHEN** the bar is active and the user presses Down
- **THEN** the active menu's pull-down opens with its first item focused, and further Down presses move the focus down the item list

#### Scenario: Enter opens and activates
- **WHEN** the user presses Enter with no pull-down open
- **THEN** the active menu's pull-down opens, and pressing Enter again activates the focused item

#### Scenario: Escape unwinds the menu
- **WHEN** a pull-down is open and the user presses Escape
- **THEN** the pull-down closes and the bar stays active, and pressing Escape again deactivates the bar

### Requirement: Pull-down items reuse the function-button action dispatch
Each top-level menu SHALL open a pull-down list of items. Activating a Files item — View, Edit, Copy, RenMov, Mkdir, or Delete — SHALL dispatch the same action as the matching function button (3View through 8Delete). Every other pull-down item SHALL be a placeholder that reports its own label through the existing mock report and SHALL NOT perform a real file operation or modify any file, directory, panel data, or application state. Activating any item SHALL close the pull-down and deactivate the bar.

#### Scenario: Files item dispatches its function-button action
- **WHEN** the user activates the View item under Files
- **THEN** the same action as the 3View function button is dispatched and no separate top-menu-only action runs

#### Scenario: Placeholder item reports its label
- **WHEN** the user activates a placeholder item such as Commands > Find file
- **THEN** the mock report names that item's label and no file or directory is created, changed, or removed

#### Scenario: Focused item does not select the frame
- **WHEN** a pull-down item is focused
- **THEN** the cursor style covers only the item's text, while the box borders keep the pull-down frame style

#### Scenario: Activation closes the menu
- **WHEN** the user activates any pull-down item
- **THEN** the pull-down closes and the bar becomes inactive

### Requirement: Mouse activates captions and pull-down items
The system SHALL accept a mouse click on a rendered caption as activating the bar and opening that menu, and a mouse click on a rendered pull-down item as activating that item. A click on the top menu row SHALL NOT reach the panels or windows.

#### Scenario: Click on a caption opens its menu
- **WHEN** the user clicks a top-level caption
- **THEN** the bar becomes active and that menu's pull-down opens

#### Scenario: Click on an item activates it
- **WHEN** the user clicks a pull-down item
- **THEN** that item is activated as if selected by keyboard

#### Scenario: Click on the top row is consumed
- **WHEN** the user clicks anywhere on the top menu row
- **THEN** no panel or window handles the click

### Requirement: Top menu is styled from injected styles
The top menu bar SHALL render every visible part from injected styles: the bar background, the inactive caption, the active caption, the hotkey letter, the pull-down background, and the focused pull-down item. A `styles.css` override in the process working directory SHALL change any of them.

#### Scenario: Default styles applied
- **WHEN** the top menu bar is rendered with the embedded theme
- **THEN** its background, captions, hotkey letters, pull-down, and focused item use the injected top-menu styles

#### Scenario: Override restyles the menu
- **WHEN** a working-directory `styles.css` sets the top-menu classes to different colors
- **THEN** the bar renders with those colors

#### Scenario: Active caption differs from inactive
- **WHEN** a menu is active
- **THEN** its caption uses the active-caption style and the other captions use the inactive-caption style

#### Scenario: Caption padding is themed
- **WHEN** the theme sets horizontal padding on the caption classes
- **THEN** each caption renders with that padding and the active caption's highlight covers the padding as well as the caption text
