package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type CleanParams struct {
}

type CleanController interface {
	core.Controller
	Params() *CleanParams
}

type cleanController struct {
	bus    port.CleanUseCaseBus `do:""`
	params *CleanParams
}

func NewCleanController(i do.Injector) (CleanController, error) {
	bus, err := do.Invoke[port.CleanUseCaseBus](i)
	if err != nil {
		return nil, err
	}
	return &cleanController{
		bus:    bus,
		params: &CleanParams{},
	}, nil
}

func (c *cleanController) Params() *CleanParams {
	return c.params
}

func (c *cleanController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.CleanCommandUseCaseInputData{})
	return
}
