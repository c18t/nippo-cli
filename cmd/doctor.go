/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check nippo environment health",
	Long: `Check the health of your nippo environment.

This command verifies:
- Configuration file exists and is valid
- Required directories exist (config, data, cache)
- Required files are present (credentials.json, token.json)
- Template and asset files are available

Use this command to diagnose setup issues.`,
}

func init() {
	doctorCmd.RunE = createDoctorCommand()
	rootCmd.AddCommand(doctorCmd)
}
