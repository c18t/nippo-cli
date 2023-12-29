package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/spf13/cobra"
)

type BuildController interface {
	core.Controller
}

type buildController struct {
	bus port.BuildUsecaseBus
}

func NewBuildController(bus port.BuildUsecaseBus) BuildController {
	return &buildController{bus}
}

func (c *buildController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.BuildSiteUsecaseInputData{})
	return
}
