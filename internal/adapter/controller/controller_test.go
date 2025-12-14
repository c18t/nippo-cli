package controller

import (
	"testing"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// Mock implementations

type mockBuildUseCaseBus struct {
	handleCalled bool
}

func (m *mockBuildUseCaseBus) Handle(input port.BuildUseCaseInputData) {
	m.handleCalled = true
}

type mockCleanUseCaseBus struct {
	handleCalled bool
}

func (m *mockCleanUseCaseBus) Handle(input port.CleanUseCaseInputData) {
	m.handleCalled = true
}

type mockDeployUseCaseBus struct {
	handleCalled bool
}

func (m *mockDeployUseCaseBus) Handle(input port.DeployUseCaseInputData) {
	m.handleCalled = true
}

type mockUpdateUseCaseBus struct {
	handleCalled bool
}

func (m *mockUpdateUseCaseBus) Handle(input port.UpdateUseCaseInputData) {
	m.handleCalled = true
}

type mockFormatUseCaseBus struct {
	handleCalled bool
}

func (m *mockFormatUseCaseBus) Handle(input port.FormatUseCaseInputData) {
	m.handleCalled = true
}

type mockInitUseCaseBus struct {
	handleCalled bool
}

func (m *mockInitUseCaseBus) Handle(input port.InitUseCaseInputData) {
	m.handleCalled = true
}

type mockRootUseCaseBus struct {
	handleCalled bool
	input        port.RootUseCaseInputData
}

func (m *mockRootUseCaseBus) Handle(input port.RootUseCaseInputData) {
	m.handleCalled = true
	m.input = input
}

type mockAuthUseCaseBus struct {
	handleCalled bool
}

func (m *mockAuthUseCaseBus) Handle(input *port.AuthUseCaseInputData) {
	m.handleCalled = true
}

type mockDoctorUseCase struct {
	handleCalled bool
}

func (m *mockDoctorUseCase) Handle(input *port.DoctorUseCaseInputData) {
	m.handleCalled = true
}

// Tests for BuildController

func TestNewBuildController(t *testing.T) {
	mock := &mockBuildUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.BuildUseCaseBus, error) {
		return mock, nil
	})

	ctrl, err := NewBuildController(injector)
	if err != nil {
		t.Errorf("NewBuildController() error = %v", err)
	}
	if ctrl == nil {
		t.Error("NewBuildController() returned nil")
	}
}

func TestBuildController_Params(t *testing.T) {
	mock := &mockBuildUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.BuildUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewBuildController(injector)
	params := ctrl.Params()

	if params == nil {
		t.Error("Params() returned nil")
	}
}

func TestBuildController_Exec(t *testing.T) {
	mock := &mockBuildUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.BuildUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewBuildController(injector)
	cmd := &cobra.Command{}

	err := ctrl.Exec(cmd, []string{})
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
	if !mock.handleCalled {
		t.Error("Exec() did not call bus.Handle()")
	}
}

// Tests for CleanController

func TestNewCleanController(t *testing.T) {
	mock := &mockCleanUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.CleanUseCaseBus, error) {
		return mock, nil
	})

	ctrl, err := NewCleanController(injector)
	if err != nil {
		t.Errorf("NewCleanController() error = %v", err)
	}
	if ctrl == nil {
		t.Error("NewCleanController() returned nil")
	}
}

func TestCleanController_Params(t *testing.T) {
	mock := &mockCleanUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.CleanUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewCleanController(injector)
	params := ctrl.Params()

	if params == nil {
		t.Error("Params() returned nil")
	}
}

func TestCleanController_Exec(t *testing.T) {
	mock := &mockCleanUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.CleanUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewCleanController(injector)
	cmd := &cobra.Command{}

	err := ctrl.Exec(cmd, []string{})
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
	if !mock.handleCalled {
		t.Error("Exec() did not call bus.Handle()")
	}
}

// Tests for DeployController

func TestNewDeployController(t *testing.T) {
	mock := &mockDeployUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.DeployUseCaseBus, error) {
		return mock, nil
	})

	ctrl, err := NewDeployController(injector)
	if err != nil {
		t.Errorf("NewDeployController() error = %v", err)
	}
	if ctrl == nil {
		t.Error("NewDeployController() returned nil")
	}
}

func TestDeployController_Exec(t *testing.T) {
	mock := &mockDeployUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.DeployUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewDeployController(injector)
	cmd := &cobra.Command{}

	err := ctrl.Exec(cmd, []string{})
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
	if !mock.handleCalled {
		t.Error("Exec() did not call bus.Handle()")
	}
}

// Tests for UpdateController

func TestNewUpdateController(t *testing.T) {
	mock := &mockUpdateUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.UpdateUseCaseBus, error) {
		return mock, nil
	})

	ctrl, err := NewUpdateController(injector)
	if err != nil {
		t.Errorf("NewUpdateController() error = %v", err)
	}
	if ctrl == nil {
		t.Error("NewUpdateController() returned nil")
	}
}

func TestUpdateController_Exec(t *testing.T) {
	mock := &mockUpdateUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.UpdateUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewUpdateController(injector)
	cmd := &cobra.Command{}

	err := ctrl.Exec(cmd, []string{})
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
	if !mock.handleCalled {
		t.Error("Exec() did not call bus.Handle()")
	}
}

// Tests for FormatController

func TestNewFormatController(t *testing.T) {
	mock := &mockFormatUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.FormatUseCaseBus, error) {
		return mock, nil
	})

	ctrl, err := NewFormatController(injector)
	if err != nil {
		t.Errorf("NewFormatController() error = %v", err)
	}
	if ctrl == nil {
		t.Error("NewFormatController() returned nil")
	}
}

func TestFormatController_Exec(t *testing.T) {
	mock := &mockFormatUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.FormatUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewFormatController(injector)
	cmd := &cobra.Command{}

	err := ctrl.Exec(cmd, []string{})
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
	if !mock.handleCalled {
		t.Error("Exec() did not call bus.Handle()")
	}
}

// Tests for InitController

func TestNewInitController(t *testing.T) {
	mock := &mockInitUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.InitUseCaseBus, error) {
		return mock, nil
	})

	ctrl, err := NewInitController(injector)
	if err != nil {
		t.Errorf("NewInitController() error = %v", err)
	}
	if ctrl == nil {
		t.Error("NewInitController() returned nil")
	}
}

func TestInitController_Exec(t *testing.T) {
	mock := &mockInitUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.InitUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewInitController(injector)
	cmd := &cobra.Command{}

	err := ctrl.Exec(cmd, []string{})
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
	if !mock.handleCalled {
		t.Error("Exec() did not call bus.Handle()")
	}
}

// Tests for RootController

func TestNewRootController(t *testing.T) {
	mock := &mockRootUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.RootUseCaseBus, error) {
		return mock, nil
	})

	ctrl, err := NewRootController(injector)
	if err != nil {
		t.Errorf("NewRootController() error = %v", err)
	}
	if ctrl == nil {
		t.Error("NewRootController() returned nil")
	}
}

func TestRootController_Version(t *testing.T) {
	mock := &mockRootUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.RootUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewRootController(injector)

	// Set version
	ctrl.Version("v1.0.0")

	// Get version
	v := ctrl.Version()
	if v != "v1.0.0" {
		t.Errorf("Version() = %q, want %q", v, "v1.0.0")
	}
}

func TestRootController_Params(t *testing.T) {
	mock := &mockRootUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.RootUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewRootController(injector)
	params := ctrl.Params()

	if params == nil {
		t.Error("Params() returned nil")
	}
}

func TestRootController_Exec_ShowVersion(t *testing.T) {
	mock := &mockRootUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.RootUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewRootController(injector)
	ctrl.Version("v1.0.0")
	ctrl.Params().Version = true

	cmd := &cobra.Command{}
	err := ctrl.Exec(cmd, []string{})

	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
	if !mock.handleCalled {
		t.Error("Exec() did not call bus.Handle()")
	}
}

func TestRootController_Exec_ShowHelp(t *testing.T) {
	mock := &mockRootUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.RootUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewRootController(injector)
	ctrl.Params().Version = false
	ctrl.Params().LicenseNotice = false

	cmd := &cobra.Command{}
	err := ctrl.Exec(cmd, []string{})

	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
	// When neither Version nor LicenseNotice, it shows help (bus not called)
	if mock.handleCalled {
		t.Error("Exec() should not call bus.Handle() when showing help")
	}
}

func TestRootController_Exec_LicenseNotice(t *testing.T) {
	mock := &mockRootUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.RootUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewRootController(injector)
	ctrl.Params().LicenseNotice = true

	cmd := &cobra.Command{}
	err := ctrl.Exec(cmd, []string{})

	// LicenseNotice returns "not implemented" error
	if err == nil {
		t.Error("Exec() should return error for LicenseNotice (not implemented)")
	}
}

func TestRootController_RequireConfig_NoError(t *testing.T) {
	mock := &mockRootUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.RootUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewRootController(injector)
	// No InitConfig called, so initConfigErr is nil
	cmd := &cobra.Command{}

	err := ctrl.RequireConfig(cmd)
	if err != nil {
		t.Errorf("RequireConfig() error = %v, want nil", err)
	}
}

func TestRootController_InitConfigErr(t *testing.T) {
	mock := &mockRootUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.RootUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewRootController(injector)
	err := ctrl.InitConfigErr()

	if err != nil {
		t.Errorf("InitConfigErr() = %v, want nil", err)
	}
}

func TestRootController_RequireConfig_ConfigNotFound(t *testing.T) {
	mock := &mockRootUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.RootUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewRootController(injector)

	// Simulate InitConfig with non-existent config file
	ctrl.Params().ConfigFile = "/nonexistent/config.yaml"
	ctrl.InitConfig()

	cmd := &cobra.Command{}
	err := ctrl.RequireConfig(cmd)

	// Should return error since config file doesn't exist
	if err == nil {
		t.Error("RequireConfig() should return error for non-existent config")
	}
}

// Tests for AuthController

func TestNewAuthController(t *testing.T) {
	mock := &mockAuthUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.AuthUseCaseBus, error) {
		return mock, nil
	})

	ctrl, err := NewAuthController(injector)
	if err != nil {
		t.Errorf("NewAuthController() error = %v", err)
	}
	if ctrl == nil {
		t.Error("NewAuthController() returned nil")
	}
}

func TestAuthController_Params(t *testing.T) {
	mock := &mockAuthUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.AuthUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewAuthController(injector)
	params := ctrl.Params()

	if params == nil {
		t.Error("Params() returned nil")
	}
}

func TestAuthController_Exec(t *testing.T) {
	mock := &mockAuthUseCaseBus{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.AuthUseCaseBus, error) {
		return mock, nil
	})

	ctrl, _ := NewAuthController(injector)
	cmd := &cobra.Command{}

	err := ctrl.Exec(cmd, []string{})
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
	if !mock.handleCalled {
		t.Error("Exec() did not call bus.Handle()")
	}
}

// Tests for DoctorController

func TestNewDoctorController(t *testing.T) {
	mock := &mockDoctorUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.DoctorUseCase, error) {
		return mock, nil
	})

	ctrl, err := NewDoctorController(injector)
	if err != nil {
		t.Errorf("NewDoctorController() error = %v", err)
	}
	if ctrl == nil {
		t.Error("NewDoctorController() returned nil")
	}
}

func TestDoctorController_Exec(t *testing.T) {
	mock := &mockDoctorUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (port.DoctorUseCase, error) {
		return mock, nil
	})

	// Setup global config for doctor
	core.Cfg = &core.Config{}
	core.Cfg.Paths.DataDir = "/tmp"
	core.Cfg.Paths.CacheDir = "/tmp"

	ctrl, _ := NewDoctorController(injector)
	cmd := &cobra.Command{}

	err := ctrl.Exec(cmd, []string{})
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
	if !mock.handleCalled {
		t.Error("Exec() did not call useCase.Handle()")
	}
}

// Test param structs

func TestBuildParams(t *testing.T) {
	params := &BuildParams{}
	// Verify struct can be created
	_ = params
}

func TestCleanParams(t *testing.T) {
	params := &CleanParams{}
	// Verify struct can be created
	_ = params
}

func TestRootParams(t *testing.T) {
	params := &RootParams{
		ConfigFile:    "/path/to/config.yaml",
		Version:       true,
		LicenseNotice: false,
	}

	if params.ConfigFile != "/path/to/config.yaml" {
		t.Errorf("ConfigFile = %q, want %q", params.ConfigFile, "/path/to/config.yaml")
	}
	if !params.Version {
		t.Error("Version should be true")
	}
	if params.LicenseNotice {
		t.Error("LicenseNotice should be false")
	}
}

func TestAuthParams(t *testing.T) {
	params := &AuthParams{}
	// Verify struct can be created
	_ = params
}
