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
	Url          InitSettingProjectUrl
	TemplatePath InitSettingProjectTemplatePath
	AssetPath    InitSettingProjectAssetPath
}
type InitSettingProjectUrl string
type InitSettingProjectTemplatePath string
type InitSettingProjectAssetPath string

type InitSaveDriveTokenUseCaseInputData struct {
	InitUseCaseInputData
}
type InitSaveDriveTokenUsecaseOutputData struct {
	InitUsecaseOutputDataImpl
}

type InitSettingUseCase interface {
	core.UseCase
	Handle(input *InitSettingUseCaseInputData)
}

type InitSaveDriveTokenUseCase interface {
	core.UseCase
	Handle(input *InitSaveDriveTokenUseCaseInputData)
}

type InitUseCaseBus interface {
	Handle(input InitUseCaseInputData)
}
type initUseCaseBus struct {
	configure      InitSettingUseCase        `do:""`
	saveDriveToken InitSaveDriveTokenUseCase `do:""`
}

func NewInitUseCaseBus(i do.Injector) (InitUseCaseBus, error) {
	return &initUseCaseBus{
		configure:      do.MustInvoke[InitSettingUseCase](i),
		saveDriveToken: do.MustInvoke[InitSaveDriveTokenUseCase](i),
	}, nil
}

func (bus *initUseCaseBus) Handle(input InitUseCaseInputData) {
	switch data := input.(type) {
	case *InitSettingUseCaseInputData:
		bus.configure.Handle(data)
	case *InitSaveDriveTokenUseCaseInputData:
		bus.saveDriveToken.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
