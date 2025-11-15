package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/adapter/presenter/view"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

var InjectorInit = AddInitProvider()

func AddInitProvider() *do.RootScope {
	var i = GetInjector().Clone()

	// adapter/controller
	do.Provide(i, controller.NewInitController)

	// usecase/port
	do.Provide(i, port.NewInitUseCaseBus)

	// usecase/intractor
	do.Provide(i, interactor.NewInitSettingInteractor)
	do.Provide(i, interactor.NewInitSaveDriveTokenInteractor)

	// adapter/presenter
	do.Provide(i, presenter.NewInitSettingPresenter)
	do.Provide(i, presenter.NewInitSaveDriveTokenPresenter)

	// adapter/presenter/view
	do.Provide(i, view.NewInitViewProvider)
	do.Provide(i, view.NewConfigureProjectView)

	return i
}
