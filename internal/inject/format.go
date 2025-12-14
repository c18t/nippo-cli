package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// FormatPackage groups all services specific to the format command.
// Services are lazily initialized when first requested.
var FormatPackage = do.Package(
	// adapter/controller
	do.Lazy(controller.NewFormatController),

	// usecase/port
	do.Lazy(port.NewFormatUseCaseBus),

	// usecase/interactor
	do.Lazy(interactor.NewFormatCommandInteractor),

	// adapter/presenter
	do.Lazy(presenter.NewFormatCommandPresenter),
)

// InjectorFormat provides a DI container with both base and format-specific services.
var InjectorFormat = do.New(BasePackage, FormatPackage)
