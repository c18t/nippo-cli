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
	fmt.Print("update project files... ")
	err := downloadProject()
	cobra.CheckErr(err)
	fmt.Println("ok.")
	return nil
}
