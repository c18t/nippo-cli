package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type BuildParams struct {
}

type BuildController interface {
	core.Controller
	Params() *BuildParams
}

type buildController struct {
	bus    port.BuildUseCaseBus `do:""`
	params *BuildParams
}

func NewBuildController(i do.Injector) (BuildController, error) {
	return &buildController{
		bus:    do.MustInvoke[port.BuildUseCaseBus](i),
		params: &BuildParams{},
	}, nil
}

func (c *buildController) Params() *BuildParams {
	return c.params
}

func (c *buildController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.BuildCommandUseCaseInputData{})
	return
}
