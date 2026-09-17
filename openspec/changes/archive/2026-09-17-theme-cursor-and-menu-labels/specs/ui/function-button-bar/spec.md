## ADDED Requirements

### Requirement: Button labels use the injected menu-label style

The bar SHALL render each non-focused, non-pressed button's label text using the injected menu-label style, drawn on the bar's strip background, so the labels can be themed to match the cursor. The focused and pressed buttons SHALL keep their distinct active and pressed styles.

#### Scenario: Default labels use the menu-label style on the strip

- **WHEN** the bar renders a button that is neither focused nor pressed
- **THEN** its label text uses the injected menu-label style while the surrounding bar strip keeps the button-bar style

#### Scenario: Focused and pressed labels stay distinct

- **WHEN** a button is focused or pressed
- **THEN** its label uses the active or pressed style rather than the menu-label style

#### Scenario: Labels follow a theme override

- **WHEN** the injected menu-label style defines colors different from the embedded default
- **THEN** the default button labels render with those colors
