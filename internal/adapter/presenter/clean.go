package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type CleanPresenter interface {
	Progress(output port.CleanUsecaseOutputData)
	Complete(output port.CleanUsecaseOutputData)
	Suspend(err error)
}
type CleanBuildCachePresenter interface {
	CleanPresenter
}

type cleanPresenter struct {
	base ConsolePresenter
}

func NewCleanBuildCachePresenter() CleanBuildCachePresenter {
	return &cleanPresenter{&consolePresenter{}}
}

func (p *cleanPresenter) Progress(output port.CleanUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *cleanPresenter) Complete(output port.CleanUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *cleanPresenter) Suspend(err error) {
	p.base.Suspend(err)
}
