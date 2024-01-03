package interactor

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"go.uber.org/dig"
)

type deploySiteInteractor struct {
	provider  gateway.LocalFileProvider
	presenter presenter.DeploySitePresenter
}

type inDeploySiteInteractor struct {
	dig.In
	Provider  gateway.LocalFileProvider
	Presenter presenter.DeploySitePresenter
}

func NewDeploySiteInteractor(deployDeps inDeploySiteInteractor) port.DeploySiteUsecase {
	return &deploySiteInteractor{
		provider:  deployDeps.Provider,
		presenter: deployDeps.Presenter,
	}
}

func (u *deploySiteInteractor) Handle(input *port.DeploySiteUsecaseInputData) {
	output := &port.DeploySiteUsecaseOutputData{}

	output.Message = "deploying to vercel... "
	u.presenter.Progress(output)

	dataDir := filepath.Join(core.Cfg.GetDataDir(), "output")
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "output")

	files, err := u.provider.List(&repository.QueryListParam{
		Folders: []string{dataDir},
	})
	if err != nil {
		u.presenter.Suspend(err)
		return
	}
	for _, file := range files {
		err = u.provider.Copy(filepath.Join(outputDir, file.Name()), filepath.Join(dataDir, file.Name()))
		if err != nil {
			u.presenter.Suspend(err)
			return
		}
	}

	log, err := exec.Command("vercel", "--cwd", outputDir, "--archive=tgz", "--prod").Output()
	if err != nil {
		u.presenter.Suspend(fmt.Errorf("err: %v\ndeploy log:\n%v", err, log))
		return
	}

	output.Message = "ok. "
	u.presenter.Complete(output)
}
