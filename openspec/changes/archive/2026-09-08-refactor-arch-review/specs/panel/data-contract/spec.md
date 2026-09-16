## Purpose

Defines the contract between directory-data sources (`Connector.ReadDir`) and the table panel (`SetData`) so that listings render completely and correctly, and data delivery never races with rendering.

## ADDED Requirements

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
Each entry row SHALL be displayed as one table row whose cell values are the entry's attribute values in the order declared by the header row.

#### Scenario: Cell ordering follows header
- **WHEN** an entry has attributes in the header's declared order
- **THEN** its cell values appear in that same column order
