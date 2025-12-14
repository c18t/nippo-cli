package controller

import (
	"errors"
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
	InitConfigErr() error
	RequireConfig(cmd *cobra.Command) error
}

type rootController struct {
	bus           port.RootUseCaseBus
	params        *RootParams
	verison       string
	initConfigErr error
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
// Errors are stored and can be checked later with InitConfigErr().
func (c *rootController) InitConfig() {
	c.initConfigErr = core.InitConfig(c.params.ConfigFile)
	// Don't fail here - let individual commands decide how to handle
}

// InitConfigErr returns the error from InitConfig, if any.
func (c *rootController) InitConfigErr() error {
	return c.initConfigErr
}

// RequireConfig checks if the configuration was loaded successfully.
// Commands that require a config file should call this in PreRunE.
// Returns nil if config was loaded, or an error with instructions.
func (c *rootController) RequireConfig(cmd *cobra.Command) error {
	if c.initConfigErr == nil {
		return nil
	}

	var configNotFound *core.ErrConfigNotFound
	if errors.As(c.initConfigErr, &configNotFound) {
		return fmt.Errorf("configuration not found. Please run `nippo init` first")
	}

	return c.initConfigErr
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
		_ = cmd.Help()
	}
	return
}
