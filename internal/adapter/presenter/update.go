package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type UpdateCommandPresenter interface {
	Progress(output *port.UpdateCommandUseCaseOutputData)
	Complete(output *port.UpdateCommandUseCaseOutputData)
	Suspend(err error)
}

type updateCommandPresenter struct {
	base ConsolePresenter
}

func NewUpdateCommandPresenter(i do.Injector) (UpdateCommandPresenter, error) {
	return &updateCommandPresenter{&consolePresenter{}}, nil
}

func (p *updateCommandPresenter) Progress(output *port.UpdateCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *updateCommandPresenter) Complete(output *port.UpdateCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *updateCommandPresenter) Suspend(err error) {
	cobra.CheckErr(err)
}
