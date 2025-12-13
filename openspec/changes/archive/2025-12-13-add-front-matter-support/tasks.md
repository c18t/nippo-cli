# Tasks

## 1. Core Front-Matter Parsing

- [x] 1.1 Add `go.yaml.in/yaml/v3` dependency for YAML parsing
- [x] 1.2 Create `FrontMatter` struct in `internal/domain/model/nippo.go` with
      `Created` and `Updated` fields
- [x] 1.3 Implement `ParseFrontMatter(content []byte)` function to extract YAML
      front-matter and return remaining content
- [x] 1.4 Implement `GenerateFrontMatter(created time.Time)` function to create
      YAML front-matter string with empty `updated` field
- [x] 1.5 Update `Nippo` struct to include `FrontMatter` field
- [x] 1.6 Implement `HasFrontMatter(content []byte) bool` helper function
- [x] 1.7 Implement `HasNowPlaceholder(content []byte) bool` to detect
      `updated: now`
- [x] 1.8 Implement RFC 3339 date validation for `created` and `updated` fields
- [x] 1.9 Implement front-matter update that preserves unknown fields

## 2. Drive API Integration

- [x] 2.1 Update `DriveFileProvider` interface to include `Update` method for
      file uploads
- [x] 2.2 Implement `Update(fileId string, content []byte)` in
      `drive_file_provider.go`
- [x] 2.3 Ensure Drive API response includes `createdTime` and `modifiedTime`
      fields in file listing

## 3. Format Command

- [x] 3.1 Create `cmd/format.go` with Cobra command definition
- [x] 3.2 Create `cmd/format_invoker.go` for DI setup
- [x] 3.3 Create `internal/inject/format.go` for format package
- [x] 3.4 Create `internal/usecase/port/format.go` interface
- [x] 3.5 Create `internal/usecase/interactor/format.go` implementation
- [x] 3.6 Create `internal/adapter/controller/format.go` controller
- [x] 3.7 Create `internal/adapter/presenter/format.go` presenter
- [x] 3.8 Add `LastFormatTimestamp` to config for tracking format state
- [x] 3.9 Implement format logic:
  - Fetch files updated since last format (or all files on first run)
  - Check for missing front-matter, missing `created` field, or `now`
    placeholder
  - Generate/update front-matter as needed (preserve unknown fields)
  - Handle combined updates (missing `created` + `now` placeholder)
  - Compare content before/after and skip upload if no changes
  - Upload to Drive
  - Update last format timestamp only if no failures
- [x] 3.10 Implement format output with status icons (success/failure/no-change)
- [x] 3.11 Implement format summary with failed files list and counts

## 4. Build Command Updates

- [x] 4.1 Update `GetMarkdown()` to strip front-matter from returned content
- [x] 4.2 Update Nippo loading to parse front-matter when present
- [x] 4.3 Add `GetCreatedTime()` method with fallback to filename-derived date
- [x] 4.4 Add `GetUpdatedTime()` method (returns empty if not set)
- [x] 4.5 Update `buildIndexPage()` to use front-matter metadata
- [x] 4.6 Update `buildNippoPage()` to use front-matter metadata
- [x] 4.7 Update `buildFeed()` to use front-matter `created` for item date and
      `updated` for item modification date when available
- [x] 4.8 Update `buildSiteMap()` to use last modified dates when available

## 5. Testing

- [x] 5.1 Add unit tests for front-matter parsing (valid, missing, malformed)
- [x] 5.2 Add unit tests for front-matter generation
- [x] 5.3 Add unit tests for `now` placeholder detection and replacement
- [x] 5.4 Add unit tests for RFC 3339 date validation
- [x] 5.5 Add unit tests for unknown field preservation
- [x] 5.6 Add unit tests for combined update scenarios
- [x] 5.7 Add unit tests for empty front-matter (`---\n---`)
- [x] 5.8 Add integration test for format command flow
- [ ] 5.9 Test format timestamp update logic (success vs failure)
- [x] 5.10 Test backward compatibility with files without front-matter in build
