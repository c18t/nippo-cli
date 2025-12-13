package interactor

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
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

	// Check if config file already exists
	configPath := core.GetConfigFilePath()
	configExists := false
	if _, err := os.Stat(configPath); err == nil {
		configExists = true
	}

	// If config exists, ask for confirmation
	if configExists {
		output.Input = port.InitSettingConfirmOverwrite(false)
		confirmCh := make(chan interface{})
		go u.presenter.Prompt(confirmCh, output)
		switch ret := (<-confirmCh).(type) {
		case error:
			u.presenter.Suspend(ret)
			return
		case bool:
			if !ret {
				u.presenter.Suspend(fmt.Errorf("operation cancelled"))
				return
			}
		}
	}

	err = u.configureProject(input, output, configExists)
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	// Check if data_dir is under git repository (proactive security check)
	dataDir := core.Cfg.GetDataDir()
	if core.IsUnderGitRepo(dataDir) {
		output.Input = port.InitSettingConfirmGitWarning(false)
		gitConfirmCh := make(chan interface{})
		go u.presenter.Prompt(gitConfirmCh, output)
		switch ret := (<-gitConfirmCh).(type) {
		case error:
			u.presenter.Suspend(ret)
			return
		case bool:
			if !ret {
				u.presenter.Suspend(fmt.Errorf("operation cancelled"))
				return
			}
		}
	}

	// Create required directories
	output.Message = "Creating directories..."
	u.presenter.Progress(output)
	err = u.createDirectories()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}
	u.presenter.StopProgress()

	// Download and extract templates
	output.Message = "Downloading templates..."
	u.presenter.Progress(output)
	err = u.downloadProject()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}
	u.presenter.StopProgress()

	// Show completion message
	output.Message = "Initialization complete. Run `nippo auth` to authenticate with Google Drive."
	u.presenter.Complete(output)
}

func (u *initSettingInteractor) configureProject(input *port.InitSettingUseCaseInputData, output *port.InitSettingUseCaseOutputData, configExists bool) error {
	// Load existing config values for defaults if config exists
	var existingDriveFolder, existingSiteUrl, existingUrl, existingBranch, existingTemplatePath, existingAssetPath string
	if configExists {
		existingDriveFolder = core.Cfg.Project.DriveFolderId
		existingSiteUrl = core.Cfg.Project.SiteUrl
		existingUrl = core.Cfg.Project.Url
		existingBranch = core.Cfg.Project.Branch
		existingTemplatePath = core.Cfg.Project.TemplatePath
		existingAssetPath = core.Cfg.Project.AssetPath
	}

	// Drive folder input
	output.Input = port.InitSettingProjectDriveFolder("")
	driveFolderCh := make(chan interface{})
	go u.presenter.Prompt(driveFolderCh, output)
	switch ret := (<-driveFolderCh).(type) {
	case error:
		return ret
	case string:
		folderId := extractDriveFolderId(ret)
		if folderId == "" {
			folderId = ret // Use as-is if not a URL
		}
		if folderId == "" && existingDriveFolder != "" {
			folderId = existingDriveFolder
		}
		output.Project.DriveFolder = port.InitSettingProjectDriveFolder(folderId)
	}

	// Site URL input
	output.Input = port.InitSettingProjectSiteUrl("")
	siteUrlCh := make(chan interface{})
	go u.presenter.Prompt(siteUrlCh, output)
	switch ret := (<-siteUrlCh).(type) {
	case error:
		return ret
	case string:
		if ret == "" && existingSiteUrl != "" {
			ret = existingSiteUrl
		}
		output.Project.SiteUrl = port.InitSettingProjectSiteUrl(ret)
	}

	// Project URL input
	output.Input = port.InitSettingProjectUrl("")
	projectUrlCh := make(chan interface{})
	go u.presenter.Prompt(projectUrlCh, output)
	switch ret := (<-projectUrlCh).(type) {
	case error:
		return ret
	case string:
		if ret == "" && existingUrl != "" {
			ret = existingUrl
		}
		parsed, err := url.Parse(ret)
		if err != nil {
			return err
		}
		if parsed.Host == "github.com" {
			// Will construct download URL later using branch
			output.Project.Url = port.InitSettingProjectUrl(ret)
		} else {
			output.Project.Url = port.InitSettingProjectUrl(ret)
		}
	}

	// Branch input
	output.Input = port.InitSettingProjectBranch("")
	branchCh := make(chan interface{})
	go u.presenter.Prompt(branchCh, output)
	switch ret := (<-branchCh).(type) {
	case error:
		return ret
	case string:
		if ret == "" {
			if existingBranch != "" {
				ret = existingBranch
			} else {
				ret = "main"
			}
		}
		output.Project.Branch = port.InitSettingProjectBranch(ret)
	}

	// Template path input
	output.Input = port.InitSettingProjectTemplatePath("")
	templatePathCh := make(chan interface{})
	go u.presenter.Prompt(templatePathCh, output)
	switch ret := (<-templatePathCh).(type) {
	case error:
		return ret
	case string:
		if ret == "" {
			if existingTemplatePath != "" {
				ret = existingTemplatePath
			} else {
				ret = "/templates"
			}
		}
		output.Project.TemplatePath = port.InitSettingProjectTemplatePath(ret)
	}

	// Asset path input
	output.Input = port.InitSettingProjectAssetPath("")
	assetPathCh := make(chan interface{})
	go u.presenter.Prompt(assetPathCh, output)
	switch ret := (<-assetPathCh).(type) {
	case error:
		return ret
	case string:
		if ret == "" {
			if existingAssetPath != "" {
				ret = existingAssetPath
			} else {
				ret = "/dist"
			}
		}
		output.Project.AssetPath = port.InitSettingProjectAssetPath(ret)
	}

	// Save configuration
	output.Message = "Saving project config..."
	u.presenter.Progress(output)

	core.Cfg.Project.DriveFolderId = string(output.Project.DriveFolder)
	core.Cfg.Project.SiteUrl = string(output.Project.SiteUrl)
	core.Cfg.Project.Url = string(output.Project.Url)
	core.Cfg.Project.Branch = string(output.Project.Branch)
	core.Cfg.Project.TemplatePath = string(output.Project.TemplatePath)
	core.Cfg.Project.AssetPath = string(output.Project.AssetPath)

	if err := core.Cfg.SaveConfig(); err != nil {
		return err
	}

	// Stop progress to show "ok."
	u.presenter.StopProgress()

	return nil
}

func (u *initSettingInteractor) createDirectories() error {
	configDir := core.Cfg.GetConfigDir()
	dataDir := core.Cfg.GetDataDir()
	cacheDir := core.Cfg.GetCacheDir()

	dirs := []string{configDir, dataDir, cacheDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func (u *initSettingInteractor) downloadProject() error {
	// Construct download URL
	projectUrl := core.Cfg.Project.Url
	branch := core.Cfg.Project.Branch
	if branch == "" {
		branch = "main"
	}

	parsed, err := url.Parse(projectUrl)
	if err != nil {
		return err
	}

	var downloadUrl string
	if parsed.Host == "github.com" {
		downloadUrl = fmt.Sprintf("https://codeload.github.com/%s/zip/refs/heads/%s",
			strings.Trim(parsed.Path, "/"), branch)
	} else {
		downloadUrl = projectUrl
	}

	// Get directories
	cacheDir := core.Cfg.GetCacheDir()
	dataDir := core.Cfg.GetDataDir()

	// Download ZIP file
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
func (u *initSettingInteractor) unzipSelective(zipFile, destDir, branch string) error {
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
	parsed, _ := url.Parse(core.Cfg.Project.Url)
	pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	repoName := "nippo"
	if len(pathParts) >= 2 {
		repoName = pathParts[1]
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

func (u *initSettingInteractor) extractFile(f *zip.File, targetPath string) error {
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

// extractDriveFolderId extracts folder ID from various Google Drive URL formats
func extractDriveFolderId(input string) string {
	// Pattern for folder ID in URL
	patterns := []string{
		`drive\.google\.com/drive/(?:u/\d+/)?folders/([a-zA-Z0-9_-]+)`,
		`drive\.google\.com/open\?id=([a-zA-Z0-9_-]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(input)
		if len(matches) >= 2 {
			return matches[1]
		}
	}

	return ""
}
