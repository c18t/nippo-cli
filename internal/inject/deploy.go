package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

var InjectorDeploy = AddDeployProvider()

func AddDeployProvider() *do.RootScope {
	var i = GetInjector().Clone()

	// adapter/controller
	do.Provide(i, controller.NewDeployController)

	// usecase/port
	do.Provide(i, port.NewDeployUseCaseBus)

	// usecase/intractor
	do.Provide(i, interactor.NewDeployCommandInteractor)

	// adapter/presenter
	do.Provide(i, presenter.NewDeployCommandPresenter)

	return i
}
