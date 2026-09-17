## ADDED Requirements

### Requirement: Cursor uses the injected cursor style

The panel SHALL render its selection cursor using the injected cursor style rather than a built-in style, so a theme — including a `styles.css` override — controls the cursor's colors and weight. The cursor SHALL remain a full-width bar across the selected row, and a blurred panel SHALL still render no cursor highlight.

#### Scenario: Cursor follows the injected style

- **WHEN** a focused panel renders its selection cursor
- **THEN** the cursor uses the colors and weight of the injected cursor style

#### Scenario: Cursor follows a theme override

- **WHEN** the injected cursor style defines colors different from the embedded default
- **THEN** the rendered cursor uses those colors

#### Scenario: Blurred panel renders no cursor

- **WHEN** a panel is not focused
- **THEN** its selected row is rendered without the cursor highlight
