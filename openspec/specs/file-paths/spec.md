# file-paths Specification

## Purpose

TBD - created by archiving change refactor-file-structure. Update Purpose after archive.

## Requirements

### Requirement: Configurable File Paths

The system SHALL allow users to configure file paths via the `[paths]` section
in the configuration file.

#### Scenario: Custom data directory

- **WHEN** `paths.data_dir` is set in configuration
- **THEN** the system SHALL use the specified directory as the data directory
- **AND** `templates/`, `assets/`, `credentials.json`, `token.json` SHALL be
  under it
- **AND** environment variables in the path SHALL be expanded

#### Scenario: Custom cache directory

- **WHEN** `paths.cache_dir` is set in configuration
- **THEN** the system SHALL use the specified directory for markdown files and
  build output
- **AND** environment variables in the path SHALL be expanded

#### Scenario: Relative path resolution

- **WHEN** a configured path is relative (not starting with `/` or drive letter)
- **THEN** the path SHALL be resolved relative to the config file directory

#### Scenario: Tilde expansion

- **WHEN** a configured path starts with `~`
- **THEN** the tilde SHALL be expanded to the user's home directory

#### Scenario: Paths not configured

- **WHEN** `[paths]` section is not present or fields are empty
- **THEN** the system SHALL use platform-specific default paths

#### Scenario: Fixed paths under data directory

- **WHEN** the system resolves paths under data directory
- **THEN** templates SHALL always be at `{data_dir}/templates/`
- **AND** assets SHALL always be at `{data_dir}/assets/`
- **AND** credentials SHALL always be at `{data_dir}/credentials.json`
- **AND** token SHALL always be at `{data_dir}/token.json`
- **AND** `data_dir` itself is configurable via `paths.data_dir`
- **AND** subdirectory/file names under data_dir are NOT configurable

### Requirement: Selective ZIP Extraction

The system SHALL extract only required folders from the template repository ZIP.

#### Scenario: Extract templates from ZIP

- **WHEN** the update command extracts the template ZIP
- **THEN** only files under `project.template_path` (default: `/templates`) SHALL
  be extracted
- **AND** they SHALL be placed in `{data_dir}/templates/`

#### Scenario: Extract assets from ZIP

- **WHEN** the update command extracts the template ZIP
- **THEN** only files under `project.asset_path` (default: `/dist`) SHALL be
  extracted
- **AND** they SHALL be placed in `{data_dir}/assets/`

#### Scenario: Other files not extracted

- **WHEN** the update command extracts the template ZIP
- **THEN** files outside `project.template_path` and `project.asset_path` SHALL
  NOT be extracted
- **AND** user files (credentials.json, token.json) SHALL NOT be overwritten

#### Scenario: Templates not found in ZIP

- **WHEN** the update command extracts the template ZIP
- **AND** `project.template_path` does not exist in the ZIP
- **THEN** the command SHALL fail with an error
- **AND** the error message SHALL indicate that templates are required

#### Scenario: Assets not found in ZIP

- **WHEN** the update command extracts the template ZIP
- **AND** `project.asset_path` does not exist in the ZIP
- **THEN** the command SHALL create an empty `{data_dir}/assets/` directory
- **AND** the command SHALL NOT fail

#### Scenario: Zip Slip prevention

- **WHEN** the update command extracts a file from the template ZIP
- **AND** the resolved file path is outside the destination directory
- **THEN** the file SHALL NOT be extracted
- **AND** an error SHALL be logged or reported

### Requirement: Path Resolution Fallback Chain

The system SHALL resolve default paths using a three-level fallback chain.

#### Scenario: XDG environment variables set

- **WHEN** `$XDG_CONFIG_HOME` is set
- **THEN** config directory SHALL be `$XDG_CONFIG_HOME/nippo`
- **AND** this applies on all platforms including Windows

#### Scenario: Windows fallback when XDG not set

- **WHEN** running on Windows
- **AND** `$XDG_CONFIG_HOME` is not set
- **THEN** config directory SHALL fall back to `%APPDATA%\nippo`
- **AND** data directory SHALL fall back to `%LOCALAPPDATA%\nippo`
- **AND** cache directory SHALL fall back to `%LOCALAPPDATA%\nippo\cache`

#### Scenario: Default XDG paths as final fallback

- **WHEN** XDG environment variables are not set
- **AND** Windows environment variables are not available (non-Windows) or not
  set
- **THEN** config directory SHALL default to `~/.config/nippo`
- **AND** data directory SHALL default to `~/.local/share/nippo`
- **AND** cache directory SHALL default to `~/.cache/nippo`

#### Scenario: Fallback priority order

- **WHEN** resolving a default path
- **THEN** the system SHALL check in order:
  1. XDG environment variable (e.g., `$XDG_CONFIG_HOME`)
  2. Windows environment variable (e.g., `%APPDATA%`) on Windows only
  3. Default XDG path (e.g., `~/.config`)
- **AND** the first non-empty value SHALL be used

### Requirement: Diagnostic Command

The system SHALL provide a diagnostic command for health checking and
troubleshooting. The TUI implementation SHALL follow the patterns defined in
`ui-presentation` spec.

#### Scenario: Run health check

- **WHEN** user runs `nippo doctor` command
- **THEN** the system SHALL check configuration validity
- **AND** the system SHALL check project settings (drive folder, site URL, URL, branch)
- **AND** the system SHALL check all resolved directory paths for existence
- **AND** the system SHALL check required files (credentials.json, token.json,
  templates/, assets/)
- **AND** the system SHALL report cache status (nippo-template.zip, md/, output/)

#### Scenario: All checks pass

- **WHEN** all health checks pass
- **THEN** the output SHALL show ✓ indicator for each item
- **AND** the output SHALL show file modification timestamps where applicable
- **AND** the output SHALL show file counts for directories
- **AND** the output SHALL end with "All checks passed!"

#### Scenario: Issues found

- **WHEN** one or more health checks fail
- **THEN** the output SHALL show ✗ indicator for failed items
- **AND** the output SHALL list all issues at the end
- **AND** each issue SHALL include a suggested fix or command to run

#### Scenario: Directory does not exist

- **WHEN** a configured directory does not exist
- **THEN** the output SHALL indicate the directory is missing
- **AND** the suggested fix SHALL include the mkdir command

#### Scenario: credentials.json missing

- **WHEN** credentials.json is missing
- **THEN** the output SHALL indicate the file is missing
- **AND** the suggested fix SHALL guide the user to download from Google Cloud
  Console
- **AND** the suggested fix SHALL include a link to setup documentation

#### Scenario: token.json missing

- **WHEN** token.json is missing
- **THEN** the output SHALL indicate the file is missing
- **AND** the suggested fix SHALL guide the user to run `nippo auth` to
  authenticate

#### Scenario: Config file missing

- **WHEN** no configuration file exists
- **THEN** the output SHALL indicate the config file is missing
- **AND** the Configuration section SHALL show ✗ for config file
- **AND** the Paths section SHALL show default paths (not configured paths)
- **AND** the suggested fix SHALL guide the user to run `nippo init`

#### Scenario: Styled output with Lipgloss

- **WHEN** the doctor command renders output
- **THEN** success indicators (✓) SHALL be styled with success color
- **AND** failure indicators (✗) SHALL be styled with error color
- **AND** the output SHALL use View Provider architecture as per `ui-presentation`
  spec

#### Scenario: Data directory under git repository

- **WHEN** the doctor command checks the data directory
- **AND** the data directory is inside a git repository (has `.git` ancestor)
- **THEN** the output SHALL show a warning indicator (⚠)
- **AND** the warning message SHALL explain the security risk
- **AND** the suggested fix SHALL recommend moving sensitive files outside git
