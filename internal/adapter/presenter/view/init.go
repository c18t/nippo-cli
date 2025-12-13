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
	ConfigureProjectSequence_ConfirmOverwrite ConfigureProjectSequence = iota
	ConfigureProjectSequence_InputDriveFolder
	ConfigureProjectSequence_InputSiteUrl
	ConfigureProjectSequence_InputProjectUrl
	ConfigureProjectSequence_InputBranch
	ConfigureProjectSequence_SelectTemplatePath
	ConfigureProjectSequence_SelectAssetPath
	ConfigureProjectSequence_ConfirmGitWarning
)

type ConfigureProjectViewModel struct {
	core.ViewModel
	viewModel
	Sequence        ConfigureProjectSequence
	ConfigExists    bool     // Whether config file already exists
	IsUnderGit      bool     // Whether data_dir is under git repo
	DefaultValues   []string // Defaults from existing config (for merge)
}

type ConfigureProjectView interface {
	Update(vm *ConfigureProjectViewModel)
}

type configureProjectView struct {
	cachedValues        []string
	cachedConfirmations []bool
	cachedErr           error
	wizardRan           bool
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
		v.runWizard(vm)
	}

	// Return cached value or error based on sequence
	if v.cachedErr != nil {
		vm.Input <- v.cachedErr
		return
	}

	switch vm.Sequence {
	case ConfigureProjectSequence_ConfirmOverwrite:
		vm.Input <- v.cachedConfirmations[0]
	case ConfigureProjectSequence_InputDriveFolder:
		vm.Input <- v.cachedValues[0]
	case ConfigureProjectSequence_InputSiteUrl:
		vm.Input <- v.cachedValues[1]
	case ConfigureProjectSequence_InputProjectUrl:
		vm.Input <- v.cachedValues[2]
	case ConfigureProjectSequence_InputBranch:
		vm.Input <- v.cachedValues[3]
	case ConfigureProjectSequence_SelectTemplatePath:
		vm.Input <- v.cachedValues[4]
	case ConfigureProjectSequence_SelectAssetPath:
		vm.Input <- v.cachedValues[5]
	case ConfigureProjectSequence_ConfirmGitWarning:
		vm.Input <- v.cachedConfirmations[1]
	}
}

func (v *configureProjectView) runWizard(vm *ConfigureProjectViewModel) {
	v.wizardRan = true
	v.cachedConfirmations = []bool{true, true} // Default to true if no confirmation needed

	// Handle overwrite confirmation if config exists
	if vm.ConfigExists {
		confirmed, err := tui.RunConfirm("Configuration file already exists. Comments may be lost. Continue?", false)
		if err != nil {
			v.cachedErr = err
			return
		}
		v.cachedConfirmations[0] = confirmed
		if !confirmed {
			v.cachedErr = fmt.Errorf("operation cancelled by user")
			return
		}
	}

	// Placeholder hints (always shown as hints)
	placeholders := []string{
		"https://drive.google.com/drive/folders/xxx",
		"https://nippo.example.com",
		"https://github.com/user/repo",
		"main",
		"/templates",
		"/dist",
	}

	// Get existing config values (will be pre-filled if non-empty)
	existingValues := make([]string, 6)
	if len(vm.DefaultValues) >= 6 {
		copy(existingValues, vm.DefaultValues)
	}

	steps := []tui.WizardStep{
		{
			Label:        "input Google Drive folder URL or ID",
			Placeholder:  placeholders[0],
			InitialValue: existingValues[0],
		},
		{
			Label:        "input site URL",
			Placeholder:  placeholders[1],
			InitialValue: existingValues[1],
		},
		{
			Label:        "input nippo project repository URL",
			Placeholder:  placeholders[2],
			InitialValue: existingValues[2],
		},
		{
			Label:        "input project branch name",
			Placeholder:  placeholders[3],
			InitialValue: existingValues[3],
		},
		{
			Label:        "input template path in ZIP",
			Placeholder:  placeholders[4],
			InitialValue: existingValues[4],
		},
		{
			Label:        "input asset path in ZIP",
			Placeholder:  placeholders[5],
			InitialValue: existingValues[5],
		},
	}

	values, err := tui.RunWizard(steps)
	if err != nil {
		v.cachedErr = err
		return
	}

	v.cachedValues = values

	// Handle git warning if data_dir is under git repo
	if vm.IsUnderGit {
		confirmed, err := tui.RunConfirm("Warning: Data directory is under git repository. Credentials may be tracked. Continue?", false)
		if err != nil {
			v.cachedErr = err
			return
		}
		v.cachedConfirmations[1] = confirmed
		if !confirmed {
			v.cachedErr = fmt.Errorf("operation cancelled by user")
			return
		}
	}
}
