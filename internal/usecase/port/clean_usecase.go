package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"go.uber.org/dig"
)

type CleanUsecaseInputData interface{}
type CleanUsecaseOutputData interface{}
type CleanBuildCacheUsecaseInputData struct {
	CleanUsecaseInputData
}
type CleanBuildCacheUsecaseOutputData struct {
	CleanUsecaseOutputData
	Message string
}
type CleanBuildCacheUsecase interface {
	core.Usecase
	Handle(input *CleanBuildCacheUsecaseInputData)
}
type CleanUsecaseBus interface {
	Handle(input CleanUsecaseInputData)
}
type cleanUsecaseBus struct {
	cleanBuildCache CleanBuildCacheUsecase
}
type inCleanUsecaseBus struct {
	dig.In
	CleanBuildCache CleanBuildCacheUsecase
}

func NewCleanUsecaseBus(bus inCleanUsecaseBus) CleanUsecaseBus {
	return &cleanUsecaseBus{
		cleanBuildCache: bus.CleanBuildCache,
	}
}

func (bus *cleanUsecaseBus) Handle(input CleanUsecaseInputData) {
	switch data := input.(type) {
	case *CleanBuildCacheUsecaseInputData:
		bus.cleanBuildCache.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
