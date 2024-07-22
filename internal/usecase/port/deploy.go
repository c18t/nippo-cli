package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/samber/do/v2"
)

type DeployUseCaseInputData interface{}
type DeployUseCaseOutputData interface{}

type DeployCommandUseCaseInputData struct {
	DeployUseCaseInputData
}
type DeployCommandUseCaseOutputData struct {
	DeployUseCaseOutputData
	Message string
}
type DeployCommandUseCase interface {
	core.UseCase
	Handle(input *DeployCommandUseCaseInputData)
}

type DeployUseCaseBus interface {
	Handle(input DeployUseCaseInputData)
}
type deployUseCaseBus struct {
	command DeployCommandUseCase `do:""`
}

func NewDeployUseCaseBus(i do.Injector) (DeployUseCaseBus, error) {
	return &deployUseCaseBus{
		command: do.MustInvoke[DeployCommandUseCase](i),
	}, nil
}

func (bus *deployUseCaseBus) Handle(input DeployUseCaseInputData) {
	switch data := input.(type) {
	case *DeployCommandUseCaseInputData:
		bus.command.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
