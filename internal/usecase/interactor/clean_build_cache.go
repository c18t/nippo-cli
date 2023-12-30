package interactor

import (
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"go.uber.org/dig"
)

type cleanBuildCacheInteractor struct {
	repository repository.AssetRepository
	presenter  presenter.CleanBuildCachePresenter
}
type inCleanBuildCacheInteractor struct {
	dig.In
	Repository repository.AssetRepository
	Presenter  presenter.CleanBuildCachePresenter
}

func NewCleanBuildCacheInteractor(cleanDeps inCleanBuildCacheInteractor) port.CleanBuildCacheUsecase {
	return &cleanBuildCacheInteractor{
		repository: cleanDeps.Repository,
		presenter:  cleanDeps.Presenter,
	}
}

func (u *cleanBuildCacheInteractor) Handle(input *port.CleanBuildCacheUsecaseInputData) {
	output := &port.CleanBuildCacheUsecaseOutputData{}

	output.Message = "clean cache files... "
	u.presenter.Progress(output)

	u.repository.CleanNippoCache()
	u.repository.CleanBuildCache()

	output.Message = "ok. "
	u.presenter.Complete(output)
}
