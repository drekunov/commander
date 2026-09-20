## ADDED Requirements

### Requirement: File-kind theme classes

The styling system SHALL expose a style class per file kind and load each into its own field on the injected `Styles` value: `.file-directory`, `.file-executable`, `.file-symlink`, `.file-image`, `.file-archive`, `.file-source`, and `.file-config`. No file-kind class SHALL be required for the others to work, and the embedded Norton Commander Classic theme SHALL define all seven with distinct foreground colors so the kinds are visually distinguishable on the panel background.

#### Scenario: File-kind selectors load into style fields

- **WHEN** the stylesheet defines the file-kind classes
- **THEN** `LoadStyles` populates the matching file-kind style fields

#### Scenario: Working-directory override restyles file kinds

- **WHEN** a working-directory `styles.css` sets a file-kind class to a color different from the embedded theme
- **THEN** the loaded style for that kind reflects the override

#### Scenario: Missing class falls back

- **WHEN** a stylesheet omits one or more file-kind classes
- **THEN** the remaining file-kind classes still load and the panel stays renderable

#### Scenario: Embedded theme distinguishes kinds

- **WHEN** the embedded theme is loaded
- **THEN** each file-kind class defines a foreground color and the classes are distinct from one another
