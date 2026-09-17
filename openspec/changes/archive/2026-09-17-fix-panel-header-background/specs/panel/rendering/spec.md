## Purpose

Defines how the file panel presents its listing visually, so the panel reads as one continuous themed surface rather than a patchwork of differently styled regions.

## ADDED Requirements

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
