## ADDED Requirements

### Requirement: Buttons have equal width
The bar SHALL lay its ten buttons out with equal width. Each button SHALL render its function-key number in a two-cell field followed by its name in a six-cell field, and any columns left over after the ten buttons SHALL be split into equal gaps between adjacent buttons so the row spans the full terminal width. When the terminal is too narrow for ten eight-cell buttons, the buttons SHALL shrink to equal widths and labels SHALL be truncated to fit.

#### Scenario: All buttons share one width
- **WHEN** the bar renders at a width that fits ten eight-cell buttons
- **THEN** every button occupies the same number of columns, with its number in a two-cell field and its name in a six-cell field

#### Scenario: Leftover columns become equal gaps
- **WHEN** the terminal width leaves columns unused after the ten buttons
- **THEN** those columns are split into equal gaps between adjacent buttons and the row still spans the full width

#### Scenario: Narrow terminal shrinks buttons equally
- **WHEN** the terminal is too narrow for ten eight-cell buttons
- **THEN** the buttons shrink to equal widths and each label is truncated to fit
