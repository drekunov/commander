## ADDED Requirements

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
