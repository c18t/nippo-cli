package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

type DeployCommandPresenter interface {
	Progress(output *port.DeployCommandUseCaseOutputData)
	StopProgress()
	Complete(output *port.DeployCommandUseCaseOutputData)
	Suspend(err error)
}

type deployCommandPresenter struct {
	base ConsolePresenter
}

func NewDeployCommandPresenter(i do.Injector) (DeployCommandPresenter, error) {
	base, err := do.Invoke[ConsolePresenter](i)
	if err != nil {
		return nil, err
	}
	return &deployCommandPresenter{base}, nil
}

func (p *deployCommandPresenter) Progress(output *port.DeployCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *deployCommandPresenter) StopProgress() {
	p.base.StopProgress()
}

func (p *deployCommandPresenter) Complete(output *port.DeployCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *deployCommandPresenter) Suspend(err error) {
	p.base.Suspend(err)
}
