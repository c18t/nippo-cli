package presenter

import (
	"errors"
	"testing"

	"github.com/samber/do/v2"
)

func TestNewConsolePresenter(t *testing.T) {
	injector := do.New()
	presenter, err := NewConsolePresenter(injector)
	if err != nil {
		t.Errorf("NewConsolePresenter() error = %v", err)
	}
	if presenter == nil {
		t.Error("NewConsolePresenter() returned nil")
	}
}

func TestConsolePresenter_Progress(t *testing.T) {
	injector := do.New()
	p, _ := NewConsolePresenter(injector)

	// Test that Progress doesn't panic
	p.Progress("Testing progress...")
	p.StopProgress()
}

func TestConsolePresenter_StopProgress(t *testing.T) {
	injector := do.New()
	p, _ := NewConsolePresenter(injector)

	// Test that StopProgress doesn't panic when no progress is running
	p.StopProgress()
}

func TestConsolePresenter_Warning(t *testing.T) {
	injector := do.New()
	p, _ := NewConsolePresenter(injector)

	// Test that Warning doesn't panic
	p.Warning(errors.New("test warning"))
}

func TestConsolePresenter_Complete(t *testing.T) {
	injector := do.New()
	p, _ := NewConsolePresenter(injector)

	// Test that Complete doesn't panic
	p.Complete("Test completed")
}

func TestConsolePresenter_IsCancelled(t *testing.T) {
	injector := do.New()
	p, _ := NewConsolePresenter(injector)

	// Test that IsCancelled returns a boolean
	result := p.IsCancelled()
	if result != false {
		t.Errorf("IsCancelled() = %v, want false initially", result)
	}
}
