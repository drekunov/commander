# rendering Specification

## Purpose

Defines how the file panel presents its listing visually, so the panel reads as one continuous themed surface rather than a patchwork of differently styled regions.

## Requirements

### Requirement: Column header shares the panel background
The panel SHALL render its column header row with the same background as the table body, taking that background from the injected table style, so no header cell or its padding shows a background different from the listing below. This SHALL hold whether the panel is focused or unfocused. Column titles, column widths, and the header's bold text SHALL be unchanged.

#### Scenario: Header background matches the body
- **WHEN** a panel renders a listing
- **THEN** every column header cell and its padding is painted with the panel's table background, the same background as the entry rows

#### Scenario: Header follows a theme override
- **WHEN** the injected table style defines a background different from the embedded default
- **THEN** the header row uses that same background and remains continuous with the entry rows

#### Scenario: Unfocused panel header matches
- **WHEN** a panel loses focus and renders without the cursor highlight
- **THEN** its column header row still uses the panel's table background

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
