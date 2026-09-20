## ADDED Requirements

### Requirement: Name cell is colored by entry kind

The panel SHALL render the **Name** cell of each entry using the injected style for that entry's file kind, leaving the entry's other cells and the column header unchanged. The style SHALL be taken from the kind supplied on the entry as a hidden attribute; an entry whose kind is absent, empty, or has no injected style SHALL render its Name cell in the default table style. Styling the Name cell SHALL NOT change column alignment, widths, or the rest of the row.

#### Scenario: Name uses the kind style

- **WHEN** a panel renders an entry whose kind has an injected style
- **THEN** that entry's Name cell renders with the kind's style and its other cells render unchanged

#### Scenario: Kinds map to their classes

- **WHEN** a panel renders directory, executable, symlink, image, archive, source, and config entries
- **THEN** each Name cell uses the style injected for its own kind

#### Scenario: Unclassified entry keeps the default style

- **WHEN** a panel renders an entry with no kind, an unknown kind, or a kind whose class is not injected
- **THEN** its Name cell renders in the default table style

#### Scenario: Parent row is a directory

- **WHEN** a panel renders the synthetic `..` parent row
- **THEN** its Name cell uses the injected directory style

#### Scenario: Other connectors are unaffected

- **WHEN** a panel renders a listing whose entries carry no kind attribute
- **THEN** every Name cell renders in the default table style

#### Scenario: Working-directory override restyles names

- **WHEN** a working-directory `styles.css` sets a file-kind class to a color different from the embedded theme
- **THEN** the matching entries render their Name cells with that color

### Requirement: Cursor row overrides the kind color

While a panel is focused, its cursor row SHALL render with the injected cursor style uniformly, including its Name cell, so the kind color does not break the highlight. When the panel is blurred, the cursor style SHALL NOT apply and the kind colors SHALL show on every row.

#### Scenario: Focused cursor row is uniform

- **WHEN** a focused panel renders the row under its cursor
- **THEN** that row's Name cell uses the cursor style rather than the kind style

#### Scenario: Blurred panel shows kind colors

- **WHEN** a panel is not focused
- **THEN** every row, including the formerly cursor row, renders its Name cell with its kind style
