package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/adapter/presenter/view/tui"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
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
