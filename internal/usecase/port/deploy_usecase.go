package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"go.uber.org/dig"
)

type DeployUsecaseInputData interface{}
type DeployUsecaseOutputData interface{}
type DeploySiteUsecaseInputData struct {
	DeployUsecaseInputData
}
type DeploySiteUsecaseOutputData struct {
	DeployUsecaseOutputData
	Message string
}
type DeploySiteUsecase interface {
	core.Usecase
	Handle(input *DeploySiteUsecaseInputData)
}
type DeployUsecaseBus interface {
	Handle(input DeployUsecaseInputData)
}
type deployUsecaseBus struct {
	deploySite DeploySiteUsecase
}
type inDeployUsecaseBus struct {
	dig.In
	DeploySite DeploySiteUsecase
}

func NewDeployUsecaseBus(bus inDeployUsecaseBus) DeployUsecaseBus {
	return &deployUsecaseBus{
		deploySite: bus.DeploySite,
	}
}

func (bus *deployUsecaseBus) Handle(input DeployUsecaseInputData) {
	switch data := input.(type) {
	case *DeploySiteUsecaseInputData:
		bus.deploySite.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
