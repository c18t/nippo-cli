package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

type BuildCommandPresenter interface {
	Progress(output *port.BuildCommandUseCaseOutputData)
	Complete(output *port.BuildCommandUseCaseOutputData)
	Suspend(err error)
}

type buildCommandPresenter struct {
	base ConsolePresenter
}

func NewBuildCommandPresenter(i do.Injector) (BuildCommandPresenter, error) {
	return &buildCommandPresenter{&consolePresenter{}}, nil
}

func (p *buildCommandPresenter) Progress(output *port.BuildCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *buildCommandPresenter) Complete(output *port.BuildCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *buildCommandPresenter) Suspend(err error) {
	p.base.Suspend(err)
}
