## Purpose

Provides a debugging hotkey that captures the exact rendered screen to a timestamped file, so TUI layout and compositing can be inspected and tested outside the live terminal.

## ADDED Requirements

### Requirement: Debug key captures the screen to a file
Pressing the debug key while the application is running SHALL write the complete currently-rendered screen (the exact output the UI would draw, including ANSI styling, window frames, and any open dialogs or overlays) to a file on disk.

#### Scenario: Pressing the debug key
- **WHEN** the user presses the debug key while the UI is rendering
- **THEN** the full current frame is written to a dump file in the process working directory and the UI continues running without interruption

#### Scenario: Frame matches what is displayed
- **WHEN** the debug key is pressed
- **THEN** the dump file content equals the exact string the UI renders for the current frame

### Requirement: Dump file naming
Each screen dump SHALL be written to a uniquely named file whose name identifies it as a screen dump and includes the capture time, so successive dumps never overwrite each other.

#### Scenario: Timestamped unique filename
- **WHEN** a screen dump is written
- **THEN** the file is named `screen-<timestamp>.ansi` where `<timestamp>` encodes the local capture time with second precision

#### Scenario: Multiple dumps in one session
- **WHEN** the debug key is pressed more than once
- **THEN** each dump is written to a distinct file and no earlier dump is overwritten

### Requirement: Dump does not block the UI
Writing the dump file SHALL NOT block rendering or input handling; the UI stays responsive during and after the write.

#### Scenario: Dump during active UI
- **WHEN** a dump is triggered while the UI is running
- **THEN** the write happens asynchronously and key input and rendering continue without noticeable delay

### Requirement: Dump failure does not crash the application
If the dump file cannot be written (e.g. working directory is not writable), the application SHALL continue running normally.

#### Scenario: Unwritable working directory
- **WHEN** the debug key is pressed in a directory where the dump file cannot be created
- **THEN** the application keeps running and does not terminate or enter an error state
