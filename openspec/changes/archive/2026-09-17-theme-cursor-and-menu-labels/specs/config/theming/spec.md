## ADDED Requirements

### Requirement: CSS style properties include bold text

The CSS engine SHALL recognize the `font-weight: bold` declaration and apply bold text to the resulting style, so a theme can express bold through the same `styles.css` rules it uses for colors, decoration, spacing, and borders.

#### Scenario: Bold declaration applies bold text

- **WHEN** a style rule sets `font-weight: bold`
- **THEN** the applied style renders its text in bold

#### Scenario: Bold follows a working-directory override

- **WHEN** a `styles.css` override sets `font-weight: bold` on a class
- **THEN** widgets using that class render bold text
