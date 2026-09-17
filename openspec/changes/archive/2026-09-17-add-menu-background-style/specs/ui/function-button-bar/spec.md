## MODIFIED Requirements

### Requirement: Bar renders a full-width themed background
The bar SHALL paint its full-width background using the injected menu-background style across the bottom row, falling back to the button-bar style for any attribute the menu-background style leaves unset. The ten buttons SHALL be drawn on top of it.

#### Scenario: Background strip spans the bar row
- **WHEN** the bar is rendered
- **THEN** the full width of the bottom row is painted with the menu-background style and the ten buttons overlay it

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
