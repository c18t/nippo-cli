/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package cmd

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/spf13/cobra"
)

var initC controller.InitController

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize nippo command",
	Long:  ``,
}

func init() {
	initCmd.RunE = createInitCommand()
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
