package controller

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/spf13/cobra"
)

type UpdateController interface {
	core.Controller
}

type updateController struct {
}

func NewUpdateController() UpdateController {
	return &updateController{}
}

func (c *updateController) Exec(cmd *cobra.Command, args []string) error {
	fmt.Print("update project files... ")
	err := downloadProject()
	cobra.CheckErr(err)
	fmt.Println("ok.")
	return nil
}

func downloadProject() error {
	// ダウンロードするURL
	url := "https://codeload.github.com/c18t/nippo/zip/refs/heads/main"

	cacheDir := core.Cfg.GetCacheDir()
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// ダウンロードしたファイルを格納するファイル名
	filename := filepath.Base(url)

	// ダウンロード
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// XDG_CACHE_HOMEディレクトリにファイルを保存
	f, err := os.Create(filepath.Join(cacheDir, filename))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	// 展開するディレクトリを取得
	dataDir := core.Cfg.GetDataDir()
	err = os.MkdirAll(dataDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// ZIPファイルを展開
	err = unzip(filepath.Join(cacheDir, filename), dataDir)
	if err != nil {
		return err
	}

	return nil
}

func unzip(zipFile, destDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}

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
		outFile := filepath.Join(destDir, relPath)
		outDir := filepath.Dir(outFile)
		err = os.MkdirAll(outDir, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}

		// ファイルを書き込む
		dest, err := os.Create(outFile)
		if err != nil {
			return err
		}
		defer dest.Close()

		srcFile, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(dest, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}
