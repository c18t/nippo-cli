package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/adapter/presenter/view"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

type InitCommandPresenter interface {
	Progress(output port.InitUseCaseOutputData)
	Complete(output port.InitUseCaseOutputData)
	Suspend(err error)
}
type InitSettingPresenter interface {
	InitCommandPresenter
	Prompt(ch chan<- interface{}, output *port.InitSettingUseCaseOutputData)
}
type InitSaveDriveTokenPresenter interface {
	InitCommandPresenter
}

type initCommandPresenter struct {
	base         ConsolePresenter      `do:""`
	viewProvider view.InitViewProvider `do:""`
}

func NewInitSettingPresenter(i do.Injector) (InitSettingPresenter, error) {
	base, err := do.Invoke[ConsolePresenter](i)
	if err != nil {
		return nil, err
	}
	viewProvider, err := do.Invoke[view.InitViewProvider](i)
	if err != nil {
		return nil, err
	}
	return &initCommandPresenter{
		base,
		viewProvider,
	}, nil
}

func NewInitSaveDriveTokenPresenter(i do.Injector) (InitSaveDriveTokenPresenter, error) {
	base, err := do.Invoke[ConsolePresenter](i)
	if err != nil {
		return nil, err
	}
	viewProvider, err := do.Invoke[view.InitViewProvider](i)
	if err != nil {
		return nil, err
	}
	return &initCommandPresenter{
		base,
		viewProvider,
	}, nil
}

func (p *initCommandPresenter) Prompt(ch chan<- interface{}, output *port.InitSettingUseCaseOutputData) {
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

func (p *initCommandPresenter) Progress(output port.InitUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *initCommandPresenter) Complete(output port.InitUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *initCommandPresenter) Suspend(err error) {
	p.base.Suspend(err)
}
