package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type InitPresenter interface {
	Progress(output port.InitUsecaseOutputData)
	Complete(output port.InitUsecaseOutputData)
	Suspend(err error)
}
type InitDownloadProjectPresenter interface {
	InitPresenter
}
type InitSaveDriveTokenPresenter interface {
	InitPresenter
}

type initPresenter struct {
	base ConsolePresenter
}

func NewInitDownloadProjectPresenter() InitDownloadProjectPresenter {
	return &initPresenter{&consolePresenter{}}
}
func NewInitSaveDriveTokenPresenter() InitSaveDriveTokenPresenter {
	return &initPresenter{&consolePresenter{}}
}

func (p *initPresenter) Progress(output port.InitUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *initPresenter) Complete(output port.InitUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *initPresenter) Suspend(err error) {
	p.base.Suspend(err)
}
