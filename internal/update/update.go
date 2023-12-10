package update

import (
	"fmt"

	"github.com/spf13/cobra"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

func CreateCmdFunc() RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("update called")
		return nil
	}
}
