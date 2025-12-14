# oauth-flow Specification

## Purpose

Define the OAuth authorization flow for Google Drive access, including
automatic callback handling via local HTTP server and browser integration.

## Requirements

### Requirement: Automatic OAuth Callback Handling

The system SHALL automatically receive OAuth authorization callbacks via a
temporary local HTTP server, eliminating the need for manual code entry.

#### Scenario: Successful automatic authorization

- **GIVEN** the user runs `nippo init`
- **AND** `http://localhost/callback` is configured in Google Cloud Console
- **WHEN** OAuth authorization is required
- **THEN** a local HTTP server SHALL start on an available localhost port
  (21660-21669)
- **AND** the user's browser SHALL automatically open to the Google
  authorization page
- **AND** after user authorization, the callback SHALL be received
  automatically
- **AND** the authorization code SHALL be extracted from the callback
- **AND** the OAuth token SHALL be exchanged and saved
- **AND** the HTTP server SHALL be shut down

#### Scenario: Browser auto-open fails

- **GIVEN** the user runs `nippo init`
- **WHEN** the browser fails to open automatically
- **THEN** the authorization URL SHALL be displayed in the terminal
- **AND** instructions SHALL be shown to manually open the URL
- **AND** the callback server SHALL continue waiting for the callback

#### Scenario: Port fallback when preferred port is in use

- **GIVEN** port 21660 is already in use
- **WHEN** the callback server attempts to start
- **THEN** the server SHALL try the next available port (21661, 21662, etc.)
- **AND** the server SHALL successfully start on the first available port
- **AND** the browser SHALL open with the authorization URL using the selected
  port
- **AND** the callback SHALL be received on the selected port

#### Scenario: All ports exhausted

- **GIVEN** all ports 21660-21669 are already in use
- **WHEN** the callback server attempts to start
- **THEN** all port attempts SHALL fail
- **AND** an error message SHALL be displayed indicating all ports are in use
- **AND** the command SHALL exit with a non-zero status code

#### Scenario: Authorization timeout

- **GIVEN** the callback server is waiting for authorization
- **WHEN** 2 minutes pass without receiving a callback
- **THEN** the server SHALL shut down gracefully
- **AND** a timeout error SHALL be displayed
- **AND** instructions to retry SHALL be shown

#### Scenario: State parameter mismatch (CSRF protection)

- **GIVEN** the callback server is running
- **WHEN** a callback is received with incorrect state parameter
- **THEN** the callback SHALL be rejected
- **AND** an error page SHALL be shown in the browser
- **AND** an error SHALL be reported in the terminal

### Requirement: Browser Integration

The system SHALL attempt to automatically open the user's default browser to
the authorization URL.

#### Scenario: Browser opens on macOS

- **GIVEN** the system is running on macOS
- **WHEN** the authorization URL needs to be opened
- **THEN** the `open` command SHALL be executed with the URL

#### Scenario: Browser opens on Linux

- **GIVEN** the system is running on Linux
- **WHEN** the authorization URL needs to be opened
- **THEN** the `xdg-open` command SHALL be executed with the URL

#### Scenario: Browser opens on Windows

- **GIVEN** the system is running on Windows
- **WHEN** the authorization URL needs to be opened
- **THEN** the `start` command SHALL be executed via cmd with the URL

#### Scenario: Unsupported platform

- **GIVEN** the system is running on an unsupported platform
- **WHEN** automatic browser opening is attempted
- **THEN** the operation SHALL fail gracefully
- **AND** the authorization URL SHALL be displayed for manual opening

### Requirement: Success Feedback

The system SHALL display a success page in the user's browser after successful
authorization.

#### Scenario: Success page displayed

- **GIVEN** the callback server receives a valid authorization code
- **WHEN** the callback request is processed
- **THEN** an HTML success page SHALL be returned to the browser
- **AND** the page SHALL indicate successful authorization
- **AND** the page SHALL instruct the user to return to the terminal

### Requirement: Concurrent Authorization Attempts

The system SHALL support multiple simultaneous authorization attempts by using
different ports.

#### Scenario: Multiple concurrent authorization sessions

- **GIVEN** the callback server is already running on port 21660 in another
  terminal
- **WHEN** a user attempts to run `nippo init` in a second terminal
- **THEN** the second session SHALL automatically use port 21661
- **AND** both sessions SHALL proceed independently
- **AND** each session SHALL complete successfully on its respective port

#### Scenario: Maximum concurrent sessions reached

- **GIVEN** 10 callback servers are already running on ports 21660-21669
- **WHEN** a user attempts to run `nippo init` in an 11th terminal
- **THEN** all port attempts SHALL fail
- **AND** an error message SHALL be displayed indicating all ports are in use
- **AND** the command SHALL exit with a non-zero status code
