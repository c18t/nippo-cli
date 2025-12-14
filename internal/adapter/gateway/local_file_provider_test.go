package gateway

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/samber/do/v2"
)

func TestNewLocalFileProvider(t *testing.T) {
	injector := do.New()
	provider, err := NewLocalFileProvider(injector)
	if err != nil {
		t.Errorf("NewLocalFileProvider() error = %v", err)
	}
	if provider == nil {
		t.Error("NewLocalFileProvider() returned nil")
	}
}

func TestLocalFileProvider_List(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gateway_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create test files
	testFiles := []string{"file1.md", "file2.md", "file3.txt"}
	for _, name := range testFiles {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create a subdirectory (should be skipped)
	if err := os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}

	provider := &localFileProvider{}

	tests := []struct {
		name           string
		param          *repository.QueryListParam
		expectedCount  int
		expectErr      bool
	}{
		{
			name: "list all files",
			param: &repository.QueryListParam{
				Folders: []string{tmpDir},
			},
			expectedCount: 3,
			expectErr:     false,
		},
		{
			name: "list only md files",
			param: &repository.QueryListParam{
				Folders:        []string{tmpDir},
				FileExtensions: []string{"md"},
			},
			expectedCount: 2,
			expectErr:     false,
		},
		{
			name: "list txt files only",
			param: &repository.QueryListParam{
				Folders:        []string{tmpDir},
				FileExtensions: []string{"txt"},
			},
			expectedCount: 1,
			expectErr:     false,
		},
		{
			name: "non-existent directory",
			param: &repository.QueryListParam{
				Folders: []string{"/nonexistent/path"},
			},
			expectedCount: 0,
			expectErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := provider.List(tt.param)
			if tt.expectErr {
				if err == nil {
					t.Error("List() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("List() error = %v", err)
				return
			}
			if len(files) != tt.expectedCount {
				t.Errorf("List() returned %d files, want %d", len(files), tt.expectedCount)
			}
		})
	}
}

func TestLocalFileProvider_Write(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gateway_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	provider := &localFileProvider{}

	tests := []struct {
		name      string
		filePath  string
		content   []byte
		expectErr bool
	}{
		{
			name:      "write new file",
			filePath:  filepath.Join(tmpDir, "newfile.txt"),
			content:   []byte("test content"),
			expectErr: false,
		},
		{
			name:      "write to nested directory",
			filePath:  filepath.Join(tmpDir, "nested", "dir", "file.txt"),
			content:   []byte("nested content"),
			expectErr: false,
		},
		{
			name:      "overwrite existing file",
			filePath:  filepath.Join(tmpDir, "newfile.txt"),
			content:   []byte("updated content"),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provider.Write(tt.filePath, tt.content)
			if tt.expectErr {
				if err == nil {
					t.Error("Write() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Write() error = %v", err)
				return
			}

			// Verify file was written
			content, readErr := os.ReadFile(tt.filePath)
			if readErr != nil {
				t.Errorf("Failed to read written file: %v", readErr)
				return
			}
			if string(content) != string(tt.content) {
				t.Errorf("File content = %q, want %q", string(content), string(tt.content))
			}
		})
	}
}

func TestLocalFileProvider_Copy(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gateway_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	srcContent := []byte("source content")
	if err := os.WriteFile(srcPath, srcContent, 0644); err != nil {
		t.Fatal(err)
	}

	provider := &localFileProvider{}

	tests := []struct {
		name      string
		baseDir   string
		destPath  string
		srcPath   string
		expectErr bool
	}{
		{
			name:      "copy file to same directory",
			baseDir:   tmpDir,
			destPath:  filepath.Join(tmpDir, "dest.txt"),
			srcPath:   srcPath,
			expectErr: false,
		},
		{
			name:      "copy file to nested directory",
			baseDir:   tmpDir,
			destPath:  filepath.Join(tmpDir, "nested", "dest.txt"),
			srcPath:   srcPath,
			expectErr: false,
		},
		{
			name:      "copy non-existent file",
			baseDir:   tmpDir,
			destPath:  filepath.Join(tmpDir, "dest2.txt"),
			srcPath:   filepath.Join(tmpDir, "nonexistent.txt"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provider.Copy(tt.baseDir, tt.destPath, tt.srcPath)
			if tt.expectErr {
				if err == nil {
					t.Error("Copy() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Copy() error = %v", err)
				return
			}

			// Verify file was copied
			content, readErr := os.ReadFile(tt.destPath)
			if readErr != nil {
				t.Errorf("Failed to read copied file: %v", readErr)
				return
			}
			if string(content) != string(srcContent) {
				t.Errorf("Copied file content = %q, want %q", string(content), string(srcContent))
			}
		})
	}
}

func TestLocalFileProvider_Read(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gateway_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create test file
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("test content")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatal(err)
	}

	provider := &localFileProvider{}

	tests := []struct {
		name      string
		baseDir   string
		filePath  string
		expectErr bool
	}{
		{
			name:      "read existing file",
			baseDir:   tmpDir,
			filePath:  testFile,
			expectErr: false,
		},
		{
			name:      "read non-existent file",
			baseDir:   tmpDir,
			filePath:  filepath.Join(tmpDir, "nonexistent.txt"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := provider.Read(tt.baseDir, tt.filePath)
			if tt.expectErr {
				if err == nil {
					t.Error("Read() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Read() error = %v", err)
			}
		})
	}
}
