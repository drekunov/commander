## MODIFIED Requirements

### Requirement: Button labels use the injected menu-label style
The bar SHALL render each button's label split into its numeric prefix and its name. For a button that is neither focused nor pressed, the numeric prefix SHALL use the injected menu-number style and the name SHALL use the injected menu-label style, both drawn on the bar's strip background. When the button is focused or pressed, the numeric prefix SHALL keep the menu-number style while the name SHALL use the active or pressed style. The focused and pressed buttons SHALL keep their distinct active and pressed styles for the name.

#### Scenario: Default labels use the menu-label style on the strip
- **WHEN** the bar renders a button that is neither focused nor pressed
- **THEN** its name uses the menu-label style and its numeric prefix uses the menu-number style, while the surrounding bar strip keeps the button-bar style

#### Scenario: Focused and pressed labels stay distinct
- **WHEN** a button is focused or pressed
- **THEN** its numeric prefix still uses the menu-number style and its name uses the active or pressed style rather than the menu-label style

#### Scenario: Labels follow a theme override
- **WHEN** the injected menu-number or menu-label style defines colors different from the embedded default
- **THEN** the number and name render with those colors
