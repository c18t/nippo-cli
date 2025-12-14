package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

// TestEnv holds the test environment configuration
type TestEnv struct {
	TmpDir    string
	ConfigDir string
	DataDir   string
	CacheDir  string
	cleanup   func()
}

// SetupTestEnv creates an isolated test environment that prevents
// tests from reading or writing to the user's actual config files.
// It returns a TestEnv with paths to temporary directories and a cleanup function.
//
// Usage:
//
//	func TestSomething(t *testing.T) {
//	    env := core.SetupTestEnv(t)
//	    defer env.Cleanup()
//	    // ... test code ...
//	}
func SetupTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	// Reset viper state to prevent contamination from previous tests
	viper.Reset()

	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "nippo_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	configDir := filepath.Join(tmpDir, ".config", "nippo")
	dataDir := filepath.Join(tmpDir, ".local", "share", "nippo")
	cacheDir := filepath.Join(tmpDir, ".cache", "nippo")

	// Create directories
	for _, dir := range []string{configDir, dataDir, cacheDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			_ = os.RemoveAll(tmpDir)
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Set XDG environment variables to point to temp directories
	// t.Setenv automatically restores the original value after the test
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, ".config"))
	t.Setenv("XDG_DATA_HOME", filepath.Join(tmpDir, ".local", "share"))
	t.Setenv("XDG_CACHE_HOME", filepath.Join(tmpDir, ".cache"))

	// Initialize a fresh config
	Cfg = &Config{}
	Cfg.Paths.DataDir = dataDir
	Cfg.Paths.CacheDir = cacheDir

	env := &TestEnv{
		TmpDir:    tmpDir,
		ConfigDir: configDir,
		DataDir:   dataDir,
		CacheDir:  cacheDir,
	}

	env.cleanup = func() {
		viper.Reset()
		Cfg = nil
		_ = os.RemoveAll(tmpDir)
	}

	return env
}

// Cleanup removes the temporary test environment.
// This should be called with defer after SetupTestEnv.
func (e *TestEnv) Cleanup() {
	if e.cleanup != nil {
		e.cleanup()
	}
}

// CreateConfigFile creates a test config file with the given content.
func (e *TestEnv) CreateConfigFile(t *testing.T, content string) string {
	t.Helper()
	configPath := filepath.Join(e.ConfigDir, "nippo.toml")
	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	return configPath
}

// CreateCredentialsFile creates a test credentials.json file.
func (e *TestEnv) CreateCredentialsFile(t *testing.T, content string) string {
	t.Helper()
	credPath := filepath.Join(e.DataDir, "credentials.json")
	if err := os.WriteFile(credPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create credentials file: %v", err)
	}
	return credPath
}

// CreateTokenFile creates a test token.json file.
func (e *TestEnv) CreateTokenFile(t *testing.T, content string) string {
	t.Helper()
	tokenPath := filepath.Join(e.DataDir, "token.json")
	if err := os.WriteFile(tokenPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create token file: %v", err)
	}
	return tokenPath
}
