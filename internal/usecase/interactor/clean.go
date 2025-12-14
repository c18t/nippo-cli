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
	repo, err := do.Invoke[repository.AssetRepository](i)
	if err != nil {
		return nil, err
	}
	p, err := do.Invoke[presenter.CleanCommandPresenter](i)
	if err != nil {
		return nil, err
	}
	return &cleanCommandInteractor{
		repository: repo,
		presenter:  p,
	}, nil
}

func (u *cleanCommandInteractor) Handle(input *port.CleanCommandUseCaseInputData) {
	output := &port.CleanCommandUseCaseOutputData{}

	output.Message = "cleaning cache files..."
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

	// Progress() で開始したスピナーは自動的に "ok." が付く
	u.presenter.StopProgress()
}
