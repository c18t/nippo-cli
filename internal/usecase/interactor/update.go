package interactor

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
	output.Message = "updating project files..."
	u.presenter.Progress(output)

	err := u.downloadProject()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	// Progress() で開始したスピナーは自動的に "ok." が付く
	u.presenter.StopProgress()
}

func (u *updateCommandInteractor) downloadProject() error {
	// Construct download URL from config
	projectUrl := core.Cfg.Project.Url
	branch := core.Cfg.Project.Branch
	if branch == "" {
		branch = "main"
	}

	var downloadUrl string
	if projectUrl != "" {
		parsed, err := url.Parse(projectUrl)
		if err != nil {
			return err
		}
		if parsed.Host == "github.com" {
			downloadUrl = fmt.Sprintf("https://codeload.github.com/%s/zip/refs/heads/%s",
				strings.Trim(parsed.Path, "/"), branch)
		} else {
			downloadUrl = projectUrl
		}
	} else {
		// Default fallback
		downloadUrl = fmt.Sprintf("https://codeload.github.com/c18t/nippo/zip/refs/heads/%s", branch)
	}

	// 展開するディレクトリを取得
	cacheDir := core.Cfg.GetCacheDir()
	dataDir := core.Cfg.GetDataDir()

	// ダウンロード
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer func() { io.Copy(io.Discard, resp.Body); resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: HTTP %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Save with standard filename
	zipFilePath := filepath.Join(cacheDir, "nippo-template.zip")
	if err := u.provider.Write(zipFilePath, content); err != nil {
		return err
	}

	// Extract ZIP with selective extraction
	err = u.unzipSelective(zipFilePath, dataDir, branch)
	if err != nil {
		return err
	}

	return nil
}

// unzipSelective extracts only template_path and asset_path from ZIP
func (u *updateCommandInteractor) unzipSelective(zipFile, destDir, branch string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	templatePath := core.Cfg.Project.TemplatePath
	assetPath := core.Cfg.Project.AssetPath
	if templatePath == "" {
		templatePath = "/templates"
	}
	if assetPath == "" {
		assetPath = "/dist"
	}

	// Remove leading slash for path matching
	templatePath = strings.TrimPrefix(templatePath, "/")
	assetPath = strings.TrimPrefix(assetPath, "/")

	// Extract repo name from project URL for ZIP prefix
	repoName := "nippo"
	if core.Cfg.Project.Url != "" {
		parsed, _ := url.Parse(core.Cfg.Project.Url)
		if parsed != nil {
			pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
			if len(pathParts) >= 2 {
				repoName = pathParts[1]
			}
		}
	}
	zipPrefix := fmt.Sprintf("%s-%s/", repoName, branch)

	templateFound := false
	templatesDir := filepath.Join(destDir, "templates")
	assetsDir := filepath.Join(destDir, "assets")

	// Create output directories
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		// Remove ZIP prefix (e.g., "nippo-main/")
		name := strings.TrimPrefix(f.Name, zipPrefix)

		var targetDir string
		var relPath string

		// Check if file is under template_path
		if strings.HasPrefix(name, templatePath+"/") {
			relPath = strings.TrimPrefix(name, templatePath+"/")
			targetDir = templatesDir
			templateFound = true
		} else if strings.HasPrefix(name, assetPath+"/") {
			// Check if file is under asset_path
			relPath = strings.TrimPrefix(name, assetPath+"/")
			targetDir = assetsDir
		} else {
			// Skip files not in template_path or asset_path
			continue
		}

		// Construct target path
		targetPath := filepath.Join(targetDir, relPath)

		// Zip Slip prevention: validate path is within destination
		if !core.IsPathSafe(targetDir, targetPath) {
			continue // Skip potentially malicious paths
		}

		// Extract file
		err = u.extractFile(f, targetPath)
		if err != nil {
			return err
		}
	}

	// Error if templates not found (required)
	if !templateFound {
		return fmt.Errorf("template path '%s' not found in ZIP archive", templatePath)
	}

	return nil
}

func (u *updateCommandInteractor) extractFile(f *zip.File, targetPath string) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	content, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	return u.provider.Write(targetPath, content)
}
