package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type InitParams struct {
}

type InitController interface {
	core.Controller
	Params() *InitParams
}

type initController struct {
	bus    port.InitUseCaseBus `do:""`
	params *InitParams
}

func NewInitController(i do.Injector) (InitController, error) {
	bus, err := do.Invoke[port.InitUseCaseBus](i)
	if err != nil {
		return nil, err
	}
	return &initController{
		bus:    bus,
		params: &InitParams{},
	}, nil
}

func (c *initController) Params() *InitParams {
	return c.params
}

func (c *initController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.InitSettingUseCaseInputData{})
	c.bus.Handle(&port.InitSaveDriveTokenUseCaseInputData{})
	return
}
