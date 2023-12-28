package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/domain/service"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"go.uber.org/dig"
)

var Container = NewContainer()

func NewContainer() *dig.Container {
	container := dig.New()

	// adapter/controller
	container.Provide(controller.NewBuildController)
	container.Provide(controller.NewCleanController)
	container.Provide(controller.NewDeployController)
	container.Provide(controller.NewInitController)
	container.Provide(controller.NewRootController)
	container.Provide(controller.NewUpdateController)

	// adapter/presenter
	container.Provide(presenter.NewRootVersionPresenter)
	container.Provide(presenter.NewInitDownloadProjectPresenter)
	container.Provide(presenter.NewInitSaveDriveTokenPresenter)

	// usecase/port
	container.Provide(port.NewRootUsecaseBus)
	container.Provide(port.NewInitUsecaseBus)

	// usecase/intractor
	container.Provide(interactor.NewRootVersionInteractor)
	container.Provide(interactor.NewInitDownloadProjectInteractor)
	container.Provide(interactor.NewInitSaveDriveTokenInteractor)

	// domain/service
	container.Provide(service.NewTemplateService)
	return container
}
