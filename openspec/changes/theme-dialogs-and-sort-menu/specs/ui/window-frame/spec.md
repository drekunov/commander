## Purpose

Renders the shared window frame — border, background, title, and resize grip — from the injected theme, so panels and dialogs share consistent, overridable chrome.

## ADDED Requirements

### Requirement: Window frame chrome uses injected styles
A window SHALL render its frame border and background, its title in the focused and unfocused states, and its resize grip using injected theme styles. The focused title SHALL be visually distinct from the unfocused title, and no frame chrome color SHALL be hardcoded.

#### Scenario: Focused title uses the focused-title style
- **WHEN** a window has focus
- **THEN** its title renders with the injected focused-title style

#### Scenario: Unfocused title uses the unfocused-title style
- **WHEN** a window does not have focus
- **THEN** its title renders with the injected unfocused-title style, distinct from the focused-title style

#### Scenario: Frame border and background use the frame style
- **WHEN** a window renders
- **THEN** its border characters and background use the injected frame style

#### Scenario: Resize grip uses the grip style
- **WHEN** a window renders its bottom border
- **THEN** the resize grip uses the injected grip style
