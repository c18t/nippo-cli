package service

import (
	"testing"

	"github.com/c18t/nippo-cli/internal/domain/model"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	ds "github.com/c18t/nippo-cli/internal/domain/service"
	"github.com/samber/do/v2"
	"google.golang.org/api/drive/v3"
)

// Mock implementations
type mockRemoteNippoQuery struct {
	nippos    []model.Nippo
	listErr   error
	downloadErr error
	updateErr error
}

func (m *mockRemoteNippoQuery) List(param *repository.QueryListParam, option *repository.QueryListOption) ([]model.Nippo, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.nippos, nil
}

func (m *mockRemoteNippoQuery) Download(nippo *model.Nippo) error {
	if m.downloadErr != nil {
		return m.downloadErr
	}
	nippo.Content = []byte("downloaded content")
	return nil
}

func (m *mockRemoteNippoQuery) Update(nippo *model.Nippo, content []byte) error {
	return m.updateErr
}

type mockLocalNippoQuery struct{}

func (m *mockLocalNippoQuery) Exist(date *model.NippoDate) bool { return true }
func (m *mockLocalNippoQuery) Find(date *model.NippoDate) (*model.Nippo, error) {
	return nil, nil
}
func (m *mockLocalNippoQuery) List(param *repository.QueryListParam, option *repository.QueryListOption) ([]model.Nippo, error) {
	return nil, nil
}
func (m *mockLocalNippoQuery) Load(nippo *model.Nippo) error { return nil }

type mockLocalNippoCommand struct {
	createErr error
	deleteErr error
}

func (m *mockLocalNippoCommand) Create(nippo *model.Nippo) error { return m.createErr }
func (m *mockLocalNippoCommand) Delete(nippo *model.Nippo) error { return m.deleteErr }

func TestNewNippoFacade(t *testing.T) {
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (repository.RemoteNippoQuery, error) {
		return &mockRemoteNippoQuery{}, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoQuery, error) {
		return &mockLocalNippoQuery{}, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoCommand, error) {
		return &mockLocalNippoCommand{}, nil
	})

	facade, err := NewNippoFacade(injector)
	if err != nil {
		t.Errorf("NewNippoFacade() error = %v", err)
	}
	if facade == nil {
		t.Error("NewNippoFacade() returned nil")
	}
}

func TestNippoFacade_Send_Search(t *testing.T) {
	remoteQuery := &mockRemoteNippoQuery{
		nippos: []model.Nippo{
			{Date: model.NewNippoDate("2024-01-15.md"), RemoteFile: &drive.File{Id: "1", Name: "2024-01-15.md"}},
			{Date: model.NewNippoDate("2024-01-16.md"), RemoteFile: &drive.File{Id: "2", Name: "2024-01-16.md"}},
		},
	}

	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (repository.RemoteNippoQuery, error) {
		return remoteQuery, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoQuery, error) {
		return &mockLocalNippoQuery{}, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoCommand, error) {
		return &mockLocalNippoCommand{}, nil
	})

	facade, _ := NewNippoFacade(injector)

	resp, err := facade.Send(&ds.NippoFacadeRequest{
		Action: ds.NippoFacadeActionSearch,
		Query:  &repository.QueryListParam{},
		Option: &repository.QueryListOption{},
	}, nil)

	if err != nil {
		t.Errorf("Send() error = %v", err)
	}
	if resp == nil {
		t.Fatal("Send() returned nil response")
	}
	// When Action is Search only (no Download or Cache), the content is empty
	// The count reflects the nippoList which comes from request.Content when no Cache action
	if resp.Result == nil {
		t.Error("Result should not be nil")
	}
}

func TestNippoFacade_Send_Download(t *testing.T) {
	remoteQuery := &mockRemoteNippoQuery{
		nippos: []model.Nippo{
			{Date: model.NewNippoDate("2024-01-15.md"), RemoteFile: &drive.File{Id: "1", Name: "2024-01-15.md"}},
		},
	}

	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (repository.RemoteNippoQuery, error) {
		return remoteQuery, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoQuery, error) {
		return &mockLocalNippoQuery{}, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoCommand, error) {
		return &mockLocalNippoCommand{}, nil
	})

	facade, _ := NewNippoFacade(injector)

	resp, err := facade.Send(&ds.NippoFacadeRequest{
		Action: ds.NippoFacadeActionDownload,
		Query:  &repository.QueryListParam{},
		Option: &repository.QueryListOption{},
	}, nil)

	if err != nil {
		t.Errorf("Send() error = %v", err)
	}
	if resp == nil {
		t.Fatal("Send() returned nil response")
	}
}

func TestNippoFacade_Send_Cache(t *testing.T) {
	remoteQuery := &mockRemoteNippoQuery{
		nippos: []model.Nippo{
			{Date: model.NewNippoDate("2024-01-15.md"), RemoteFile: &drive.File{Id: "1", Name: "2024-01-15.md"}},
		},
	}

	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (repository.RemoteNippoQuery, error) {
		return remoteQuery, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoQuery, error) {
		return &mockLocalNippoQuery{}, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoCommand, error) {
		return &mockLocalNippoCommand{}, nil
	})

	facade, _ := NewNippoFacade(injector)

	resp, err := facade.Send(&ds.NippoFacadeRequest{
		Action: ds.NippoFacadeActionDownload | ds.NippoFacadeActionCache,
		Query:  &repository.QueryListParam{},
		Option: &repository.QueryListOption{},
	}, nil)

	if err != nil {
		t.Errorf("Send() error = %v", err)
	}
	if resp == nil {
		t.Fatal("Send() returned nil response")
	}
}

func TestNippoFacade_Send_WithProgressCallback(t *testing.T) {
	remoteQuery := &mockRemoteNippoQuery{
		nippos: []model.Nippo{
			{Date: model.NewNippoDate("2024-01-15.md"), RemoteFile: &drive.File{Id: "1", Name: "2024-01-15.md"}},
		},
	}

	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (repository.RemoteNippoQuery, error) {
		return remoteQuery, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoQuery, error) {
		return &mockLocalNippoQuery{}, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoCommand, error) {
		return &mockLocalNippoCommand{}, nil
	})

	facade, _ := NewNippoFacade(injector)

	progressCalled := false
	resp, err := facade.Send(&ds.NippoFacadeRequest{
		Action: ds.NippoFacadeActionDownload,
		Query:  &repository.QueryListParam{},
		Option: &repository.QueryListOption{},
	}, &ds.NippoFacadeOption{
		OnProgress: func(filename string, fileId string, current int, total int) bool {
			progressCalled = true
			return true
		},
	})

	if err != nil {
		t.Errorf("Send() error = %v", err)
	}
	if resp == nil {
		t.Fatal("Send() returned nil response")
	}
	if !progressCalled {
		t.Error("Progress callback was not called")
	}
}

func TestNippoFacade_Send_CancelledByCallback(t *testing.T) {
	remoteQuery := &mockRemoteNippoQuery{
		nippos: []model.Nippo{
			{Date: model.NewNippoDate("2024-01-15.md"), RemoteFile: &drive.File{Id: "1", Name: "2024-01-15.md"}},
		},
	}

	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (repository.RemoteNippoQuery, error) {
		return remoteQuery, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoQuery, error) {
		return &mockLocalNippoQuery{}, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.LocalNippoCommand, error) {
		return &mockLocalNippoCommand{}, nil
	})

	facade, _ := NewNippoFacade(injector)

	_, err := facade.Send(&ds.NippoFacadeRequest{
		Action: ds.NippoFacadeActionDownload,
		Query:  &repository.QueryListParam{},
		Option: &repository.QueryListOption{},
	}, &ds.NippoFacadeOption{
		OnProgress: func(filename string, fileId string, current int, total int) bool {
			return false // Cancel
		},
	})

	if err != ds.ErrCancelled {
		t.Errorf("Send() error = %v, want ErrCancelled", err)
	}
}
