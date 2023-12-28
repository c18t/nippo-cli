package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type DeployPresenter interface {
	Progress(output port.DeployUsecaseOutputData)
	Complete(output port.DeployUsecaseOutputData)
	Suspend(err error)
}
type DeploySitePresenter interface {
	DeployPresenter
}

type deployPresenter struct {
	base ConsolePresenter
}

func NewDeploySitePresenter() DeploySitePresenter {
	return &deployPresenter{&consolePresenter{}}
}

func (p *deployPresenter) Progress(output port.DeployUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *deployPresenter) Complete(output port.DeployUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *deployPresenter) Suspend(err error) {
	p.base.Suspend(err)
}
