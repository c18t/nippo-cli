# oauth-credentials Specification

## Purpose

TBD - created by archiving change refactor-init-command. Update Purpose after archive.

## Requirements

### Requirement: Data Directory Creation

The init command SHALL create the data directory before attempting to read or
write OAuth credentials files. File providers SHALL NOT create directories.

#### Scenario: Data directory created on init

**Given** the data directory does not exist
**When** the init command attempts to read credentials.json
**Then** the init command SHALL create the data directory with permissions 0755
**And** subsequent file operations SHALL proceed

#### Scenario: Data directory already exists

**Given** the data directory exists
**When** the init command attempts to read credentials.json
**Then** no error SHALL be raised for directory creation
**And** file operations SHALL proceed normally

#### Scenario: Drive provider does not create directories

**Given** the data directory does not exist
**When** the drive file provider attempts to read credentials.json
**Then** the drive file provider SHALL NOT create the data directory
**And** an error message SHALL be displayed instructing the user to run `nippo
init`

### Requirement: Helpful credentials.json Error Messages

The system SHALL provide clear, actionable error messages when credentials.json
is not found, including instructions for obtaining the file.

#### Scenario: credentials.json not found in init command

**Given** credentials.json does not exist in the data directory
**When** the init command attempts to read the file
**Then** an error message SHALL be displayed
**And** the error message SHALL include "credentials.json not found"
**And** the error message SHALL include step-by-step instructions for obtaining
OAuth 2.0 credentials from Google Cloud Console
**And** the error message SHALL include the exact file path where the file
should be placed
**And** the instructions SHALL include:

- Google Cloud Console URL
- Credential type (OAuth 2.0 Client ID)
- Application type (Desktop app)
- Download and save instructions

#### Scenario: credentials.json not found in drive provider

**Given** credentials.json does not exist in the data directory
**When** the drive file provider attempts to read the file
**Then** an error message SHALL be displayed
**And** the error message SHALL include the same helpful instructions as the
init command
**And** the error message SHALL include a note to run `nippo init` to set up
the environment

#### Scenario: credentials.json exists but cannot be read

**Given** credentials.json exists in the data directory
**And** the file cannot be read (e.g., permission denied)
**When** the system attempts to read the file
**Then** a standard file read error SHALL be returned
**And** the error SHALL indicate the specific read failure (not "file not
found")

#### Scenario: Error distinguishes file not found from other errors

**Given** a file operation fails
**When** the error is checked
**Then** the system SHALL use `os.IsNotExist(err)` to distinguish "not found"
from other errors
**And** only "not found" errors SHALL show the helpful setup instructions
**And** other errors SHALL show standard error messages
