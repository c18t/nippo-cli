# Tasks: Improve OAuth authentication flow

## Implementation Tasks

1. **Implement HTTP callback server** ✓
   - Create `setupCallbackServer` function in `init_save_drive_token.go`
   - Handle `/callback` endpoint
   - Verify state parameter for CSRF protection
   - Extract authorization code from query parameters
   - Return success HTML page
   - Use channels to communicate code/errors to main flow
   - **Validation**: Server accepts callback and extracts code

2. **Implement state parameter generation** ✓
   - Create `generateRandomState` function using `crypto/rand`
   - Generate 16-byte random value
   - Encode as base64 URL-safe string
   - **Validation**: Each call generates unique random state
   - **Note**: Task 1 depends on this task being completed first

3. **Implement browser opening** ✓
   - Create `openBrowser` function with platform detection
   - Support macOS (`open` command)
   - Support Linux (`xdg-open` command)
   - Support Windows (`cmd /c start` command)
   - Return error for unsupported platforms
   - **Validation**: Browser opens on supported platforms

4. **Implement server startup with port fallback** ✓
   - Attempt to start server on ports 21660-21669 sequentially
   - Try each port until one is available
   - If all ports are in use, return error indicating all ports are occupied
   - No manual code entry fallback (OAuth callback is required)
   - **Validation**: Server starts on first available port in range or exits
     with clear error message

5. **Refactor `getTokenFromWeb` function** ✓
   - Change signature to return `(*oauth2.Token, error)`
   - Implement port selection loop (21660-21669)
   - Configure redirect URI dynamically based on selected port
   - Use buffered `errChan := make(chan error, 1)` to prevent goroutine leak
   - Open browser automatically with the selected port's URL
   - Wait for callback with timeout (2 minutes)
   - Handle success/error/timeout cases
   - If all ports fail, exit with error (no manual fallback)
   - Add proper error handling and cleanup (defer server shutdown)
   - In callback handler, flush response and add 100ms delay before sending
     code to channel
   - **Validation**: Function works end-to-end with automatic callback and port
     fallback

6. **Add necessary imports** ✓
   - `net/http` for HTTP server
   - `crypto/rand` for random state generation
   - `encoding/base64` for state encoding
   - `os/exec` for browser command execution
   - `runtime` for platform detection
   - `time` for timeout
   - **Validation**: No import errors during compilation

## Testing Tasks

1. **Test happy path**
   - Run `nippo init` with no existing token
   - Verify browser opens automatically
   - Complete authorization in browser
   - Verify automatic redirect to localhost
   - Verify success page displays
   - Verify token is saved
   - **Validation**: Complete flow works without manual intervention

2. **Test browser auto-open failure**
   - Mock browser opening to fail
   - Verify URL is displayed in terminal
   - Verify instructions for manual opening
   - Complete authorization manually
   - **Validation**: Flow completes with manual browser opening

3. **Test port fallback**
   - Start a service on port 21660
   - Run `nippo init`
   - Verify command uses port 21661 instead
   - Complete authorization successfully on port 21661
   - Stop the service and verify next `nippo init` uses port 21660 again
   - **Validation**: Port fallback works correctly

4. **Test timeout**
   - Run `nippo init`
   - Wait without completing authorization
   - Verify timeout occurs after 2 minutes
   - Verify server shuts down
   - Verify error message is clear
   - **Validation**: Timeout handling works correctly

5. **Test state mismatch**
   - Manually craft callback with wrong state parameter
   - Verify callback is rejected
   - Verify error page in browser
   - Verify error in terminal
   - **Validation**: CSRF protection works

6. **Test concurrent authorization sessions**
   - Run `nippo init` in first terminal (should use port 21660)
   - While first is waiting, run `nippo init` in second terminal (should use
     port 21661)
   - Verify both sessions proceed independently
   - Complete authorization in both terminals
   - Verify both sessions complete successfully
   - **Validation**: Multiple concurrent sessions work correctly

7. **Test all ports exhausted**
   - Start services on all ports 21660-21669
   - Run `nippo init`
   - Verify command exits with error indicating all ports are in use
   - **Validation**: Clear error when all ports are occupied

8. **Test platform-specific browser opening**
   - Test on macOS (uses `open`)
   - Test on Linux (uses `xdg-open`)
   - Test on Windows (uses `cmd /c start`)
   - **Validation**: Browser opens on each platform

9. **Test success page rendering**
   - Complete authorization
   - Verify HTML success page displays in browser
   - Verify page content is readable
   - Verify styling is applied
   - **Validation**: Success page looks good

## Documentation Tasks

1. **Update README if needed**
   - Document the improved OAuth flow
   - Mention automatic browser opening
   - Note that ports 21660-21669 are used with automatic fallback
   - Document that manual code entry is no longer supported
   - Mention support for concurrent authorization sessions
   - **Validation**: Documentation reflects new behavior

2. **Update Google Cloud Console setup instructions**
   - Document that users must configure `http://localhost/callback` (without
     port number) **before** using the improved flow
   - Include step-by-step instructions in error messages or README
   - Explain that Google OAuth accepts localhost without specific port, enabling
     dynamic port selection
   - Explain that without this configuration, the command will fail
   - Document that at least one port in range 21660-21669 must be available
   - **Validation**: Instructions are complete and accurate

## Dependencies

- Task 6 (imports) must complete first
- Tasks 2, 3, 4 can be implemented in parallel
- Task 1 (setupCallbackServer) depends on task 2 (generateRandomState)
- Task 5 (refactor getTokenFromWeb) depends on tasks 1-4 being complete
- Testing tasks depend on all implementation tasks being complete
- Documentation tasks can run in parallel with testing
