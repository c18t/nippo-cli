package controller

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type AuthParams struct{}

type AuthController interface {
	core.Controller
	Params() *AuthParams
}

type authController struct {
	bus    port.AuthUseCaseBus
	params *AuthParams
}

func NewAuthController(i do.Injector) (AuthController, error) {
	bus, err := do.Invoke[port.AuthUseCaseBus](i)
	if err != nil {
		return nil, err
	}
	return &authController{
		bus:    bus,
		params: &AuthParams{},
	}, nil
}

func (c *authController) Params() *AuthParams {
	return c.params
}

func (c *authController) Exec(cmd *cobra.Command, args []string) (err error) {
	c.bus.Handle(&port.AuthUseCaseInputData{})
	return
}
