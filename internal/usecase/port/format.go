package port

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/samber/do/v2"
)

type FormatUseCaseInputData interface{}
type FormatUseCaseOutputData interface{}

type FormatCommandUseCaseInputData struct {
	FormatUseCaseInputData
}

// FormatFileStatus represents the result status of processing a file
type FormatFileStatus int

const (
	FormatFileStatusSuccess   FormatFileStatus = iota // Successfully updated
	FormatFileStatusNoChange                          // No changes needed
	FormatFileStatusFailed                            // Failed to process/upload
)

type FormatCommandUseCaseOutputData struct {
	FormatUseCaseOutputData
	Message  string
	Filename string
	FileId   string
	Status   FormatFileStatus
	Error    error
}

type FormatCommandUseCase interface {
	core.UseCase
	Handle(input *FormatCommandUseCaseInputData)
}

type FormatUseCaseBus interface {
	Handle(input FormatUseCaseInputData)
}
type formatUseCaseBus struct {
	command FormatCommandUseCase `do:""`
}

func NewFormatUseCaseBus(i do.Injector) (FormatUseCaseBus, error) {
	command, err := do.Invoke[FormatCommandUseCase](i)
	if err != nil {
		return nil, err
	}
	return &formatUseCaseBus{
		command: command,
	}, nil
}

func (bus *formatUseCaseBus) Handle(input FormatUseCaseInputData) {
	switch data := input.(type) {
	case *FormatCommandUseCaseInputData:
		bus.command.Handle(data)
	default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}
