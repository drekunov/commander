# data-contract Specification

## Purpose

Defines the contract between directory-data sources (`Connector.ReadDir`) and the table panel (`SetData`) so that listings render completely and correctly, and data delivery never races with rendering.

## Requirements

### Requirement: ReadDir returns header row before entry rows
A directory connector's `ReadDir` SHALL return a sequence of attribute rows in which the first row is the column header (attribute names) and every subsequent row is one directory entry. Connectors MUST NOT omit the header row.

#### Scenario: Header row precedes entry rows
- **WHEN** a connector reads a directory containing entries
- **THEN** the returned sequence begins with a header row whose attribute names become column titles, followed by one row per entry

#### Scenario: Empty directory
- **WHEN** a connector reads an empty directory
- **THEN** it returns no rows and no header row, and the caller renders an empty table without error

### Requirement: Panel renders the complete listing
The panel SHALL build column titles from the header row and table rows from all remaining rows, so every directory entry including the first is displayed.

#### Scenario: Full listing rendered
- **WHEN** a listing of N entries is provided with a header row
- **THEN** the table displays exactly N entry rows with the first entry present

#### Scenario: No data provided
- **WHEN** no data rows are provided
- **THEN** the panel leaves the table unchanged and does not panic

### Requirement: Race-free data delivery
The panel SHALL accept listing updates concurrently with UI rendering without exposing partially-updated or corrupted table state, and without data races detectable by the race detector.

#### Scenario: Update while rendering
- **WHEN** a new listing is delivered while the UI is rendering the current table
- **THEN** the panel either renders the previous or the new listing, never a torn mix of both, and no data race is reported

### Requirement: Entry attribute values become cells
Each entry row SHALL be displayed as one table row whose cell values are the entry's visible attribute values in the order declared by the header row.

#### Scenario: Cell ordering follows header
- **WHEN** an entry has visible attributes in the header's declared order
- **THEN** its cell values appear in that same column order

#### Scenario: Hidden attributes are omitted
- **WHEN** an entry carries an attribute marked hidden
- **THEN** no cell is rendered for it and the visible cells stay aligned with the columns

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

### Requirement: Flexible column fills the remaining width
The header row MAY mark one column flexible. The panel SHALL size the flexible column to the width left after the fixed columns and the per-cell padding, and SHALL shrink it when the panel is narrower, down to a minimum of eight cells. Fixed columns SHALL keep their declared widths. When no column is marked flexible, the panel SHALL use every column's declared width.

#### Scenario: Flexible column grows to fill the panel
- **WHEN** the panel is wider than the fixed columns plus padding
- **THEN** the flexible column grows so the row spans the panel width

#### Scenario: Flexible column shrinks on a narrow panel
- **WHEN** the panel is narrower than the columns need
- **THEN** the flexible column shrinks, never below eight cells, while the fixed columns keep their widths

#### Scenario: No flexible column
- **WHEN** no header attribute is marked flexible
- **THEN** every column uses its declared width
