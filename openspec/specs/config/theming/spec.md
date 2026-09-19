# theming Specification

## Purpose

Provides theme styling to widgets through explicit injection instead of a process-wide mutable singleton, so styles are testable, replaceable, and instance-scoped.

## Requirements

### Requirement: Styles are injected into widgets
Every styled widget SHALL receive its styles through its constructor or an explicit setter; widgets MUST NOT read style state from a process-global mutable variable.

#### Scenario: Widget styled via injection
- **WHEN** a widget is constructed with a given set of styles
- **THEN** the widget renders using exactly those styles and no other source of style state

#### Scenario: Two widgets with different styles
- **WHEN** two widgets are constructed with different style sets
- **THEN** each renders with its own styles and neither affects the other

### Requirement: No ambient mutable style state
The styling system SHALL NOT expose a mutable, globally-assignable style variable that any package can read or overwrite.

#### Scenario: Style state cannot be mutated globally
- **WHEN** one widget is restyled
- **THEN** no other widget's rendering changes as a side effect

### Requirement: Startup override from working directory
The styling system SHALL continue to support overriding the embedded default theme with a `styles.css` file in the process working directory, resolved once at startup.

#### Scenario: Override file present at startup
- **WHEN** the process starts in a directory containing a `styles.css` file
- **THEN** the styles used by the application reflect that file's rules

#### Scenario: No override file
- **WHEN** the process starts without a `styles.css` file in the working directory
- **THEN** the embedded default theme is used

### Requirement: CSS style properties include bold text
The CSS engine SHALL recognize the `font-weight: bold` declaration and apply bold text to the resulting style, so a theme can express bold through the same `styles.css` rules it uses for colors, decoration, spacing, and borders.

#### Scenario: Bold declaration applies bold text
- **WHEN** a style rule sets `font-weight: bold`
- **THEN** the applied style renders its text in bold

#### Scenario: Bold follows a working-directory override
- **WHEN** a `styles.css` override sets `font-weight: bold` on a class
- **THEN** widgets using that class render bold text

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

### Requirement: Theme classes cover the dialog set and frame chrome
The styling system SHALL expose theme classes for the dialog body text, the dialog input, the dialog buttons, the dialog menu options (selected and unselected), and the window frame chrome (frame border and background, focused title, unfocused title, and resize grip), so a working-directory `styles.css` can restyle every visible part of a dialog.

#### Scenario: Override restyles the dialog body
- **WHEN** a working-directory `styles.css` sets a style on the dialog body class
- **THEN** the dialog set renders its body with that style

#### Scenario: Override restyles the menu options
- **WHEN** a working-directory `styles.css` sets styles on the selected and unselected menu option classes
- **THEN** the sort window renders its option rows with those styles

#### Scenario: Override restyles the frame chrome
- **WHEN** a working-directory `styles.css` sets styles on the frame, title, or grip classes
- **THEN** windows render their frame, title, and resize grip with those styles
