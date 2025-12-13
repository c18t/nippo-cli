package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// AuthPackage groups all services specific to the auth command.
// Services are lazily initialized when first requested.
var AuthPackage = do.Package(
	// adapter/controller
	do.Lazy(controller.NewAuthController),

	// usecase/port
	do.Lazy(port.NewAuthUseCaseBus),

	// usecase/interactor
	do.Lazy(interactor.NewAuthInteractor),

	// adapter/presenter
	do.Lazy(presenter.NewAuthPresenter),
)

// InjectorAuth provides a DI container with both base and auth-specific services.
var InjectorAuth = do.New(BasePackage, AuthPackage)
