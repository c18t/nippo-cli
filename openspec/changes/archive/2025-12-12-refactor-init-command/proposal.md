# Proposal: Fix directory creation and improve error messages

## Why

The `init` command fails on fresh systems due to missing directories and
provides unhelpful error messages:

1. **Config directory issue**: When `~/.config/nippo/` doesn't exist,
   `viper.SafeWriteConfig()` fails (line 109 in `internal/core/config.go`)
2. **Data directory issue**: When `~/.local/share/nippo/` doesn't exist, file
   operations fail
3. **Missing credentials.json**: When `credentials.json` is not found, the
   error message "unable to read client secret file" doesn't explain how to
   obtain it

## What Changes

### Root Cause

#### Config Directory

In `internal/core/config.go`, the `LoadConfig` method:

1. Calls `viper.ReadInConfig()` (line 105)
2. On `ConfigFileNotFoundError`, calls `viper.SafeWriteConfig()` (line 109)
3. `SafeWriteConfig()` requires the config directory to exist, but no code
   ensures this directory is created before the write operation

#### Data Directory and credentials.json

In `internal/usecase/interactor/init_save_drive_token.go`:

1. Calls `os.ReadFile(filepath.Join(dataDir, "credentials.json"))` (line 37)
2. If file doesn't exist, shows generic error (line 39)
3. No directory creation or helpful guidance for obtaining credentials

### Proposed Solution

1. **Ensure directories exist**: Create config and data directories before file
   operations, following the pattern used elsewhere in the codebase (e.g.,
   `local_file_provider.go:71`, `template_service.go:24`)

2. **Improve credentials.json error message**: When `credentials.json` is
   missing, provide clear instructions on how to obtain it from Google Cloud
   Console

### Scope

- **Capabilities**:
  - `config-initialization` (new)
  - `oauth-credentials` (new)
- **Modified files**:
  - `internal/core/config.go` - Config directory creation
  - `internal/usecase/interactor/init_save_drive_token.go` - Data directory
    creation and helpful error message
  - `internal/adapter/gateway/drive_file_provider.go` - Helpful error message
    only (no directory creation, as it's not the provider's responsibility)

## Non-Goals

- Changing the config file format or structure
- Modifying config file path resolution logic
- Adding config migration functionality

## Success Criteria

- Config and data directories are automatically created with appropriate
  permissions (0755) when needed
- When `credentials.json` is missing, a helpful error message is displayed that
  includes:
  - Clear explanation of what the file is
  - Instructions on how to obtain it from Google Cloud Console
  - Expected file location path
- Users can proceed with `init` command after placing credentials.json in the
  correct location
- Existing functionality remains unchanged when directories already exist
