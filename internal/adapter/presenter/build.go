package presenter

import (
	"fmt"
	"reflect"

	"github.com/c18t/nippo-cli/internal/adapter/presenter/view/tui"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// Status icons for build command output
const (
	BuildIconSuccess = "✓"
	BuildIconFailed  = "✗"
)

type BuildCommandPresenter interface {
	Progress(output *port.BuildCommandUseCaseOutputData)
	StopProgress()
	Complete(output *port.BuildCommandUseCaseOutputData)
	Suspend(err error)
	StartBuildProgress(totalFiles int)
	UpdateBuildProgress(filename string, fileId string)
	StopBuildProgress()
	IsBuildCancelled() bool
	Summary(downloadedFiles []FileInfo, failedFiles []FileInfo, buildError error)
}

// FileInfo holds file name and ID for summary display
type FileInfo struct {
	Name string
	Id   string
}

type buildCommandPresenter struct {
	base             ConsolePresenter
	buildProgressCtl *tui.BuildProgressController
}

func NewBuildCommandPresenter(i do.Injector) (BuildCommandPresenter, error) {
	base, err := do.Invoke[ConsolePresenter](i)
	if err != nil {
		return nil, err
	}
	return &buildCommandPresenter{
		base:             base,
		buildProgressCtl: tui.NewBuildProgressController(),
	}, nil
}

func (p *buildCommandPresenter) Progress(output *port.BuildCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *buildCommandPresenter) StopProgress() {
	p.base.StopProgress()
}

func (p *buildCommandPresenter) Complete(output *port.BuildCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *buildCommandPresenter) Suspend(err error) {
	p.StopBuildProgress()
	p.base.Suspend(err)
}

func (p *buildCommandPresenter) StartBuildProgress(totalFiles int) {
	p.buildProgressCtl.Start(totalFiles)
}

func (p *buildCommandPresenter) UpdateBuildProgress(filename string, fileId string) {
	p.buildProgressCtl.UpdateFile(filename, fileId)
}

func (p *buildCommandPresenter) StopBuildProgress() {
	p.buildProgressCtl.Stop()
}

func (p *buildCommandPresenter) IsBuildCancelled() bool {
	return p.buildProgressCtl.IsCancelled()
}

func (p *buildCommandPresenter) Summary(downloadedFiles []FileInfo, failedFiles []FileInfo, buildError error) {
	if len(downloadedFiles) > 0 {
		tui.Println("")
		tui.Println(tui.SuccessStyle.Render("Downloaded files:"))
		for _, f := range downloadedFiles {
			tui.Println(fmt.Sprintf("  %s %s (%s)",
				tui.SuccessStyle.Render(BuildIconSuccess),
				f.Name,
				tui.DimStyle.Render(f.Id),
			))
		}
	}

	if len(failedFiles) > 0 {
		tui.Println("")
		tui.Println(tui.ErrorStyle.Render("Failed files:"))
		for _, f := range failedFiles {
			tui.Println(fmt.Sprintf("  %s %s (%s)",
				tui.ErrorStyle.Render(BuildIconFailed),
				f.Name,
				tui.DimStyle.Render(f.Id),
			))
		}
	}

	tui.Println("")

	// Show build status
	if buildError != nil {
		tui.Println(fmt.Sprintf("Build failed: %s", tui.ErrorStyle.Render(buildError.Error())))
	} else {
		tui.Println(fmt.Sprintf("Build complete: %s downloaded, %s failed",
			tui.SuccessStyle.Render(fmt.Sprintf("%d", len(downloadedFiles))),
			tui.ErrorStyle.Render(fmt.Sprintf("%d", len(failedFiles))),
		))
	}
}
