package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/samber/do/v2"
)

func TestNewTemplateService(t *testing.T) {
	injector := do.New()
	service, err := NewTemplateService(injector)
	if err != nil {
		t.Errorf("NewTemplateService() error = %v", err)
	}
	if service == nil {
		t.Error("NewTemplateService() returned nil")
	}
}

func TestTemplateService_SaveTo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create template directory with templates
	templateDir := filepath.Join(tmpDir, "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create layout template
	layoutContent := `{{define "layout"}}<!DOCTYPE html><html>{{template "content" .}}</html>{{end}}`
	if err := os.WriteFile(filepath.Join(templateDir, "layout.html"), []byte(layoutContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create test template
	testContent := `{{define "test"}}Hello, {{.Name}}!{{end}}`
	if err := os.WriteFile(filepath.Join(templateDir, "test.html"), []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Set up global config
	core.Cfg = &core.Config{}
	core.Cfg.Paths.DataDir = tmpDir

	injector := do.New()
	service, _ := NewTemplateService(injector)

	outputPath := filepath.Join(tmpDir, "output", "test.html")
	data := struct {
		Name string
	}{
		Name: "World",
	}

	err = service.SaveTo(outputPath, "test", data)
	if err != nil {
		t.Errorf("SaveTo() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("SaveTo() did not create output file")
	}
}

func TestTemplateService_SaveTo_CreateDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create template directory with templates
	templateDir := filepath.Join(tmpDir, "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatal(err)
	}

	layoutContent := `{{define "layout"}}{{template "content" .}}{{end}}`
	if err := os.WriteFile(filepath.Join(templateDir, "layout.html"), []byte(layoutContent), 0644); err != nil {
		t.Fatal(err)
	}

	testContent := `{{define "test"}}Test{{end}}`
	if err := os.WriteFile(filepath.Join(templateDir, "test.html"), []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Set up global config
	core.Cfg = &core.Config{}
	core.Cfg.Paths.DataDir = tmpDir

	injector := do.New()
	service, _ := NewTemplateService(injector)

	// Create output in nested directory that doesn't exist
	outputPath := filepath.Join(tmpDir, "deep", "nested", "output.html")
	err = service.SaveTo(outputPath, "test", nil)
	if err != nil {
		t.Errorf("SaveTo() error = %v", err)
	}

	// Verify directory was created
	dir := filepath.Dir(outputPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("SaveTo() did not create output directory")
	}
}
