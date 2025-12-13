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
