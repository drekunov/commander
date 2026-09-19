## MODIFIED Requirements

### Requirement: Activated buttons give visible feedback
When a button is activated, the system SHALL render it in its pressed state at the moment of activation and SHALL clear that pressed state when the activation is dispatched, so the button returns to its normal style before or as the consequence of the action appears. A button whose action opens a dialog SHALL NOT remain pressed while the dialog is open.

#### Scenario: Pressed state shown on activation
- **WHEN** a button is activated by key or mouse
- **THEN** the button is shown in its pressed styling at the moment of activation

#### Scenario: Pressed state ends when the action is dispatched
- **WHEN** a button is activated and its action has been dispatched
- **THEN** the button is rendered in its normal style and is no longer shown as pressed

#### Scenario: Button does not stay pressed while its dialog is open
- **WHEN** a button such as Menu opens a dialog
- **THEN** the button is rendered unpressed for as long as the dialog remains open
