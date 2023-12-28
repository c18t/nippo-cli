package interactor

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type cleanBuildCacheInteractor struct {
	presenter presenter.CleanBuildCachePresenter
}

func NewCleanBuildCacheInteractor(presenter presenter.CleanBuildCachePresenter) port.CleanBuildCacheUsecase {
	return &cleanBuildCacheInteractor{presenter}
}

func (u *cleanBuildCacheInteractor) Handle(input *port.CleanBuildCacheUsecaseInputData) {
	output := &port.CleanBuildCacheUsecaseOutputData{}

	output.Message = "clean cache files... "
	u.presenter.Progress(output)

	clearNippoCache()
	clearBuildCache()

	output.Message = "ok. "
	u.presenter.Complete(output)
}

func clearNippoCache() error {
	outputDir := path.Join(core.Cfg.GetCacheDir(), "md")
	files, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		err = os.Remove(path.Join(outputDir, fileName))
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func clearBuildCache() error {
	outputDir := path.Join(core.Cfg.GetCacheDir(), "output")
	files, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".html") {
			err = os.Remove(path.Join(outputDir, fileName))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
