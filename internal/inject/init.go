package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/adapter/presenter/view"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// InitPackage groups all services specific to the init command.
// Services are lazily initialized when first requested.
var InitPackage = do.Package(
	// adapter/controller
	do.Lazy(controller.NewInitController),

	// usecase/port
	do.Lazy(port.NewInitUseCaseBus),

	// usecase/interactor
	do.Lazy(interactor.NewInitSettingInteractor),

	// adapter/presenter
	do.Lazy(presenter.NewInitSettingPresenter),

	// adapter/presenter/view
	do.Lazy(view.NewInitViewProvider),
	do.Lazy(view.NewConfigureProjectView),
)

// InjectorInit provides a DI container with both base and init-specific services.
var InjectorInit = do.New(BasePackage, InitPackage)
