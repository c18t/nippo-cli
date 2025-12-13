# Front-Matter Parsing

This capability enables parsing and management of YAML front-matter in Markdown
files for nippo content.

## ADDED Requirements

### Requirement: YAML Front-Matter Parsing

The system SHALL parse YAML front-matter from Markdown files when present. The
front-matter block MUST be delimited by `---` markers at the beginning of the
file.

#### Scenario: Valid front-matter is parsed

- **WHEN** a Markdown file contains valid YAML front-matter between `---`
  markers
- **THEN** the system SHALL extract the front-matter fields into the Nippo model
- **AND** the remaining content SHALL be available as the Markdown body

#### Scenario: File without front-matter

- **WHEN** a Markdown file does not contain `---` markers at the beginning
- **THEN** the system SHALL treat the entire file as Markdown content
- **AND** front-matter fields SHALL be derived from fallback sources

#### Scenario: Malformed front-matter in build

- **WHEN** the YAML between `---` markers is invalid during build
- **THEN** the system SHALL log a warning
- **AND** treat the file as if it has no front-matter

#### Scenario: Malformed front-matter in format

- **WHEN** the YAML between `---` markers is invalid during format
- **THEN** the system SHALL log an error with the filename
- **AND** skip processing that file
- **AND** continue with the next file

### Requirement: Front-Matter Data Structure

The system SHALL support the following front-matter fields:

- `created`: RFC 3339 datetime when the nippo was originally created
- `updated`: RFC 3339 datetime when the nippo was last modified (optional)

The system SHALL preserve any unknown fields added by the user.

#### Scenario: All fields present

- **WHEN** front-matter contains both `created` and `updated` fields
- **THEN** the system SHALL use these values for the Nippo model

#### Scenario: Only created field present

- **WHEN** front-matter contains only `created` field
- **THEN** the system SHALL use the `created` value
- **AND** `updated` SHALL remain empty (no automatic derivation)

#### Scenario: Updated field with special value now

- **WHEN** front-matter contains `updated: now` (unquoted string)
- **THEN** the format command SHALL replace it with Drive's `modifiedTime`

#### Scenario: Unknown fields are preserved

- **WHEN** front-matter contains fields other than `created` and `updated`
- **THEN** the system SHALL preserve these fields when updating front-matter

#### Scenario: Invalid date format

- **WHEN** front-matter contains `created` or `updated` with non-RFC 3339 value
- **THEN** the system SHALL treat this as malformed front-matter

### Requirement: Front-Matter Format

The system SHALL generate front-matter in YAML format following this structure:

```yaml
---
created: 2024-01-15T09:30:00+09:00
---
```

#### Scenario: Generated front-matter format

- **WHEN** the system generates front-matter
- **THEN** datetime values SHALL be in RFC 3339 format with local timezone
- **AND** the block SHALL start and end with `---` on separate lines
- **AND** `updated` field SHALL NOT be included (omit the field entirely)

### Requirement: Format Command

The system SHALL provide a `format` command to manage front-matter for files
on Google Drive.

A file requires front-matter updates if any of the following conditions are met:

1. The file has no front-matter
2. The file has front-matter but no `created` field
3. The file has `updated: now` placeholder

#### Scenario: Format fetches recently updated files

- **WHEN** running `nippo format`
- **THEN** the system SHALL fetch files updated since the last format timestamp
  stored in config
- **AND** process only files that need front-matter updates

#### Scenario: Format first run

- **WHEN** running `nippo format` for the first time (no last format timestamp)
- **THEN** the system SHALL fetch all files from Google Drive
- **AND** process all files that need front-matter updates

#### Scenario: Format adds front-matter to file without it

- **WHEN** processing a file without front-matter
- **THEN** the system SHALL generate front-matter with:
  - `created` set to Drive file's `createdTime` (local timezone)
  - `updated` field omitted
- **AND** prepend the front-matter to the file content
- **AND** upload the modified content to Google Drive

#### Scenario: Format updates now placeholder

- **WHEN** processing a file with `updated: now` in front-matter
- **THEN** the system SHALL replace `now` with Drive file's `modifiedTime`
- **AND** preserve all other fields in front-matter
- **AND** upload the modified content to Google Drive

#### Scenario: Format adds created to front-matter missing it

- **WHEN** processing a file with front-matter but no `created` field
- **THEN** the system SHALL add `created` set to Drive file's `createdTime`
- **AND** preserve all other fields in front-matter
- **AND** upload the modified content to Google Drive

#### Scenario: Format handles combined updates

- **WHEN** processing a file that requires multiple updates (e.g., missing
  `created` AND has `updated: now`)
- **THEN** the system SHALL apply all updates in a single operation
- **AND** upload the modified content once

#### Scenario: Format skips files with complete front-matter

- **WHEN** processing a file with valid front-matter including `created` (no
  `now` placeholder)
- **THEN** the system SHALL NOT modify the file
- **AND** NOT upload to Google Drive
- **AND** log with "no change" indicator

#### Scenario: Format skips upload when no changes

- **WHEN** processing results in no actual changes to file content
- **THEN** the system SHALL NOT upload to Google Drive
- **AND** log with "no change" indicator

#### Scenario: Format records last format timestamp on success

- **WHEN** format completes without any upload failures
- **THEN** the system SHALL record the current timestamp
- **AND** use this timestamp to filter files in subsequent format runs

#### Scenario: Format does not update timestamp on failure

- **WHEN** format completes with one or more upload failures
- **THEN** the system SHALL NOT update the last format timestamp
- **AND** failed files will be reprocessed in the next format run

### Requirement: Build Command Metadata Handling

The build command SHALL use front-matter metadata when available, with fallback
to filename-derived date for files without front-matter.

#### Scenario: Build with front-matter

- **WHEN** building a file with valid front-matter
- **THEN** the system SHALL use front-matter `created` for the date
- **AND** use front-matter `updated` when present
- **AND** NOT modify the source file on Google Drive

#### Scenario: Build without front-matter

- **WHEN** building a file without front-matter
- **THEN** the system SHALL derive the date from the filename pattern
  (`YYYY-MM-DD`)
- **AND** NOT modify the source file on Google Drive
- **AND** NOT generate front-matter automatically

#### Scenario: Feed generation with front-matter

- **WHEN** generating feed.xml for a file with front-matter
- **THEN** the system SHALL use `created` field for the item's publication date
- **AND** use `updated` field for the item's modification date when present

#### Scenario: Feed generation without front-matter

- **WHEN** generating feed.xml for a file without front-matter
- **THEN** the system SHALL use filename-derived date for the item's date

### Requirement: Drive Upload for Format Command

The format command SHALL upload modified content back to Google Drive.

#### Scenario: Successful upload after front-matter injection

- **WHEN** front-matter is generated or updated
- **THEN** the system SHALL update the file on Google Drive
- **AND** preserve the original file ID
- **AND** log with success indicator icon

#### Scenario: Upload failure handling

- **WHEN** uploading the modified content to Google Drive fails
- **THEN** the system SHALL log with failure indicator icon
- **AND** continue processing remaining files

### Requirement: Format Command Output

The format command SHALL provide clear visual feedback similar to build command.

#### Scenario: Progress logging with status icons

- **WHEN** processing files
- **THEN** the system SHALL log each file with a status indicator icon:
  - Success icon for successfully updated files
  - Failure icon for files that failed to upload
  - No-change icon for files that required no updates

#### Scenario: Summary with failed files list

- **WHEN** format completes with one or more failures
- **THEN** the system SHALL display a summary listing all failed filenames
- **AND** indicate the total count of successful, failed, and unchanged files
