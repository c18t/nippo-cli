package presenter

import (
	"fmt"

	"github.com/spf13/cobra"
)

type ConsolePresenter interface {
	Progress(message string)
	Warning(err error)
	Complete(message string)
	Suspend(err error)
}

type consolePresenter struct{}

func (presenter *consolePresenter) Progress(message string) {
	fmt.Print(message)
}

func (presenter *consolePresenter) Warning(err error) {
	fmt.Print(err)
}

func (presenter *consolePresenter) Complete(message string) {
	fmt.Println(message)
}

func (presenter *consolePresenter) Suspend(err error) {
	cobra.CheckErr(err)
}
