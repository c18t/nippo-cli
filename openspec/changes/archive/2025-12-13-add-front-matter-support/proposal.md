# Change: Add Front-Matter Support for Nippo Markdown Files

## Why

Currently, nippo-cli derives metadata (dates, titles) solely from file names.
This approach limits flexibility and prevents storing additional metadata like
creation and modification timestamps directly within the content. By supporting
YAML front-matter, users gain explicit control over metadata while enabling
richer content generation (accurate timestamps in feeds, OGP data, etc.).

## What Changes

- **Parse YAML front-matter** from Markdown files (delimited by `---`)
- **Add `format` command** to manage front-matter:
  - Fetch files updated since last format run
  - Generate front-matter with `created` (from Drive's `createdTime`) and
    `updated` (empty by default, or Drive's `modifiedTime` if `now`)
  - Upload modified content back to Google Drive
- **Update `build` command** behavior:
  - Parse front-matter when present
  - Fall back to Drive's `createdTime` for files without front-matter
  - **No write operations to Google Drive** (read-only)
- **Update Nippo model** to store parsed front-matter metadata

## Impact

- Affected specs: `ui-presentation` (progress bar and summary display)
- Affected code:
  - `internal/domain/model/nippo.go` - Add FrontMatter struct and parsing
  - `cmd/format.go` - New format command
  - `internal/usecase/interactor/format.go` - New format interactor
  - `internal/adapter/gateway/drive_file_provider.go` - Add upload capability
  - `internal/usecase/interactor/build.go` - Use front-matter metadata
  - `internal/adapter/presenter/view/tui/format_progress.go` - Format progress TUI
  - `internal/adapter/presenter/view/tui/build_progress.go` - Build progress updates
  - `internal/adapter/presenter/format.go` - Format presenter with summary
  - `internal/adapter/presenter/build.go` - Build presenter with summary
