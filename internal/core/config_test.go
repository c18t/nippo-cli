package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestConfig_GetDataDir(t *testing.T) {
	homeDir := os.Getenv("HOME")

	tests := []struct {
		name       string
		config     Config
		envVars    map[string]string
		expected   string
		skipOnGOOS string
	}{
		{
			name: "configured absolute path",
			config: Config{
				Paths: ConfigPaths{
					DataDir: "/custom/data/dir",
				},
			},
			expected: "/custom/data/dir",
		},
		{
			name: "configured tilde path",
			config: Config{
				Paths: ConfigPaths{
					DataDir: "~/nippo-data",
				},
			},
			expected: filepath.Join(homeDir, "nippo-data"),
		},
		{
			name: "configured relative path resolves to config dir",
			config: Config{
				configDir: "/config/dir",
				Paths: ConfigPaths{
					DataDir: "data",
				},
			},
			expected: "/config/dir/data",
		},
		{
			name: "fallback to XDG_DATA_HOME",
			config: Config{
				Paths: ConfigPaths{
					DataDir: "", // not configured
				},
			},
			envVars: map[string]string{
				"XDG_DATA_HOME": "/xdg/data",
			},
			expected: "/xdg/data/nippo",
		},
		{
			name: "fallback to default XDG path when nothing configured",
			config: Config{
				Paths: ConfigPaths{
					DataDir: "", // not configured
				},
			},
			envVars: map[string]string{
				"XDG_DATA_HOME": "", // clear XDG
			},
			expected: filepath.Join(homeDir, ".local", "share", "nippo"),
		},
		{
			name: "cached value is returned",
			config: Config{
				dataDir: "/cached/data/dir",
				Paths: ConfigPaths{
					DataDir: "/should/be/ignored",
				},
			},
			expected: "/cached/data/dir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipOnGOOS != "" && tt.skipOnGOOS == os.Getenv("GOOS") {
				t.Skip("Skipping on " + tt.skipOnGOOS)
			}

			// Set test environment variables (t.Setenv auto-restores after test)
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg := tt.config
			result := cfg.GetDataDir()
			if result != tt.expected {
				t.Errorf("GetDataDir() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestConfig_GetCacheDir(t *testing.T) {
	homeDir := os.Getenv("HOME")

	tests := []struct {
		name       string
		config     Config
		envVars    map[string]string
		expected   string
		skipOnGOOS string
	}{
		{
			name: "configured absolute path",
			config: Config{
				Paths: ConfigPaths{
					CacheDir: "/custom/cache/dir",
				},
			},
			expected: "/custom/cache/dir",
		},
		{
			name: "configured tilde path",
			config: Config{
				Paths: ConfigPaths{
					CacheDir: "~/nippo-cache",
				},
			},
			expected: filepath.Join(homeDir, "nippo-cache"),
		},
		{
			name: "configured relative path resolves to config dir",
			config: Config{
				configDir: "/config/dir",
				Paths: ConfigPaths{
					CacheDir: "cache",
				},
			},
			expected: "/config/dir/cache",
		},
		{
			name: "fallback to XDG_CACHE_HOME",
			config: Config{
				Paths: ConfigPaths{
					CacheDir: "", // not configured
				},
			},
			envVars: map[string]string{
				"XDG_CACHE_HOME": "/xdg/cache",
			},
			expected: "/xdg/cache/nippo",
		},
		{
			name: "fallback to default XDG path when nothing configured",
			config: Config{
				Paths: ConfigPaths{
					CacheDir: "", // not configured
				},
			},
			envVars: map[string]string{
				"XDG_CACHE_HOME": "", // clear XDG
			},
			expected: filepath.Join(homeDir, ".cache", "nippo"),
		},
		{
			name: "cached value is returned",
			config: Config{
				cacheDir: "/cached/cache/dir",
				Paths: ConfigPaths{
					CacheDir: "/should/be/ignored",
				},
			},
			expected: "/cached/cache/dir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipOnGOOS != "" && tt.skipOnGOOS == os.Getenv("GOOS") {
				t.Skip("Skipping on " + tt.skipOnGOOS)
			}

			// Set test environment variables (t.Setenv auto-restores after test)
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg := tt.config
			result := cfg.GetCacheDir()
			if result != tt.expected {
				t.Errorf("GetCacheDir() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestConfig_GetConfigDir(t *testing.T) {
	homeDir := os.Getenv("HOME")

	tests := []struct {
		name     string
		config   Config
		envVars  map[string]string
		expected string
	}{
		{
			name: "cached value is returned",
			config: Config{
				configDir: "/cached/config/dir",
			},
			expected: "/cached/config/dir",
		},
		{
			name: "fallback to XDG_CONFIG_HOME",
			config: Config{
				configDir: "", // not cached
			},
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "/xdg/config",
			},
			expected: "/xdg/config/nippo",
		},
		{
			name: "fallback to default XDG path when nothing configured",
			config: Config{
				configDir: "", // not cached
			},
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "", // clear XDG
			},
			expected: filepath.Join(homeDir, ".config", "nippo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test environment variables (t.Setenv auto-restores after test)
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg := tt.config
			result := cfg.GetConfigDir()
			if result != tt.expected {
				t.Errorf("GetConfigDir() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestErrConfigNotFound(t *testing.T) {
	err := &ErrConfigNotFound{Path: "/path/to/config.toml"}
	expected := "configuration file not found: /path/to/config.toml"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestConfig_EnvironmentVariableExpansion(t *testing.T) {
	// Set test environment variable (t.Setenv auto-restores after test)
	t.Setenv("NIPPO_TEST_VAR", "/expanded/path")

	cfg := Config{
		Paths: ConfigPaths{
			DataDir: "$NIPPO_TEST_VAR/data",
		},
	}

	result := cfg.GetDataDir()
	expected := "/expanded/path/data"
	if result != expected {
		t.Errorf("GetDataDir() with env var = %q, want %q", result, expected)
	}
}

func TestConfig_ResetLastUpdateCheckTimestamp(t *testing.T) {
	cfg := &Config{}
	cfg.LastUpdateCheckTimestamp = cfg.getDefaultLastUpdateCheckTimestamp().Add(24 * 60 * 60 * 1000000000)

	cfg.ResetLastUpdateCheckTimestamp()

	expected := cfg.getDefaultLastUpdateCheckTimestamp()
	if !cfg.LastUpdateCheckTimestamp.Equal(expected) {
		t.Errorf("ResetLastUpdateCheckTimestamp() = %v, want %v", cfg.LastUpdateCheckTimestamp, expected)
	}
}

func TestConfig_SaveConfig(t *testing.T) {
	// Reset viper state before and after test
	defer func() {
		viper.Reset()
	}()
	viper.Reset()

	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a config file first
	configPath := filepath.Join(tmpDir, "nippo.toml")
	configContent := `[project]
url = "https://github.com/c18t/nippo"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		configDir: tmpDir,
		Project: ConfigProject{
			Url:           "https://github.com/test/nippo",
			DriveFolderId: "test-folder",
			SiteUrl:       "https://test.example.com",
		},
	}

	// Load config first to set viper config file
	err = cfg.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Update config
	cfg.Project.Url = "https://github.com/updated/nippo"

	err = cfg.SaveConfig()
	if err != nil {
		t.Errorf("SaveConfig() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("SaveConfig() did not create config file")
	}
}

func TestConfig_configFieldMap(t *testing.T) {
	cfg := &Config{
		Project: ConfigProject{
			Url:           "https://github.com/test/nippo",
			DriveFolderId: "test-folder",
		},
		Paths: ConfigPaths{
			DataDir:  "/data",
			CacheDir: "/cache",
		},
	}

	cMap, err := cfg.configFieldMap(map[string]any{}, *cfg, "")
	if err != nil {
		t.Errorf("configFieldMap() error = %v", err)
	}

	if cMap["project.url"] != "https://github.com/test/nippo" {
		t.Errorf("configFieldMap() project.url = %v, want %v", cMap["project.url"], "https://github.com/test/nippo")
	}
	if cMap["project.drive_folder_id"] != "test-folder" {
		t.Errorf("configFieldMap() project.drive_folder_id = %v, want %v", cMap["project.drive_folder_id"], "test-folder")
	}
	if cMap["path.data_dir"] != "/data" {
		t.Errorf("configFieldMap() path.data_dir = %v, want %v", cMap["path.data_dir"], "/data")
	}
}

func TestGetConfigFilePath(t *testing.T) {
	result := GetConfigFilePath()
	if result == "" {
		t.Error("GetConfigFilePath() should return non-empty path")
	}
}
