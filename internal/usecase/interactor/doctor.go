package interactor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

type doctorInteractor struct {
	presenter presenter.DoctorPresenter `do:""`
}

func NewDoctorInteractor(i do.Injector) (port.DoctorUseCase, error) {
	p, err := do.Invoke[presenter.DoctorPresenter](i)
	if err != nil {
		return nil, err
	}
	return &doctorInteractor{
		presenter: p,
	}, nil
}

func (u *doctorInteractor) Handle(input *port.DoctorUseCaseInputData) {
	output := &port.DoctorUseCaseOutputData{
		Checks: []port.DoctorCheck{},
	}

	// Check configuration
	u.checkConfiguration(output)

	// Check directories
	u.checkDirectories(output)

	// Check required files
	u.checkRequiredFiles(output)

	// Check cache status
	u.checkCacheStatus(output)

	// Security check
	u.checkSecurity(output)

	u.presenter.Show(output)
}

func (u *doctorInteractor) checkConfiguration(output *port.DoctorUseCaseOutputData) {
	configPath := core.GetConfigFilePath()

	// Check config file
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category:   "Configuration",
			Item:       "Config file",
			Status:     port.DoctorCheckStatusFail,
			Message:    "Not found: " + configPath,
			Suggestion: "Run `nippo init` to create configuration",
		})
		return
	}

	output.Checks = append(output.Checks, port.DoctorCheck{
		Category: "Configuration",
		Item:     "Config file",
		Status:   port.DoctorCheckStatusPass,
		Message:  configPath,
	})

	// Check drive folder
	if core.Cfg.Project.DriveFolderId == "" {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category:   "Configuration",
			Item:       "Drive folder",
			Status:     port.DoctorCheckStatusWarn,
			Message:    "Not configured",
			Suggestion: "Run `nippo init` to configure",
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Configuration",
			Item:     "Drive folder",
			Status:   port.DoctorCheckStatusPass,
			Message:  core.Cfg.Project.DriveFolderId,
		})
	}

	// Check site URL
	if core.Cfg.Project.SiteUrl == "" {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category:   "Configuration",
			Item:       "Site URL",
			Status:     port.DoctorCheckStatusWarn,
			Message:    "Not configured (using default)",
			Suggestion: "Run `nippo init` to configure",
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Configuration",
			Item:     "Site URL",
			Status:   port.DoctorCheckStatusPass,
			Message:  core.Cfg.Project.SiteUrl,
		})
	}

	// Check project URL
	if core.Cfg.Project.Url == "" {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category:   "Configuration",
			Item:       "Project URL",
			Status:     port.DoctorCheckStatusWarn,
			Message:    "Not configured (using default)",
			Suggestion: "Run `nippo init` to configure",
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Configuration",
			Item:     "Project URL",
			Status:   port.DoctorCheckStatusPass,
			Message:  core.Cfg.Project.Url,
		})
	}

	// Check branch
	branch := core.Cfg.Project.Branch
	if branch == "" {
		branch = "main (default)"
	}
	output.Checks = append(output.Checks, port.DoctorCheck{
		Category: "Configuration",
		Item:     "Project branch",
		Status:   port.DoctorCheckStatusPass,
		Message:  branch,
	})

	// Check template path
	templatePath := core.Cfg.Project.TemplatePath
	if templatePath == "" {
		templatePath = "/templates (default)"
	}
	output.Checks = append(output.Checks, port.DoctorCheck{
		Category: "Configuration",
		Item:     "Project template path",
		Status:   port.DoctorCheckStatusPass,
		Message:  templatePath,
	})

	// Check asset path
	assetPath := core.Cfg.Project.AssetPath
	if assetPath == "" {
		assetPath = "/dist (default)"
	}
	output.Checks = append(output.Checks, port.DoctorCheck{
		Category: "Configuration",
		Item:     "Project asset path",
		Status:   port.DoctorCheckStatusPass,
		Message:  assetPath,
	})
}

func (u *doctorInteractor) checkDirectories(output *port.DoctorUseCaseOutputData) {
	// Use slice to maintain stable order
	dirs := []struct {
		name string
		path string
	}{
		{"Config directory", core.Cfg.GetConfigDir()},
		{"Data directory", core.Cfg.GetDataDir()},
		{"Cache directory", core.Cfg.GetCacheDir()},
	}

	for _, dir := range dirs {
		if info, err := os.Stat(dir.path); err == nil && info.IsDir() {
			output.Checks = append(output.Checks, port.DoctorCheck{
				Category: "Paths",
				Item:     dir.name,
				Status:   port.DoctorCheckStatusPass,
				Message:  dir.path,
			})
		} else {
			output.Checks = append(output.Checks, port.DoctorCheck{
				Category:   "Paths",
				Item:       dir.name,
				Status:     port.DoctorCheckStatusFail,
				Message:    "Not found: " + dir.path,
				Suggestion: "mkdir -p " + dir.path,
			})
		}
	}
}

func (u *doctorInteractor) checkRequiredFiles(output *port.DoctorUseCaseOutputData) {
	dataDir := core.Cfg.GetDataDir()

	// Check credentials.json
	credPath := filepath.Join(dataDir, "credentials.json")
	if _, err := os.Stat(credPath); os.IsNotExist(err) {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category:   "Required Files",
			Item:       "credentials.json",
			Status:     port.DoctorCheckStatusFail,
			Message:    "Not found: " + credPath,
			Suggestion: "Download from Google Cloud Console: https://console.cloud.google.com/apis/credentials",
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Required Files",
			Item:     "credentials.json",
			Status:   port.DoctorCheckStatusPass,
			Message:  credPath,
		})
	}

	// Check token.json
	tokenPath := filepath.Join(dataDir, "token.json")
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category:   "Required Files",
			Item:       "token.json",
			Status:     port.DoctorCheckStatusFail,
			Message:    "Not found: " + tokenPath,
			Suggestion: "Run `nippo auth` to authenticate with Google Drive",
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Required Files",
			Item:     "token.json",
			Status:   port.DoctorCheckStatusPass,
			Message:  tokenPath,
		})
	}

	// Check templates directory
	templatesDir := filepath.Join(dataDir, "templates")
	if info, err := os.Stat(templatesDir); err == nil && info.IsDir() {
		files, _ := os.ReadDir(templatesDir)
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Required Files",
			Item:     "templates/",
			Status:   port.DoctorCheckStatusPass,
			Message:  templatesDir + " (" + itoa(len(files)) + " files)",
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category:   "Required Files",
			Item:       "templates/",
			Status:     port.DoctorCheckStatusFail,
			Message:    "Not found: " + templatesDir,
			Suggestion: "Run `nippo init` or `nippo update` to download templates",
		})
	}

	// Check assets directory
	assetsDir := filepath.Join(dataDir, "assets")
	if info, err := os.Stat(assetsDir); err == nil && info.IsDir() {
		files, _ := os.ReadDir(assetsDir)
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Required Files",
			Item:     "assets/",
			Status:   port.DoctorCheckStatusPass,
			Message:  assetsDir + " (" + itoa(len(files)) + " files)",
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category:   "Required Files",
			Item:       "assets/",
			Status:     port.DoctorCheckStatusWarn,
			Message:    "Not found: " + assetsDir,
			Suggestion: "Run `nippo init` or `nippo update` to download assets",
		})
	}
}

func (u *doctorInteractor) checkCacheStatus(output *port.DoctorUseCaseOutputData) {
	cacheDir := core.Cfg.GetCacheDir()

	// Check nippo-template.zip
	zipPath := filepath.Join(cacheDir, "nippo-template.zip")
	if info, err := os.Stat(zipPath); err == nil {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Cache Status",
			Item:     "nippo-template.zip",
			Status:   port.DoctorCheckStatusPass,
			Message:  info.ModTime().Format("2006-01-02 15:04:05"),
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Cache Status",
			Item:     "nippo-template.zip",
			Status:   port.DoctorCheckStatusWarn,
			Message:  "Not found (will be downloaded on init/update)",
		})
	}

	// Check md/ directory
	mdDir := filepath.Join(cacheDir, "md")
	if info, err := os.Stat(mdDir); err == nil && info.IsDir() {
		files, _ := os.ReadDir(mdDir)
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Cache Status",
			Item:     "md/",
			Status:   port.DoctorCheckStatusPass,
			Message:  itoa(len(files)) + " files",
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Cache Status",
			Item:     "md/",
			Status:   port.DoctorCheckStatusWarn,
			Message:  "Not found (will be created on build)",
		})
	}

	// Check output/ directory
	outputDir := filepath.Join(cacheDir, "output")
	if info, err := os.Stat(outputDir); err == nil && info.IsDir() {
		files, _ := os.ReadDir(outputDir)
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Cache Status",
			Item:     "output/",
			Status:   port.DoctorCheckStatusPass,
			Message:  itoa(len(files)) + " files",
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Cache Status",
			Item:     "output/",
			Status:   port.DoctorCheckStatusWarn,
			Message:  "Not found (will be created on build)",
		})
	}
}

func (u *doctorInteractor) checkSecurity(output *port.DoctorUseCaseOutputData) {
	dataDir := core.Cfg.GetDataDir()

	if core.IsUnderGitRepo(dataDir) {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category:   "Security",
			Item:       "Git repository",
			Status:     port.DoctorCheckStatusWarn,
			Message:    "Data directory is under git repository",
			Suggestion: "Credentials may be tracked. Add to .gitignore or move data_dir",
		})
	} else {
		output.Checks = append(output.Checks, port.DoctorCheck{
			Category: "Security",
			Item:     "Git repository",
			Status:   port.DoctorCheckStatusPass,
			Message:  "Data directory is not under git repository",
		})
	}
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
