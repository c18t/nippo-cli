package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/spf13/cobra"
)

type DeployController interface {
	core.Controller
}

type deployController struct {
	bus port.DeployUsecaseBus
}

func NewDeployController(bus port.DeployUsecaseBus) DeployController {
	return &deployController{bus}
}

func (c *deployController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.DeploySiteUsecaseInputData{})
	return
}
