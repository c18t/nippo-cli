package interactor

import (
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/usecase/port"
)

type updateProjectDataInteractor struct {
	presenter presenter.UpdateProjectDataPresenter
}

func NewUpdateProjectDataInteractor(presenter presenter.UpdateProjectDataPresenter) port.UpdateProjectDataUsecase {
	return &updateProjectDataInteractor{presenter}
}

func (u *updateProjectDataInteractor) Handle(input *port.UpdateProjectDataUsecaseInputData) {
	output := &port.UpdateProjectDataUsecaseOutputData{}

	output.Message = "update project files... "
	u.presenter.Progress(output)

	err := downloadProject()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	output.Message = "ok."
	u.presenter.Complete(output)
}
