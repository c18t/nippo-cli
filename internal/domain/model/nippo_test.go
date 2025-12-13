package model

import (
	"strings"
	"testing"
	"time"
)

func TestHasFrontMatter(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "valid front-matter",
			content:  "---\ncreated: 2024-01-15T09:30:00+09:00\n---\n# Content",
			expected: true,
		},
		{
			name:     "no front-matter",
			content:  "# Content without front-matter",
			expected: false,
		},
		{
			name:     "empty content",
			content:  "",
			expected: false,
		},
		{
			name:     "only delimiter",
			content:  "---",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasFrontMatter([]byte(tt.content))
			if result != tt.expected {
				t.Errorf("HasFrontMatter(%q) = %v, want %v", tt.content, result, tt.expected)
			}
		})
	}
}

func TestHasNowPlaceholder(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "has updated: now",
			content:  "---\ncreated: 2024-01-15T09:30:00+09:00\nupdated: now\n---\n# Content",
			expected: true,
		},
		{
			name:     "has updated:now (no space)",
			content:  "---\ncreated: 2024-01-15T09:30:00+09:00\nupdated:now\n---\n# Content",
			expected: true,
		},
		{
			name:     "has updated: 'now' (single quotes)",
			content:  "---\ncreated: 2024-01-15T09:30:00+09:00\nupdated: 'now'\n---\n# Content",
			expected: true,
		},
		{
			name:     "has updated: \"now\" (double quotes)",
			content:  "---\ncreated: 2024-01-15T09:30:00+09:00\nupdated: \"now\"\n---\n# Content",
			expected: true,
		},
		{
			name:     "no now placeholder",
			content:  "---\ncreated: 2024-01-15T09:30:00+09:00\n---\n# Content",
			expected: false,
		},
		{
			name:     "no front-matter",
			content:  "# Content",
			expected: false,
		},
		{
			name:     "now in body, not front-matter",
			content:  "---\ncreated: 2024-01-15T09:30:00+09:00\n---\nupdated: now",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasNowPlaceholder([]byte(tt.content))
			if result != tt.expected {
				t.Errorf("HasNowPlaceholder(%q) = %v, want %v", tt.content, result, tt.expected)
			}
		})
	}
}

func TestParseFrontMatter(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectFM      bool
		expectCreated time.Time
		expectUpdated time.Time
		expectBody    string
		expectErr     bool
	}{
		{
			name:          "valid front-matter with created",
			content:       "---\ncreated: 2024-01-15T09:30:00+09:00\n---\n# Content",
			expectFM:      true,
			expectCreated: time.Date(2024, 1, 15, 9, 30, 0, 0, time.FixedZone("", 9*60*60)),
			expectBody:    "# Content",
		},
		{
			name:          "valid front-matter with created and updated",
			content:       "---\ncreated: 2024-01-15T09:30:00+09:00\nupdated: 2024-01-16T10:00:00+09:00\n---\n# Content",
			expectFM:      true,
			expectCreated: time.Date(2024, 1, 15, 9, 30, 0, 0, time.FixedZone("", 9*60*60)),
			expectUpdated: time.Date(2024, 1, 16, 10, 0, 0, 0, time.FixedZone("", 9*60*60)),
			expectBody:    "# Content",
		},
		{
			name:       "no front-matter",
			content:    "# Content without front-matter",
			expectFM:   false,
			expectBody: "# Content without front-matter",
		},
		{
			name:       "empty front-matter",
			content:    "---\n---\n# Content",
			expectFM:   true,
			expectBody: "# Content",
		},
		{
			name:       "unclosed front-matter",
			content:    "---\ncreated: 2024-01-15T09:30:00+09:00\n# No closing delimiter",
			expectFM:   false,
			expectBody: "---\ncreated: 2024-01-15T09:30:00+09:00\n# No closing delimiter",
		},
		{
			name:      "malformed YAML",
			content:   "---\n: invalid yaml\n---\n# Content",
			expectFM:  false,
			expectErr: true,
		},
		{
			name:          "date only format (YAML parses as time.Time)",
			content:       "---\ncreated: 2024-01-15\n---\n# Content",
			expectFM:      true,
			expectCreated: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expectBody:    "# Content",
		},
		{
			name:          "now placeholder preserved",
			content:       "---\ncreated: 2024-01-15T09:30:00+09:00\nupdated: now\n---\n# Content",
			expectFM:      true,
			expectCreated: time.Date(2024, 1, 15, 9, 30, 0, 0, time.FixedZone("", 9*60*60)),
			expectBody:    "# Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body, err := ParseFrontMatter([]byte(tt.content))

			if tt.expectErr {
				if err == nil {
					t.Errorf("ParseFrontMatter(%q) expected error, got nil", tt.content)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseFrontMatter(%q) unexpected error: %v", tt.content, err)
				return
			}

			if tt.expectFM && fm == nil {
				t.Errorf("ParseFrontMatter(%q) expected front-matter, got nil", tt.content)
				return
			}

			if !tt.expectFM && fm != nil {
				t.Errorf("ParseFrontMatter(%q) expected no front-matter, got %v", tt.content, fm)
				return
			}

			if fm != nil {
				if !tt.expectCreated.IsZero() && !fm.Created.Equal(tt.expectCreated) {
					t.Errorf("ParseFrontMatter(%q) created = %v, want %v", tt.content, fm.Created, tt.expectCreated)
				}
				if !tt.expectUpdated.IsZero() && !fm.Updated.Equal(tt.expectUpdated) {
					t.Errorf("ParseFrontMatter(%q) updated = %v, want %v", tt.content, fm.Updated, tt.expectUpdated)
				}
			}

			if string(body) != tt.expectBody {
				t.Errorf("ParseFrontMatter(%q) body = %q, want %q", tt.content, string(body), tt.expectBody)
			}
		})
	}
}

func TestGenerateFrontMatter(t *testing.T) {
	created := time.Date(2024, 1, 15, 9, 30, 0, 0, time.FixedZone("JST", 9*60*60))
	result := GenerateFrontMatter(created)

	expected := "---\ncreated: 2024-01-15T09:30:00+09:00\n---\n\n"
	if result != expected {
		t.Errorf("GenerateFrontMatter() = %q, want %q", result, expected)
	}
}

func TestUpdateFrontMatter(t *testing.T) {
	created := time.Date(2024, 1, 15, 9, 30, 0, 0, time.FixedZone("JST", 9*60*60))
	updated := time.Date(2024, 1, 16, 10, 0, 0, 0, time.FixedZone("JST", 9*60*60))

	tests := []struct {
		name       string
		content    string
		created    time.Time
		updated    time.Time
		replaceNow bool
		expectErr  bool
		validate   func(result []byte) bool
	}{
		{
			name:    "add front-matter to content without it",
			content: "# Content",
			created: created,
			validate: func(result []byte) bool {
				return HasFrontMatter(result)
			},
		},
		{
			name:    "no changes when no front-matter and no created",
			content: "# Content",
			validate: func(result []byte) bool {
				return string(result) == "# Content"
			},
		},
		{
			name:       "replace now placeholder",
			content:    "---\ncreated: 2024-01-15T09:30:00+09:00\nupdated: now\n---\n# Content",
			updated:    updated,
			replaceNow: true,
			validate: func(result []byte) bool {
				if HasNowPlaceholder(result) {
					return false
				}
				// Verify timestamp is not quoted
				resultStr := string(result)
				// Should contain unquoted timestamp like "updated: 2024-01-16T..."
				// Should NOT contain quoted timestamp like 'updated: "2024-01-16T...'
				return strings.Contains(resultStr, "updated: 2024-01-16T") &&
					!strings.Contains(resultStr, `updated: "`)
			},
		},
		{
			name:    "preserve unknown fields",
			content: "---\ncreated: 2024-01-15T09:30:00+09:00\ntags:\n  - diary\n  - work\n---\n# Content",
			created: created,
			validate: func(result []byte) bool {
				fm, _, _ := ParseFrontMatter(result)
				if fm == nil {
					return false
				}
				tags, ok := fm.Raw["tags"]
				return ok && tags != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UpdateFrontMatter([]byte(tt.content), tt.created, tt.updated, tt.replaceNow)

			if tt.expectErr {
				if err == nil {
					t.Errorf("UpdateFrontMatter() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateFrontMatter() unexpected error: %v", err)
				return
			}

			if !tt.validate(result) {
				t.Errorf("UpdateFrontMatter() validation failed, result = %q", string(result))
			}
		})
	}
}

func TestUnknownFieldsPreservation(t *testing.T) {
	content := `---
created: 2024-01-15T09:30:00+09:00
author: test
tags:
  - diary
  - work
custom_field: custom_value
---
# Content`

	fm, body, err := ParseFrontMatter([]byte(content))
	if err != nil {
		t.Fatalf("ParseFrontMatter() error = %v", err)
	}

	if fm == nil {
		t.Fatal("ParseFrontMatter() returned nil front-matter")
	}

	// Check unknown fields are preserved
	if _, ok := fm.Raw["author"]; !ok {
		t.Error("author field not preserved")
	}
	if _, ok := fm.Raw["tags"]; !ok {
		t.Error("tags field not preserved")
	}
	if _, ok := fm.Raw["custom_field"]; !ok {
		t.Error("custom_field not preserved")
	}

	if string(body) != "# Content" {
		t.Errorf("body = %q, want %q", string(body), "# Content")
	}
}

func TestCombinedUpdateScenario(t *testing.T) {
	// Test case: front-matter has tags but missing created, and has updated: now
	content := `---
tags:
  - diary
updated: now
---
# Content`

	created := time.Date(2024, 1, 15, 9, 30, 0, 0, time.FixedZone("JST", 9*60*60))
	updated := time.Date(2024, 1, 16, 10, 0, 0, 0, time.FixedZone("JST", 9*60*60))

	result, err := UpdateFrontMatter([]byte(content), created, updated, true)
	if err != nil {
		t.Fatalf("UpdateFrontMatter() error = %v", err)
	}

	// Parse result to verify
	fm, _, err := ParseFrontMatter(result)
	if err != nil {
		t.Fatalf("ParseFrontMatter() error = %v", err)
	}

	// Verify created was added
	if fm.Created.IsZero() {
		t.Error("created field was not added")
	}

	// Verify updated was replaced (not "now" anymore)
	if HasNowPlaceholder(result) {
		t.Error("now placeholder was not replaced")
	}

	// Verify tags were preserved
	if _, ok := fm.Raw["tags"]; !ok {
		t.Error("tags field was not preserved")
	}
}

func TestEmptyFrontMatter(t *testing.T) {
	content := "---\n---\n# Content"

	fm, body, err := ParseFrontMatter([]byte(content))
	if err != nil {
		t.Fatalf("ParseFrontMatter() error = %v", err)
	}

	if fm == nil {
		t.Fatal("ParseFrontMatter() returned nil for empty front-matter")
	}

	if fm.Raw == nil {
		t.Error("Raw map should not be nil for empty front-matter")
	}

	if string(body) != "# Content" {
		t.Errorf("body = %q, want %q", string(body), "# Content")
	}
}

func TestRFC3339DateValidation(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		expectErr bool
	}{
		{
			name:      "valid RFC 3339",
			content:   "---\ncreated: 2024-01-15T09:30:00+09:00\n---\n",
			expectErr: false,
		},
		{
			name:      "valid RFC 3339 with Z timezone",
			content:   "---\ncreated: 2024-01-15T00:30:00Z\n---\n",
			expectErr: false,
		},
		{
			name:      "date only (YAML parses as time.Time)",
			content:   "---\ncreated: 2024-01-15\n---\n",
			expectErr: false, // YAML library parses this as time.Time
		},
		{
			name:      "datetime without timezone (treated as string, fails RFC 3339)",
			content:   "---\ncreated: 2024-01-15T09:30:00\n---\n",
			expectErr: true, // YAML treats this as string, RFC 3339 parse fails
		},
		{
			name:      "invalid format",
			content:   "---\ncreated: January 15, 2024\n---\n",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ParseFrontMatter([]byte(tt.content))
			if tt.expectErr && err == nil {
				t.Errorf("ParseFrontMatter(%q) expected error, got nil", tt.content)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("ParseFrontMatter(%q) unexpected error: %v", tt.content, err)
			}
		})
	}
}
