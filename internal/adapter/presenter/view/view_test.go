package view

import (
	"testing"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/samber/do/v2"
)

// Mock ConfigureProjectView
type mockConfigureProjectView struct {
	updateCalled bool
}

func (m *mockConfigureProjectView) Update(vm *ConfigureProjectViewModel) {
	m.updateCalled = true
	if vm.Input != nil {
		vm.Input <- "test"
	}
}

// Tests for InitViewProvider

func TestNewInitViewProvider(t *testing.T) {
	mockView := &mockConfigureProjectView{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConfigureProjectView, error) {
		return mockView, nil
	})

	vp, err := NewInitViewProvider(injector)
	if err != nil {
		t.Errorf("NewInitViewProvider() error = %v", err)
	}
	if vp == nil {
		t.Error("NewInitViewProvider() returned nil")
	}
}

func TestInitViewProvider_Handle(t *testing.T) {
	mockView := &mockConfigureProjectView{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConfigureProjectView, error) {
		return mockView, nil
	})

	vp, _ := NewInitViewProvider(injector)

	ch := make(chan interface{}, 1)
	vm := &ConfigureProjectViewModel{
		Sequence: ConfigureProjectSequence_InputDriveFolder,
	}
	vm.Input = ch

	vp.Handle(vm)

	if !mockView.updateCalled {
		t.Error("Handle() did not call configureProjectView.Update()")
	}
}

func TestInitViewProvider_Handle_Panic(t *testing.T) {
	mockView := &mockConfigureProjectView{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConfigureProjectView, error) {
		return mockView, nil
	})

	vp, _ := NewInitViewProvider(injector)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Handle() should panic for unknown view model type")
		}
	}()

	// Pass an unrecognized type
	vp.Handle("unknown type")
}

// Tests for ConfigureProjectView

func TestNewConfigureProjectView(t *testing.T) {
	injector := do.New()
	view, err := NewConfigureProjectView(injector)
	if err != nil {
		t.Errorf("NewConfigureProjectView() error = %v", err)
	}
	if view == nil {
		t.Error("NewConfigureProjectView() returned nil")
	}
}

// Tests for ConfigureProjectSequence

func TestConfigureProjectSequence(t *testing.T) {
	tests := []struct {
		sequence ConfigureProjectSequence
		want     int
	}{
		{ConfigureProjectSequence_ConfirmOverwrite, 0},
		{ConfigureProjectSequence_InputDriveFolder, 1},
		{ConfigureProjectSequence_InputSiteUrl, 2},
		{ConfigureProjectSequence_InputProjectUrl, 3},
		{ConfigureProjectSequence_InputBranch, 4},
		{ConfigureProjectSequence_SelectTemplatePath, 5},
		{ConfigureProjectSequence_SelectAssetPath, 6},
		{ConfigureProjectSequence_ConfirmGitWarning, 7},
	}

	for _, tt := range tests {
		if int(tt.sequence) != tt.want {
			t.Errorf("ConfigureProjectSequence = %d, want %d", tt.sequence, tt.want)
		}
	}
}

// Tests for ConfigureProjectViewModel

func TestConfigureProjectViewModel(t *testing.T) {
	vm := &ConfigureProjectViewModel{
		Sequence:      ConfigureProjectSequence_InputDriveFolder,
		ConfigExists:  true,
		IsUnderGit:    false,
		DefaultValues: []string{"folder_id", "https://example.com", "https://github.com/user/repo", "main", "/templates", "/dist"},
	}

	if vm.Sequence != ConfigureProjectSequence_InputDriveFolder {
		t.Errorf("Sequence = %d, want %d", vm.Sequence, ConfigureProjectSequence_InputDriveFolder)
	}
	if !vm.ConfigExists {
		t.Error("ConfigExists should be true")
	}
	if vm.IsUnderGit {
		t.Error("IsUnderGit should be false")
	}
	if len(vm.DefaultValues) != 6 {
		t.Errorf("len(DefaultValues) = %d, want 6", len(vm.DefaultValues))
	}
}

// Tests for viewModel struct

func TestViewModel(t *testing.T) {
	ch := make(chan interface{}, 1)
	vm := viewModel{
		Input:  ch,
		Output: "test output",
	}

	if vm.Input == nil {
		t.Error("Input should not be nil")
	}
	if vm.Output != "test output" {
		t.Errorf("Output = %v, want %q", vm.Output, "test output")
	}
}

// Tests for message function

func TestMessage_WithOutput(t *testing.T) {
	// message function prints output, we just test it returns true
	result := message("test")
	if !result {
		t.Error("message() should return true when output is not nil")
	}
}

func TestMessage_WithNil(t *testing.T) {
	result := message(nil)
	if result {
		t.Error("message() should return false when output is nil")
	}
}

// Test that ViewModel interface is properly implemented
func TestConfigureProjectViewModel_ImplementsInterface(t *testing.T) {
	var _ core.ViewModel = &ConfigureProjectViewModel{}
}
