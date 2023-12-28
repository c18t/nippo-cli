package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type UpdatePresenter interface {
	Progress(output port.UpdateUsecaseOutputData)
	Complete(output port.UpdateUsecaseOutputData)
	Suspend(err error)
}
type UpdateProjectDataPresenter interface {
	UpdatePresenter
}

type updatePresenter struct {
	base ConsolePresenter
}

func NewUpdateProjectDataPresenter() UpdateProjectDataPresenter {
	return &updatePresenter{&consolePresenter{}}
}

func (p *updatePresenter) Progress(output port.UpdateUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *updatePresenter) Complete(output port.UpdateUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *updatePresenter) Suspend(err error) {
	p.base.Suspend(err)
}
