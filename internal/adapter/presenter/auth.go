package presenter

import (
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

type AuthPresenter interface {
	Progress(output *port.AuthUseCaseOutputData)
	Complete(output *port.AuthUseCaseOutputData)
	Suspend(err error)
	IsCancelled() bool
}

type authPresenter struct {
	base ConsolePresenter
}

func NewAuthPresenter(i do.Injector) (AuthPresenter, error) {
	base, err := do.Invoke[ConsolePresenter](i)
	if err != nil {
		return nil, err
	}
	return &authPresenter{
		base: base,
	}, nil
}

func (p *authPresenter) Progress(output *port.AuthUseCaseOutputData) {
	p.base.Progress(output.Message)
}

func (p *authPresenter) Complete(output *port.AuthUseCaseOutputData) {
	p.base.Complete(output.Message)
}

func (p *authPresenter) Suspend(err error) {
	p.base.Suspend(err)
}

func (p *authPresenter) IsCancelled() bool {
	return p.base.IsCancelled()
}
