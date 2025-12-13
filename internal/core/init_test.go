package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitConfig_NotFound(t *testing.T) {
	// Create a temp directory
	tmpDir, err := os.MkdirTemp("", "init_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Set XDG_CONFIG_HOME to temp directory (config shouldn't exist)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create nippo config directory (but not the config file)
	nippoDir := filepath.Join(tmpDir, "nippo")
	if err := os.MkdirAll(nippoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Don't pass explicit config path - let it use XDG_CONFIG_HOME fallback
	// This tests the ConfigFileNotFoundError path
	err = InitConfig("")
	if err == nil {
		t.Error("InitConfig should return error for non-existent config")
	}

	// Check error type
	_, ok := err.(*ErrConfigNotFound)
	if !ok {
		t.Errorf("InitConfig error should be *ErrConfigNotFound, got %T", err)
	}

	// Cfg should still be set (with defaults)
	if Cfg == nil {
		t.Error("Cfg should be set even when config not found")
	}
}

func TestInitConfig_Valid(t *testing.T) {
	// Create a temp directory with valid config
	tmpDir, err := os.MkdirTemp("", "init_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create nippo directory
	nippoDir := filepath.Join(tmpDir, "nippo")
	if err := os.MkdirAll(nippoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a valid config file
	configPath := filepath.Join(nippoDir, "nippo.toml")
	configContent := `
[project]
url = "https://github.com/c18t/nippo"
drive_folder_id = "test-folder-id"
site_url = "https://nippo.example.com"
branch = "main"
template_path = "/templates"
asset_path = "/dist"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	err = InitConfig(configPath)
	if err != nil {
		t.Errorf("InitConfig should not return error for valid config: %v", err)
	}

	if Cfg == nil {
		t.Fatal("Cfg should be set")
	}

	// Verify config values
	if Cfg.Project.Url != "https://github.com/c18t/nippo" {
		t.Errorf("Project.Url = %q, want %q", Cfg.Project.Url, "https://github.com/c18t/nippo")
	}
	if Cfg.Project.DriveFolderId != "test-folder-id" {
		t.Errorf("Project.DriveFolderId = %q, want %q", Cfg.Project.DriveFolderId, "test-folder-id")
	}
	if Cfg.Project.SiteUrl != "https://nippo.example.com" {
		t.Errorf("Project.SiteUrl = %q, want %q", Cfg.Project.SiteUrl, "https://nippo.example.com")
	}
}

func TestInitConfig_BackwardCompatibility(t *testing.T) {
	// Test that config without [path] section works
	tmpDir, err := os.MkdirTemp("", "init_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create nippo directory
	nippoDir := filepath.Join(tmpDir, "nippo")
	if err := os.MkdirAll(nippoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create config WITHOUT [path] section (legacy format)
	configPath := filepath.Join(nippoDir, "nippo.toml")
	configContent := `
[project]
url = "https://github.com/c18t/nippo"
drive_folder_id = "test-folder-id"
site_url = "https://nippo.example.com"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	err = InitConfig(configPath)
	if err != nil {
		t.Errorf("InitConfig should work without [path] section: %v", err)
	}

	if Cfg == nil {
		t.Fatal("Cfg should be set")
	}

	// Verify that Paths are empty (will use fallback)
	if Cfg.Paths.DataDir != "" {
		t.Errorf("Paths.DataDir should be empty, got %q", Cfg.Paths.DataDir)
	}
	if Cfg.Paths.CacheDir != "" {
		t.Errorf("Paths.CacheDir should be empty, got %q", Cfg.Paths.CacheDir)
	}

	// GetDataDir and GetCacheDir should still work (using fallbacks)
	dataDir := Cfg.GetDataDir()
	if dataDir == "" {
		t.Error("GetDataDir should return non-empty path using fallback")
	}

	cacheDir := Cfg.GetCacheDir()
	if cacheDir == "" {
		t.Error("GetCacheDir should return non-empty path using fallback")
	}
}

func TestInitConfig_WithCustomPaths(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "init_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create nippo directory
	nippoDir := filepath.Join(tmpDir, "nippo")
	if err := os.MkdirAll(nippoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create config WITH [path] section
	configPath := filepath.Join(nippoDir, "nippo.toml")
	configContent := `
[project]
url = "https://github.com/c18t/nippo"
drive_folder_id = "test-folder-id"
site_url = "https://nippo.example.com"

[path]
data_dir = "/custom/data/dir"
cache_dir = "/custom/cache/dir"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	err = InitConfig(configPath)
	if err != nil {
		t.Errorf("InitConfig should work with [path] section: %v", err)
	}

	if Cfg == nil {
		t.Fatal("Cfg should be set")
	}

	// Verify that custom paths are used
	if Cfg.GetDataDir() != "/custom/data/dir" {
		t.Errorf("GetDataDir = %q, want %q", Cfg.GetDataDir(), "/custom/data/dir")
	}
	if Cfg.GetCacheDir() != "/custom/cache/dir" {
		t.Errorf("GetCacheDir = %q, want %q", Cfg.GetCacheDir(), "/custom/cache/dir")
	}
}
