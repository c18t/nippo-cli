package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type DeployParams struct {
}

type DeployController interface {
	core.Controller
	Params() *DeployParams
}

type deployController struct {
	bus    port.DeployUseCaseBus `do:""`
	params *DeployParams
}

func NewDeployController(i do.Injector) (DeployController, error) {
	return &deployController{
		bus:    do.MustInvoke[port.DeployUseCaseBus](i),
		params: &DeployParams{},
	}, nil
}

func (c *deployController) Params() *DeployParams {
	return c.params
}

func (c *deployController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.DeployCommandUseCaseInputData{})
	return
}
