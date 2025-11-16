package view

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/manifoldco/promptui"
	"github.com/samber/do/v2"
)

type InitViewProvider interface {
	Handle(vm core.ViewModel)
}
type initViewProvider struct {
	configureProjectView ConfigureProjectView
}

func NewInitViewProvider(i do.Injector) (InitViewProvider, error) {
	configureProjectView, err := do.Invoke[ConfigureProjectView](i)
	if err != nil {
		return nil, err
	}
	return &initViewProvider{
		configureProjectView: configureProjectView,
	}, nil
}

func (vp *initViewProvider) Handle(vm core.ViewModel) {
	switch data := vm.(type) {
	case *ConfigureProjectViewModel:
		vp.configureProjectView.Update(data)
	default:
		panic(fmt.Errorf("view provier for '%T' is not implemented", data))
	}
}

type ConfigureProjectSequence int

const (
	ConfigureProjectSequence_InputProjectUrl ConfigureProjectSequence = iota
	ConfigureProjectSequence_SelectTemplatePath
	ConfigureProjectSequence_SelectAssetPath
)

type ConfigureProjectViewModel struct {
	core.ViewModel
	viewModel
	Sequence ConfigureProjectSequence
}

type ConfigureProjectView interface {
	Update(vm *ConfigureProjectViewModel)
}
type configureProjectView struct{}

func NewConfigureProjectView(_ do.Injector) (ConfigureProjectView, error) {
	return &configureProjectView{}, nil
}

func (v *configureProjectView) Update(vm *ConfigureProjectViewModel) {
	if message(vm.Output) {
		return
	}

	switch vm.Sequence {
	case ConfigureProjectSequence_InputProjectUrl:
		prompt := promptui.Prompt{
			Label:   "input your nippo project repository url",
			Default: "https://github.com/c18t/nippo",
		}
		vm.Input <- either2(prompt.Run())
	case ConfigureProjectSequence_SelectTemplatePath:
		prompt := promptui.Prompt{
			Label:   "input project template path",
			Default: "/templates",
		}
		vm.Input <- either2(prompt.Run())
	case ConfigureProjectSequence_SelectAssetPath:
		prompt := promptui.Prompt{
			Label:   "input project asset path",
			Default: "/output",
		}
		vm.Input <- either2(prompt.Run())
	}
}
