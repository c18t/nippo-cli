package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// RootPackage groups all services specific to the root command.
// Services are lazily initialized when first requested.
var RootPackage = do.Package(
	// adapter/controller
	do.Lazy(controller.NewRootController),

	// usecase/port
	do.Lazy(port.NewRootUseCaseBus),

	// usecase/interactor
	do.Lazy(interactor.NewRootCommandInteractor),

	// adapter/presenter
	do.Lazy(presenter.NewRootCommandPresenter),
)

// InjectorRoot provides a DI container with both base and root-specific services.
var InjectorRoot = do.New(BasePackage, RootPackage)
