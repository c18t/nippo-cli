package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

var InjectorRoot = AddRootProvider()

func AddRootProvider() *do.RootScope {
	var i = GetInjector().Clone()

	// adapter/controller
	do.Provide(i, controller.NewRootController)

	// usecase/port
	do.Provide(i, port.NewRootUseCaseBus)

	// usecase/intractor
	do.Provide(i, interactor.NewRootCommandInteractor)

	// adapter/presenter
	do.Provide(i, presenter.NewRootCommandPresenter)

	return i
}
