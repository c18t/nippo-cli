package presenter

import (
	"fmt"
	"reflect"

	"github.com/c18t/nippo-cli/internal/adapter/presenter/view/tui"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// Status icons for format command output
const (
	IconSuccess  = "✓"
	IconFailed   = "✗"
	IconNoChange = "○"
)

type FormatCommandPresenter interface {
	Progress(output *port.FormatCommandUseCaseOutputData)
	StopProgress()
	Complete(output *port.FormatCommandUseCaseOutputData)
	Suspend(err error)
	StartFormatProgress(totalFiles int)
	UpdateFormatProgress(output *port.FormatCommandUseCaseOutputData)
	StopFormatProgress()
	IsFormatCancelled() bool
	Summary(successCount, noChangeCount, failedCount int, updatedFiles, failedFiles []FileInfo)
}

type formatCommandPresenter struct {
	base              ConsolePresenter
	formatProgressCtl *tui.FormatProgressController
}

func NewFormatCommandPresenter(i do.Injector) (FormatCommandPresenter, error) {
	base, err := do.Invoke[ConsolePresenter](i)
	if err != nil {
		return nil, err
	}
	return &formatCommandPresenter{
		base:              base,
		formatProgressCtl: tui.NewFormatProgressController(),
	}, nil
}

func (p *formatCommandPresenter) Progress(output *port.FormatCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *formatCommandPresenter) StopProgress() {
	p.base.StopProgress()
}

func (p *formatCommandPresenter) Complete(output *port.FormatCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *formatCommandPresenter) Suspend(err error) {
	p.StopFormatProgress()
	p.base.Suspend(err)
}

func (p *formatCommandPresenter) StartFormatProgress(totalFiles int) {
	p.formatProgressCtl.Start(totalFiles)
}

func (p *formatCommandPresenter) UpdateFormatProgress(output *port.FormatCommandUseCaseOutputData) {
	status := convertStatus(output.Status)
	message := ""
	if output.Error != nil {
		message = output.Error.Error()
	}
	p.formatProgressCtl.UpdateFile(output.Filename, output.FileId, status, message)
}

func (p *formatCommandPresenter) StopFormatProgress() {
	p.formatProgressCtl.Stop()
}

func (p *formatCommandPresenter) IsFormatCancelled() bool {
	return p.formatProgressCtl.IsCancelled()
}

func (p *formatCommandPresenter) Summary(successCount, noChangeCount, failedCount int, updatedFiles, failedFiles []FileInfo) {
	if len(updatedFiles) > 0 {
		tui.Println("")
		tui.Println(tui.SuccessStyle.Render("Updated files:"))
		for _, f := range updatedFiles {
			tui.Println(fmt.Sprintf("  %s %s (%s)",
				tui.SuccessStyle.Render(IconSuccess),
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
				tui.ErrorStyle.Render(IconFailed),
				f.Name,
				tui.DimStyle.Render(f.Id),
			))
		}
	}

	tui.Println("")
	tui.Println(fmt.Sprintf("Format complete: %s updated, %s unchanged, %s failed",
		tui.SuccessStyle.Render(fmt.Sprintf("%d", successCount)),
		tui.DimStyle.Render(fmt.Sprintf("%d", noChangeCount)),
		tui.ErrorStyle.Render(fmt.Sprintf("%d", failedCount)),
	))
}

// convertStatus converts port.FormatFileStatus to tui.FormatFileStatus
func convertStatus(status port.FormatFileStatus) tui.FormatFileStatus {
	switch status {
	case port.FormatFileStatusSuccess:
		return tui.FormatFileStatusSuccess
	case port.FormatFileStatusNoChange:
		return tui.FormatFileStatusNoChange
	case port.FormatFileStatusFailed:
		return tui.FormatFileStatusFailed
	default:
		return tui.FormatFileStatusNoChange
	}
}
