package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/spf13/cobra"
)

type CleanController interface {
	core.Controller
}
type cleanController struct {
	bus port.CleanUsecaseBus
}

func NewCleanController(bus port.CleanUsecaseBus) CleanController {
	return &cleanController{bus}
}

func (c *cleanController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.CleanBuildCacheUsecaseInputData{})
	return
}
