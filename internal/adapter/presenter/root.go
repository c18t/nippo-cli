package presenter

import (
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type RootCommandPresenter interface {
	Complete(output *port.RootCommandUseCaseOutputData)
	Suspend(err error)
}

type rootCommandPresenter struct {
	base ConsolePresenter
}

func NewRootCommandPresenter(i do.Injector) (RootCommandPresenter, error) {
	return &rootCommandPresenter{&consolePresenter{}}, nil
}

func (p *rootCommandPresenter) Complete(output *port.RootCommandUseCaseOutputData) {
	p.base.Complete(output.Message)
}

func (p *rootCommandPresenter) Suspend(err error) {
	cobra.CheckErr(err)
}
