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
	"github.com/samber/do/v2"
)

type deployCommandInteractor struct {
	provider  gateway.LocalFileProvider        `do:""`
	presenter presenter.DeployCommandPresenter `do:""`
}

func NewDeployCommandInteractor(i do.Injector) (port.DeployCommandUseCase, error) {
	provider, err := do.Invoke[gateway.LocalFileProvider](i)
	if err != nil {
		return nil, err
	}
	p, err := do.Invoke[presenter.DeployCommandPresenter](i)
	if err != nil {
		return nil, err
	}
	return &deployCommandInteractor{
		provider:  provider,
		presenter: p,
	}, nil
}

func (u *deployCommandInteractor) Handle(input *port.DeployCommandUseCaseInputData) {
	output := &port.DeployCommandUseCaseOutputData{}
	output.Message = "deploying to vercel..."
	u.presenter.Progress(output)

	// Static assets from template repo (mapped from project.asset_path)
	assetsDir := filepath.Join(core.Cfg.GetDataDir(), "assets")
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "output")

	files, err := u.provider.List(&repository.QueryListParam{
		Folders: []string{assetsDir},
	})
	if err != nil {
		u.presenter.Suspend(err)
		return
	}
	for _, file := range files {
		err = u.provider.Copy(assetsDir, filepath.Join(outputDir, file.Name()), filepath.Join(assetsDir, file.Name()))
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

	// Progress() で開始したスピナーは自動的に "ok." が付く
	u.presenter.StopProgress()
}
