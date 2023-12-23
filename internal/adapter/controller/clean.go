package controller

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/spf13/cobra"
)

type CleanController interface {
	core.Controller
}
type cleanController struct{}

func NewCleanController() CleanController {
	return &cleanController{}
}

func (c *cleanController) Exec(cmd *cobra.Command, args []string) error {
	fmt.Println("clean called")
	return nil
}
