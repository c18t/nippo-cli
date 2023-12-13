package init

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

func CreateCmdFunc() RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		// ダウンロードするURL
		url := "https://codeload.github.com/c18t/nippo/zip/refs/heads/main"

		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// XDG_CACHE_HOMEディレクトリを取得
		defaultCacheDir := path.Join(home, ".cache")
		cacheDir := os.Getenv("XDG_CACHE_HOME")
		if cacheDir == "" || !path.IsAbs(cacheDir) {
			cacheDir = defaultCacheDir
		}
		cacheDir = path.Join(cacheDir, "nippo")
		err = os.MkdirAll(cacheDir, 0755)
		if err != nil && !os.IsExist(err) {
			fmt.Println(err)
			return nil
		}

		// ダウンロードしたファイルを格納するファイル名
		filename := filepath.Base(url)

		// ダウンロード
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		defer resp.Body.Close()

		// XDG_CACHE_HOMEディレクトリにファイルを保存
		f, err := os.Create(filepath.Join(cacheDir, filename))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer f.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		// 展開するディレクトリを取得
		defaultDataDir := path.Join(home, ".local", "share")
		dataDir := os.Getenv("XDG_DATA_HOME")
		if dataDir == "" || !path.IsAbs(dataDir) {
			dataDir = defaultDataDir
		}
		dataDir = path.Join(dataDir, "nippo")
		err = os.MkdirAll(dataDir, 0755)
		if err != nil && !os.IsExist(err) {
			fmt.Println(err)
			return nil
		}

		// ZIPファイルを展開
		err = unzip(filepath.Join(cacheDir, filename), dataDir)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		fmt.Println("ダウンロードと展開が完了しました。")
		return nil
	}
}

// ZIPファイルを展開する関数
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
		outFile := filepath.Join(destDir, filepath.Base(f.Name))

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
