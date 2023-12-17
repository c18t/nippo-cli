/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package cmd

import (
	"github.com/c18t/nippo-cli/internal/cmd/clean"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean built nippo site files",
	Long:  ``,
	RunE:  clean.CreateCmdFunc(),
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
