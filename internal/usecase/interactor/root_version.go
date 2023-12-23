package interactor

import (
	"fmt"
	"runtime/debug"

	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type rootVersionInteractor struct{}

func NewRootVersionInteractor() port.RootVersionUsecase {
	return &rootVersionInteractor{}
}

func (u *rootVersionInteractor) Handle(input *port.RootVersionUsecaseInputData) {
	if input.Version != "" {
		// go build -ldflags "-X 'main.version=vx.x.x'"
		fmt.Println(input.Version)
	} else if buildInfo, ok := debug.ReadBuildInfo(); ok {
		// go install version tag
		fmt.Println(buildInfo.Main.Version)
	} else {
		// unknown version
		fmt.Println("(unknown)")
	}
}
