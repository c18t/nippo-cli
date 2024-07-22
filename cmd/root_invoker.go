package cmd

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/inject"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

func createRootCommand() core.RunEFunc {
	cmd, err := do.Invoke[controller.RootController](inject.InjectorRoot)
	cobra.CheckErr(err)
	root = cmd
	return cmd.Exec
}
