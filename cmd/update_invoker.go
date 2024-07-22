package cmd

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/inject"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

func createUpdateCommand() core.RunEFunc {
	cmd, err := do.Invoke[controller.UpdateController](inject.InjectorUpdate)
	cobra.CheckErr(err)
	return cmd.Exec
}
