## ADDED Requirements

### Requirement: CSS style classes cover the top menu
The CSS engine SHALL expose style classes for the top menu bar so a theme can restyle it: a class for the bar background, a class for inactive captions, a class for the active caption, a class for the hotkey letter, a class for the pull-down background, and a class for the focused pull-down item. Each class SHALL be loaded into its own field on the injected `Styles` value, and no top-menu class SHALL be required for the others to work.

#### Scenario: Top-menu selectors load into style fields
- **WHEN** the stylesheet defines the top-menu classes
- **THEN** `LoadStyles` populates the matching top-menu style fields

#### Scenario: Working-directory override restyles the top menu
- **WHEN** a working-directory `styles.css` sets the top-menu classes to colors different from the embedded theme
- **THEN** the loaded styles reflect the override and the top menu renders with those colors

#### Scenario: Missing class falls back
- **WHEN** a stylesheet omits one or more top-menu classes
- **THEN** the remaining top-menu classes still load and the menu stays renderable
