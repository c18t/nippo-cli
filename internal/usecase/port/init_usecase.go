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
type InitDownloadProjectUsecaseInputData struct {
	InitUsecaseInputData
}
type InitDownloadProjectUsecaseOutpuData struct {
	InitUsecaseOutputDataImpl
}
type InitDownloadProjectUsecase interface {
	core.Usecase
	Handle(input *InitDownloadProjectUsecaseInputData)
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
	downloadProject InitDownloadProjectUsecase
	saveDriveToken  InitSaveDriveTokenUsecase
}
type inInitUsecaseBus struct {
	dig.In
	DownloadProject InitDownloadProjectUsecase
	SaveDriveToken  InitSaveDriveTokenUsecase
}

func NewInitUsecaseBus(bus inInitUsecaseBus) InitUsecaseBus {
	return &initUsecaseBus{
		downloadProject: bus.DownloadProject,
		saveDriveToken:  bus.SaveDriveToken,
	}
}

func (bus *initUsecaseBus) Handle(input InitUsecaseInputData) {
	switch data := input.(type) {
	case *InitDownloadProjectUsecaseInputData:
		bus.downloadProject.Handle(data)
	case *InitSaveDriveTokenUsecaseInputData:
		bus.saveDriveToken.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
