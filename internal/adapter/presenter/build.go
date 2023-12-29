package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type BuildPresenter interface {
	Progress(output port.BuildUsecaseOutputData)
	Complete(output port.BuildUsecaseOutputData)
	Suspend(err error)
}
type BuildSitePresenter interface {
	BuildPresenter
}

type buildPresenter struct {
	base ConsolePresenter
}

func NewBuildSitePresenter() BuildSitePresenter {
	return &buildPresenter{&consolePresenter{}}
}

func (p *buildPresenter) Progress(output port.BuildUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *buildPresenter) Complete(output port.BuildUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *buildPresenter) Suspend(err error) {
	p.base.Suspend(err)
}
