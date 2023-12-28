package interactor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type deploySiteInteractor struct {
	presenter presenter.DeploySitePresenter
}

func NewDeploySiteInteractor(presenter presenter.DeploySitePresenter) port.DeploySiteUsecase {
	return &deploySiteInteractor{presenter}
}

func (u *deploySiteInteractor) Handle(input *port.DeploySiteUsecaseInputData) {
	output := &port.DeploySiteUsecaseOutputData{}

	output.Message = "deploy to vercel... "
	u.presenter.Progress(output)

	dataDir := path.Join(core.Cfg.GetDataDir(), "output")
	outputDir := path.Join(core.Cfg.GetCacheDir(), "output")

	files, err := os.ReadDir(dataDir)
	if err != nil {
		u.presenter.Suspend(err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		src, err := os.Open(path.Join(dataDir, file.Name()))
		if err != nil {
			u.presenter.Suspend(err)
			return
		}
		defer src.Close()
		dest, err := os.Create(path.Join(outputDir, file.Name()))
		if err != nil {
			u.presenter.Suspend(err)
			return
		}
		defer dest.Close()
		_, err = io.Copy(dest, src)
		if err != nil {
			u.presenter.Suspend(err)
			return
		}
	}

	log, err := exec.Command("vercel", "--cwd", outputDir, "--prod").Output()
	if err != nil {
		u.presenter.Suspend(fmt.Errorf("err: %v\ndeploy log:\n%v", err, log))
		return
	}

	output.Message = "ok. "
	u.presenter.Complete(output)
}
