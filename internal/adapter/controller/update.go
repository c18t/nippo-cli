package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/spf13/cobra"
)

type UpdateController interface {
	core.Controller
}

type updateController struct {
	bus port.UpdateUsecaseBus
}

func NewUpdateController(bus port.UpdateUsecaseBus) UpdateController {
	return &updateController{bus}
}

func (c *updateController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.UpdateProjectDataUsecaseInputData{})
	return
}
