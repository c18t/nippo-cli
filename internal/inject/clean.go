package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// CleanPackage groups all services specific to the clean command.
// Services are lazily initialized when first requested.
var CleanPackage = do.Package(
	// adapter/controller
	do.Lazy(controller.NewCleanController),

	// usecase/port
	do.Lazy(port.NewCleanUseCaseBus),

	// usecase/interactor
	do.Lazy(interactor.NewCleanCommandInteractor),

	// adapter/presenter
	do.Lazy(presenter.NewCleanCommandPresenter),
)

// InjectorClean provides a DI container with both base and clean-specific services.
var InjectorClean = do.New(BasePackage, CleanPackage)
