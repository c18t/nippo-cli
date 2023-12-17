/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package cmd

import (
	"github.com/c18t/nippo-cli/internal/cmd/deploy"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy nippo site",
	Long:  ``,
	RunE:  deploy.CreateCmdFunc(),
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
