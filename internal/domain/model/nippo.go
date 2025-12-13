package model

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"go.yaml.in/yaml/v3"
	"google.golang.org/api/drive/v3"
)

// FrontMatter represents the YAML front-matter metadata in a nippo file
type FrontMatter struct {
	Created time.Time `yaml:"created,omitempty"`
	Updated time.Time `yaml:"updated,omitempty"`
	Raw     map[string]interface{}
}

var (
	frontMatterDelimiter = []byte("---")
	ErrInvalidDateFormat = errors.New("invalid date format: expected RFC 3339")
	ErrMalformedYAML     = errors.New("malformed YAML in front-matter")
)

type Nippo struct {
	Date        NippoDate
	FilePath    string
	Content     []byte
	FrontMatter *FrontMatter
	RemoteFile  *drive.File
}

type NippoDate interface {
	PathString() string
	FileString() string
	TitleString() string

	Year() int
	Month() time.Month
	Day() int
	Weekday() time.Weekday
}
type nippoDate struct {
	time time.Time
}

func NewNippo(filePath string) (*Nippo, error) {
	nippo := &Nippo{}

	if err := checkNippoIsExist(filePath); err != nil {
		return nippo, err
	}

	nippo.Date = NewNippoDate(filePath)
	nippo.FilePath = filePath
	return nippo, nil
}

func NewNippoDate(filePath string) NippoDate {
	date, err := time.Parse("2006-01-02", filepath.Base(filePath)[:10])
	if err != nil {
		panic(err)
	}
	return &nippoDate{date}
}

func (n *Nippo) GetMarkdown() ([]byte, error) {
	if len(n.Content) > 0 {
		// If content is already loaded, strip front-matter and return body
		_, body, err := ParseFrontMatter(n.Content)
		if err != nil {
			// On parse error, log warning and return content as-is
			fmt.Fprintf(os.Stderr, "Warning: malformed front-matter in %s: %v\n", n.FilePath, err)
			return n.Content, nil
		}
		return body, nil
	}
	f, err := os.Open(n.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rawContent, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// Parse front-matter and store it
	fm, body, parseErr := ParseFrontMatter(rawContent)
	if parseErr != nil {
		// Log warning for malformed front-matter
		fmt.Fprintf(os.Stderr, "Warning: malformed front-matter in %s: %v\n", n.FilePath, parseErr)
	} else if fm != nil {
		n.FrontMatter = fm
	}

	n.Content = rawContent
	return body, nil
}

func (n *Nippo) GetHtml() ([]byte, error) {
	data, err := n.GetMarkdown()
	if err != nil {
		return nil, err
	}

	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(data)

	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer), nil
}

func checkNippoIsExist(filePath string) error {
	if f, err := os.Stat(filePath); os.IsNotExist(err) || f.IsDir() {
		return fmt.Errorf("nippo not found: filePath=%v", filePath)
	}
	return nil
}

func (date *nippoDate) String() string {
	return date.PathString()
}

func (date *nippoDate) FileString() string {
	return fmt.Sprintf("%04d-%02d-%02d", date.time.Year(), date.time.Month(), date.time.Day())
}

func (date *nippoDate) PathString() string {
	return fmt.Sprintf("%04d%02d%02d", date.time.Year(), date.time.Month(), date.time.Day())
}

func (date *nippoDate) TitleString() string {
	return fmt.Sprintf("%02d/%02d %s", date.time.Month(), date.time.Day(),
		strings.ToLower(date.time.Weekday().String()[:3]))
}

func (date *nippoDate) Year() int {
	return date.time.Year()
}

func (date *nippoDate) Month() time.Month {
	return date.time.Month()
}

func (date *nippoDate) Day() int {
	return date.time.Day()
}

func (date *nippoDate) Weekday() time.Weekday {
	return date.time.Weekday()
}

// HasFrontMatter checks if content starts with front-matter delimiters
func HasFrontMatter(content []byte) bool {
	return bytes.HasPrefix(content, frontMatterDelimiter)
}

// HasNowPlaceholder checks if content contains "updated: now" in front-matter
func HasNowPlaceholder(content []byte) bool {
	if !HasFrontMatter(content) {
		return false
	}

	// Find the end of front-matter
	rest := content[3:]
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	}

	endIndex := bytes.Index(rest, []byte("\n---"))
	if endIndex == -1 {
		return false
	}

	yamlContent := rest[:endIndex]
	// Check for "updated: now" pattern (with variations)
	return bytes.Contains(yamlContent, []byte("updated: now")) ||
		bytes.Contains(yamlContent, []byte("updated:now")) ||
		bytes.Contains(yamlContent, []byte("updated: 'now'")) ||
		bytes.Contains(yamlContent, []byte("updated: \"now\""))
}

// ParseFrontMatter extracts YAML front-matter from content and returns
// the parsed FrontMatter and the remaining body content.
// Returns nil FrontMatter if no front-matter is present.
func ParseFrontMatter(content []byte) (*FrontMatter, []byte, error) {
	if !HasFrontMatter(content) {
		return nil, content, nil
	}

	// Skip the first "---"
	rest := content[3:]
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	}

	// Find the closing "---"
	// Handle empty front-matter case (---\n---) where rest starts with "---"
	var endIndex int
	var yamlContent []byte
	if bytes.HasPrefix(rest, frontMatterDelimiter) {
		// Empty front-matter
		endIndex = 0
		yamlContent = []byte{}
	} else {
		endIndex = bytes.Index(rest, []byte("\n---"))
		if endIndex == -1 {
			// No closing delimiter, treat as no front-matter
			return nil, content, nil
		}
		yamlContent = rest[:endIndex]
	}

	// Parse into map to preserve unknown fields
	var raw map[string]interface{}
	if len(yamlContent) > 0 {
		if err := yaml.Unmarshal(yamlContent, &raw); err != nil {
			return nil, content, ErrMalformedYAML
		}
	}

	// Handle empty front-matter (---\n---)
	if raw == nil {
		raw = make(map[string]interface{})
	}

	fm := &FrontMatter{Raw: raw}

	// Extract and validate created field
	if createdVal, ok := raw["created"]; ok {
		created, err := parseDateTime(createdVal)
		if err != nil {
			return nil, content, fmt.Errorf("created: %w", ErrInvalidDateFormat)
		}
		fm.Created = created
	}

	// Extract and validate updated field (skip "now" placeholder)
	if updatedVal, ok := raw["updated"]; ok {
		if str, ok := updatedVal.(string); ok && str == "now" {
			// Keep as placeholder, don't parse
		} else {
			updated, err := parseDateTime(updatedVal)
			if err != nil {
				return nil, content, fmt.Errorf("updated: %w", ErrInvalidDateFormat)
			}
			fm.Updated = updated
		}
	}

	// Extract body (after the closing "---" and newline)
	var body []byte
	if endIndex == 0 {
		// Empty front-matter: skip "---" and newline
		body = rest[3:]
	} else {
		body = rest[endIndex+4:] // Skip "\n---"
	}
	// Strip all leading newlines from body (we'll add exactly one blank line when reconstructing)
	for len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	}

	return fm, body, nil
}

// parseDateTime parses a datetime value from front-matter
func parseDateTime(val interface{}) (time.Time, error) {
	switch v := val.(type) {
	case time.Time:
		return v, nil
	case string:
		// Try RFC 3339 format
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Time{}, err
		}
		return t, nil
	default:
		return time.Time{}, fmt.Errorf("unsupported type: %T", val)
	}
}

// formatTimeValue converts a value to RFC3339 string format.
// Handles time.Time, string, and other types.
func formatTimeValue(val interface{}) string {
	switch v := val.(type) {
	case time.Time:
		return v.Format(time.RFC3339)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GenerateFrontMatter creates a YAML front-matter string with the given created time.
// The updated field is omitted by default.
// A blank line is added after the closing delimiter for readability.
func GenerateFrontMatter(created time.Time) string {
	return fmt.Sprintf("---\ncreated: %s\n---\n\n", created.Format(time.RFC3339))
}

// UpdateFrontMatter updates the front-matter in content while preserving unknown fields.
// If created is non-zero, it sets/updates the created field.
// If updated is non-zero, it sets/updates the updated field.
// If replaceNow is true and updated is non-zero, it replaces "now" placeholder.
func UpdateFrontMatter(content []byte, created, updated time.Time, replaceNow bool) ([]byte, error) {
	fm, body, err := ParseFrontMatter(content)
	if err != nil {
		return nil, err
	}

	// If no front-matter exists, create new one
	if fm == nil {
		if created.IsZero() {
			return content, nil
		}
		newFM := GenerateFrontMatter(created)
		return append([]byte(newFM), content...), nil
	}

	// Update the raw map
	if !created.IsZero() {
		fm.Raw["created"] = created.Format(time.RFC3339)
	}

	if !updated.IsZero() {
		if replaceNow {
			// Only replace if current value is "now"
			if val, ok := fm.Raw["updated"]; ok {
				if str, ok := val.(string); ok && str == "now" {
					fm.Raw["updated"] = updated.Format(time.RFC3339)
				}
			}
		} else {
			fm.Raw["updated"] = updated.Format(time.RFC3339)
		}
	}

	// Serialize back to YAML manually to avoid quoting timestamps
	var buf bytes.Buffer
	buf.WriteString("---\n")

	// Write created field first (if present)
	if val, ok := fm.Raw["created"]; ok {
		buf.WriteString(fmt.Sprintf("created: %s\n", formatTimeValue(val)))
	}

	// Write updated field second (if present)
	if val, ok := fm.Raw["updated"]; ok {
		buf.WriteString(fmt.Sprintf("updated: %s\n", formatTimeValue(val)))
	}

	// Write other fields using yaml.Marshal
	for key, val := range fm.Raw {
		if key == "created" || key == "updated" {
			continue
		}
		// Marshal single field
		fieldBytes, err := yaml.Marshal(map[string]interface{}{key: val})
		if err != nil {
			return nil, err
		}
		buf.Write(fieldBytes)
	}

	buf.WriteString("---\n\n")
	buf.Write(body)

	return buf.Bytes(), nil
}

// GetCreatedTime returns the created time from front-matter if available,
// otherwise returns the time derived from the nippo date (filename).
func (n *Nippo) GetCreatedTime() time.Time {
	if n.FrontMatter != nil && !n.FrontMatter.Created.IsZero() {
		return n.FrontMatter.Created
	}
	// Fallback to filename-derived date
	return time.Date(n.Date.Year(), n.Date.Month(), n.Date.Day(), 0, 0, 0, 0, time.Local)
}

// GetUpdatedTime returns the updated time from front-matter if available,
// otherwise returns zero time.
func (n *Nippo) GetUpdatedTime() time.Time {
	if n.FrontMatter != nil {
		return n.FrontMatter.Updated
	}
	return time.Time{}
}
