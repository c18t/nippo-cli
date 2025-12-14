# ui-presentation Specification

## Purpose

Define the terminal UI components and presentation patterns for nippo-cli,
using Bubbletea for interactive elements like spinners, progress bars, and
text inputs.

## Requirements

### Requirement: TUI Framework Integration

The system SHALL use Bubbletea as the TUI framework for rendering interactive
terminal UI components.

#### Scenario: Bubbletea program lifecycle

- **GIVEN** a CLI command is executed
- **WHEN** UI interaction is required
- **THEN** a Bubbletea program SHALL be started with appropriate model
- **AND** the program SHALL handle user input through the Update function
- **AND** the program SHALL render output through the View function

### Requirement: Spinner Component for Progress Indication

The system SHALL display a spinner component when long-running operations are
in progress and progress percentage is not available.

#### Scenario: Build command shows spinner

- **GIVEN** the user runs `nippo build`
- **WHEN** the build process is running
- **THEN** a spinner SHALL be displayed with the current operation message
- **AND** the spinner SHALL stop when the operation completes

#### Scenario: Deploy command shows spinner

- **GIVEN** the user runs `nippo deploy`
- **WHEN** the deploy process is running
- **THEN** a spinner SHALL be displayed with the current operation message

#### Scenario: Clean command shows spinner

- **GIVEN** the user runs `nippo clean`
- **WHEN** the clean process is running
- **THEN** a spinner SHALL be displayed with the current operation message

#### Scenario: Update command shows spinner

- **GIVEN** the user runs `nippo update`
- **WHEN** the update process is running
- **THEN** a spinner SHALL be displayed with the current operation message

#### Scenario: Init command shows spinner during authentication

- **GIVEN** the user runs `nippo init`
- **WHEN** the Google Drive authentication process is running
- **THEN** a spinner SHALL be displayed with the current operation message
- **AND** the spinner SHALL NOT be displayed during text input prompts

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

### Requirement: Text Input Component for User Input

The system SHALL use Bubbletea text input component for collecting user input
during interactive prompts.

#### Scenario: Init command prompts for project URL

- **GIVEN** the user runs `nippo init`
- **WHEN** the system prompts for project repository URL
- **THEN** a text input component SHALL be displayed
- **AND** the input SHALL have a default value
- **AND** the user SHALL be able to edit or accept the default

#### Scenario: Init command prompts for template path

- **GIVEN** the user is in the init flow
- **WHEN** the system prompts for template path
- **THEN** a text input component SHALL be displayed with default `/templates`

#### Scenario: Init command prompts for asset path

- **GIVEN** the user is in the init flow
- **WHEN** the system prompts for asset path
- **THEN** a text input component SHALL be displayed with default `/output`

### Requirement: Styled Output with Lipgloss

The system SHALL use Lipgloss for consistent styling of terminal output.

#### Scenario: Success messages are styled

- **GIVEN** an operation completes successfully
- **WHEN** the completion message is displayed
- **THEN** the message SHALL be styled with success colors

#### Scenario: Error messages are styled

- **GIVEN** an operation fails
- **WHEN** the error message is displayed
- **THEN** the message SHALL be styled with error colors

### Requirement: View Provider Architecture

The system SHALL maintain the existing View Provider pattern while integrating
Bubbletea components.

#### Scenario: View provider handles view model

- **GIVEN** a presenter creates a view model
- **WHEN** the view provider receives the view model
- **THEN** the appropriate Bubbletea component SHALL be rendered
- **AND** user input SHALL be communicated back through channels

### Requirement: Preserve UI Output on Termination

The system SHALL NOT clear UI output when operations complete or are
interrupted. The final state of the UI component SHALL remain visible after
program termination.

#### Scenario: Spinner shows final state on completion

- **GIVEN** a spinner is displayed during an operation
- **WHEN** the operation completes successfully
- **THEN** the spinner output SHALL NOT be cleared from the terminal
- **AND** a success message SHALL be displayed on a new line

#### Scenario: Spinner shows final state on interruption

- **GIVEN** a spinner is displayed during an operation
- **WHEN** the user interrupts with Ctrl-C
- **THEN** the spinner output SHALL NOT be cleared from the terminal
- **AND** an interruption message SHALL be displayed on a new line

#### Scenario: Progress bar shows final state on completion

- **GIVEN** a progress bar is displayed during an operation
- **WHEN** the operation completes
- **THEN** the progress bar output SHALL NOT be cleared from the terminal
- **AND** the progress bar SHALL show 100% completion

#### Scenario: Text input preserves entered value

- **GIVEN** a text input prompt is displayed
- **WHEN** the user submits a value
- **THEN** the prompt and entered value SHALL remain visible in the terminal

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
