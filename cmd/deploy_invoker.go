package cmd

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/inject"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

func createDeployCommand() core.RunEFunc {
	cmd, err := do.Invoke[controller.DeployController](inject.InjectorDeploy)
	cobra.CheckErr(err)
	return cmd.Exec
}
