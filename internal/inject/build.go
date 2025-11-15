package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

var InjectorBuild = AddBuildProvider()

func AddBuildProvider() *do.RootScope {
	var i = GetInjector().Clone()

	// adapter/controller
	do.Provide(i, controller.NewBuildController)

	// usecase/port
	do.Provide(i, port.NewBuildUseCaseBus)

	// usecase/intractor
	do.Provide(i, interactor.NewBuildCommandInteractor)

	// adapter/presenter
	do.Provide(i, presenter.NewBuildCommandPresenter)

	return i
}
