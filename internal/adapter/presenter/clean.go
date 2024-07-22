package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type CleanCommandPresenter interface {
	Progress(output *port.CleanCommandUseCaseOutputData)
	Complete(output *port.CleanCommandUseCaseOutputData)
	Suspend(err error)
}

type cleanCommandPresenter struct {
	base ConsolePresenter
}

func NewCleanCommandPresenter(i do.Injector) (CleanCommandPresenter, error) {
	return &cleanCommandPresenter{&consolePresenter{}}, nil
}

func (p *cleanCommandPresenter) Progress(output *port.CleanCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *cleanCommandPresenter) Complete(output *port.CleanCommandUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *cleanCommandPresenter) Suspend(err error) {
	cobra.CheckErr(err)
}
