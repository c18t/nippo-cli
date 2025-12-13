package core

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ResolveConfigDir resolves the configuration directory using fallback chain:
// 1. XDG_CONFIG_HOME environment variable (all platforms)
// 2. Windows APPDATA (Windows only, when XDG not set)
// 3. Default XDG path (~/.config)
func ResolveConfigDir() string {
	// Priority 1: XDG environment variable
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" && filepath.IsAbs(xdgConfig) {
		return filepath.Join(xdgConfig, "nippo")
	}

	// Priority 2: Windows environment variable (Windows only)
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "nippo")
		}
	}

	// Priority 3: Default XDG path
	home := homeDir()
	return filepath.Join(home, ".config", "nippo")
}

// ResolveDataDir resolves the data directory using fallback chain:
// 1. XDG_DATA_HOME environment variable (all platforms)
// 2. Windows LOCALAPPDATA (Windows only, when XDG not set)
// 3. Default XDG path (~/.local/share)
func ResolveDataDir() string {
	// Priority 1: XDG environment variable
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" && filepath.IsAbs(xdgData) {
		return filepath.Join(xdgData, "nippo")
	}

	// Priority 2: Windows environment variable (Windows only)
	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return filepath.Join(localAppData, "nippo")
		}
	}

	// Priority 3: Default XDG path
	home := homeDir()
	return filepath.Join(home, ".local", "share", "nippo")
}

// ResolveCacheDir resolves the cache directory using fallback chain:
// 1. XDG_CACHE_HOME environment variable (all platforms)
// 2. Windows LOCALAPPDATA\cache (Windows only, when XDG not set)
// 3. Default XDG path (~/.cache)
func ResolveCacheDir() string {
	// Priority 1: XDG environment variable
	if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" && filepath.IsAbs(xdgCache) {
		return filepath.Join(xdgCache, "nippo")
	}

	// Priority 2: Windows environment variable (Windows only)
	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return filepath.Join(localAppData, "nippo", "cache")
		}
	}

	// Priority 3: Default XDG path
	home := homeDir()
	return filepath.Join(home, ".cache", "nippo")
}

// ExpandPath expands environment variables and tilde in a path.
// - Environment variables: $VAR or %VAR% (cross-platform via os.ExpandEnv)
// - Tilde: ~ expands to user's home directory
func ExpandPath(path string) string {
	if path == "" {
		return path
	}

	// Expand environment variables (handles both $VAR and %VAR%)
	path = os.ExpandEnv(path)

	// Expand tilde
	if strings.HasPrefix(path, "~") {
		home := homeDir()
		if path == "~" {
			return home
		}
		if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
			return filepath.Join(home, path[2:])
		}
	}

	return path
}

// ResolvePath resolves a path that may be relative to a base directory.
// - Absolute paths are returned as-is (after expansion)
// - Relative paths are resolved relative to baseDir
// - Environment variables and tilde are expanded first
func ResolvePath(path, baseDir string) string {
	if path == "" {
		return path
	}

	// First expand environment variables and tilde
	expanded := ExpandPath(path)

	// If already absolute, return as-is
	if filepath.IsAbs(expanded) {
		return expanded
	}

	// Resolve relative to base directory
	return filepath.Join(baseDir, expanded)
}

// IsUnderGitRepo checks if the given path is inside a git repository
// by traversing parent directories looking for .git directory.
func IsUnderGitRepo(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	for {
		gitPath := filepath.Join(absPath, ".git")
		if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
			return true
		}
		parent := filepath.Dir(absPath)
		if parent == absPath { // reached root
			return false
		}
		absPath = parent
	}
}

// IsPathSafe checks if targetPath is safely contained within basePath.
// This prevents Zip Slip vulnerabilities by ensuring the target doesn't
// escape the base directory via ../ or similar path traversal.
func IsPathSafe(basePath, targetPath string) bool {
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return false
	}
	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return false
	}
	// Ensure target starts with base path followed by separator
	return strings.HasPrefix(absTarget, absBase+string(filepath.Separator)) || absTarget == absBase
}

// homeDir returns the user's home directory
func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}
