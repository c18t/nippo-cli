package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type DoctorController interface {
	core.Controller
}

type doctorController struct {
	useCase port.DoctorUseCase `do:""`
}

func NewDoctorController(i do.Injector) (DoctorController, error) {
	useCase, err := do.Invoke[port.DoctorUseCase](i)
	if err != nil {
		return nil, err
	}
	return &doctorController{
		useCase: useCase,
	}, nil
}

func (c *doctorController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.useCase.Handle(&port.DoctorUseCaseInputData{})
	return nil
}
