package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/samber/do/v2"
)

type BuildUseCaseInputData interface{}
type BuildUseCaseOutputData interface{}

type BuildCommandUseCaseInputData struct {
	BuildUseCaseInputData
}
type BuildCommandUseCaseOutputData struct {
	BuildUseCaseOutputData
	Message string
}
type BuildCommandUseCase interface {
	core.UseCase
	Handle(input *BuildCommandUseCaseInputData)
}

type BuildUseCaseBus interface {
	Handle(input BuildUseCaseInputData)
}
type buildUseCaseBus struct {
	command BuildCommandUseCase `do:""`
}

func NewBuildUseCaseBus(i do.Injector) (BuildUseCaseBus, error) {
	return &buildUseCaseBus{
		command: do.MustInvoke[BuildCommandUseCase](i),
	}, nil
}

func (bus *buildUseCaseBus) Handle(input BuildUseCaseInputData) {
	switch data := input.(type) {
	case *BuildCommandUseCaseInputData:
		bus.command.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
