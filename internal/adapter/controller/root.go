package controller

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type RootParams struct {
	ConfigFile    string
	Version       bool
	LicenseNotice bool
}

type RootController interface {
	core.Controller
	Params() *RootParams
	Version(...string) string

	InitConfig()
}

type rootController struct {
	bus     port.RootUseCaseBus
	params  *RootParams
	verison string
}

func NewRootController(i do.Injector) (RootController, error) {
	bus, err := do.Invoke[port.RootUseCaseBus](i)
	if err != nil {
		return nil, err
	}
	return &rootController{
		bus:    bus,
		params: &RootParams{},
	}, nil
}

func (c *rootController) Version(v ...string) string {
	if len(v) > 0 {
		c.verison = v[0]
	}
	return c.verison
}

func (c *rootController) Params() *RootParams {
	return c.params
}

// initConfig reads in config file and ENV variables if set.
func (c *rootController) InitConfig() {
	core.Cfg = &core.Config{}
	err := core.Cfg.LoadConfig(c.params.ConfigFile)
	cobra.CheckErr(err)
}

func (c *rootController) Exec(cmd *cobra.Command, args []string) (err error) {
	if c.params.Version {
		// show nippo-cli version
		c.bus.Handle(&port.RootCommandUseCaseInputData{Version: c.verison})
	} else if c.params.LicenseNotice {
		// show license notice
		fmt.Println("set --license-notice")
		err = fmt.Errorf("not implemented")
	} else {
		// show help
		cmd.Help()
	}
	return
}
