package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/controller"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/interactor"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/samber/do/v2"
)

// DoctorPackage groups all services specific to the doctor command.
var DoctorPackage = do.Package(
	// adapter/controller
	do.Lazy(controller.NewDoctorController),

	// usecase/port
	do.Lazy(func(i do.Injector) (port.DoctorUseCase, error) {
		return interactor.NewDoctorInteractor(i)
	}),

	// adapter/presenter
	do.Lazy(presenter.NewDoctorPresenter),
)

// InjectorDoctor provides a DI container with both base and doctor-specific services.
var InjectorDoctor = do.New(BasePackage, DoctorPackage)
