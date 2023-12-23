package core

import "github.com/spf13/cobra"

type RunEFunc func(cmd *cobra.Command, args []string) error

type Controller interface {
	Exec(cmd *cobra.Command, args []string) error
}

type Usecase interface{}
