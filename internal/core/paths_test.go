package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde expansion",
			input:    "~/test",
			expected: filepath.Join(os.Getenv("HOME"), "test"),
		},
		{
			name:     "no expansion needed",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "relative path",
			input:    "relative/path",
			expected: "relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.input)
			if result != tt.expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsUnderGitRepo(t *testing.T) {
	// Create temp directory structure
	tmpDir, err := os.MkdirTemp("", "paths_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a .git directory in tmpDir
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "directory with .git",
			path:     tmpDir,
			expected: true,
		},
		{
			name:     "subdirectory of git repo",
			path:     subDir,
			expected: true,
		},
		{
			name:     "non-existent path",
			path:     "/nonexistent/path",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUnderGitRepo(tt.path)
			if result != tt.expected {
				t.Errorf("IsUnderGitRepo(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsPathSafe(t *testing.T) {
	tests := []struct {
		name       string
		basePath   string
		targetPath string
		expected   bool
	}{
		{
			name:       "safe path within base",
			basePath:   "/base/dir",
			targetPath: "/base/dir/file.txt",
			expected:   true,
		},
		{
			name:       "safe nested path",
			basePath:   "/base/dir",
			targetPath: "/base/dir/sub/file.txt",
			expected:   true,
		},
		{
			name:       "unsafe path escaping base",
			basePath:   "/base/dir",
			targetPath: "/base/other/file.txt",
			expected:   false,
		},
		{
			name:       "unsafe path with parent traversal",
			basePath:   "/base/dir",
			targetPath: "/base/dir/../other/file.txt",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPathSafe(tt.basePath, tt.targetPath)
			if result != tt.expected {
				t.Errorf("IsPathSafe(%q, %q) = %v, want %v", tt.basePath, tt.targetPath, result, tt.expected)
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	homeDir := os.Getenv("HOME")

	tests := []struct {
		name     string
		path     string
		baseDir  string
		expected string
	}{
		{
			name:     "absolute path",
			path:     "/absolute/path",
			baseDir:  "/base",
			expected: "/absolute/path",
		},
		{
			name:     "relative path",
			path:     "relative/path",
			baseDir:  "/base",
			expected: "/base/relative/path",
		},
		{
			name:     "tilde path",
			path:     "~/mydir",
			baseDir:  "/base",
			expected: filepath.Join(homeDir, "mydir"),
		},
		{
			name:     "empty path",
			path:     "",
			baseDir:  "/base",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolvePath(tt.path, tt.baseDir)
			if result != tt.expected {
				t.Errorf("ResolvePath(%q, %q) = %q, want %q", tt.path, tt.baseDir, result, tt.expected)
			}
		})
	}
}

func TestExpandPath_TildeOnly(t *testing.T) {
	homeDir := os.Getenv("HOME")
	result := ExpandPath("~")
	if result != homeDir {
		t.Errorf("ExpandPath(\"~\") = %q, want %q", result, homeDir)
	}
}

func TestExpandPath_EmptyPath(t *testing.T) {
	result := ExpandPath("")
	if result != "" {
		t.Errorf("ExpandPath(\"\") = %q, want %q", result, "")
	}
}

func TestResolveConfigDir(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		contains string
	}{
		{
			name: "XDG_CONFIG_HOME set",
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "/xdg/config",
			},
			contains: "nippo",
		},
		{
			name: "fallback to default",
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "",
			},
			contains: "nippo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			result := ResolveConfigDir()
			if result == "" {
				t.Error("ResolveConfigDir() should return non-empty path")
			}
			if !filepath.IsAbs(result) {
				t.Errorf("ResolveConfigDir() = %q, should be absolute path", result)
			}
		})
	}
}

func TestResolveDataDir(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		contains string
	}{
		{
			name: "XDG_DATA_HOME set",
			envVars: map[string]string{
				"XDG_DATA_HOME": "/xdg/data",
			},
			contains: "nippo",
		},
		{
			name: "fallback to default",
			envVars: map[string]string{
				"XDG_DATA_HOME": "",
			},
			contains: "nippo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			result := ResolveDataDir()
			if result == "" {
				t.Error("ResolveDataDir() should return non-empty path")
			}
		})
	}
}

func TestResolveCacheDir(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
	}{
		{
			name: "XDG_CACHE_HOME set",
			envVars: map[string]string{
				"XDG_CACHE_HOME": "/xdg/cache",
			},
		},
		{
			name: "fallback to default",
			envVars: map[string]string{
				"XDG_CACHE_HOME": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			result := ResolveCacheDir()
			if result == "" {
				t.Error("ResolveCacheDir() should return non-empty path")
			}
		})
	}
}

func TestSafeOpen(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "paths_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		basePath  string
		filePath  string
		expectErr bool
	}{
		{
			name:      "safe path",
			basePath:  tmpDir,
			filePath:  testFile,
			expectErr: false,
		},
		{
			name:      "path traversal attack",
			basePath:  tmpDir,
			filePath:  filepath.Join(tmpDir, "..", "outside.txt"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := SafeOpen(tt.basePath, tt.filePath)
			if tt.expectErr {
				if err == nil {
					t.Error("SafeOpen() should return error for unsafe path")
				}
			} else {
				if err != nil {
					t.Errorf("SafeOpen() error = %v", err)
				}
				if f != nil {
					_ = f.Close()
				}
			}
		})
	}
}

func TestSafeOpenFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "paths_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tests := []struct {
		name      string
		basePath  string
		filePath  string
		expectErr bool
	}{
		{
			name:      "safe path for new file",
			basePath:  tmpDir,
			filePath:  filepath.Join(tmpDir, "newfile.txt"),
			expectErr: false,
		},
		{
			name:      "path traversal attack",
			basePath:  tmpDir,
			filePath:  filepath.Join(tmpDir, "..", "outside.txt"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := SafeOpenFile(tt.basePath, tt.filePath, os.O_CREATE|os.O_WRONLY, 0644)
			if tt.expectErr {
				if err == nil {
					t.Error("SafeOpenFile() should return error for unsafe path")
				}
			} else {
				if err != nil {
					t.Errorf("SafeOpenFile() error = %v", err)
				}
				if f != nil {
					_ = f.Close()
				}
			}
		})
	}
}
