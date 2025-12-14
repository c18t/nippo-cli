# config-initialization Spec Delta

## MODIFIED Requirements

### Requirement: Required Directory Creation

The system SHALL create required directories (config and data) before
attempting file operations that depend on them. Additionally, the system SHALL
support platform-specific default directories.

#### Scenario: Init creates config directory if missing

**Given** the user runs `nippo init`
**And** the configuration directory does not exist
**When** the init command proceeds to save configuration
**Then** the configuration directory must be created with permissions 0755
**And** a new configuration file must be written to the directory

#### Scenario: Init creates config file in existing directory

**Given** the user runs `nippo init`
**And** the configuration directory exists
**And** the configuration file does not exist
**When** the init command proceeds to save configuration
**Then** a new configuration file must be written to the existing directory

#### Scenario: Init fails when directory creation fails

**Given** the user runs `nippo init`
**And** the configuration directory does not exist
**And** directory creation will fail (e.g., permission denied)
**When** the init command attempts to save configuration
**Then** an error must be returned
**And** the error message must indicate directory creation failure
**And** the error must include the underlying system error

#### Scenario: Init fails when config file write fails

**Given** the user runs `nippo init`
**And** the configuration directory was successfully created
**And** writing the config file will fail (e.g., disk full)
**When** the init command attempts to save configuration
**Then** an error must be returned
**And** the error message must indicate config write failure
**And** the error must include the underlying write error

#### Scenario: Config directory with XDG environment variable

**Given** `XDG_CONFIG_HOME` is set to a valid path
**When** the application determines the config directory
**Then** the config directory SHALL be `$XDG_CONFIG_HOME/nippo`
**And** this applies on all platforms including Windows

#### Scenario: Config directory fallback on Windows

**Given** the application is running on Windows
**And** `XDG_CONFIG_HOME` is not set
**And** `%APPDATA%` is set
**When** the application determines the config directory
**Then** the config directory SHALL be `%APPDATA%\nippo`

#### Scenario: Config directory default XDG path fallback

**Given** `XDG_CONFIG_HOME` is not set
**And** Windows environment variables are not available or not set
**When** the application determines the config directory
**Then** the config directory SHALL be `~/.config/nippo`

#### Scenario: Command requires initialization

**Given** the user has not run `nippo init`
**And** no configuration file exists
**When** the user runs a command other than `init`, `doctor`, `help`, or `version`
**Then** the command SHALL fail with an error
**And** the error message SHALL instruct the user to run `nippo init` first
**And** no configuration file SHALL be created automatically

#### Scenario: Doctor command works without initialization

**Given** the user has not run `nippo init`
**And** no configuration file exists
**When** the user runs `nippo doctor`
**Then** the command SHALL run successfully
**And** the output SHALL indicate that configuration file is missing
**And** the suggested fix SHALL guide the user to run `nippo init`

## ADDED Requirements

### Requirement: Project Configuration Enhancements

The init command SHALL allow users to configure additional project settings
including Git branch and Google Drive folder ID.

#### Scenario: Prompt for Google Drive folder during init

**Given** the user runs `nippo init`
**When** the init command proceeds to project configuration
**Then** the user SHALL be prompted for the Google Drive folder (URL or ID) first
**And** the default value SHALL be
`https://drive.google.com/drive/folders/1HNSRS2tJI2t7DKP_8XQJ2NTleSH-rs4y`

#### Scenario: Prompt for site URL during init

**Given** the user runs `nippo init`
**And** Google Drive folder is configured
**When** the init command proceeds to project configuration
**Then** the user SHALL be prompted for the site URL
**And** the default value SHALL be `https://nippo.c18t.me`

#### Scenario: Prompt for project URL during init

**Given** the user runs `nippo init`
**And** Google Drive folder and site URL are configured
**When** the init command proceeds to project configuration
**Then** the user SHALL be prompted for the template project URL
**And** the default value SHALL be `https://github.com/c18t/nippo`

#### Scenario: Extract folder ID from Google Drive URL

**Given** the user enters a Google Drive URL in one of these formats:

- `https://drive.google.com/drive/u/0/folders/{folder_id}`
- `https://drive.google.com/drive/folders/{folder_id}`
- `https://drive.google.com/open?id={folder_id}`

**When** the init command processes the input
**Then** the folder ID SHALL be extracted from the URL
**And** the extracted ID SHALL be saved to `project.drive_folder_id`

#### Scenario: Accept folder ID directly

**Given** the user enters a folder ID directly (not a URL)
**When** the init command processes the input
**Then** the input SHALL be used as `project.drive_folder_id`

#### Scenario: Prompt for branch name during init

**Given** the user runs `nippo init`
**And** drive folder ID, site URL, and project URL are configured
**When** the init command proceeds to branch configuration
**Then** the user SHALL be prompted for an optional branch name
**And** the default value SHALL be `main`

#### Scenario: Prompt for template path during init

**Given** the user runs `nippo init`
**And** drive folder, site URL, project URL, and branch are configured
**When** the init command proceeds to template path configuration
**Then** the user SHALL be prompted for the template path in ZIP
**And** the default value SHALL be `/templates`

#### Scenario: Prompt for asset path during init

**Given** the user runs `nippo init`
**And** drive folder, site URL, project URL, branch, and template path are
configured
**When** the init command proceeds to asset path configuration
**Then** the user SHALL be prompted for the asset path in ZIP
**And** the default value SHALL be `/dist`

#### Scenario: Branch used in URL generation

**Given** `project.branch` is set to a value (e.g., `develop`)
**When** the system generates the download URL for a GitHub repository
**Then** the URL SHALL use the configured branch instead of `main`
**And** the URL format SHALL be
`https://codeload.github.com/{owner}/{repo}/zip/refs/heads/{branch}`

#### Scenario: Update command uses configured branch

**Given** `project.branch` is set in configuration
**When** the user runs `nippo update`
**Then** the update command SHALL download from the configured branch

#### Scenario: Build command uses configured drive folder ID

**Given** `project.drive_folder_id` is set in configuration
**When** the user runs `nippo build`
**Then** the build command SHALL use the configured folder ID to locate nippo files

#### Scenario: Format command uses configured drive folder ID

**Given** `project.drive_folder_id` is set in configuration
**When** the user runs `nippo format`
**Then** the format command SHALL use the configured folder ID to locate nippo files

#### Scenario: Build command uses configured site URL

**Given** `project.site_url` is set in configuration
**When** the user runs `nippo build`
**Then** the build command SHALL use the configured site URL for:

- HTML meta tags (canonical URL, og:url)
- OGP image URLs
- RSS feed links and IDs
- Sitemap URLs

### Requirement: Simplified Init Command

The init command SHALL focus on essential project settings only and SHALL NOT
prompt for path customization. The init command SHALL write a `[paths]` section
with commented-out defaults for discoverability.

#### Scenario: Init does not prompt for path configuration

**Given** the user runs `nippo init`
**And** drive folder, site URL, project URL, and branch configuration is complete
**When** the init command finishes configuration prompts
**Then** the user SHALL NOT be prompted for data directory
**And** the user SHALL NOT be prompted for cache directory
**And** the system SHALL use platform-specific default paths

#### Scenario: Init writes commented paths section

**Given** the user runs `nippo init`
**And** no existing configuration file is present
**When** the init command writes the configuration file
**Then** the config file SHALL contain a `[paths]` section
**And** the `data_dir` option SHALL be commented out with its default value
**And** the `cache_dir` option SHALL be commented out with its default value
**And** a comment SHALL explain how to customize paths

#### Scenario: Init preserves existing paths section

**Given** the user runs `nippo init`
**And** an existing configuration file is present with `[paths]` section
**When** the init command writes the configuration file
**Then** the existing `[paths]` section SHALL be preserved

#### Scenario: Init creates required directories

**Given** the user runs `nippo init`
**And** the data directory does not exist
**When** the init command proceeds to save files
**Then** the config directory SHALL be created if missing
**And** the data directory SHALL be created if missing
**And** the cache directory SHALL be created if missing
**And** directories SHALL be created with permissions 0755

#### Scenario: Init warns if data directory is under git repository

**Given** the user runs `nippo init`
**And** the resolved data directory is inside a git repository
**When** the init command determines the data directory path
**Then** a warning SHALL be shown explaining the security risk
**And** the user SHALL be prompted for confirmation (y/N) to proceed
**And** if the user declines, the command SHALL exit without changes

#### Scenario: Init shows auth instruction on completion

**Given** the user runs `nippo init`
**When** the init command completes successfully
**Then** a message SHALL be shown instructing the user to run `nippo auth`
**And** the init command SHALL NOT perform OAuth authentication

#### Scenario: Init confirms before updating existing config

**Given** the user runs `nippo init`
**And** a configuration file already exists
**When** the init command starts
**Then** a warning SHALL be shown: "Configuration file already exists. Comments
may be lost."
**And** the user SHALL be prompted for confirmation (y/N)
**And** if the user declines, the command SHALL exit without changes

#### Scenario: Init uses existing values as prompt defaults

**Given** the user runs `nippo init`
**And** the user confirms to update existing configuration
**And** an existing configuration file has `project.drive_folder_id` set
**When** the init command prompts for drive folder
**Then** the existing value SHALL be shown as the default
**And** pressing Enter without input SHALL keep the existing value

#### Scenario: Init only overwrites explicitly entered values

**Given** the user runs `nippo init`
**And** an existing configuration file is present
**When** the user presses Enter without entering a new value
**Then** the existing value SHALL be preserved
**And** the value SHALL NOT be overwritten with the hardcoded default

#### Scenario: Init uses existing paths for file operations

**Given** the user runs `nippo init`
**And** an existing configuration file has `paths.data_dir` set
**When** the init command extracts templates from ZIP
**Then** the templates SHALL be saved to the configured `data_dir/templates/`
**And** the assets SHALL be saved to the configured `data_dir/assets/`

#### Scenario: Path customization via manual config edit

**Given** the user wants to customize file paths
**When** the user edits the configuration file directly
**And** uncomments `data_dir` or `cache_dir` in the `[paths]` section
**Then** the system SHALL use the configured paths
**And** `nippo doctor` SHALL show the configured paths

### Requirement: Auth Command

The system SHALL provide an `auth` command for Google Drive authentication,
separated from the init command.

#### Scenario: Auth command checks credentials.json

**Given** the user runs `nippo auth`
**And** credentials.json does not exist in data directory
**When** the auth command attempts to start OAuth flow
**Then** the command SHALL show an error message
**And** the error SHALL include instructions to download credentials.json
**And** the error SHALL include a link to setup documentation
**And** the OAuth flow SHALL NOT proceed

#### Scenario: Auth command performs OAuth flow

**Given** the user runs `nippo auth`
**And** credentials.json exists in data directory
**When** the auth command starts OAuth flow
**Then** the OAuth flow SHALL be performed
**And** token.json SHALL be saved to the data directory

#### Scenario: Auth command can refresh authentication

**Given** the user runs `nippo auth`
**And** token.json already exists
**When** the auth command completes OAuth flow
**Then** the existing token.json SHALL be overwritten with new token

#### Scenario: Auth command requires config file

**Given** the user runs `nippo auth`
**And** no configuration file exists
**When** the auth command starts
**Then** the command SHALL fail with an error
**And** the error message SHALL instruct the user to run `nippo init` first

### Requirement: Path Configuration Loading

The system SHALL load path configuration from the `[paths]` section and resolve
paths appropriately.

#### Scenario: Load paths section from config

**Given** the configuration file contains a `[paths]` section
**When** the configuration is loaded
**Then** the path values SHALL be parsed and stored
**And** empty or missing values SHALL use platform defaults

#### Scenario: Environment variable expansion in paths

**Given** a path configuration contains environment variables
**When** the path is resolved
**Then** environment variables SHALL be expanded using `os.ExpandEnv()`
**And** both Unix (`$VAR`) and Windows (`%VAR%`) formats SHALL be supported

#### Scenario: Relative path resolution

**Given** a path configuration contains a relative path
**When** the path is resolved
**Then** the path SHALL be resolved relative to the config file directory
**And** the resolved path SHALL be an absolute path
