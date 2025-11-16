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
	"github.com/samber/do/v2"
)

type updateCommandInteractor struct {
	provider  gateway.LocalFileProvider        `do:""`
	presenter presenter.UpdateCommandPresenter `do:""`
}

func NewUpdateCommandInteractor(i do.Injector) (port.UpdateCommandUseCase, error) {
	provider, err := do.Invoke[gateway.LocalFileProvider](i)
	if err != nil {
		return nil, err
	}
	p, err := do.Invoke[presenter.UpdateCommandPresenter](i)
	if err != nil {
		return nil, err
	}
	return &updateCommandInteractor{
		provider:  provider,
		presenter: p,
	}, nil
}

func (u *updateCommandInteractor) Handle(input *port.UpdateCommandUseCaseInputData) {
	output := &port.UpdateCommandUseCaseOutputData{}
	output.Message = "updating project files... "
	u.presenter.Progress(output)

	err := u.downloadProject()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	output.Message = "ok."
	u.presenter.Complete(output)
}

func (u *updateCommandInteractor) downloadProject() error {
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
func (u *updateCommandInteractor) unzip(zipFile, destDir string) error {
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
