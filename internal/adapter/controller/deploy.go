package controller

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/spf13/cobra"
)

type DeployController interface {
	core.Controller
}

type deployController struct{}

func NewDeployController() DeployController {
	return &deployController{}
}

func (c *deployController) Exec(cmd *cobra.Command, args []string) error {
	fmt.Println("deploy called")
	return nil
}
