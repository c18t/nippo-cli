package gateway

import (
	"testing"
	"time"

	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/samber/do/v2"
)

func TestNewDriveFileTimestamp(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC)
	ts := NewDriveFileTimestamp(testTime)

	expected := "2024-01-15T09:30:00Z"
	if ts.String() != expected {
		t.Errorf("String() = %q, want %q", ts.String(), expected)
	}
}

func TestDriveFileTimestamp_String(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "UTC time",
			time:     time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC),
			expected: "2024-01-15T09:30:00Z",
		},
		{
			name:     "with timezone offset",
			time:     time.Date(2024, 1, 15, 18, 30, 0, 0, time.FixedZone("JST", 9*60*60)),
			expected: "2024-01-15T18:30:00+09:00",
		},
		{
			name:     "midnight",
			time:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "2024-01-01T00:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewDriveFileTimestamp(tt.time)
			if ts.String() != tt.expected {
				t.Errorf("String() = %q, want %q", ts.String(), tt.expected)
			}
		})
	}
}

func TestNewDriveFileProvider(t *testing.T) {
	injector := do.New()
	provider, err := NewDriveFileProvider(injector)
	if err != nil {
		t.Errorf("NewDriveFileProvider() error = %v", err)
	}
	if provider == nil {
		t.Error("NewDriveFileProvider() returned nil")
	}
}

func TestDriveFileProvider_Shutdown(t *testing.T) {
	provider := &driveFileProvider{}
	err := provider.Shutdown()
	if err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}

func TestDriveFileProvider_queryBuilder(t *testing.T) {
	provider := &driveFileProvider{}

	tests := []struct {
		name     string
		param    *repository.QueryListParam
		contains []string
	}{
		{
			name: "with folders",
			param: &repository.QueryListParam{
				Folders: []string{"folder1"},
			},
			contains: []string{"parents in", "'folder1'"},
		},
		{
			name: "with file extensions",
			param: &repository.QueryListParam{
				FileExtensions: []string{"md", "txt"},
			},
			contains: []string{"fileExtension in"},
		},
		{
			name: "with updated time",
			param: &repository.QueryListParam{
				UpdatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			contains: []string{"modifiedTime >="},
		},
		{
			name: "with multiple folders",
			param: &repository.QueryListParam{
				Folders: []string{"folder1", "folder2"},
			},
			contains: []string{"'folder1'", "'folder2'"},
		},
		{
			name: "empty param",
			param: &repository.QueryListParam{},
			contains: []string{"trashed != true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.queryBuilder(tt.param)
			for _, substr := range tt.contains {
				if !containsString(result, substr) {
					t.Errorf("queryBuilder() = %q, want to contain %q", result, substr)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestDriveFolderMimeType(t *testing.T) {
	expected := "application/vnd.google-apps.folder"
	if DriveFolderMimeType != expected {
		t.Errorf("DriveFolderMimeType = %q, want %q", DriveFolderMimeType, expected)
	}
}
