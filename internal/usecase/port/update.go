package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/samber/do/v2"
)

type UpdateUseCaseInputData interface{}
type UpdateUseCaseOutputData interface{}

type UpdateCommandUseCaseInputData struct {
	UpdateUseCaseInputData
}
type UpdateCommandUseCaseOutputData struct {
	UpdateUseCaseOutputData
	Message string
}
type UpdateCommandUseCase interface {
	core.UseCase
	Handle(input *UpdateCommandUseCaseInputData)
}

type UpdateUseCaseBus interface {
	Handle(input UpdateUseCaseInputData)
}
type updateUseCaseBus struct {
	command UpdateCommandUseCase `do:""`
}

func NewUpdateUseCaseBus(i do.Injector) (UpdateUseCaseBus, error) {
	command, err := do.Invoke[UpdateCommandUseCase](i)
	if err != nil {
		return nil, err
	}
	return &updateUseCaseBus{
		command: command,
	}, nil
}

func (bus *updateUseCaseBus) Handle(input UpdateUseCaseInputData) {
	switch data := input.(type) {
	case *UpdateCommandUseCaseInputData:
		bus.command.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
