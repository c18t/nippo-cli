package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/samber/do/v2"
)

// mockDirEntry implements os.DirEntry for testing
type mockDirEntry struct {
	name  string
	isDir bool
}

func (m *mockDirEntry) Name() string               { return m.name }
func (m *mockDirEntry) IsDir() bool                { return m.isDir }
func (m *mockDirEntry) Type() os.FileMode          { return 0 }
func (m *mockDirEntry) Info() (os.FileInfo, error) { return nil, nil }

func TestNewAssetRepository(t *testing.T) {
	mock := &mockLocalFileProvider{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (gateway.LocalFileProvider, error) {
		return mock, nil
	})

	repo, err := NewAssetRepository(injector)
	if err != nil {
		t.Errorf("NewAssetRepository() error = %v", err)
	}
	if repo == nil {
		t.Error("NewAssetRepository() returned nil")
	}
}

func TestAssetRepository_CleanNippoCache(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "asset_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create cache directory and test files
	cacheDir := filepath.Join(tmpDir, "md")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFiles := []string{"2024-01-15.md", "2024-01-16.md"}
	for _, name := range testFiles {
		if err := os.WriteFile(filepath.Join(cacheDir, name), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Set up global config
	core.Cfg = &core.Config{}
	core.Cfg.Paths.CacheDir = tmpDir

	// Create mock that returns the actual files
	mock := &mockLocalFileProvider{
		entries: []os.DirEntry{
			&mockDirEntry{name: "2024-01-15.md"},
			&mockDirEntry{name: "2024-01-16.md"},
		},
	}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (gateway.LocalFileProvider, error) {
		return mock, nil
	})

	repo, _ := NewAssetRepository(injector)
	err = repo.CleanNippoCache()
	if err != nil {
		t.Errorf("CleanNippoCache() error = %v", err)
	}

	// Verify files were deleted
	for _, name := range testFiles {
		if _, err := os.Stat(filepath.Join(cacheDir, name)); !os.IsNotExist(err) {
			t.Errorf("File %s should have been deleted", name)
		}
	}
}

func TestAssetRepository_CleanBuildCache(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "asset_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create output directory and test files
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFiles := []string{"index.html", "page.html"}
	for _, name := range testFiles {
		if err := os.WriteFile(filepath.Join(outputDir, name), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Set up global config
	core.Cfg = &core.Config{}
	core.Cfg.Paths.CacheDir = tmpDir

	// Create mock that returns the actual files
	mock := &mockLocalFileProvider{
		entries: []os.DirEntry{
			&mockDirEntry{name: "index.html"},
			&mockDirEntry{name: "page.html"},
		},
	}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (gateway.LocalFileProvider, error) {
		return mock, nil
	})

	repo, _ := NewAssetRepository(injector)
	err = repo.CleanBuildCache()
	if err != nil {
		t.Errorf("CleanBuildCache() error = %v", err)
	}

	// Verify files were deleted
	for _, name := range testFiles {
		if _, err := os.Stat(filepath.Join(outputDir, name)); !os.IsNotExist(err) {
			t.Errorf("File %s should have been deleted", name)
		}
	}
}

func TestAssetRepository_CleanEmptyDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "asset_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Set up global config
	core.Cfg = &core.Config{}
	core.Cfg.Paths.CacheDir = tmpDir

	// Create mock with no files
	mock := &mockLocalFileProvider{
		entries: []os.DirEntry{},
	}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (gateway.LocalFileProvider, error) {
		return mock, nil
	})

	repo, _ := NewAssetRepository(injector)
	err = repo.CleanNippoCache()
	if err != nil {
		t.Errorf("CleanNippoCache() with empty dir error = %v", err)
	}
}
