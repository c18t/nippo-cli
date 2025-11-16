package interactor

import (
	"runtime/debug"

	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

type rootCommandInteractor struct {
	presenter presenter.RootCommandPresenter `do:""`
}

func NewRootCommandInteractor(i do.Injector) (port.RootCommandUseCase, error) {
	p, err := do.Invoke[presenter.RootCommandPresenter](i)
	if err != nil {
		return nil, err
	}
	return &rootCommandInteractor{
		presenter: p,
	}, nil
}

func (u *rootCommandInteractor) Handle(input *port.RootCommandUseCaseInputData) {
	output := &port.RootCommandUseCaseOutputData{}
	if input.Version != "" {
		// go build -ldflags "-X 'main.version=vx.x.x'"
		output.Message = input.Version
	} else if buildInfo, ok := debug.ReadBuildInfo(); ok {
		// go install version tag
		output.Message = buildInfo.Main.Version
	} else {
		// unknown version
		output.Message = "(unknown)"
	}

	u.presenter.Complete(output)
}
