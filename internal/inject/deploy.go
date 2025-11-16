package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// DeployPackage groups all services specific to the deploy command.
// Services are lazily initialized when first requested.
var DeployPackage = do.Package(
	// adapter/controller
	do.Lazy(controller.NewDeployController),

	// usecase/port
	do.Lazy(port.NewDeployUseCaseBus),

	// usecase/interactor
	do.Lazy(interactor.NewDeployCommandInteractor),

	// adapter/presenter
	do.Lazy(presenter.NewDeployCommandPresenter),
)

// InjectorDeploy provides a DI container with both base and deploy-specific services.
var InjectorDeploy = do.New(BasePackage, DeployPackage)
