package presenter

import (
	"reflect"

	"github.com/c18t/nippo-cli/internal/adapter/presenter/view"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

type InitCommandPresenter interface {
	Progress(output port.InitUseCaseOutputData)
	StopProgress()
	Complete(output port.InitUseCaseOutputData)
	Suspend(err error)
	IsCancelled() bool
}
type InitSettingPresenter interface {
	InitCommandPresenter
	Prompt(ch chan<- interface{}, output *port.InitSettingUseCaseOutputData)
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

func (p *initCommandPresenter) Prompt(ch chan<- interface{}, output *port.InitSettingUseCaseOutputData) {
	vm := &view.ConfigureProjectViewModel{}
	vm.Input = ch

	// Populate default values from existing config if available
	if core.Cfg != nil {
		vm.DefaultValues = []string{
			core.Cfg.Project.DriveFolderId,
			core.Cfg.Project.SiteUrl,
			core.Cfg.Project.Url,
			core.Cfg.Project.Branch,
			core.Cfg.Project.TemplatePath,
			core.Cfg.Project.AssetPath,
		}
	}

	switch output.Input.(type) {
	case port.InitSettingConfirmOverwrite:
		vm.Sequence = view.ConfigureProjectSequence_ConfirmOverwrite
		vm.ConfigExists = true
	case port.InitSettingProjectDriveFolder:
		vm.Sequence = view.ConfigureProjectSequence_InputDriveFolder
	case port.InitSettingProjectSiteUrl:
		vm.Sequence = view.ConfigureProjectSequence_InputSiteUrl
	case port.InitSettingProjectUrl:
		vm.Sequence = view.ConfigureProjectSequence_InputProjectUrl
	case port.InitSettingProjectBranch:
		vm.Sequence = view.ConfigureProjectSequence_InputBranch
	case port.InitSettingProjectTemplatePath:
		vm.Sequence = view.ConfigureProjectSequence_SelectTemplatePath
	case port.InitSettingProjectAssetPath:
		vm.Sequence = view.ConfigureProjectSequence_SelectAssetPath
	case port.InitSettingConfirmGitWarning:
		vm.Sequence = view.ConfigureProjectSequence_ConfirmGitWarning
		vm.IsUnderGit = true
	}

	p.viewProvider.Handle(vm)
}

func (p *initCommandPresenter) Progress(output port.InitUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Progress(v.String())
}

func (p *initCommandPresenter) StopProgress() {
	p.base.StopProgress()
}

func (p *initCommandPresenter) Complete(output port.InitUseCaseOutputData) {
	v := reflect.Indirect(reflect.ValueOf(output)).FieldByName("Message")
	p.base.Complete(v.String())
}

func (p *initCommandPresenter) Suspend(err error) {
	p.base.Suspend(err)
}

func (p *initCommandPresenter) IsCancelled() bool {
	return p.base.IsCancelled()
}
