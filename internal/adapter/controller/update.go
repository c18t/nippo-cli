package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type UpdateParams struct {
}

type UpdateController interface {
	core.Controller
	Params() *UpdateParams
}

type updateController struct {
	bus    port.UpdateUseCaseBus `do:""`
	params *UpdateParams
}

func NewUpdateController(i do.Injector) (UpdateController, error) {
	bus, err := do.Invoke[port.UpdateUseCaseBus](i)
	if err != nil {
		return nil, err
	}
	return &updateController{
		bus:    bus,
		params: &UpdateParams{},
	}, nil
}

func (c *updateController) Params() *UpdateParams {
	return c.params
}

func (c *updateController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.UpdateCommandUseCaseInputData{})
	return
}
