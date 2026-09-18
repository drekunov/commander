## ADDED Requirements

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
