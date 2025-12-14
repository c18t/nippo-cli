package port

import "github.com/c18t/nippo-cli/internal/core"

type DoctorUseCaseInputData struct{}

type DoctorUseCaseOutputData struct {
	Checks []DoctorCheck
}

type DoctorCheck struct {
	Category   string
	Item       string
	Status     DoctorCheckStatus
	Message    string
	Suggestion string
}

type DoctorCheckStatus int

const (
	DoctorCheckStatusPass DoctorCheckStatus = iota
	DoctorCheckStatusFail
	DoctorCheckStatusWarn
)

type DoctorUseCase interface {
	core.UseCase
	Handle(input *DoctorUseCaseInputData)
}
