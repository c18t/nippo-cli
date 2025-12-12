# Proposal: Improve OAuth authentication flow with automatic callback

## Why

The current OAuth authentication flow for Google Drive integration requires
manual intervention that creates friction and confusion:

1. **Manual code extraction**: After user authorizes the app in their browser,
   they are redirected to a non-existent localhost URL. Users must manually
   copy the authorization code from the URL and paste it into the terminal.
2. **Poor user experience**: The process is error-prone (typos, copying wrong
   part of URL) and confusing for non-technical users.
3. **No browser automation**: Users must manually open the authorization URL in
   their browser.

## What Changes

### Root Cause

In `internal/usecase/interactor/init_save_drive_token.go`, the
`getTokenFromWeb` function:

1. Prints the authorization URL to the terminal (line 68-70)
2. Waits for manual input using `fmt.Scan` (line 72-75)
3. Does not start a local HTTP server to receive the OAuth callback
4. OAuth config uses a redirect URI that points to a non-existent localhost
   server

### Proposed Solution

Implement automatic OAuth callback handling by:

1. **Start local HTTP server**: Launch a temporary HTTP server on localhost
   (ports 21660-21669) to receive the OAuth callback
2. **Auto-open browser**: Automatically open the authorization URL in the
   user's default browser (with fallback to manual URL display)
3. **Receive callback automatically**: Extract the authorization code from the
   callback request
4. **Display success page**: Show a simple success page in the browser after
   successful authentication
5. **Graceful timeout**: Implement timeout (e.g., 2 minutes) with appropriate
   error handling
6. **Dynamic port selection**: Try ports 21660-21669 sequentially to support
   concurrent sessions

This pattern is commonly used by CLI tools like `gh` (GitHub CLI), `gcloud`
(Google Cloud SDK), and others.

### Scope

- **Capability**: `oauth-flow` (new)
- **Modified files**:
  - `internal/usecase/interactor/init_save_drive_token.go` - Implement
    localhost callback server
- **New dependencies**: Uses standard library only (`net/http` for HTTP server,
  `os/exec` for browser opening)
- **Prerequisites**: Users must configure `http://localhost/callback` (without
  a specific port number) as the authorized redirect URI in their Google Cloud
  Console OAuth 2.0 Client ID settings. Google OAuth accepts `http://localhost`
  without a port, enabling dynamic port selection at runtime.

## Non-Goals

- Implementing OAuth for other services (only Google Drive)
- Changing the OAuth scopes or permissions
- Adding OAuth token refresh logic (already handled by oauth2 library)
- Supporting OAuth flows other than authorization code flow

## Success Criteria

- User runs `nippo init` and their browser automatically opens to Google
  authorization page
- After user authorizes, they are redirected to localhost and see a success
  message
- Authorization code is automatically extracted and token is saved
- If browser auto-open fails, clear instructions are shown with the URL
- If preferred port (21660) is in use, system automatically tries ports
  21661-21669
- Multiple concurrent `nippo init` sessions work independently on different
  ports
- If all ports (21660-21669) are in use, command exits with clear error message
- Process times out gracefully after 2 minutes if user doesn't complete
  authorization
