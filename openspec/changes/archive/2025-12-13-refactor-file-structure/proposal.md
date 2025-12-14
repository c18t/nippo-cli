# Change: Refactor File Structure and Path Configuration

## Why

The current file structure follows XDG Base Directory specification but has
several limitations:

1. File paths (credentials.json, token.json, templates/, md/, output/) are
   hardcoded and cannot be customized
2. Windows support is limited (XDG fallbacks use Unix-style paths)
3. The semantic categorization of files (config vs data vs cache) could be
   improved for better clarity
4. Google Drive folder structure has too many year-month folders at the top
   level, making manual navigation difficult

This change will make file paths configurable, improve cross-platform support,
clarify the purpose of each file location, and document the planned Google
Drive folder restructuring.

## What Changes

### 1. Configurable File Paths

- **ADDED**: Configuration options for customizing file paths
  - `paths.data_dir` - Data directory (templates/, assets/, credentials, token)
  - `paths.cache_dir` - Cache directory (md/, output/, nippo-template.zip)

### 2. Cross-Platform Default Paths

- **MODIFIED**: Default path resolution with fallback chain
  - Priority 1: XDG environment variables (`$XDG_CONFIG_HOME`, etc.)
  - Priority 2: Windows environment variables (`%APPDATA%`, `%LOCALAPPDATA%`)
  - Priority 3: Default XDG paths (`~/.config`, `~/.local/share`, `~/.cache`)

### 3. File Location Semantics

- **MODIFIED**: Clarify file categorization
  - Config dir: Application configuration only (`nippo.toml`)
  - Data dir: User-specific persistent data (`credentials.json`, `token.json`,
    `templates/`, `assets/`)
  - Cache dir: Regenerable/temporary data (`md/`, `output/`,
    `nippo-template.zip`)
- **ADDED**: New `assets/` directory for static assets extracted from template
  - Static assets (CSS, JS, images) are stored in `{data_dir}/assets/`
  - Build output (HTML) remains in `{cache_dir}/output/`
  - Note: Template repo uses `dist/` folder; mapped via `project.asset_path`

### 3.1. Selective ZIP Extraction

- **MODIFIED**: Extract only required folders from template ZIP
  - Current behavior: Extracts ALL files from ZIP to data directory
  - New behavior: Extracts only `templates/` and `assets/` folders
  - `project.template_path` - Source path in ZIP (default: `/templates`)
  - `project.asset_path` - Source path in ZIP (default: `/dist`)
  - Destination paths are fixed subdirectories under `data_dir`:
    - `{data_dir}/templates/` for templates
    - `{data_dir}/assets/` for static assets
  - Error handling:
    - If `project.template_path` not found in ZIP → **error** (templates required)
    - If `project.asset_path` not found in ZIP → create empty `{data_dir}/assets/`
  - Benefits:
    - Prevents dev files (.github, src/, etc.) from being extracted
    - Protects user files (credentials.json, token.json) from overwrite
    - Reduces disk usage

### 4. Require Initialization

- **MODIFIED**: Commands other than `init`, `doctor`, `help`, `version` require
  config
  - If config file not found, show error: "Please run `nippo init` first"
  - No automatic config file creation
  - `doctor` command works without config to help diagnose setup issues

### 4.1. Simplified Init Command

- **MODIFIED**: Init command prompts for all project settings
  - Prompts for: Drive folder, site URL, project URL, branch, template_path,
    asset_path
  - All `[project]` values are explicitly written (no fallback logic needed)
  - Does NOT prompt for path customization (data_dir, cache_dir)
  - Writes `[paths]` section with commented-out defaults (for discoverability)
  - Path customization: uncomment and modify values in config file
- **MODIFIED**: Init merges with existing config if present
  - If config file exists, show warning and ask for confirmation (y/N)
  - Warning: "Configuration file already exists. Comments may be lost. Continue?"
  - If user declines, exit without changes
  - If user confirms, load existing config and use values as prompt defaults
  - Preserve existing `[paths]` section (do not overwrite)
  - Enables both initial setup and configuration updates via `nippo init`
- **MODIFIED**: Init creates required directories
  - Create config directory if missing
  - Create data directory if missing (for token.json, templates/, assets/)
  - Create cache directory if missing (for md/, output/)
- **ADDED**: Init warns if data directory is under git repository
  - Check before creating directories or extracting templates
  - Show warning: "Data directory is under git repository. Credentials may be
    committed accidentally."
  - Ask for confirmation (y/N) to proceed
  - Proactive security check (same as `doctor` command)
- **MODIFIED**: Init does NOT perform OAuth authentication
  - Authentication is separated into `nippo auth` command
  - After init completes, show message: "Run `nippo auth` to authenticate"

### 4.1.1. Auth Command (NEW)

- **ADDED**: `nippo auth` command for Google Drive authentication
  - Check credentials.json exists before starting OAuth flow
  - If missing, show error with download instructions and documentation link
  - Perform OAuth flow and save token.json
  - Can be re-run anytime to refresh authentication
  - Works independently of init (only requires config file to exist)

### 4.2. Full Path Customization

- **NOTE**: Users can fully customize all file locations
  - Config file location: `nippo --config /path/to/nippo.toml`
  - Data/cache directories: `[paths]` section in config file
  - Combined: All nippo files can be placed in user-defined locations
  - Use case: Portable installations, custom directory structures, shared configs

### 4.3. Update Command Changes

- **MODIFIED**: `nippo update` uses selective ZIP extraction
  - Same extraction logic as `init` command
  - Downloads ZIP from configured `project.branch` (not hardcoded `main`)
  - Extracts only `project.template_path` and `project.asset_path`
  - Saves to `{data_dir}/templates/` and `{data_dir}/assets/`
  - Renames downloaded file from `main` to `nippo-template.zip`
  - Validates paths to prevent Zip Slip vulnerability

### 5. Project Configuration Enhancements

- **ADDED**: Configurable Google Drive folder ID
  - `project.drive_folder_id` - Root folder ID for nippo files on Google Drive
  - Currently hardcoded as `1HNSRS2tJI2t7DKP_8XQJ2NTleSH-rs4y`
  - Required for `build` and `format` commands to locate nippo files
- **ADDED**: Configurable site URL
  - `project.site_url` - Base URL for the generated site
  - Currently hardcoded as `https://nippo.c18t.net` (incorrect, should be `.me`)
  - Default: `https://nippo.c18t.me`
  - Used in HTML meta tags, OGP, RSS feeds, and sitemaps
- **ADDED**: Configurable branch for template repository
  - `project.branch` - Git branch to download (default: `main`)
  - Allows testing with feature branches during development
  - `nippo update` also uses configured branch

### 6. Diagnostic Command

- **ADDED**: `nippo doctor` command for health check and troubleshooting
  - Checks configuration validity and required files
  - Shows resolved paths with existence status
  - Reports issues and suggests fixes
  - **Security check**: Warns if data directory is under git repository
    - Prevents accidental commit of credentials.json/token.json
  - TUI implementation SHALL follow `ui-presentation` spec patterns
    - Use Lipgloss for styled output (success/error/warning colors)
    - Use View Provider architecture for rendering
  - Example output:

    ```text
    nippo doctor - Health Check

    Configuration
      Config file:     /Users/xxx/.config/nippo/nippo.toml ✓
      Drive folder:    https://drive.google.com/drive/folders/1HNSR...rs4y ✓
      Site URL:        https://nippo.c18t.me ✓
      Project URL:     https://github.com/c18t/nippo ✓
      Project Branch:  main

    Paths
      Config:      /Users/xxx/.config/nippo           ✓
      Data:        /Users/xxx/.local/share/nippo      ✓
      Cache:       /Users/xxx/.cache/nippo            ✓

    Required Files
      credentials.json    ✓  2023-12-15 12:00:00
      token.json          ✓  2025-12-13 13:05:00
      templates/          ✓  6 files
      assets/             ✓  5 files

    Cache Status
      nippo-template.zip  ✓  2025-12-13 13:05:00
      md/                 ✓  1553 files
      output/             ✓  1648 files

    All checks passed!
    ```

  - Example output with issues:

    ```text
    nippo doctor - Health Check

    Configuration
      Config file:     /Users/xxx/.config/nippo/nippo.toml ✓
      Drive folder:    https://drive.google.com/drive/folders/1HNSR...rs4y ✓
      Site URL:        https://nippo.c18t.me ✓
      Project URL:     https://github.com/c18t/nippo ✓
      Project Branch:  main

    Paths
      Config:      /Users/xxx/.config/nippo           ✓
      Data:        /Users/xxx/.local/share/nippo      ✓
      Cache:       /Users/xxx/.cache/nippo            ✗ missing

    Required Files
      credentials.json    ✗  missing
      token.json          ✗  missing
      templates/          ✓  6 files
      assets/             ✓  5 files

    Issues Found:
      ✗ Cache directory does not exist
        Run: mkdir -p /Users/xxx/.cache/nippo
      ✗ credentials.json not found
        Download from Google Cloud Console and place in data directory
        See: https://github.com/c18t/nippo-cli#setup
      ✗ token.json not found
        Run: nippo auth (to authenticate with Google Drive)
    ```

### 7. Google Drive Folder Structure (Manual Migration)

- **MODIFIED**: Reorganize nippo folder structure in Google Drive
  - Before: `YYYYMM/YYYY-MM-DD.md` (e.g., `201811/2018-11-30.md`)
  - After: `YYYY/MM/YYYY-MM-DD.md` (e.g., `2018/11/2018-11-30.md`)
  - Reduces top-level folder count from ~80+ year-month folders to ~7 year
    folders
  - **No code changes required**: The program uses folder IDs (not folder
    names) and extracts dates from filenames only

#### Code Impact Analysis

The following code was analyzed to confirm no changes are needed:

1. `internal/adapter/gateway/drive_file_provider.go` - Uses folder IDs for
   recursive file listing, not folder name patterns
2. `internal/usecase/interactor/build.go:127` - Specifies root folder by ID:
   `Folders: []string{"1HNSRS2tJI2t7DKP_8XQJ2NTleSH-rs4y"}`
3. `internal/domain/model/nippo.go:68` - Parses date from filename only:
   `filepath.Base(filePath)[:10]` extracts `YYYY-MM-DD`
4. `internal/domain/logic/repository/nippo.go` - Handles folder traversal via
   IDs, agnostic to folder naming conventions

#### Migration Steps (Manual)

1. Create year folders (`2018`, `2019`, ..., `2024`, `2025`) in Google Drive
2. Move month folders into corresponding year folders
3. Optionally rename month folders from `YYYYMM` to `MM` for consistency

## Impact

- Affected specs:
  - `config-initialization` - Add path configuration loading, auth command
  - `file-paths` (NEW) - Define path resolution behavior
  - `oauth-credentials` - Reference configurable paths
- Affected code:
  - `cmd/auth.go` (NEW) - Auth subcommand
  - `cmd/doctor.go` (NEW) - Doctor subcommand
  - `internal/core/config.go` - Path configuration and resolution
  - `internal/adapter/controller/auth_controller.go` (NEW) - OAuth flow logic
  - `internal/adapter/controller/doctor_controller.go` (NEW) - Health check logic
  - `internal/adapter/gateway/drive_file_provider.go` - Configurable paths
  - `internal/domain/logic/service/template_service.go` - Configurable paths
  - `internal/usecase/interactor/build.go` - Configurable paths, folder ID
  - `internal/usecase/interactor/format.go` - Configurable folder ID
  - `internal/usecase/interactor/init_save_drive_token.go` - Refactor/remove
    (OAuth logic moved to auth_controller.go)
