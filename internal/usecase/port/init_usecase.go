package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"go.uber.org/dig"
)

type InitUsecaseInputData interface{}
type InitUsecaseOutputData interface{}

type InitUsecaseOutputDataImpl struct {
	InitUsecaseOutputData
	Message string
}
type InitSettingUsecaseInputData struct {
	InitUsecaseInputData
}
type InitSettingUsecaseOutputData struct {
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

type InitSettingUsecase interface {
	core.Usecase
	Handle(input *InitSettingUsecaseInputData)
}

type InitSaveDriveTokenUsecaseInputData struct {
	InitUsecaseInputData
}
type InitSaveDriveTokenUsecaseOutputData struct {
	InitUsecaseOutputDataImpl
}
type InitSaveDriveTokenUsecase interface {
	core.Usecase
	Handle(input *InitSaveDriveTokenUsecaseInputData)
}

type InitUsecaseBus interface {
	Handle(input InitUsecaseInputData)
}
type initUsecaseBus struct {
	configure      InitSettingUsecase
	saveDriveToken InitSaveDriveTokenUsecase
}
type inInitUsecaseBus struct {
	dig.In
	Configure      InitSettingUsecase
	SaveDriveToken InitSaveDriveTokenUsecase
}

func NewInitUsecaseBus(bus inInitUsecaseBus) InitUsecaseBus {
	return &initUsecaseBus{
		configure:      bus.Configure,
		saveDriveToken: bus.SaveDriveToken,
	}
}

func (bus *initUsecaseBus) Handle(input InitUsecaseInputData) {
	switch data := input.(type) {
	case *InitSettingUsecaseInputData:
		bus.configure.Handle(data)
	case *InitSaveDriveTokenUsecaseInputData:
		bus.saveDriveToken.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
