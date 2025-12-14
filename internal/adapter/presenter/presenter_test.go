package presenter

import (
	"testing"

	"github.com/c18t/nippo-cli/internal/adapter/presenter/view"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// Mock for ConsolePresenter
type mockConsolePresenter struct {
	progressMessage   string
	completeMessage   string
	warningErr        error
	suspendErr        error
	cancelledReturns  bool
	progressCalled    bool
	stopProgressCalled bool
	completeCalled    bool
	warningCalled     bool
	suspendCalled     bool
}

func (m *mockConsolePresenter) Progress(message string) {
	m.progressCalled = true
	m.progressMessage = message
}

func (m *mockConsolePresenter) StopProgress() {
	m.stopProgressCalled = true
}

func (m *mockConsolePresenter) Complete(message string) {
	m.completeCalled = true
	m.completeMessage = message
}

func (m *mockConsolePresenter) Warning(err error) {
	m.warningCalled = true
	m.warningErr = err
}

func (m *mockConsolePresenter) Suspend(err error) {
	m.suspendCalled = true
	m.suspendErr = err
}

func (m *mockConsolePresenter) IsCancelled() bool {
	return m.cancelledReturns
}

// Mock for InitViewProvider
type mockInitViewProvider struct {
	handleCalled bool
}

func (m *mockInitViewProvider) Handle(vm core.ViewModel) {
	m.handleCalled = true
}

// Tests for RootCommandPresenter

func TestNewRootCommandPresenter(t *testing.T) {
	injector := do.New()
	p, err := NewRootCommandPresenter(injector)
	if err != nil {
		t.Errorf("NewRootCommandPresenter() error = %v", err)
	}
	if p == nil {
		t.Error("NewRootCommandPresenter() returned nil")
	}
}

// Tests for CleanCommandPresenter

func TestNewCleanCommandPresenter(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, err := NewCleanCommandPresenter(injector)
	if err != nil {
		t.Errorf("NewCleanCommandPresenter() error = %v", err)
	}
	if p == nil {
		t.Error("NewCleanCommandPresenter() returned nil")
	}
}

func TestCleanCommandPresenter_Progress(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, _ := NewCleanCommandPresenter(injector)
	p.Progress(&port.CleanCommandUseCaseOutputData{Message: "test message"})

	if !mockBase.progressCalled {
		t.Error("Progress() did not call base.Progress()")
	}
	if mockBase.progressMessage != "test message" {
		t.Errorf("Progress message = %q, want %q", mockBase.progressMessage, "test message")
	}
}

func TestCleanCommandPresenter_StopProgress(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, _ := NewCleanCommandPresenter(injector)
	p.StopProgress()

	if !mockBase.stopProgressCalled {
		t.Error("StopProgress() did not call base.StopProgress()")
	}
}

func TestCleanCommandPresenter_Complete(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, _ := NewCleanCommandPresenter(injector)
	p.Complete(&port.CleanCommandUseCaseOutputData{Message: "done"})

	if !mockBase.completeCalled {
		t.Error("Complete() did not call base.Complete()")
	}
	if mockBase.completeMessage != "done" {
		t.Errorf("Complete message = %q, want %q", mockBase.completeMessage, "done")
	}
}

// Tests for DeployCommandPresenter

func TestNewDeployCommandPresenter(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, err := NewDeployCommandPresenter(injector)
	if err != nil {
		t.Errorf("NewDeployCommandPresenter() error = %v", err)
	}
	if p == nil {
		t.Error("NewDeployCommandPresenter() returned nil")
	}
}

func TestDeployCommandPresenter_Progress(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, _ := NewDeployCommandPresenter(injector)
	p.Progress(&port.DeployCommandUseCaseOutputData{Message: "deploying"})

	if !mockBase.progressCalled {
		t.Error("Progress() did not call base.Progress()")
	}
}

// Tests for UpdateCommandPresenter

func TestNewUpdateCommandPresenter(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, err := NewUpdateCommandPresenter(injector)
	if err != nil {
		t.Errorf("NewUpdateCommandPresenter() error = %v", err)
	}
	if p == nil {
		t.Error("NewUpdateCommandPresenter() returned nil")
	}
}

func TestUpdateCommandPresenter_Complete(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, _ := NewUpdateCommandPresenter(injector)
	p.Complete(&port.UpdateCommandUseCaseOutputData{Message: "updated"})

	if !mockBase.completeCalled {
		t.Error("Complete() did not call base.Complete()")
	}
}

// Tests for AuthPresenter

func TestNewAuthPresenter(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, err := NewAuthPresenter(injector)
	if err != nil {
		t.Errorf("NewAuthPresenter() error = %v", err)
	}
	if p == nil {
		t.Error("NewAuthPresenter() returned nil")
	}
}

func TestAuthPresenter_IsCancelled(t *testing.T) {
	mockBase := &mockConsolePresenter{cancelledReturns: true}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, _ := NewAuthPresenter(injector)
	result := p.IsCancelled()

	if !result {
		t.Error("IsCancelled() should return true")
	}
}

// Tests for DoctorPresenter

func TestNewDoctorPresenter(t *testing.T) {
	injector := do.New()
	p, err := NewDoctorPresenter(injector)
	if err != nil {
		t.Errorf("NewDoctorPresenter() error = %v", err)
	}
	if p == nil {
		t.Error("NewDoctorPresenter() returned nil")
	}
}

func TestDoctorPresenter_Show(t *testing.T) {
	injector := do.New()
	p, _ := NewDoctorPresenter(injector)

	output := &port.DoctorUseCaseOutputData{
		Checks: []port.DoctorCheck{
			{Category: "Test", Item: "item1", Status: port.DoctorCheckStatusPass, Message: "ok"},
			{Category: "Test", Item: "item2", Status: port.DoctorCheckStatusFail, Message: "failed", Suggestion: "fix it"},
			{Category: "Test", Item: "item3", Status: port.DoctorCheckStatusWarn, Message: "warning"},
		},
	}

	// Just verify it doesn't panic
	p.Show(output)
}

// Tests for BuildCommandPresenter

func TestNewBuildCommandPresenter(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, err := NewBuildCommandPresenter(injector)
	if err != nil {
		t.Errorf("NewBuildCommandPresenter() error = %v", err)
	}
	if p == nil {
		t.Error("NewBuildCommandPresenter() returned nil")
	}
}

func TestBuildCommandPresenter_Summary(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, _ := NewBuildCommandPresenter(injector)

	downloaded := []FileInfo{{Name: "file1.md", Id: "123"}}
	failed := []FileInfo{}

	// Just verify it doesn't panic
	p.Summary(downloaded, failed, nil)
}

func TestBuildCommandPresenter_SummaryWithError(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, _ := NewBuildCommandPresenter(injector)

	downloaded := []FileInfo{}
	failed := []FileInfo{{Name: "fail.md", Id: "456"}}

	// Just verify it doesn't panic with error
	p.Summary(downloaded, failed, nil)
}

// Tests for FormatCommandPresenter

func TestNewFormatCommandPresenter(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, err := NewFormatCommandPresenter(injector)
	if err != nil {
		t.Errorf("NewFormatCommandPresenter() error = %v", err)
	}
	if p == nil {
		t.Error("NewFormatCommandPresenter() returned nil")
	}
}

func TestFormatCommandPresenter_Summary(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})

	p, _ := NewFormatCommandPresenter(injector)

	updated := []FileInfo{{Name: "file1.md", Id: "123"}}
	failed := []FileInfo{}

	// Just verify it doesn't panic
	p.Summary(1, 0, 0, updated, failed)
}

func TestConvertStatus(t *testing.T) {
	tests := []struct {
		input port.FormatFileStatus
		want  int
	}{
		{port.FormatFileStatusSuccess, 0},
		{port.FormatFileStatusNoChange, 1},
		{port.FormatFileStatusFailed, 2},
	}

	for _, tt := range tests {
		result := convertStatus(tt.input)
		if int(result) != tt.want {
			t.Errorf("convertStatus(%d) = %d, want %d", tt.input, result, tt.want)
		}
	}
}

// Tests for InitSettingPresenter

func TestNewInitSettingPresenter(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	mockViewProvider := &mockInitViewProvider{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})
	do.Provide(injector, func(_ do.Injector) (view.InitViewProvider, error) {
		return mockViewProvider, nil
	})

	p, err := NewInitSettingPresenter(injector)
	if err != nil {
		t.Errorf("NewInitSettingPresenter() error = %v", err)
	}
	if p == nil {
		t.Error("NewInitSettingPresenter() returned nil")
	}
}

func TestInitSettingPresenter_Progress(t *testing.T) {
	mockBase := &mockConsolePresenter{}
	mockViewProvider := &mockInitViewProvider{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})
	do.Provide(injector, func(_ do.Injector) (view.InitViewProvider, error) {
		return mockViewProvider, nil
	})

	p, _ := NewInitSettingPresenter(injector)
	output := &port.InitSettingUseCaseOutputData{}
	output.Message = "test progress"
	p.Progress(output)

	if !mockBase.progressCalled {
		t.Error("Progress() did not call base.Progress()")
	}
}

func TestInitSettingPresenter_IsCancelled(t *testing.T) {
	mockBase := &mockConsolePresenter{cancelledReturns: true}
	mockViewProvider := &mockInitViewProvider{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (ConsolePresenter, error) {
		return mockBase, nil
	})
	do.Provide(injector, func(_ do.Injector) (view.InitViewProvider, error) {
		return mockViewProvider, nil
	})

	p, _ := NewInitSettingPresenter(injector)
	result := p.IsCancelled()

	if !result {
		t.Error("IsCancelled() should return true")
	}
}

// Test FileInfo struct

func TestFileInfo(t *testing.T) {
	info := FileInfo{
		Name: "test.md",
		Id:   "file123",
	}

	if info.Name != "test.md" {
		t.Errorf("Name = %q, want %q", info.Name, "test.md")
	}
	if info.Id != "file123" {
		t.Errorf("Id = %q, want %q", info.Id, "file123")
	}
}
