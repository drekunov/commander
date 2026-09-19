## MODIFIED Requirements

### Requirement: Select dialog cursor uses the cursor style
The Select dialog SHALL render its focused option row with the injected dialog cursor style and every other option row with the injected dialog option style. The dialog cursor style SHALL be a themeable style distinct from the file panel cursor, so the sort window can be themed independently of the panels.

#### Scenario: Cursor row uses the cursor style
- **WHEN** the Select dialog renders its focused option
- **THEN** that row uses the injected dialog cursor style, distinct from the unfocused options

#### Scenario: Unselected options use the option style
- **WHEN** the Select dialog renders an option that is not focused
- **THEN** that row uses the injected dialog option style

## ADDED Requirements

### Requirement: Dialog body is painted with the dialog background
Every dialog SHALL render its body text, inputs, options, and controls on the injected dialog frame background, so the dialog content is continuous with its frame rather than showing the terminal's default background behind text.

#### Scenario: Dialog content shares the frame background
- **WHEN** an Info, Input, Confirm, or Select dialog renders
- **THEN** each body line draws the dialog frame background behind its content

### Requirement: Select options are left-aligned with row padding
The Select dialog SHALL render each option row starting at the left edge of the dialog body, with no `>` cursor prefix. Each option row SHALL use the injected option padding, so the option text is inset from the row edges. The focused option SHALL be indicated only by the dialog cursor style. In multi-select mode each option SHALL begin with its `[x]`/`[ ]` checkbox before the option text.

#### Scenario: No cursor marker
- **WHEN** the Select dialog renders its options
- **THEN** no option row begins with a `>` marker

#### Scenario: Options are flushed left
- **WHEN** the Select dialog renders a single-select option
- **THEN** the option row starts at the left edge of the dialog body

#### Scenario: Option rows use the injected padding
- **WHEN** the injected option style sets padding
- **THEN** each option row renders with that padding on its left and right, and the padding carries the row's background

#### Scenario: Multi-select keeps its checkbox
- **WHEN** the Select dialog renders a multi-select option
- **THEN** the option row begins with its `[x]` or `[ ]` checkbox before the option text

#### Scenario: Focus is shown by style
- **WHEN** an option is focused
- **THEN** the row is distinguished by the dialog cursor style rather than by a marker character
