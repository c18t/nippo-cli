package controller

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/spf13/cobra"
)

type UpdateController interface {
	core.Controller
}

type updateController struct {
}

func NewUpdateController() UpdateController {
	return &updateController{}
}

func (c *updateController) Exec(cmd *cobra.Command, args []string) error {
	fmt.Println("update called")
	return nil
}
