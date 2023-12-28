package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/spf13/cobra"
)

type InitParams struct{}

type InitController interface {
	core.Controller
}
type initController struct {
	bus    port.InitUsecaseBus
	params *InitParams
}

func NewInitController(bus port.InitUsecaseBus) InitController {
	return &initController{bus: bus, params: &InitParams{}}
}

func (c *initController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.InitDownloadProjectUsecaseInputData{})
	c.bus.Handle(&port.InitSaveDriveTokenUsecaseInputData{})
	return
}
