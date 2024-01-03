package presenter

import (
	"fmt"
	"reflect"

	"github.com/c18t/nippo-cli/internal/adapter/presenter/view"
	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type InitPresenter interface {
	Progress(output port.InitUsecaseOutputData)
	Complete(output port.InitUsecaseOutputData)
	Suspend(err error)
}
type InitSettingPresenter interface {
	InitPresenter
	Prompt(ch chan<- interface{}, output *port.InitSettingUsecaseOutputData)
}
type InitSaveDriveTokenPresenter interface {
	InitPresenter
}

type initPresenter struct {
	base         ConsolePresenter
	viewProvider view.InitViewProvider
}

func NewInitSettingPresenter(vp view.InitViewProvider) InitSettingPresenter {
	return &initPresenter{&consolePresenter{}, vp}
}
func NewInitSaveDriveTokenPresenter(vp view.InitViewProvider) InitSaveDriveTokenPresenter {
	return &initPresenter{&consolePresenter{}, vp}
}

func (p *initPresenter) Prompt(ch chan<- interface{}, output *port.InitSettingUsecaseOutputData) {
	switch output.Input.(type) {
	case port.InitSettingProjectUrl:
		vm := &view.ConfigureProjectViewModel{Sequence: view.ConfigureProjectSequence_InputProjectUrl}
		vm.Input = ch
		p.viewProvider.Handle(vm)
	case port.InitSettingProjectTemplatePath:
		vm := &view.ConfigureProjectViewModel{Sequence: view.ConfigureProjectSequence_SelectTemplatePath}
		vm.Input = ch
		p.viewProvider.Handle(vm)
	case port.InitSettingProjectAssetPath:
		vm := &view.ConfigureProjectViewModel{Sequence: view.ConfigureProjectSequence_SelectAssetPath}
		vm.Input = ch
		p.viewProvider.Handle(vm)
	}
}

func (p *initPresenter) Progress(output port.InitUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	vm := &view.ConfigureProjectViewModel{}
	vm.Output = fmt.Sprint(v.String())
	p.viewProvider.Handle(vm)
}

func (p *initPresenter) Complete(output port.InitUsecaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	vm := &view.ConfigureProjectViewModel{}
	vm.Output = fmt.Sprintln(v.String())
	p.viewProvider.Handle(vm)
}

func (p *initPresenter) Suspend(err error) {
	p.base.Suspend(err)
}
