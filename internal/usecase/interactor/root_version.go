package interactor

import (
	"runtime/debug"

	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type rootVersionInteractor struct {
	presenter presenter.RootVersionPresenter
}

func NewRootVersionInteractor(presenter presenter.RootVersionPresenter) port.RootVersionUsecase {
	return &rootVersionInteractor{presenter}
}

func (u *rootVersionInteractor) Handle(input *port.RootVersionUsecaseInputData) {
	output := &port.RootVersionUsecaseOutpuData{}
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
