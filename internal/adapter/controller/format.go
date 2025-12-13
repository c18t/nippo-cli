package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type FormatParams struct {
}

type FormatController interface {
	core.Controller
	Params() *FormatParams
}

type formatController struct {
	bus    port.FormatUseCaseBus `do:""`
	params *FormatParams
}

func NewFormatController(i do.Injector) (FormatController, error) {
	bus, err := do.Invoke[port.FormatUseCaseBus](i)
	if err != nil {
		return nil, err
	}
	return &formatController{
		bus:    bus,
		params: &FormatParams{},
	}, nil
}

func (c *formatController) Params() *FormatParams {
	return c.params
}

func (c *formatController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.FormatCommandUseCaseInputData{})
	return
}
