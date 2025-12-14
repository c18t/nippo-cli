/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Manage front-matter in nippo files",
	Long: `Format command manages YAML front-matter in nippo Markdown files on Google Drive.

This command:
- Adds front-matter to files that don't have it
- Adds 'created' field if missing (using Drive's createdTime)
- Replaces 'updated: now' placeholder with Drive's modifiedTime

Files are only uploaded when changes are made. The command tracks the last format
timestamp and only processes files modified since then.`,
}

func init() {
	formatCmd.RunE = createFormatCommand()
	rootCmd.AddCommand(formatCmd)
}
