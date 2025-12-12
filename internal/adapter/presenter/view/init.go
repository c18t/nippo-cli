package view

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/adapter/presenter/view/tui"
	"github.com/c18t/nippo-cli/internal/core"
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

type configureProjectView struct {
	cachedValues []string
	cachedErr    error
	wizardRan    bool
}

func NewConfigureProjectView(_ do.Injector) (ConfigureProjectView, error) {
	return &configureProjectView{}, nil
}

func (v *configureProjectView) Update(vm *ConfigureProjectViewModel) {
	if message(vm.Output) {
		return
	}

	// Run wizard on first call and cache all results
	if !v.wizardRan {
		v.runWizard()
	}

	// Return cached value or error based on sequence
	if v.cachedErr != nil {
		vm.Input <- v.cachedErr
		return
	}

	switch vm.Sequence {
	case ConfigureProjectSequence_InputProjectUrl:
		vm.Input <- v.cachedValues[0]
	case ConfigureProjectSequence_SelectTemplatePath:
		vm.Input <- v.cachedValues[1]
	case ConfigureProjectSequence_SelectAssetPath:
		vm.Input <- v.cachedValues[2]
	}
}

func (v *configureProjectView) runWizard() {
	v.wizardRan = true

	steps := []tui.WizardStep{
		{
			Label:        "input your nippo project repository url",
			DefaultValue: "https://github.com/c18t/nippo",
		},
		{
			Label:        "input project template path",
			DefaultValue: "/templates",
		},
		{
			Label:        "input project asset path",
			DefaultValue: "/output",
		},
	}

	values, err := tui.RunWizard(steps)
	if err != nil {
		v.cachedErr = err
		return
	}

	v.cachedValues = values
}
