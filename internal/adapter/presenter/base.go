package presenter

import (
	"fmt"
	"strings"

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
	IsCancelled() bool
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
	// Manually finish the previous message before starting a new one
	prevMsg := p.spinner.GetCurrentMessage()
	p.spinner.Stop()
	if strings.TrimSpace(prevMsg) != "" {
		tui.Print(fmt.Sprintf("%s %s\n", prevMsg, tui.SuccessStyle.Render("ok.")))
	}
	p.spinner.Start(message)
}

func (p *consolePresenter) StopProgress() {
	// Finish the current message if any
	prevMsg := p.spinner.GetCurrentMessage()
	p.spinner.Stop()
	if strings.TrimSpace(prevMsg) != "" {
		tui.Print(fmt.Sprintf("%s %s\n", prevMsg, tui.SuccessStyle.Render("ok.")))
	}
}

func (p *consolePresenter) Warning(err error) {
	p.spinner.Stop()
	tui.PrintWarning(err.Error())
}

func (p *consolePresenter) Complete(message string) {
	// Finish the previous message if any
	prevMsg := p.spinner.GetCurrentMessage()
	p.spinner.Stop()
	if strings.TrimSpace(prevMsg) != "" {
		tui.Print(fmt.Sprintf("%s %s\n", prevMsg, tui.SuccessStyle.Render("ok.")))
	}
	tui.PrintSuccess(message)
}

func (p *consolePresenter) Suspend(err error) {
	p.spinner.Stop()
	cobra.CheckErr(err)
}

func (p *consolePresenter) IsCancelled() bool {
	return p.spinner.IsCancelled()
}
