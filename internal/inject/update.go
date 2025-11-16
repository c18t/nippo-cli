package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// UpdatePackage groups all services specific to the update command.
// Services are lazily initialized when first requested.
var UpdatePackage = do.Package(
	// adapter/controller
	do.Lazy(controller.NewUpdateController),

	// usecase/port
	do.Lazy(port.NewUpdateUseCaseBus),

	// usecase/interactor
	do.Lazy(interactor.NewUpdateCommandInteractor),

	// adapter/presenter
	do.Lazy(presenter.NewUpdateCommandPresenter),
)

// InjectorUpdate provides a DI container with both base and update-specific services.
var InjectorUpdate = do.New(BasePackage, UpdatePackage)
