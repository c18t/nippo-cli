package interactor

import (
	"bytes"
	"fmt"
	"time"

	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/model"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

type formatCommandInteractor struct {
	remoteNippoQuery repository.RemoteNippoQuery     `do:""`
	presenter        presenter.FormatCommandPresenter `do:""`
}

func NewFormatCommandInteractor(i do.Injector) (port.FormatCommandUseCase, error) {
	remoteNippoQuery, err := do.Invoke[repository.RemoteNippoQuery](i)
	if err != nil {
		return nil, err
	}
	p, err := do.Invoke[presenter.FormatCommandPresenter](i)
	if err != nil {
		return nil, err
	}
	return &formatCommandInteractor{
		remoteNippoQuery: remoteNippoQuery,
		presenter:        p,
	}, nil
}

func (u *formatCommandInteractor) Handle(input *port.FormatCommandUseCaseInputData) {
	// Show progress while fetching file list
	u.presenter.Progress(&port.FormatCommandUseCaseOutputData{Message: "Fetching file list from Google Drive..."})

	// Fetch files updated since last format timestamp
	nippoList, err := u.fetchFiles()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	u.presenter.StopProgress()

	if len(nippoList) == 0 {
		u.presenter.Complete(&port.FormatCommandUseCaseOutputData{Message: "No files to process."})
		return
	}

	// Track results
	var successCount, noChangeCount, failedCount int
	var updatedFiles, failedFiles []presenter.FileInfo
	hasFailure := false

	// Start format progress TUI
	u.presenter.StartFormatProgress(len(nippoList))

	// Process each file
	for i := range nippoList {
		// Check for cancellation
		if u.presenter.IsFormatCancelled() {
			break
		}

		result := u.processFile(&nippoList[i])
		u.presenter.UpdateFormatProgress(result)

		switch result.Status {
		case port.FormatFileStatusSuccess:
			successCount++
			updatedFiles = append(updatedFiles, presenter.FileInfo{Name: result.Filename, Id: result.FileId})
		case port.FormatFileStatusNoChange:
			noChangeCount++
		case port.FormatFileStatusFailed:
			failedCount++
			failedFiles = append(failedFiles, presenter.FileInfo{Name: result.Filename, Id: result.FileId})
			hasFailure = true
		}
	}

	// Stop format progress TUI
	u.presenter.StopFormatProgress()

	// Show summary
	u.presenter.Summary(successCount, noChangeCount, failedCount, updatedFiles, failedFiles)

	// Only update timestamp if no failures
	if !hasFailure {
		core.Cfg.LastFormatTimestamp = time.Now()
		if err := core.Cfg.SaveConfig(); err != nil {
			u.presenter.Suspend(err)
			return
		}
	}
}

func (u *formatCommandInteractor) fetchFiles() ([]model.Nippo, error) {
	// Use configured drive folder ID
	driveFolderId := core.Cfg.Project.DriveFolderId
	if driveFolderId == "" {
		return nil, fmt.Errorf("drive folder ID is not configured. Run `nippo init` to configure")
	}

	// Build query parameters
	param := &repository.QueryListParam{
		Folders:        []string{driveFolderId},
		FileExtensions: []string{"md"},
		OrderBy:        "name",
	}

	// If we have a last format timestamp, only fetch files modified since then
	if !core.Cfg.LastFormatTimestamp.IsZero() {
		param.UpdatedAt = core.Cfg.LastFormatTimestamp
	}

	// Use RemoteNippoQuery with recursive option to fetch all files including subfolders
	nippoList, err := u.remoteNippoQuery.List(param, &repository.QueryListOption{
		Recursive:   true,
		WithContent: true,
	})
	if err != nil {
		return nil, err
	}

	return nippoList, nil
}

func (u *formatCommandInteractor) processFile(nippo *model.Nippo) *port.FormatCommandUseCaseOutputData {
	result := &port.FormatCommandUseCaseOutputData{
		Filename: nippo.RemoteFile.Name,
		FileId:   nippo.RemoteFile.Id,
	}

	// Check if file needs processing
	needsUpdate, reason, parseErr := u.needsUpdate(nippo.Content)
	if parseErr != nil {
		// Malformed front-matter - log error and skip
		result.Status = port.FormatFileStatusFailed
		result.Error = parseErr
		result.Message = "malformed front-matter: " + parseErr.Error()
		return result
	}
	if !needsUpdate {
		result.Status = port.FormatFileStatusNoChange
		result.Message = "No changes needed"
		return result
	}

	// Process the file based on what's needed
	newContent, err := u.updateContent(nippo)
	if err != nil {
		result.Status = port.FormatFileStatusFailed
		result.Error = err
		result.Message = err.Error()
		return result
	}

	// Check if content actually changed
	if bytes.Equal(nippo.Content, newContent) {
		result.Status = port.FormatFileStatusNoChange
		result.Message = "No changes needed"
		return result
	}

	// Upload to Drive
	if err := u.remoteNippoQuery.Update(nippo, newContent); err != nil {
		result.Status = port.FormatFileStatusFailed
		result.Error = err
		result.Message = err.Error()
		return result
	}

	result.Status = port.FormatFileStatusSuccess
	result.Message = reason
	return result
}

func (u *formatCommandInteractor) needsUpdate(content []byte) (bool, string, error) {
	// Check if file has front-matter
	if !model.HasFrontMatter(content) {
		return true, "Added front-matter", nil
	}

	// Parse front-matter to check for missing created or now placeholder
	fm, _, err := model.ParseFrontMatter(content)
	if err != nil {
		// Malformed front-matter - return error to log it
		return false, "", err
	}

	// Check if created field is missing
	if fm != nil && fm.Created.IsZero() {
		if _, hasCreated := fm.Raw["created"]; !hasCreated {
			return true, "Added created field", nil
		}
	}

	// Check for "now" placeholder
	if model.HasNowPlaceholder(content) {
		return true, "Replaced updated: now", nil
	}

	return false, "", nil
}

func (u *formatCommandInteractor) updateContent(nippo *model.Nippo) ([]byte, error) {
	content := nippo.Content

	// Get timestamps from RemoteFile
	createdTime, _ := time.Parse(time.RFC3339, nippo.RemoteFile.CreatedTime)
	modifiedTime, _ := time.Parse(time.RFC3339, nippo.RemoteFile.ModifiedTime)

	// If no front-matter, generate new one
	if !model.HasFrontMatter(content) {
		created := createdTime.Local()
		newFM := model.GenerateFrontMatter(created)
		return append([]byte(newFM), content...), nil
	}

	// Parse existing front-matter
	fm, _, err := model.ParseFrontMatter(content)
	if err != nil {
		return nil, err
	}

	// Determine what to update
	var created, updated time.Time
	var replaceNow bool

	// Add created if missing
	if fm != nil {
		if _, hasCreated := fm.Raw["created"]; !hasCreated {
			created = createdTime.Local()
		}
	}

	// Replace "now" placeholder with modifiedTime
	if model.HasNowPlaceholder(content) {
		updated = modifiedTime.Local()
		replaceNow = true
	}

	// If nothing to update, return original content
	if created.IsZero() && !replaceNow {
		return content, nil
	}

	return model.UpdateFrontMatter(content, created, updated, replaceNow)
}
