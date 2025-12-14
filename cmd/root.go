/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package cmd

import (
	"os"

	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/inject"
	"github.com/spf13/cobra"
)

var Version string
var root controller.RootController

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nippo",
	Short: "nippo - The tool to power my nippo.",
	Long:  ``,
}

// commandsWithoutConfig lists commands that can run without a config file
var commandsWithoutConfig = map[string]bool{
	"init":    true,
	"doctor":  true,
	"help":    true,
	"version": true,
	"nippo":   true, // root command (shows help or version)
}

func init() {
	rootCmd.RunE = createRootCommand()
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		root.Version(Version)
	}

	// PersistentPreRunE checks if config is required for the command
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Skip config check for commands that don't require it
		if commandsWithoutConfig[cmd.Name()] {
			return nil
		}
		return root.RequireConfig(cmd)
	}

	cobra.OnInitialize(root.InitConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&root.Params().ConfigFile, "config", "", "config file (default is $XDG_CONFIG_HOME/nippo/nippo.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&root.Params().Version, "version", "V", false, "show nippo-cli version")
	rootCmd.Flags().BoolVarP(&root.Params().LicenseNotice, "license-notice", "", false, "show copyright notices and license texts of third-party library that nippo-cli depends on")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Ensure graceful shutdown on exit
	defer func() {
		// Ignore shutdown errors since the main operation may have already completed
		_ = inject.GetInjector().Shutdown()
	}()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
