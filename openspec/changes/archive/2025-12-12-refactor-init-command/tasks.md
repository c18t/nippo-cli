# Tasks: Fix directory creation and improve error messages

## Implementation Tasks

1. **Modify config loading to create config directory** ✓
   - Update `LoadConfig()` in `internal/core/config.go`
   - Add `os.MkdirAll(c.GetConfigDir(), 0755)` before `SafeWriteConfig()`
   - Replace silent error ignoring with proper error handling
   - Wrap errors with descriptive messages using `fmt.Errorf("%w")`
   - **Validation**: Code compiles without errors

2. **Add data directory creation in init command** ✓
   - Update `Handle()` in
     `internal/usecase/interactor/init_save_drive_token.go`
   - Add `os.MkdirAll(dataDir, 0755)` before reading credentials.json
   - Add error handling for directory creation
   - **Validation**: Code compiles without errors

3. **Improve credentials.json error message in init command** ✓
   - Update error handling in `init_save_drive_token.go`
   - Use `os.IsNotExist(err)` to detect missing file
   - Show helpful multi-line error with setup instructions
   - Include Google Cloud Console URL and steps
   - Show exact file path where credentials should be placed
   - **Validation**: Code compiles, error message is clear and actionable

4. **Improve credentials.json error message in drive provider** ✓
   - Update error handling in
     `internal/adapter/gateway/drive_file_provider.go`
   - Use same helpful error message pattern (with note to run `nippo init`)
   - Do NOT add directory creation (not provider's responsibility)
   - **Validation**: Code compiles, error message is helpful and consistent

5. **Add import for os package if needed** ✓
   - Verify `os` package is imported in all modified files
   - Add import if missing
   - **Validation**: No import errors during compilation

## Testing Tasks

1. **Test fresh system initialization**
   - Remove existing directories: `rm -rf ~/.config/nippo ~/.local/share/nippo`
   - Run: `go run nippo/nippo.go init`
   - **Validation**: Helpful error message about missing credentials.json with
     setup instructions

2. **Test with existing directories but no credentials**
   - Create directories: `mkdir -p ~/.config/nippo ~/.local/share/nippo`
   - Run: `go run nippo/nippo.go init`
   - **Validation**: Helpful error message showing exact path for
     credentials.json

3. **Test with credentials.json**
   - Place valid credentials.json in data directory
   - Run: `go run nippo/nippo.go init`
   - **Validation**: OAuth flow proceeds normally

4. **Test error message content**
   - Remove credentials.json
   - Run init command and capture error
   - **Validation**: Error includes:
     - "credentials.json not found"
     - Google Cloud Console URL
     - Step-by-step instructions
     - Exact file path

5. **Test drive provider error message**
   - Remove credentials.json
   - Trigger drive provider (e.g., deploy command)
   - **Validation**: Helpful error message with note to run `nippo init`

6. **Test permission errors**
   - Create credentials.json with no read permission
   - Run init command
   - **Validation**: Error indicates permission problem, not "file not found"

## Documentation Tasks

1. **Update project documentation if needed**
   - Review if any documentation mentions config initialization behavior
   - Update if necessary
   - Consider adding credentials.json setup to README
   - **Validation**: Documentation accurately reflects new behavior

## Dependencies

- Implementation task 5 must complete before implementation tasks 1-4
- Implementation tasks 1-4 can run in parallel after implementation task 5
  completes
- Testing tasks can run in parallel after implementation tasks complete
- Documentation tasks can run in parallel with testing tasks
