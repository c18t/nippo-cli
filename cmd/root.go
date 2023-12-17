/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package cmd

import (
	"os"

	"github.com/c18t/nippo-cli/internal/cmd/root"
	"github.com/spf13/cobra"
)

var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nippo",
	Short: "nippo - The tool to power my nippo.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: root.CreateCmdFunc(),
	PreRun: func(cmd *cobra.Command, args []string) {
		root.Version = Version
	},
}

func init() {
	cobra.OnInitialize(root.InitConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&root.Params.ConfigFile, "config", "", "config file (default is $XDG_CONFIG_HOME/nippo/.nippo.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&root.Params.Version, "version", "V", false, "show nippo-cli version")
	rootCmd.Flags().BoolVarP(&root.Params.LicenseNotice, "license-notice", "", false, "show copyright notices and license texts of third-party library that nippo-cli depends on")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
