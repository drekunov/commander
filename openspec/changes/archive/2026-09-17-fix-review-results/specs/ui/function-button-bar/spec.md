## ADDED Requirements

### Requirement: Bar renders a full-width themed background
The bar SHALL paint its full-width background strip using the configured button-bar style across the bottom row, with the ten buttons drawn on top of it.

#### Scenario: Background strip spans the bar row
- **WHEN** the bar is rendered
- **THEN** the full width of the bottom row is painted with the bar's background style and the ten buttons overlay it

### Requirement: Mock activation does not stack dialogs
Repeatedly activating non-quit buttons while a mock result dialog is already open SHALL NOT open additional dialogs, and SHALL NOT accumulate idle goroutines.

#### Scenario: Repeated activation shows a single mock dialog
- **WHEN** the user activates a non-quit button several times in quick succession
- **THEN** at most one mock result dialog is open at a time and no dialog or goroutine accumulates
