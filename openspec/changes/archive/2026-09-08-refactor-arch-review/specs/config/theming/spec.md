## Purpose

Provides theme styling to widgets through explicit injection instead of a process-wide mutable singleton, so styles are testable, replaceable, and instance-scoped.

## ADDED Requirements

### Requirement: Styles are injected into widgets
Every styled widget SHALL receive its styles through its constructor or an explicit setter; widgets MUST NOT read style state from a process-global mutable variable.

#### Scenario: Widget styled via injection
- **WHEN** a widget is constructed with a given set of styles
- **THEN** the widget renders using exactly those styles and no other source of style state

#### Scenario: Two widgets with different styles
- **WHEN** two widgets are constructed with different style sets
- **THEN** each renders with its own styles and neither affects the other

### Requirement: No ambient mutable style state
The styling system SHALL NOT expose a mutable, globally-assignable style variable that any package can read or overwrite.

#### Scenario: Style state cannot be mutated globally
- **WHEN** one widget is restyled
- **THEN** no other widget's rendering changes as a side effect

### Requirement: Startup override from working directory
The styling system SHALL continue to support overriding the embedded default theme with a `styles.css` file in the process working directory, resolved once at startup.

#### Scenario: Override file present at startup
- **WHEN** the process starts in a directory containing a `styles.css` file
- **THEN** the styles used by the application reflect that file's rules

#### Scenario: No override file
- **WHEN** the process starts without a `styles.css` file in the working directory
- **THEN** the embedded default theme is used
