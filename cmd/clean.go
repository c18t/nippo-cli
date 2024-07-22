/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package cmd

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/spf13/cobra"
)

var clean controller.CleanController

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean built nippo site files",
	Long:  ``,
}

func init() {
	cleanCmd.RunE = createCleanCommand()
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
