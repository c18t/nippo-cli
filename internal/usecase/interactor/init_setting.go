package interactor

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

type initSettingInteractor struct {
	provider  gateway.LocalFileProvider
	presenter presenter.InitSettingPresenter
}

func NewInitSettingInteractor(i do.Injector) (port.InitSettingUseCase, error) {
	provider, err := do.Invoke[gateway.LocalFileProvider](i)
	if err != nil {
		return nil, err
	}
	p, err := do.Invoke[presenter.InitSettingPresenter](i)
	if err != nil {
		return nil, err
	}
	return &initSettingInteractor{
		provider:  provider,
		presenter: p,
	}, nil
}

func (u *initSettingInteractor) Handle(input *port.InitSettingUseCaseInputData) {
	var err error
	output := &port.InitSettingUseCaseOutputData{
		Project: port.InitSettingProject{},
	}

	err = u.configureProject(input, output)
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	err = u.downloadProject()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}
}

func (u *initSettingInteractor) configureProject(input *port.InitSettingUseCaseInputData, output *port.InitSettingUseCaseOutputData) error {
	// プロジェクトURLの入力
	output.Input = port.InitSettingProjectUrl("")
	projectUrl := make(chan interface{})
	go u.presenter.Prompt(projectUrl, output)
	switch ret := (<-projectUrl).(type) {
	case error:
		return ret
	case string:
		parsed, err := url.Parse(ret)
		if err != nil {
			return err
		}
		if parsed.Host == "github.com" {
			ret = fmt.Sprintf("https://codeload.github.com/%s/zip/refs/heads/main", strings.Trim(parsed.Path, "/"))
		}
		output.Project.Url = port.InitSettingProjectUrl(ret)
	}

	// テンプレートパスの入力
	output.Input = port.InitSettingProjectTemplatePath("")
	projectTemplatePath := make(chan interface{})
	go u.presenter.Prompt(projectTemplatePath, output)
	switch ret := (<-projectTemplatePath).(type) {
	case error:
		return ret
	case string:
		output.Project.TemplatePath = port.InitSettingProjectTemplatePath(ret)
	}

	// アセットパスの入力
	output.Input = port.InitSettingProjectAssetPath("")
	projectAssetPath := make(chan interface{})
	go u.presenter.Prompt(projectAssetPath, output)
	switch ret := (<-projectAssetPath).(type) {
	case error:
		return ret
	case string:
		output.Project.AssetPath = port.InitSettingProjectAssetPath(ret)
	}

	// 設定の更新
	output.Message = "Saving project config..."
	u.presenter.Progress(output)
	core.Cfg.Project.Url = string(output.Project.Url)
	core.Cfg.Project.TemplatePath = string(output.Project.TemplatePath)
	core.Cfg.Project.AssetPath = string(output.Project.AssetPath)
	if err := core.Cfg.SaveConfig(); err != nil {
		return err
	}
	// Progress() で開始したスピナーは次の Progress() または Complete() で自動的に "ok." が付く
	output.Message = ""

	return nil
}

func (u *initSettingInteractor) downloadProject() error {
	// ダウンロードするURL
	url := core.Cfg.Project.Url

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
func (u *initSettingInteractor) unzip(zipFile, destDir string) error {
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
