package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"go.uber.org/dig"
)

type BuildUsecaseInputData interface{}
type BuildUsecaseOutputData interface{}
type BuildSiteUsecaseInputData struct {
	BuildUsecaseInputData
}
type BuildSiteUsecaseOutputData struct {
	BuildUsecaseOutputData
	Message string
}
type BuildSiteUsecase interface {
	core.Usecase
	Handle(input *BuildSiteUsecaseInputData)
}
type BuildUsecaseBus interface {
	Handle(input BuildUsecaseInputData)
}
type buildUsecaseBus struct {
	buildSite BuildSiteUsecase
}
type inBuildUsecaseBus struct {
	dig.In
	BuildSite BuildSiteUsecase
}

func NewBuildUsecaseBus(bus inBuildUsecaseBus) BuildUsecaseBus {
	return &buildUsecaseBus{
		buildSite: bus.BuildSite,
	}
}

func (bus *buildUsecaseBus) Handle(input BuildUsecaseInputData) {
	switch data := input.(type) {
	case *BuildSiteUsecaseInputData:
		bus.buildSite.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
