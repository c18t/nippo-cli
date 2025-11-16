package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// BuildPackage groups all services specific to the build command.
// Services are lazily initialized when first requested.
var BuildPackage = do.Package(
	// adapter/controller
	do.Lazy(controller.NewBuildController),

	// usecase/port
	do.Lazy(port.NewBuildUseCaseBus),

	// usecase/interactor
	do.Lazy(interactor.NewBuildCommandInteractor),

	// adapter/presenter
	do.Lazy(presenter.NewBuildCommandPresenter),
)

// InjectorBuild provides a DI container with both base and build-specific services.
var InjectorBuild = do.New(BasePackage, BuildPackage)
