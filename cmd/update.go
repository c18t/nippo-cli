/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package cmd

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/inject"
	"github.com/spf13/cobra"
)

var update controller.CleanController

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Download latest nippo source project",
	Long:  ``,
}

func init() {
	updateCmd.RunE = createUpdateCommand()
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createUpdateCommand() core.RunEFunc {
	_ = inject.Container.Invoke(func(c controller.UpdateController) error {
		update = c
		return nil
	})
	return update.Exec
}
