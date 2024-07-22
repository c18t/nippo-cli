package interactor

import (
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

type cleanCommandInteractor struct {
	repository repository.AssetRepository      `do:""`
	presenter  presenter.CleanCommandPresenter `do:""`
}

func NewCleanCommandInteractor(i do.Injector) (port.CleanCommandUseCase, error) {
	return &cleanCommandInteractor{
		repository: do.MustInvoke[repository.AssetRepository](i),
		presenter:  do.MustInvoke[presenter.CleanCommandPresenter](i),
	}, nil
}

func (u *cleanCommandInteractor) Handle(input *port.CleanCommandUseCaseInputData) {
	output := &port.CleanCommandUseCaseOutputData{}

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
