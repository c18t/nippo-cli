package interactor

import (
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
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

	output.Message = "cleaning cache files... "
	u.presenter.Progress(output)

	if err := u.repository.CleanNippoCache(); err != nil {
		u.presenter.Suspend(err)
		return
	}

	core.Cfg.ResetLastUpdateCheckTimestamp()
	if err := core.Cfg.SaveConfig(); err != nil {
		u.presenter.Suspend(err)
		return
	}

	if err := u.repository.CleanBuildCache(); err != nil {
		u.presenter.Suspend(err)
		return
	}

	output.Message = "ok. "
	u.presenter.Complete(output)
}
