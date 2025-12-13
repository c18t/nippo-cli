/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Google Drive",
	Long: `Authenticate with Google Drive using OAuth 2.0.

This command requires credentials.json to be present in the data directory.
If not present, instructions will be shown to download it from Google Cloud Console.

The command will:
1. Check for credentials.json
2. Open a browser for Google OAuth authentication
3. Save the token to token.json in the data directory

You can re-run this command anytime to refresh your authentication.`,
}

func init() {
	authCmd.RunE = createAuthCommand()
	rootCmd.AddCommand(authCmd)
}
