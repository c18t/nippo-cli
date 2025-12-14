# Design: Improve OAuth authentication flow with automatic callback

## Overview

This change implements an automatic OAuth callback mechanism using a temporary
local HTTP server, improving the user experience by eliminating manual
authorization code copying.

## Current Behavior

`internal/usecase/interactor/init_save_drive_token.go` handles OAuth as
follows:

```go
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the "+
        "authorization code: \n%v\n", authURL)

    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
        fmt.Printf("Unable to read authorization code %v\n", err)
    }

    tok, err := config.Exchange(context.TODO(), authCode)
    // ...
    return tok
}
```

Problems:

1. User must manually copy/paste authorization code from URL
2. No automatic browser opening
3. Error-prone and confusing UX

## Proposed Implementation

### Architecture

```text
┌─────────────┐
│   User CLI  │
└──────┬──────┘
       │ 1. Start local server
       ▼
┌──────────────────┐
│ HTTP Server      │◄────┐
│ localhost:21660  │     │ 4. GET /callback?code=xxx
└──────────────────┘     │
       │                 │
       │ 2. Open browser │
       ▼                 │
┌──────────────────┐     │
│   Browser        │─────┘
│ Google Auth Page │ 3. User authorizes
└──────────────────┘
```

### Core Implementation

```go
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
    // 1. Setup callback server
    state := generateRandomState()
    codeChan := make(chan string)
    errChan := make(chan error, 1)  // Buffered to prevent goroutine leak

    server := setupCallbackServer(state, codeChan, errChan)

    // 2. Configure redirect URI
    config.RedirectURL = "http://localhost:21660/callback"
    authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)

    // 3. Start server in background
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            errChan <- err
        }
    }()
    defer server.Shutdown(context.Background())

    // 4. Open browser
    fmt.Println("Opening browser for authorization...")
    if err := openBrowser(authURL); err != nil {
        fmt.Printf("Failed to open browser. Please visit:\n%s\n", authURL)
    }

    // 5. Wait for callback or timeout
    select {
    case code := <-codeChan:
        return config.Exchange(context.Background(), code)
    case err := <-errChan:
        return nil, fmt.Errorf("callback server error: %w", err)
    case <-time.After(2 * time.Minute):
        return nil, fmt.Errorf("authorization timeout")
    }
}

func setupCallbackServer(expectedState string, codeChan chan string,
                        errChan chan error) *http.Server {
    mux := http.NewServeMux()

    mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
        // Verify state to prevent CSRF
        if r.URL.Query().Get("state") != expectedState {
            errChan <- fmt.Errorf("state mismatch")
            http.Error(w, "Invalid state", http.StatusBadRequest)
            return
        }

        code := r.URL.Query().Get("code")
        if code == "" {
            errChan <- fmt.Errorf("no authorization code")
            http.Error(w, "No code", http.StatusBadRequest)
            return
        }

        // Send success page
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprintf(w, successHTML)

        // Flush response to ensure browser receives success page
        if f, ok := w.(http.Flusher); ok {
            f.Flush()
        }
        // Brief delay to ensure browser receives the page
        time.Sleep(100 * time.Millisecond)

        // Send code to main goroutine
        codeChan <- code
    })

    return &http.Server{
        Addr:    "127.0.0.1:21660",  // Bind to localhost only for security
        Handler: mux,
    }
}

func openBrowser(url string) error {
    var cmd string
    var args []string

    switch runtime.GOOS {
    case "darwin":
        cmd = "open"
        args = []string{url}
    case "linux":
        cmd = "xdg-open"
        args = []string{url}
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start", url}
    default:
        return fmt.Errorf("unsupported platform")
    }

    return exec.Command(cmd, args...).Start()
}

func generateRandomState() string {
    b := make([]byte, 16)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
```

### Success HTML Page

```go
const successHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Authorization Successful</title>
    <style>
        body { font-family: sans-serif; text-align: center; padding: 50px; }
        h1 { color: #4CAF50; }
    </style>
</head>
<body>
    <h1>✓ Authorization Successful</h1>
    <p>You can close this window and return to the terminal.</p>
</body>
</html>`
```

## Prerequisites

Users must configure their Google Cloud Console OAuth 2.0 Client ID settings
**before** using the improved OAuth flow:

1. Go to Google Cloud Console > APIs & Services > Credentials
2. Select or create an OAuth 2.0 Client ID (Application type: Desktop app)
3. Add `http://localhost/callback` to "Authorized redirect URIs"
   - Note: Google OAuth accepts `http://localhost` without a port number,
     which allows the application to use any port dynamically
4. Download the updated credentials JSON file
5. Save it as `~/.local/share/nippo/credentials.json`

**Important**: The redirect URI should be configured as `http://localhost/callback`
(without a specific port number). The application will automatically select an
available port from 21660-21669.

## Design Decisions

### Port Selection

**Decision**: Try ports 21660-21669 sequentially until an available port is
found

**Rationale**:

- Port 21660 is the preferred starting port, derived from "nippo" branding
  (ni=2, i=1, roku=6, roku=6, zero=0)
- Google OAuth accepts `http://localhost/callback` without a specific port
  number, allowing dynamic port selection at runtime
- Ports 21660-21669 are in the user/private port range and unlikely to conflict
  with common services
- Sequential port probing (up to 10 ports) handles edge cases where the
  preferred port is in use
- If all ports 21660-21669 are in use, the command exits with a clear error
  message
- This approach allows multiple concurrent `nippo init` sessions without
  conflicts

**Alternative Considered**: Use only fixed port 21660

- **Rejected**: Would prevent concurrent authorization sessions and require
  users to manually resolve port conflicts. The flexibility of trying multiple
  ports provides better user experience with minimal complexity.

### State Parameter for CSRF Protection

**Decision**: Generate random state parameter and verify it in callback

**Rationale**:

- Standard OAuth 2.0 security practice
- Prevents CSRF attacks
- Simple to implement with crypto/rand

### Browser Opening

**Decision**: Attempt automatic open with fallback to manual URL display

**Rationale**:

- Best UX when automatic opening works
- Graceful degradation when it doesn't
- Users can still complete flow manually

**Alternative Considered**: Always require manual URL opening

- **Rejected**: Poor UX, defeats purpose of improvement

### Timeout Duration

**Decision**: 2-minute timeout for authorization

**Rationale**:

- Long enough for user to read, authorize, and complete flow
- Short enough to not leave server running indefinitely
- Matches common CLI tool timeouts

### Server Cleanup

**Decision**: Use defer to ensure server shutdown

**Rationale**:

- Prevents resource leaks
- Ensures port is released
- Works even if authorization fails

### Error Handling for Server Failures

**Decision**: If all ports (21660-21669) are in use, exit with clear error
message

**Rationale**:

- At least one port in the range 21660-21669 must be available for OAuth
  callback
- Manual code entry is no longer supported (callback URL is required)
- Clear error message informs user that all ports are in use
- Extremely rare edge case (requires 10 concurrent sessions or port conflicts)

## Error Handling

1. **All ports in use**: Exit with error message if all ports (21660-21669) are
   occupied (e.g., "Unable to start callback server: all ports (21660-21669)
   are in use. Please stop any conflicting processes and try again.")
2. **Browser open fails**: Display URL for manual opening, callback server
   continues waiting on the selected port
3. **State mismatch**: Reject callback, show error in browser and terminal
4. **Timeout**: Close server, show error with instructions to retry
5. **Exchange fails**: Show error with token exchange failure details

## Security Considerations

1. **CSRF Protection**: State parameter prevents cross-site request forgery
2. **Localhost only**: Server only binds to localhost, not accessible from
   network
3. **Timeout**: Limited time window for authorization
4. **HTTPS not required**: OAuth library handles token security, localhost
   callback is local-only

## Testing Strategy

1. **Happy path**: Start server on port 21660, authorize, receive callback,
   exchange token
2. **Browser open failure**: Verify URL displayed for manual opening, server
   continues waiting
3. **Timeout**: Verify server closes and error is shown after 2 minutes
4. **State mismatch**: Verify rejection of invalid state parameter
5. **Port fallback**: Verify command uses port 21661 when port 21660 is in use
6. **Concurrent sessions**: Verify multiple `nippo init` sessions work
   simultaneously on different ports
7. **All ports exhausted**: Verify command exits with clear error when all ports
   (21660-21669) are in use

## Impact Analysis

- **Files Modified**: `internal/usecase/interactor/init_save_drive_token.go`
- **Breaking Changes**:
  - Internal breaking change - `getTokenFromWeb` function signature changes
    from `(*oauth2.Config) *oauth2.Token` to
    `(*oauth2.Config) (*oauth2.Token, error)`. This is a private function, so
    no external API impact.
  - Manual authorization code entry is removed; callback server is required
  - Users must update their Google Cloud Console OAuth settings with the new
    redirect URI
- **Behavioral Changes**:
  - Automatic browser opening and callback handling
  - Dynamic port selection (21660-21669) with automatic fallback
  - Supports concurrent authorization sessions
  - No fallback to manual code entry
- **Performance Impact**: Minimal (local HTTP server for ~2 minutes, plus up to
  500ms for port probing in worst case)
- **Dependencies**: No new external dependencies (uses standard library only)
- **User Prerequisites**: Users must configure `http://localhost/callback` in
  Google Cloud Console before using the improved flow

## Platform Compatibility

- **macOS**: Uses `open` command
- **Linux**: Uses `xdg-open` command
- **Windows**: Uses `start` command via `cmd`
- **Other**: Falls back to manual URL display

## Rollback Plan

If issues arise, revert the changes to `getTokenFromWeb` function and restore
the manual authorization code entry flow. Users will need to revert their
Google Cloud Console redirect URI settings to the previous configuration.
