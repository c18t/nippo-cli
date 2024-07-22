package core

import "github.com/spf13/cobra"

type RunEFunc func(cmd *cobra.Command, args []string) (err error)

type Controller interface {
	Exec(cmd *cobra.Command, args []string) (err error)
}

type UseCase interface{}

type ViewModel interface{}
