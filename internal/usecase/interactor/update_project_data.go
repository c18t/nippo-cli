package interactor

import (
	"archive/zip"
	"io"
	"net/http"
	"path/filepath"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"go.uber.org/dig"
)

type updateProjectDataInteractor struct {
	provider  gateway.LocalFileProvider
	presenter presenter.UpdateProjectDataPresenter
}

type inUpdateProjectDataInteractor struct {
	dig.In
	Provider  gateway.LocalFileProvider
	Presenter presenter.UpdateProjectDataPresenter
}

func NewUpdateProjectDataInteractor(updateDeps inUpdateProjectDataInteractor) port.UpdateProjectDataUsecase {
	return &updateProjectDataInteractor{
		provider:  updateDeps.Provider,
		presenter: updateDeps.Presenter,
	}
}

func (u *updateProjectDataInteractor) Handle(input *port.UpdateProjectDataUsecaseInputData) {
	output := &port.UpdateProjectDataUsecaseOutputData{}

	output.Message = "update project files... "
	u.presenter.Progress(output)

	err := u.downloadProject()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	output.Message = "ok."
	u.presenter.Complete(output)
}

func (u *updateProjectDataInteractor) downloadProject() error {
	// ダウンロードするURL
	url := "https://codeload.github.com/c18t/nippo/zip/refs/heads/main"

	// 展開するディレクトリを取得
	cacheDir := core.Cfg.GetCacheDir()
	dataDir := core.Cfg.GetDataDir()

	// ダウンロードしたファイルを格納するファイル名
	filename := filepath.Base(url)

	// ダウンロード
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { io.Copy(io.Discard, resp.Body); resp.Body.Close() }()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zipFilePath := filepath.Join(cacheDir, filename)
	u.provider.Write(zipFilePath, content)

	// ZIPファイルを展開
	err = u.unzip(zipFilePath, dataDir)
	if err != nil {
		return err
	}

	return nil
}

// ZIPファイルを展開する関数
func (u *updateProjectDataInteractor) unzip(zipFile, destDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// ディレクトリの場合はスキップ
		if f.FileInfo().IsDir() {
			continue
		}

		// 出力先ファイル名を生成
		relPath, err := filepath.Rel("nippo-main", f.Name)
		if err != nil {
			return err
		}

		// zip内ファイルを開く
		err = func() error {
			src, err := f.Open()
			if err != nil {
				return err
			}
			defer src.Close()
			content, err := io.ReadAll(src)
			if err != nil {
				return err
			}

			u.provider.Write(filepath.Join(destDir, relPath), content)
			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
