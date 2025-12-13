# ui-presentation Spec Changes

## ADDED Requirements

### Requirement: Format Command Progress Display

The system SHALL display a progress bar with file processing status during the
format command execution.

#### Scenario: Format command shows spinner during file fetching

- **GIVEN** the user runs `nippo format`
- **WHEN** the system is fetching the file list from Google Drive
- **THEN** a spinner SHALL be displayed with the message
  "Fetching file list from Google Drive..."

#### Scenario: Format command shows progress bar during processing

- **GIVEN** files are being processed by the format command
- **WHEN** each file is processed
- **THEN** a progress bar SHALL be displayed showing current/total files
- **AND** the elapsed time SHALL be displayed in parentheses (e.g., "(1m30s)")
- **AND** a spinner SHALL indicate the current file being processed

#### Scenario: Format command shows recent files list

- **GIVEN** files are being processed by the format command
- **WHEN** files complete processing
- **THEN** the last 10 processed files SHALL be displayed
- **AND** each file SHALL show a status icon (success, no-change, or failed)
- **AND** each file SHALL display its Google Drive File ID

#### Scenario: Format command supports cancellation

- **GIVEN** the format command is running
- **WHEN** the user presses Ctrl+C
- **THEN** the operation SHALL stop gracefully
- **AND** "(interrupted)" SHALL be displayed in the progress line

### Requirement: Command Summary Display

The system SHALL display a summary after command completion with file lists
followed by aggregate counts.

#### Scenario: Format command shows summary

- **GIVEN** the format command has completed processing
- **WHEN** the summary is displayed
- **THEN** updated files list SHALL be displayed first (if any)
- **AND** failed files list SHALL be displayed second (if any)
- **AND** the aggregate summary line SHALL be displayed last
- **AND** each file in the lists SHALL include its Google Drive File ID

#### Scenario: Build command shows summary

- **GIVEN** the build command has completed processing
- **WHEN** the summary is displayed
- **THEN** downloaded files list SHALL be displayed first (if any)
- **AND** failed files list SHALL be displayed second (if any)
- **AND** the aggregate summary line SHALL be displayed last
- **AND** each file in the lists SHALL include its Google Drive File ID

## MODIFIED Requirements

### Requirement: Progress Bar Component for Trackable Operations

The system SHALL display a progress bar component when operations can report
progress percentage.

#### Scenario: Progress bar shows percentage

- **GIVEN** a long-running operation reports progress percentage
- **WHEN** the operation is running
- **THEN** a progress bar SHALL be displayed with the current percentage
- **AND** the progress bar SHALL update as the operation progresses

#### Scenario: Progress bar shows elapsed time

- **GIVEN** a progress bar is displayed during an operation
- **WHEN** the operation is running
- **THEN** the elapsed time SHALL be displayed in parentheses
- **AND** the elapsed time SHALL update every second

#### Scenario: Progress bar uses muted green color

- **GIVEN** a progress bar is displayed
- **WHEN** the progress bar is rendered
- **THEN** the filled portion SHALL use OliveDrab color (#6B8E23)

#### Scenario: Fallback to spinner when progress unavailable

- **GIVEN** an operation does not report progress percentage
- **WHEN** the operation is running
- **THEN** a spinner SHALL be displayed instead of a progress bar
