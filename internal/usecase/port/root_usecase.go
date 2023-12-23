package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"go.uber.org/dig"
)

type RootUsecaseInputData interface{}
type RootUsecaseOutputData interface{}

type RootVersionUsecaseInputData struct {
	RootUsecaseInputData
	Version string
}
type RootVersionUsecaseOutpuData struct {
	RootUsecaseOutputData
}
type RootVersionUsecase interface {
	core.Usecase
	Handle(input *RootVersionUsecaseInputData)
}

type RootUsecaseBus interface {
	Handle(input RootUsecaseInputData)
}
type rootUsecaseBus struct {
	version RootVersionUsecase
}
type inRootUsecaseBus struct {
	dig.In
	Version RootVersionUsecase
}

func NewRootUsecaseBus(bus inRootUsecaseBus) RootUsecaseBus {
	return &rootUsecaseBus{
		version: bus.Version,
	}
}

func (bus *rootUsecaseBus) Handle(input RootUsecaseInputData) {
	var err error
	switch data := input.(type) {
	case *RootVersionUsecaseInputData:
		bus.version.Handle(data)
	default:
		err = fmt.Errorf("not imprement")
	}
	if err != nil {
		panic(err)
	}
}
