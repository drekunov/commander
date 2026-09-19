## ADDED Requirements

### Requirement: Hidden attributes are not columns
A connector MAY mark an attribute hidden. The panel SHALL exclude hidden attributes from the table columns and from the rendered rows, while retaining them in the entry data so behavior such as directory detection still works.

#### Scenario: Hidden attribute omitted from the table
- **WHEN** the header row marks an attribute hidden
- **THEN** the table shows no column for it and no cell is rendered for it

#### Scenario: Hidden attribute still drives behavior
- **WHEN** an entry carries a hidden attribute
- **THEN** the panel can still read it to decide behavior such as entering a directory

### Requirement: Header may declare column widths
The header row MAY declare a width for an attribute. The panel SHALL use that width for the attribute's column, falling back to its default width when none is declared.

#### Scenario: Declared width is used
- **WHEN** the header declares a width for a column
- **THEN** the table uses that width for the column

#### Scenario: Default width when none declared
- **WHEN** the header declares no width for a column
- **THEN** the table uses the panel's default column width

## MODIFIED Requirements

### Requirement: Entry attribute values become cells
Each entry row SHALL be displayed as one table row whose cell values are the entry's visible attribute values in the order declared by the header row.

#### Scenario: Cell ordering follows header
- **WHEN** an entry has visible attributes in the header's declared order
- **THEN** its cell values appear in that same column order

#### Scenario: Hidden attributes are omitted
- **WHEN** an entry carries an attribute marked hidden
- **THEN** no cell is rendered for it and the visible cells stay aligned with the columns
