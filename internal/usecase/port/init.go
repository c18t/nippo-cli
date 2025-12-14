package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/samber/do/v2"
)

type InitUseCaseInputData interface{}
type InitUseCaseOutputData interface{}
type InitUsecaseOutputDataImpl struct {
	InitUseCaseOutputData
	Message string
}

type InitCommandUseCaseInputData struct {
	InitUseCaseInputData
}
type InitCommandUseCaseOutputData struct {
	InitUsecaseOutputDataImpl
}
type InitSettingUseCaseInputData struct {
	InitUseCaseInputData
}
type InitSettingUseCaseOutputData struct {
	InitUsecaseOutputDataImpl
	Input             interface{}
	Project           InitSettingProject
	ProjectConfigured bool
}
type InitSettingProject struct {
	DriveFolder  InitSettingProjectDriveFolder
	SiteUrl      InitSettingProjectSiteUrl
	Url          InitSettingProjectUrl
	Branch       InitSettingProjectBranch
	TemplatePath InitSettingProjectTemplatePath
	AssetPath    InitSettingProjectAssetPath
}
type InitSettingProjectDriveFolder string
type InitSettingProjectSiteUrl string
type InitSettingProjectUrl string
type InitSettingProjectBranch string
type InitSettingProjectTemplatePath string
type InitSettingProjectAssetPath string

// Confirmation prompts
type InitSettingConfirmOverwrite bool
type InitSettingConfirmGitWarning bool

type InitSettingUseCase interface {
	core.UseCase
	Handle(input *InitSettingUseCaseInputData)
}

type InitUseCaseBus interface {
	Handle(input InitUseCaseInputData)
}
type initUseCaseBus struct {
	configure InitSettingUseCase `do:""`
}

func NewInitUseCaseBus(i do.Injector) (InitUseCaseBus, error) {
	configure, err := do.Invoke[InitSettingUseCase](i)
	if err != nil {
		return nil, err
	}
	return &initUseCaseBus{
		configure: configure,
	}, nil
}

func (bus *initUseCaseBus) Handle(input InitUseCaseInputData) {
	switch data := input.(type) {
	case *InitSettingUseCaseInputData:
		bus.configure.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
