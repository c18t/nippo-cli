package interactor

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type initSaveDriveTokenInteractor struct {
	presenter presenter.InitSaveDriveTokenPresenter
}

func NewInitSaveDriveTokenInteractor(i do.Injector) (port.InitSaveDriveTokenUseCase, error) {
	p, err := do.Invoke[presenter.InitSaveDriveTokenPresenter](i)
	if err != nil {
		return nil, err
	}
	return &initSaveDriveTokenInteractor{
		presenter: p,
	}, nil
}

func (u *initSaveDriveTokenInteractor) Handle(input *port.InitSaveDriveTokenUseCaseInputData) {
	output := &port.InitSaveDriveTokenUsecaseOutputData{}

	dataDir := core.Cfg.GetDataDir()

	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		u.presenter.Suspend(fmt.Errorf("failed to create data directory: %w", err))
		return
	}

	credPath := filepath.Join(dataDir, "credentials.json")
	b, err := os.ReadFile(credPath)
	if err != nil {
		if os.IsNotExist(err) {
			u.presenter.Suspend(fmt.Errorf(`credentials.json not found

Please download the OAuth 2.0 Client ID credentials from Google Cloud Console:

1. Go to https://console.cloud.google.com/apis/credentials
2. Create OAuth 2.0 Client ID (Application type: Desktop app)
3. Download the credentials JSON file
4. Save it to: %s`, credPath))
		} else {
			u.presenter.Suspend(fmt.Errorf("unable to read credentials file: %w", err))
		}
		return
	}

	oauthConfig, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		u.presenter.Suspend(fmt.Errorf("unable to parse client secret file to config: %v", err))
		return
	}

	tok, err := getTokenFromWeb(oauthConfig, u.presenter, output)
	if err != nil {
		u.presenter.Suspend(fmt.Errorf("unable to get token from web: %w", err))
		return
	}

	// Save token (spinner continues during save)
	output.Message = "Saving credentials..."
	u.presenter.Progress(output)
	if err := saveToken(filepath.Join(dataDir, "token.json"), tok); err != nil {
		u.presenter.Suspend(fmt.Errorf("unable to save token: %w", err))
		return
	}

	output.Message = "Successfully authenticated with Google Drive"
	u.presenter.Complete(output)
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}
	defer func() { _ = f.Close() }()
	if err := json.NewEncoder(f).Encode(token); err != nil {
		return fmt.Errorf("unable to encode token: %w", err)
	}
	return nil
}

// generateRandomState generates a random state parameter for CSRF protection.
func generateRandomState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// openBrowser attempts to open the URL in the user's default browser.
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

// setupCallbackServer sets up an HTTP server to receive the OAuth callback.
func setupCallbackServer(port int, expectedState string, codeChan chan string, errChan chan error) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Verify state to prevent CSRF
		if r.URL.Query().Get("state") != expectedState {
			errChan <- fmt.Errorf("state mismatch")
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no authorization code")
			http.Error(w, "No authorization code received", http.StatusBadRequest)
			return
		}

		// Send success page
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprint(w, successHTML)

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
		Addr:    fmt.Sprintf("127.0.0.1:%d", port), // Bind to localhost only for security
		Handler: mux,
	}
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config, presenter presenter.InitSaveDriveTokenPresenter, output *port.InitSaveDriveTokenUsecaseOutputData) (*oauth2.Token, error) {
	// Try ports 21660-21669 to find an available one
	const basePort = 21660
	const maxAttempts = 10

	state := generateRandomState()
	codeChan := make(chan string)
	errChan := make(chan error, 1) // Buffered to prevent goroutine leak

	var server *http.Server
	var actualPort int
	var serverStarted bool

	// Try to find an available port
	for i := 0; i < maxAttempts; i++ {
		port := basePort + i
		server = setupCallbackServer(port, state, codeChan, errChan)

		// Try to start server
		go func() {
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				errChan <- err
			}
		}()

		// Give the server a moment to start or fail
		time.Sleep(50 * time.Millisecond)

		// Check if server started successfully
		select {
		case <-errChan:
			// Port is in use, try next one
			if i == maxAttempts-1 {
				return nil, fmt.Errorf("unable to start callback server: all ports (21660-21669) are in use")
			}
			continue
		default:
			// Server started successfully
			actualPort = port
			serverStarted = true
		}

		// Break out of loop if server started
		if serverStarted {
			break
		}
	}

	if !serverStarted {
		return nil, fmt.Errorf("unable to start callback server")
	}

	defer func() {
		_ = server.Shutdown(context.Background())
	}()

	// Configure redirect URI with the actual port
	config.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", actualPort)
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Update spinner message with port number
	output.Message = fmt.Sprintf("Opening browser for authorization (using port %d)...", actualPort)
	presenter.Progress(output)

	// Open browser
	if err := openBrowser(authURL); err != nil {
		return nil, fmt.Errorf("failed to open browser automatically, please visit:\n%s", authURL)
	}

	// Wait for callback, timeout, or cancellation
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(2 * time.Minute)

	for {
		select {
		case code := <-codeChan:
			tok, err := config.Exchange(context.Background(), code)
			if err != nil {
				return nil, fmt.Errorf("unable to exchange authorization code: %w", err)
			}
			return tok, nil
		case err := <-errChan:
			return nil, fmt.Errorf("callback server error: %w", err)
		case <-timeout:
			return nil, fmt.Errorf("authorization timeout (2 minutes)")
		case <-ticker.C:
			// Check if user cancelled
			if presenter.IsCancelled() {
				return nil, fmt.Errorf("authorization cancelled by user")
			}
		}
	}
}
