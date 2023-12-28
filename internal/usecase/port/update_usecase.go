package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"go.uber.org/dig"
)

type UpdateUsecaseInputData interface{}
type UpdateUsecaseOutputData interface{}
type UpdateProjectDataUsecaseInputData struct {
	UpdateUsecaseInputData
}
type UpdateProjectDataUsecaseOutputData struct {
	UpdateUsecaseOutputData
	Message string
}
type UpdateProjectDataUsecase interface {
	core.Usecase
	Handle(input *UpdateProjectDataUsecaseInputData)
}
type UpdateUsecaseBus interface {
	Handle(input UpdateUsecaseInputData)
}
type updateUsecaseBus struct {
	updateProjectData UpdateProjectDataUsecase
}
type inUpdateUsecaseBus struct {
	dig.In
	UpdateProjectData UpdateProjectDataUsecase
}

func NewUpdateUsecaseBus(bus inUpdateUsecaseBus) UpdateUsecaseBus {
	return &updateUsecaseBus{
		updateProjectData: bus.UpdateProjectData,
	}
}

func (bus *updateUsecaseBus) Handle(input UpdateUsecaseInputData) {
	switch data := input.(type) {
	case *UpdateProjectDataUsecaseInputData:
		bus.updateProjectData.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
