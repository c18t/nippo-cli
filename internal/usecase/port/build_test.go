package port

import (
	"testing"

	"github.com/samber/do/v2"
)

// Mock BuildCommandUseCase
type mockBuildCommandUseCase struct {
	handleCalled bool
}

func (m *mockBuildCommandUseCase) Handle(input *BuildCommandUseCaseInputData) {
	m.handleCalled = true
}

func TestNewBuildUseCaseBus(t *testing.T) {
	mock := &mockBuildCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (BuildCommandUseCase, error) {
		return mock, nil
	})

	bus, err := NewBuildUseCaseBus(injector)
	if err != nil {
		t.Errorf("NewBuildUseCaseBus() error = %v", err)
	}
	if bus == nil {
		t.Error("NewBuildUseCaseBus() returned nil")
	}
}

func TestBuildUseCaseBus_Handle(t *testing.T) {
	mock := &mockBuildCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (BuildCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewBuildUseCaseBus(injector)
	bus.Handle(&BuildCommandUseCaseInputData{})

	if !mock.handleCalled {
		t.Error("Handle() did not call command.Handle()")
	}
}

func TestBuildUseCaseBus_Handle_Panic(t *testing.T) {
	mock := &mockBuildCommandUseCase{}
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (BuildCommandUseCase, error) {
		return mock, nil
	})

	bus, _ := NewBuildUseCaseBus(injector)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Handle() should panic for unknown input type")
		}
	}()

	// Pass an unknown type to trigger panic
	bus.Handle("unknown type")
}

func TestBuildCommandUseCaseOutputData(t *testing.T) {
	output := &BuildCommandUseCaseOutputData{
		Message: "test message",
	}

	if output.Message != "test message" {
		t.Errorf("Message = %q, want %q", output.Message, "test message")
	}
}
