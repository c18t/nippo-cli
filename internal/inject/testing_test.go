package inject

import (
	"os"
	"testing"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/model"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/c18t/nippo-cli/internal/domain/service"
	"github.com/samber/do/v2"
	"google.golang.org/api/drive/v3"
)

func TestNewTestInjector_NilOptions(t *testing.T) {
	injector := NewTestInjector(nil)
	if injector == nil {
		t.Error("NewTestInjector(nil) should return non-nil injector")
	}
}

func TestNewTestInjector_EmptyOptions(t *testing.T) {
	injector := NewTestInjector(&TestBasePackageOptions{})
	if injector == nil {
		t.Error("NewTestInjector with empty options should return non-nil injector")
	}
}

// Mock implementations for testing
type mockDriveFileProvider struct{}

func (m *mockDriveFileProvider) List(param *repository.QueryListParam) (*drive.FileList, error) {
	return nil, nil
}
func (m *mockDriveFileProvider) Download(id string) ([]byte, error)    { return nil, nil }
func (m *mockDriveFileProvider) Update(fileId string, content []byte) error { return nil }
func (m *mockDriveFileProvider) Shutdown() error                       { return nil }
func (m *mockDriveFileProvider) HealthCheck() error                    { return nil }

type mockLocalFileProvider struct{}

func (m *mockLocalFileProvider) List(param *repository.QueryListParam) ([]os.DirEntry, error) {
	return nil, nil
}
func (m *mockLocalFileProvider) Read(baseDir, filePath string) ([]byte, error) { return nil, nil }
func (m *mockLocalFileProvider) Write(filePath string, content []byte) error   { return nil }
func (m *mockLocalFileProvider) Copy(baseDir, destPath, srcPath string) error  { return nil }

type mockRemoteNippoQuery struct{}

func (m *mockRemoteNippoQuery) List(param *repository.QueryListParam, option *repository.QueryListOption) ([]model.Nippo, error) {
	return nil, nil
}
func (m *mockRemoteNippoQuery) Download(nippo *model.Nippo) error               { return nil }
func (m *mockRemoteNippoQuery) Update(nippo *model.Nippo, content []byte) error { return nil }

type mockLocalNippoQuery struct{}

func (m *mockLocalNippoQuery) Exist(date *model.NippoDate) bool { return false }
func (m *mockLocalNippoQuery) Find(date *model.NippoDate) (*model.Nippo, error) {
	return nil, nil
}
func (m *mockLocalNippoQuery) List(param *repository.QueryListParam, option *repository.QueryListOption) ([]model.Nippo, error) {
	return nil, nil
}
func (m *mockLocalNippoQuery) Load(nippo *model.Nippo) error { return nil }

type mockLocalNippoCommand struct{}

func (m *mockLocalNippoCommand) Create(nippo *model.Nippo) error { return nil }
func (m *mockLocalNippoCommand) Delete(nippo *model.Nippo) error { return nil }

type mockAssetRepository struct{}

func (m *mockAssetRepository) CleanNippoCache() error { return nil }
func (m *mockAssetRepository) CleanBuildCache() error { return nil }

type mockNippoFacade struct{}

func (m *mockNippoFacade) Send(request *service.NippoFacadeRequest, option *service.NippoFacadeOption) (*service.NippoFacadeReponse, error) {
	return nil, nil
}

type mockTemplateService struct{}

func (m *mockTemplateService) SaveTo(filePath string, templateName string, data any) error {
	return nil
}

func TestNewTestInjector_WithDriveFileProvider(t *testing.T) {
	mock := &mockDriveFileProvider{}
	injector := NewTestInjector(&TestBasePackageOptions{
		DriveFileProvider: mock,
	})
	if injector == nil {
		t.Error("NewTestInjector with DriveFileProvider should return non-nil injector")
	}

	// Verify the mock is returned
	provider, err := do.Invoke[gateway.DriveFileProvider](injector)
	if err != nil {
		t.Errorf("Invoke[DriveFileProvider] error = %v", err)
	}
	if provider != mock {
		t.Error("Invoke[DriveFileProvider] should return mock")
	}
}

func TestNewTestInjector_WithLocalFileProvider(t *testing.T) {
	mock := &mockLocalFileProvider{}
	injector := NewTestInjector(&TestBasePackageOptions{
		LocalFileProvider: mock,
	})
	if injector == nil {
		t.Error("NewTestInjector with LocalFileProvider should return non-nil injector")
	}

	provider, err := do.Invoke[gateway.LocalFileProvider](injector)
	if err != nil {
		t.Errorf("Invoke[LocalFileProvider] error = %v", err)
	}
	if provider != mock {
		t.Error("Invoke[LocalFileProvider] should return mock")
	}
}

func TestNewTestInjector_WithAllMocks(t *testing.T) {
	driveProvider := &mockDriveFileProvider{}
	localProvider := &mockLocalFileProvider{}

	injector := NewTestInjector(&TestBasePackageOptions{
		DriveFileProvider: driveProvider,
		LocalFileProvider: localProvider,
	})
	if injector == nil {
		t.Error("NewTestInjector with all mocks should return non-nil injector")
	}
}

func TestTestBasePackageOptions_AllFields(t *testing.T) {
	opts := &TestBasePackageOptions{
		Config:            &core.Config{},
		DriveFileProvider: &mockDriveFileProvider{},
		LocalFileProvider: &mockLocalFileProvider{},
		RemoteNippoQuery:  &mockRemoteNippoQuery{},
		LocalNippoQuery:   &mockLocalNippoQuery{},
		LocalNippoCommand: &mockLocalNippoCommand{},
		AssetRepository:   &mockAssetRepository{},
		NippoFacade:       &mockNippoFacade{},
		TemplateService:   &mockTemplateService{},
	}

	if opts.Config == nil {
		t.Error("Config should not be nil")
	}
	if opts.DriveFileProvider == nil {
		t.Error("DriveFileProvider should not be nil")
	}
}

// Test that repository interfaces can be overridden
func TestNewTestInjector_WithRepositoryMocks(t *testing.T) {
	injector := NewTestInjector(&TestBasePackageOptions{
		RemoteNippoQuery:  &mockRemoteNippoQuery{},
		LocalNippoQuery:   &mockLocalNippoQuery{},
		LocalNippoCommand: &mockLocalNippoCommand{},
		AssetRepository:   &mockAssetRepository{},
	})
	if injector == nil {
		t.Error("NewTestInjector with repository mocks should return non-nil injector")
	}

	// Verify mocks are returned
	remoteQuery, err := do.Invoke[repository.RemoteNippoQuery](injector)
	if err != nil {
		t.Errorf("Invoke[RemoteNippoQuery] error = %v", err)
	}
	if remoteQuery == nil {
		t.Error("RemoteNippoQuery should not be nil")
	}

	localQuery, err := do.Invoke[repository.LocalNippoQuery](injector)
	if err != nil {
		t.Errorf("Invoke[LocalNippoQuery] error = %v", err)
	}
	if localQuery == nil {
		t.Error("LocalNippoQuery should not be nil")
	}

	localCommand, err := do.Invoke[repository.LocalNippoCommand](injector)
	if err != nil {
		t.Errorf("Invoke[LocalNippoCommand] error = %v", err)
	}
	if localCommand == nil {
		t.Error("LocalNippoCommand should not be nil")
	}

	assetRepo, err := do.Invoke[repository.AssetRepository](injector)
	if err != nil {
		t.Errorf("Invoke[AssetRepository] error = %v", err)
	}
	if assetRepo == nil {
		t.Error("AssetRepository should not be nil")
	}
}

// Test service overrides
func TestNewTestInjector_WithServiceMocks(t *testing.T) {
	injector := NewTestInjector(&TestBasePackageOptions{
		NippoFacade:     &mockNippoFacade{},
		TemplateService: &mockTemplateService{},
	})
	if injector == nil {
		t.Error("NewTestInjector with service mocks should return non-nil injector")
	}

	facade, err := do.Invoke[service.NippoFacade](injector)
	if err != nil {
		t.Errorf("Invoke[NippoFacade] error = %v", err)
	}
	if facade == nil {
		t.Error("NippoFacade should not be nil")
	}

	templateSvc, err := do.Invoke[service.TemplateService](injector)
	if err != nil {
		t.Errorf("Invoke[TemplateService] error = %v", err)
	}
	if templateSvc == nil {
		t.Error("TemplateService should not be nil")
	}
}

// Tests for container.go

func TestGetInjector(t *testing.T) {
	injector := GetInjector()
	if injector == nil {
		t.Error("GetInjector() should return non-nil injector")
	}
}

func TestGetInjector_Singleton(t *testing.T) {
	injector1 := GetInjector()
	injector2 := GetInjector()
	if injector1 != injector2 {
		t.Error("GetInjector() should return the same instance (singleton)")
	}
}

func TestBasePackage_NotNil(t *testing.T) {
	if BasePackage == nil {
		t.Error("BasePackage should not be nil")
	}
}

// Tests for doctor.go

func TestDoctorPackage_NotNil(t *testing.T) {
	if DoctorPackage == nil {
		t.Error("DoctorPackage should not be nil")
	}
}

func TestInjectorDoctor_NotNil(t *testing.T) {
	if InjectorDoctor == nil {
		t.Error("InjectorDoctor should not be nil")
	}
}
