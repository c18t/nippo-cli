package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/domain/logic/service"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"go.uber.org/dig"
)

var Container = NewContainer()

func NewContainer() *dig.Container {
	container := dig.New()

	// adapter/controller
	container.Provide(controller.NewRootController)
	container.Provide(controller.NewInitController)
	container.Provide(controller.NewUpdateController)
	container.Provide(controller.NewBuildController)
	container.Provide(controller.NewCleanController)
	container.Provide(controller.NewDeployController)

	// adapter/presenter
	container.Provide(presenter.NewRootVersionPresenter)
	container.Provide(presenter.NewInitDownloadProjectPresenter)
	container.Provide(presenter.NewInitSaveDriveTokenPresenter)
	container.Provide(presenter.NewUpdateProjectDataPresenter)
	container.Provide(presenter.NewBuildSitePresenter)
	container.Provide(presenter.NewCleanBuildCachePresenter)
	container.Provide(presenter.NewDeploySitePresenter)

	// usecase/port
	container.Provide(port.NewRootUsecaseBus)
	container.Provide(port.NewInitUsecaseBus)
	container.Provide(port.NewUpdateUsecaseBus)
	container.Provide(port.NewBuildUsecaseBus)
	container.Provide(port.NewCleanUsecaseBus)
	container.Provide(port.NewDeployUsecaseBus)

	// usecase/intractor
	container.Provide(interactor.NewRootVersionInteractor)
	container.Provide(interactor.NewInitDownloadProjectInteractor)
	container.Provide(interactor.NewInitSaveDriveTokenInteractor)
	container.Provide(interactor.NewUpdateProjectDataInteractor)
	container.Provide(interactor.NewBuildSiteInteractor)
	container.Provide(interactor.NewCleanBuildCacheInteractor)
	container.Provide(interactor.NewDeploySiteInteractor)

	// domain/service
	container.Provide(service.NewTemplateService)
	return container
}
