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
	var i = GetInjector().Clone()

	// adapter/controller
	do.Provide(i, controller.NewCleanController)

	// usecase/port
	do.Provide(i, port.NewCleanUseCaseBus)

	// usecase/intractor
	do.Provide(i, interactor.NewCleanCommandInteractor)

	// adapter/presenter
	do.Provide(i, presenter.NewCleanCommandPresenter)

	return i
}
