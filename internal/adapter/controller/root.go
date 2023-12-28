package controller

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
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
	bus     port.RootUsecaseBus
	params  *RootParams
	verison string
}

func NewRootController(bus port.RootUsecaseBus) RootController {
	return &rootController{bus: bus, params: &RootParams{}}
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
	core.Cfg.LoadConfig(c.params.ConfigFile)
}

func (c *rootController) Exec(cmd *cobra.Command, args []string) (err error) {
	if c.params.Version {
		// show nippo-cli version
		c.bus.Handle(&port.RootVersionUsecaseInputData{Version: c.verison})
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
