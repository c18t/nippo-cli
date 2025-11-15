package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/samber/do/v2"
)

type CleanUseCaseInputData interface{}
type CleanUseCaseOutputData interface{}

type CleanCommandUseCaseInputData struct {
	CleanUseCaseInputData
}
type CleanCommandUseCaseOutputData struct {
	CleanUseCaseOutputData
	Message string
}
type CleanCommandUseCase interface {
	core.UseCase
	Handle(input *CleanCommandUseCaseInputData)
}

type CleanUseCaseBus interface {
	Handle(input CleanUseCaseInputData)
}
type cleanUseCaseBus struct {
	command CleanCommandUseCase `do:""`
}

func NewCleanUseCaseBus(i do.Injector) (CleanUseCaseBus, error) {
	command, err := do.Invoke[CleanCommandUseCase](i)
	if err != nil {
		return nil, err
	}
	return &cleanUseCaseBus{
		command: command,
	}, nil
}

func (bus *cleanUseCaseBus) Handle(input CleanUseCaseInputData) {
	switch data := input.(type) {
	case *CleanCommandUseCaseInputData:
		bus.command.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
