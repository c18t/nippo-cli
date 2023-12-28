package presenter

import (
	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type RootVersionPresenter interface {
	Complete(output *port.RootVersionUsecaseOutpuData)
	Suspend(err error)
}

type rootVersionPresenter struct {
	base ConsolePresenter
}

func NewRootVersionPresenter() RootVersionPresenter {
	return &rootVersionPresenter{&consolePresenter{}}
}

func (p *rootVersionPresenter) Complete(output *port.RootVersionUsecaseOutpuData) {
	p.base.Complete(output.Message)
}

func (p *rootVersionPresenter) Suspend(err error) {
	p.base.Suspend(err)
}
