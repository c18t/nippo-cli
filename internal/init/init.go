package init

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

func CreateCmdFunc() RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		var err error

		err = downloadProject()
		cobra.CheckErr(err)

		err = saveDriveToken()
		cobra.CheckErr(err)

		return nil
	}
}

// プロジェクトのダウンロード
func downloadProject() error {

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

func saveDriveToken() error {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	defaultDataDir := path.Join(home, ".local", "share")
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" || !path.IsAbs(dataDir) {
		dataDir = defaultDataDir
	}
	dataDir = path.Join(dataDir, "nippo")

	b, err := os.ReadFile(path.Join(dataDir, "credentials.json"))
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v\n", err)
		return nil
	}

	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v\n", err)
		return nil
	}

	tok := getTokenFromWeb(config)
	saveToken(path.Join(dataDir, "token.json"), tok)
	return nil
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Unable to cache oauth token: %v\n", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Printf("Unable to read authorization code %v\n", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		fmt.Printf("Unable to retrieve token from web %v\n", err)
	}
	return tok
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
