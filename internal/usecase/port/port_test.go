package port

import (
	"testing"

	"github.com/samber/do/v2"
)

// Mock implementations

type mockAuthUseCase struct {
	handleCalled bool
}

func (m *mockAuthUseCase) Handle(input *AuthUseCaseInputData) {
	m.handleCalled = true
}

type mockCleanCommandUseCase struct {
	handleCalled bool
}

func (m *mockCleanCommandUseCase) Handle(input *CleanCommandUseCaseInputData) {
	m.handleCalled = true
}

type mockDeployCommandUseCase struct {
	handleCalled bool
}

func (m *mockDeployCommandUseCase) Handle(input *DeployCommandUseCaseInputData) {
	m.handleCalled = true
}

type mockFormatCommandUseCase struct {
	handleCalled bool
}

func (m *mockFormatCommandUseCase) Handle(input *FormatCommandUseCaseInputData) {
	m.handleCalled = true
}

type mockInitSettingUseCase struct {
	handleCalled bool
}

func (m *mockInitSettingUseCase) Handle(input *InitSettingUseCaseInputData) {
	m.handleCalled = true
}

type mockRootCommandUseCase struct {
	handleCalled bool
}

func (m *mockRootCommandUseCase) Handle(input *RootCommandUseCaseInputData) {
	m.handleCalled = true
}

type mockUpdateCommandUseCase struct {
	handleCalled bool
}

func (m *mockUpdateCommandUseCase) Handle(input *UpdateCommandUseCaseInputData) {
	m.handleCalled = true
}

// Auth tests

func TestNewAuthUseCaseBus(t *testing.T) {
	mock := &mockAuthUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (AuthUseCase, error) {
		return mock, nil
	})

	bus, err := NewAuthUseCaseBus(injector)
	if err != nil {
		t.Errorf("NewAuthUseCaseBus() error = %v", err)
	}
	if bus == nil {
		t.Error("NewAuthUseCaseBus() returned nil")
	}
}

func TestAuthUseCaseBus_Handle(t *testing.T) {
	mock := &mockAuthUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (AuthUseCase, error) {
		return mock, nil
	})

	bus, _ := NewAuthUseCaseBus(injector)
	bus.Handle(&AuthUseCaseInputData{})

	if !mock.handleCalled {
		t.Error("Handle() did not call auth.Handle()")
	}
}

func TestAuthUseCaseOutputData(t *testing.T) {
	output := &AuthUseCaseOutputData{
		Message: "test message",
	}
	if output.Message != "test message" {
		t.Errorf("Message = %q, want %q", output.Message, "test message")
	}
}

// Clean tests

func TestNewCleanUseCaseBus(t *testing.T) {
	mock := &mockCleanCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (CleanCommandUseCase, error) {
		return mock, nil
	})

	bus, err := NewCleanUseCaseBus(injector)
	if err != nil {
		t.Errorf("NewCleanUseCaseBus() error = %v", err)
	}
	if bus == nil {
		t.Error("NewCleanUseCaseBus() returned nil")
	}
}

func TestCleanUseCaseBus_Handle(t *testing.T) {
	mock := &mockCleanCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (CleanCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewCleanUseCaseBus(injector)
	bus.Handle(&CleanCommandUseCaseInputData{})

	if !mock.handleCalled {
		t.Error("Handle() did not call command.Handle()")
	}
}

func TestCleanUseCaseBus_Handle_Panic(t *testing.T) {
	mock := &mockCleanCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (CleanCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewCleanUseCaseBus(injector)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Handle() should panic for unknown input type")
		}
	}()

	bus.Handle("unknown type")
}

func TestCleanCommandUseCaseOutputData(t *testing.T) {
	output := &CleanCommandUseCaseOutputData{
		Message: "test message",
	}
	if output.Message != "test message" {
		t.Errorf("Message = %q, want %q", output.Message, "test message")
	}
}

// Deploy tests

func TestNewDeployUseCaseBus(t *testing.T) {
	mock := &mockDeployCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (DeployCommandUseCase, error) {
		return mock, nil
	})

	bus, err := NewDeployUseCaseBus(injector)
	if err != nil {
		t.Errorf("NewDeployUseCaseBus() error = %v", err)
	}
	if bus == nil {
		t.Error("NewDeployUseCaseBus() returned nil")
	}
}

func TestDeployUseCaseBus_Handle(t *testing.T) {
	mock := &mockDeployCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (DeployCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewDeployUseCaseBus(injector)
	bus.Handle(&DeployCommandUseCaseInputData{})

	if !mock.handleCalled {
		t.Error("Handle() did not call command.Handle()")
	}
}

func TestDeployUseCaseBus_Handle_Panic(t *testing.T) {
	mock := &mockDeployCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (DeployCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewDeployUseCaseBus(injector)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Handle() should panic for unknown input type")
		}
	}()

	bus.Handle("unknown type")
}

func TestDeployCommandUseCaseOutputData(t *testing.T) {
	output := &DeployCommandUseCaseOutputData{
		Message: "test message",
	}
	if output.Message != "test message" {
		t.Errorf("Message = %q, want %q", output.Message, "test message")
	}
}

// Doctor tests

func TestDoctorCheckStatus(t *testing.T) {
	tests := []struct {
		status DoctorCheckStatus
		want   int
	}{
		{DoctorCheckStatusPass, 0},
		{DoctorCheckStatusFail, 1},
		{DoctorCheckStatusWarn, 2},
	}

	for _, tt := range tests {
		if int(tt.status) != tt.want {
			t.Errorf("DoctorCheckStatus = %d, want %d", tt.status, tt.want)
		}
	}
}

func TestDoctorCheck(t *testing.T) {
	check := DoctorCheck{
		Category:   "Config",
		Item:       "credentials",
		Status:     DoctorCheckStatusPass,
		Message:    "OK",
		Suggestion: "",
	}

	if check.Category != "Config" {
		t.Errorf("Category = %q, want %q", check.Category, "Config")
	}
	if check.Item != "credentials" {
		t.Errorf("Item = %q, want %q", check.Item, "credentials")
	}
	if check.Status != DoctorCheckStatusPass {
		t.Errorf("Status = %d, want %d", check.Status, DoctorCheckStatusPass)
	}
}

func TestDoctorUseCaseOutputData(t *testing.T) {
	output := &DoctorUseCaseOutputData{
		Checks: []DoctorCheck{
			{Category: "Test", Item: "item1", Status: DoctorCheckStatusPass},
			{Category: "Test", Item: "item2", Status: DoctorCheckStatusFail},
		},
	}

	if len(output.Checks) != 2 {
		t.Errorf("len(Checks) = %d, want 2", len(output.Checks))
	}
}

// Format tests

func TestNewFormatUseCaseBus(t *testing.T) {
	mock := &mockFormatCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (FormatCommandUseCase, error) {
		return mock, nil
	})

	bus, err := NewFormatUseCaseBus(injector)
	if err != nil {
		t.Errorf("NewFormatUseCaseBus() error = %v", err)
	}
	if bus == nil {
		t.Error("NewFormatUseCaseBus() returned nil")
	}
}

func TestFormatUseCaseBus_Handle(t *testing.T) {
	mock := &mockFormatCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (FormatCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewFormatUseCaseBus(injector)
	bus.Handle(&FormatCommandUseCaseInputData{})

	if !mock.handleCalled {
		t.Error("Handle() did not call command.Handle()")
	}
}

func TestFormatUseCaseBus_Handle_Panic(t *testing.T) {
	mock := &mockFormatCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (FormatCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewFormatUseCaseBus(injector)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Handle() should panic for unknown input type")
		}
	}()

	bus.Handle("unknown type")
}

func TestFormatFileStatus(t *testing.T) {
	tests := []struct {
		status FormatFileStatus
		want   int
	}{
		{FormatFileStatusSuccess, 0},
		{FormatFileStatusNoChange, 1},
		{FormatFileStatusFailed, 2},
	}

	for _, tt := range tests {
		if int(tt.status) != tt.want {
			t.Errorf("FormatFileStatus = %d, want %d", tt.status, tt.want)
		}
	}
}

func TestFormatCommandUseCaseOutputData(t *testing.T) {
	output := &FormatCommandUseCaseOutputData{
		Message:  "test message",
		Filename: "test.md",
		FileId:   "123",
		Status:   FormatFileStatusSuccess,
	}

	if output.Message != "test message" {
		t.Errorf("Message = %q, want %q", output.Message, "test message")
	}
	if output.Filename != "test.md" {
		t.Errorf("Filename = %q, want %q", output.Filename, "test.md")
	}
	if output.FileId != "123" {
		t.Errorf("FileId = %q, want %q", output.FileId, "123")
	}
	if output.Status != FormatFileStatusSuccess {
		t.Errorf("Status = %d, want %d", output.Status, FormatFileStatusSuccess)
	}
}

// Init tests

func TestNewInitUseCaseBus(t *testing.T) {
	mock := &mockInitSettingUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (InitSettingUseCase, error) {
		return mock, nil
	})

	bus, err := NewInitUseCaseBus(injector)
	if err != nil {
		t.Errorf("NewInitUseCaseBus() error = %v", err)
	}
	if bus == nil {
		t.Error("NewInitUseCaseBus() returned nil")
	}
}

func TestInitUseCaseBus_Handle(t *testing.T) {
	mock := &mockInitSettingUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (InitSettingUseCase, error) {
		return mock, nil
	})

	bus, _ := NewInitUseCaseBus(injector)
	bus.Handle(&InitSettingUseCaseInputData{})

	if !mock.handleCalled {
		t.Error("Handle() did not call configure.Handle()")
	}
}

func TestInitUseCaseBus_Handle_Panic(t *testing.T) {
	mock := &mockInitSettingUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (InitSettingUseCase, error) {
		return mock, nil
	})

	bus, _ := NewInitUseCaseBus(injector)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Handle() should panic for unknown input type")
		}
	}()

	bus.Handle("unknown type")
}

func TestInitSettingProject(t *testing.T) {
	project := InitSettingProject{
		DriveFolder:  InitSettingProjectDriveFolder("folder_id"),
		SiteUrl:      InitSettingProjectSiteUrl("https://example.com"),
		Url:          InitSettingProjectUrl("https://repo.com"),
		Branch:       InitSettingProjectBranch("main"),
		TemplatePath: InitSettingProjectTemplatePath("/templates"),
		AssetPath:    InitSettingProjectAssetPath("/assets"),
	}

	if string(project.DriveFolder) != "folder_id" {
		t.Errorf("DriveFolder = %q, want %q", project.DriveFolder, "folder_id")
	}
	if string(project.SiteUrl) != "https://example.com" {
		t.Errorf("SiteUrl = %q, want %q", project.SiteUrl, "https://example.com")
	}
	if string(project.Url) != "https://repo.com" {
		t.Errorf("Url = %q, want %q", project.Url, "https://repo.com")
	}
	if string(project.Branch) != "main" {
		t.Errorf("Branch = %q, want %q", project.Branch, "main")
	}
}

func TestInitSettingConfirmTypes(t *testing.T) {
	overwrite := InitSettingConfirmOverwrite(true)
	gitWarning := InitSettingConfirmGitWarning(false)

	if !bool(overwrite) {
		t.Error("InitSettingConfirmOverwrite should be true")
	}
	if bool(gitWarning) {
		t.Error("InitSettingConfirmGitWarning should be false")
	}
}

func TestInitSettingUseCaseOutputData(t *testing.T) {
	output := &InitSettingUseCaseOutputData{
		InitUsecaseOutputDataImpl: InitUsecaseOutputDataImpl{
			Message: "test message",
		},
		ProjectConfigured: true,
		Project: InitSettingProject{
			DriveFolder: "folder",
		},
	}

	if output.Message != "test message" {
		t.Errorf("Message = %q, want %q", output.Message, "test message")
	}
	if !output.ProjectConfigured {
		t.Error("ProjectConfigured should be true")
	}
}

// Root tests

func TestNewRootUseCaseBus(t *testing.T) {
	mock := &mockRootCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (RootCommandUseCase, error) {
		return mock, nil
	})

	bus, err := NewRootUseCaseBus(injector)
	if err != nil {
		t.Errorf("NewRootUseCaseBus() error = %v", err)
	}
	if bus == nil {
		t.Error("NewRootUseCaseBus() returned nil")
	}
}

func TestRootUseCaseBus_Handle(t *testing.T) {
	mock := &mockRootCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (RootCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewRootUseCaseBus(injector)
	bus.Handle(&RootCommandUseCaseInputData{Version: "1.0.0"})

	if !mock.handleCalled {
		t.Error("Handle() did not call command.Handle()")
	}
}

func TestRootUseCaseBus_Handle_Panic(t *testing.T) {
	mock := &mockRootCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (RootCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewRootUseCaseBus(injector)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Handle() should panic for unknown input type")
		}
	}()

	bus.Handle("unknown type")
}

func TestRootCommandUseCaseInputData(t *testing.T) {
	input := &RootCommandUseCaseInputData{
		Version: "1.2.3",
	}
	if input.Version != "1.2.3" {
		t.Errorf("Version = %q, want %q", input.Version, "1.2.3")
	}
}

func TestRootCommandUseCaseOutputData(t *testing.T) {
	output := &RootCommandUseCaseOutputData{
		Message: "test message",
	}
	if output.Message != "test message" {
		t.Errorf("Message = %q, want %q", output.Message, "test message")
	}
}

// Update tests

func TestNewUpdateUseCaseBus(t *testing.T) {
	mock := &mockUpdateCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (UpdateCommandUseCase, error) {
		return mock, nil
	})

	bus, err := NewUpdateUseCaseBus(injector)
	if err != nil {
		t.Errorf("NewUpdateUseCaseBus() error = %v", err)
	}
	if bus == nil {
		t.Error("NewUpdateUseCaseBus() returned nil")
	}
}

func TestUpdateUseCaseBus_Handle(t *testing.T) {
	mock := &mockUpdateCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (UpdateCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewUpdateUseCaseBus(injector)
	bus.Handle(&UpdateCommandUseCaseInputData{})

	if !mock.handleCalled {
		t.Error("Handle() did not call command.Handle()")
	}
}

func TestUpdateUseCaseBus_Handle_Panic(t *testing.T) {
	mock := &mockUpdateCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (UpdateCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewUpdateUseCaseBus(injector)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Handle() should panic for unknown input type")
		}
	}()

	bus.Handle("unknown type")
}

func TestUpdateCommandUseCaseOutputData(t *testing.T) {
	output := &UpdateCommandUseCaseOutputData{
		Message: "test message",
	}
	if output.Message != "test message" {
		t.Errorf("Message = %q, want %q", output.Message, "test message")
	}
}
