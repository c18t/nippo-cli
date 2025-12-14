# Tasks

## 1. Core Path Configuration

- [x] 1.1 Add `ConfigPaths` struct to `internal/core/config.go`
  - Fields: `DataDir`, `CacheDir`
  - mapstructure tags for TOML parsing
  - Note: templates/, assets/, credentials.json, token.json are under DataDir
- [x] 1.2 Add `Paths` field to `Config` struct
- [x] 1.3 Implement path resolution fallback chain
  - Create `internal/core/paths.go` with fallback logic
  - Priority 1: XDG environment variables (all platforms)
  - Priority 2: Windows environment variables (Windows only)
  - Priority 3: Default XDG paths (`~/.config`, etc.)
- [x] 1.4 Implement path resolution methods
  - `GetDataDir()` - returns resolved data directory (configurable)
  - `GetCacheDir()` - returns resolved cache directory (configurable)
  - Fixed paths under DataDir: `templates/`, `assets/`, `credentials.json`,
    `token.json`
  - Fixed paths under CacheDir: `md/`, `output/`, `nippo-template.zip`
- [x] 1.5 Implement environment variable expansion using `os.ExpandEnv()`
- [x] 1.6 Implement tilde (`~`) expansion for home directory
- [x] 1.7 Implement relative path resolution (relative to config file)

## 2. Update init Command for Project Configuration

- [x] 2.1 Add Google Drive folder configuration to init command
  - Prompt for drive folder URL or ID (first in project settings)
  - Default: `https://drive.google.com/drive/folders/1HNSRS2tJI2t7DKP_8XQJ2NTleSH-rs4y`
  - Add `project.drive_folder_id` field to `ConfigProject` struct
  - Implement folder ID extraction from URL formats:
    - `https://drive.google.com/drive/u/0/folders/{id}`
    - `https://drive.google.com/drive/folders/{id}`
    - `https://drive.google.com/open?id={id}`
  - Accept raw folder ID if not a URL
- [x] 2.2 Add site URL configuration to init command
  - Prompt for site URL (after drive folder)
  - Default: `https://nippo.c18t.me`
  - Add `project.site_url` field to `ConfigProject` struct
- [x] 2.3 Add branch configuration to init command
  - Prompt for branch name (after project URL, optional, default: `main`)
  - Add `project.branch` field to `ConfigProject` struct
  - Update URL generation to use configured branch
- [x] 2.4 Add template_path and asset_path configuration to init command
  - Prompt for template path in ZIP (after branch, default: `/templates`)
  - Prompt for asset path in ZIP (after template_path, default: `/dist`)
  - All `[project]` values are explicitly written (no fallback logic needed)
- [x] 2.5 Update `internal/usecase/interactor/init_setting.go`
  - Prompt order: drive folder, site URL, project URL, branch, template_path,
    asset_path
  - Implement selective ZIP extraction (same as update.go)
  - Note: Do NOT prompt for path configuration (data_dir, cache_dir)
  - If config file exists, show warning and ask for confirmation (y/N):
    - "Configuration file already exists. Comments may be lost. Continue?"
    - Exit without changes if user declines
  - Load existing config file if present (merge behavior):
    - Use existing values as prompt defaults
    - Only overwrite values that user explicitly enters (Enter = keep existing)
    - Preserve existing `[paths]` section
    - Use existing paths for file operations (token.json, templates, etc.)
  - Check if data_dir is under git repository (proactive security check):
    - Show warning if `.git` ancestor found
    - Ask for confirmation (y/N) to proceed
  - Write `[paths]` section with commented-out defaults (if no existing paths):
    ```toml
    [paths]
    # Uncomment and modify to customize file locations.
    # data_dir = "/Users/xxx/.local/share/nippo"
    # cache_dir = "/Users/xxx/.cache/nippo"
    ```
  - Use resolved absolute paths in comments (not environment variables)
- [x] 2.6 Create required directories during init
  - Create config directory if missing (before saving config file)
  - Create data directory if missing (before template extraction)
  - Create cache directory if missing
  - Use permissions 0755
- [x] 2.7 Show auth instruction on init completion
  - After init completes, show message: "Run `nippo auth` to authenticate"
  - Do NOT perform OAuth in init command

## 2.8. Add `nippo auth` Command

- [x] 2.8.1 Create `cmd/auth.go`
  - Add `auth` subcommand to root
- [x] 2.8.2 Create `internal/adapter/controller/auth_controller.go`
  - Implement OAuth flow logic (extract from init_save_drive_token.go)
- [x] 2.8.3 Check credentials.json before OAuth flow
  - Check credentials.json exists in data directory
  - If missing, show error with download instructions
  - Include link to setup documentation
    (<https://github.com/c18t/nippo-cli#setup>)
  - Do not proceed with OAuth until credentials.json is present
- [x] 2.8.4 Perform OAuth flow and save token.json
  - Reuse existing OAuth logic from init_save_drive_token.go
  - Save token.json to data directory
- [x] 2.8.5 Auth command requires config file
  - If no config file exists, show error: "Please run `nippo init` first"

## 3. Require Initialization Before Other Commands

- [x] 3.1 Modify `internal/core/config.go` LoadConfig behavior
  - Do not auto-create config file on ConfigFileNotFoundError
  - Return specific error type for "not initialized"
- [x] 3.2 Add initialization check in `cmd/root.go`
  - Allow `init`, `doctor`, `help`, `version` without config file
  - Other commands: show error "Please run `nippo init` first" and exit

## 4. Update Existing Code to Use Configurable Paths

- [x] 4.1 Update `internal/adapter/gateway/drive_file_provider.go`
  - Replace hardcoded paths with `GetDataDir()/credentials.json`
  - Replace hardcoded paths with `GetDataDir()/token.json`
- [x] 4.2 Update `internal/domain/logic/service/template_service.go`
  - Use `GetDataDir()/templates` for template loading (fixed path)
- [x] 4.3 Update `internal/usecase/interactor/build.go`
  - Use `GetCacheDir()` for md/ and output/ directories
  - Use `project.drive_folder_id` instead of hardcoded folder ID
  - Use `project.site_url` instead of hardcoded `https://nippo.c18t.net`
    - HTML meta tags (canonical URL, og:url)
    - OGP image URL
    - RSS feed links
    - Sitemap URLs
- [x] 4.4 Refactor `internal/usecase/interactor/init_save_drive_token.go`
  - Move OAuth logic to auth command (see 2.8.2)
  - Remove or rename this file after extraction
  - Note: This file's functionality is replaced by auth_controller.go
- [x] 4.5 Update `internal/domain/logic/repository/asset.go`
  - Use `GetCacheDir()` for cache operations
- [x] 4.6 Update `internal/usecase/interactor/format.go`
  - Use `project.drive_folder_id` instead of hardcoded folder ID
- [x] 4.7 Update `internal/usecase/interactor/update.go`
  - Change ZIP filename from `main` to `nippo-template.zip`
  - Use `project.branch` from config instead of hardcoded `main`
  - Implement selective ZIP extraction:
    - Only extract files matching `project.template_path` (default: `/templates`)
    - Only extract files matching `project.asset_path` (default: `/dist`)
    - Map ZIP paths to fixed destinations:
      - `{repo}-{branch}{template_path}/*` → `GetDataDir()/templates/`
      - `{repo}-{branch}{asset_path}/*` → `GetDataDir()/assets/`
  - Error handling:
    - If `template_path` not found in ZIP → error (templates required)
    - If `asset_path` not found in ZIP → create empty `assets/` directory
  - Zip Slip prevention:
    - Validate each file path is within destination directory
    - Skip files with `../` that escape destination
  - Benefits: Prevents overwriting user files, reduces disk usage
- [x] 4.8 Update `internal/usecase/interactor/deploy.go`
  - Change `data/output` to `data/assets` for static assets
  - Note: Template repo (c18t/nippo) uses `dist/` folder for built assets
  - `project.asset_path = "/dist"` maps to local `assets/` folder

## 5. Add `nippo doctor` Command

- [x] 5.1 Create `internal/adapter/controller/doctor_controller.go`
  - Implement health check logic
  - Check configuration file validity
  - Check directory existence
  - Check required files (credentials.json, token.json, templates/, assets/)
  - Report issues with suggested fixes
- [x] 5.2 Create `cmd/doctor.go`
  - Add `doctor` subcommand to root
- [x] 5.3 Implement checks:
  - Configuration section: config file, drive folder, site URL, project URL, branch
  - Paths section: config/data/cache directories with existence status
  - Required Files section: credentials.json, token.json, templates/, assets/
  - Cache Status section: nippo-template.zip, md/, output/ with file counts
  - Security check: data directory under git repository (warn if `.git` ancestor)
- [x] 5.4 Handle no-config case:
  - Show ✗ for config file with suggestion to run `nippo init`
  - Use default paths for remaining checks
  - Continue checking other items to provide useful diagnostics
- [x] 5.5 Format output:
  - Use ✓/✗ indicators for pass/fail
  - Show file modification timestamps where applicable
  - Show file counts for directories
  - List issues with actionable suggestions at the end:
    - credentials.json: guide to download from Google Cloud Console + doc link
    - token.json: guide to run `nippo auth` to authenticate
    - directories: show mkdir command
    - config file: guide to run `nippo init`
- [x] 5.6 TUI implementation (follow `ui-presentation` spec)
  - Create presenter and view model in `internal/adapter/presenter/`
  - Use Lipgloss for styled output (success green, error red)
  - Use View Provider architecture for rendering
  - Preserve output on completion (do not clear terminal)

## 6. Testing

- [x] 6.1 Add unit tests for path resolution in `internal/core/paths_test.go`
  - Test fallback chain (XDG → Windows → default XDG)
  - Test environment variable expansion
  - Test tilde (`~`) expansion
  - Test relative path resolution
- [x] 6.2 Add unit tests for `GetDataDir()`, `GetCacheDir()`
- [x] 6.3 Add integration tests for `nippo doctor` command
- [x] 6.4 Add integration tests for `nippo auth` command
  - Test credentials.json check (missing → error with instructions)
  - Test OAuth flow execution
  - Test token.json saving
  - Test config file requirement (missing → error)
- [x] 6.5 Test initialization requirement
  - Verify commands fail without config file
  - Verify init/doctor/help/version work without config
- [x] 6.6 Test backward compatibility
  - Existing installations without `[paths]` section continue to work

## 7. Documentation

- [x] 7.1 Update README with path configuration section
- [x] 7.2 Add example configuration with custom paths
- [x] 7.3 Document platform-specific default paths
