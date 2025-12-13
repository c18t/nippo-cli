package port

import (
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/samber/do/v2"
)

type AuthUseCaseInputData struct{}
type AuthUseCaseOutputData struct {
	Message string
}

type AuthUseCase interface {
	core.UseCase
	Handle(input *AuthUseCaseInputData)
}

type AuthUseCaseBus interface {
	Handle(input *AuthUseCaseInputData)
}

type authUseCaseBus struct {
	auth AuthUseCase
}

func NewAuthUseCaseBus(i do.Injector) (AuthUseCaseBus, error) {
	auth, err := do.Invoke[AuthUseCase](i)
	if err != nil {
		return nil, err
	}
	return &authUseCaseBus{
		auth: auth,
	}, nil
}

func (bus *authUseCaseBus) Handle(input *AuthUseCaseInputData) {
	bus.auth.Handle(input)
}
