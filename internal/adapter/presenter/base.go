package presenter

import (
	"github.com/c18t/nippo-cli/internal/adapter/presenter/view/tui"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type ConsolePresenter interface {
	Progress(message string)
	StopProgress()
	Warning(err error)
	Complete(message string)
	Suspend(err error)
}

type consolePresenter struct {
	spinner *tui.SpinnerController
}

// NewConsolePresenter creates a new ConsolePresenter for DI registration.
func NewConsolePresenter(_ do.Injector) (ConsolePresenter, error) {
	return &consolePresenter{
		spinner: tui.NewSpinnerController(),
	}, nil
}

func (p *consolePresenter) Progress(message string) {
	p.spinner.Start(message)
}

func (p *consolePresenter) StopProgress() {
	p.spinner.Stop()
}

func (p *consolePresenter) Warning(err error) {
	p.spinner.Stop()
	tui.PrintWarning(err.Error())
}

func (p *consolePresenter) Complete(message string) {
	p.spinner.Stop()
	tui.PrintSuccess(message)
}

func (p *consolePresenter) Suspend(err error) {
	p.spinner.Stop()
	cobra.CheckErr(err)
}
