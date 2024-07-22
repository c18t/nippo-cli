package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

var InjectorClean = AddCleanProvider()

func AddCleanProvider() *do.RootScope {
	// adapter/controller
	do.Provide(Injector, controller.NewCleanController)

	// usecase/port
	do.Provide(Injector, port.NewCleanUseCaseBus)

	// usecase/intractor
	do.Provide(Injector, interactor.NewCleanCommandInteractor)

	// adapter/presenter
	do.Provide(Injector, presenter.NewCleanCommandPresenter)

	return Injector
}
