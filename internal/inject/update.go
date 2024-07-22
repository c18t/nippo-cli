package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

var InjectorUpdate = AddUpdateProvider()

func AddUpdateProvider() *do.RootScope {
	var i = Injector.Clone()

	// adapter/controller
	do.Provide(i, controller.NewUpdateController)

	// usecase/port
	do.Provide(i, port.NewUpdateUseCaseBus)

	// usecase/intractor
	do.Provide(i, interactor.NewUpdateCommandInteractor)

	// adapter/presenter
	do.Provide(i, presenter.NewUpdateCommandPresenter)

	return i
}
