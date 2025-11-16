---
name: 'command'
root: '.'
output: '.'
questions:
  name: 'enter new command name:'
  usecase:
    message: 'enter a list of subcommand names, separated by spaces or commas[\s,]:'
    initial: 'command'
---

# Variables

- command_camel: `{{ inputs.name | camel }}`
- command_pascal: `{{ inputs.name | pascal }}`
- command_snake: `{{ inputs.name | snake }}`
- subcommand_list: `{{ inputs.usecase | replace "\s" "," }}`

# `cmd/{{ command_snake }}_invoker.go`

```go
package cmd

import (
	"{{ go_module }}/internal/adapter/controller"
	"{{ go_module }}/internal/core"
	"{{ go_module }}/internal/inject"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

func create{{ command_pascal }}Command() core.RunEFunc {
	ctrl, err := do.Invoke[controller.{{ command_pascal }}Controller](inject.Injector{{ command_pascal }})
	cobra.CheckErr(err)
	return ctrl.Exec
}

```

# `internal/adapter/controller/{{ command_snake }}.go`

```go
package controller

import (
	"{{ go_module }}/internal/core"
	"{{ go_module }}/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

type {{ command_pascal }}Params struct {
}

type {{ command_pascal }}Controller interface {
	core.Controller
	Params() *{{ command_pascal }}Params
}

type {{ command_camel }}Controller struct {
	bus    port.{{ command_pascal }}UseCaseBus `do:""`
	params *{{ command_pascal }}Params
}

func New{{ command_pascal }}Controller(i do.Injector) ({{ command_pascal }}Controller, error) {
	bus, err := do.Invoke[port.{{ command_pascal }}UseCaseBus](i)
	if err != nil {
		return nil, err
	}
	return &{{ command_snake }}Controller{
		bus:    bus,
		params: &{{ command_pascal }}Params{},
	}, nil
}

func (c *{{ command_snake }}Controller) Params() *{{ command_pascal }}Params {
	return c.params
}

func (c *{{ command_snake }}Controller) Exec(cmd *cobra.Command, args []string) (err error) {
	{{ for subcommand in subcommand_list | split ',' -}}
	{{ prefix := command_pascal + (subcommand | trim | pascal) -}}
	c.bus.Handle(&port.{{ prefix }}UseCaseInputData{})
	{{ end }}return
}

```

# `internal/usecase/port/{{ command_snake }}.go`

```go
package port

import (
	"fmt"

	"{{ go_module }}/internal/core"
	"github.com/samber/do/v2"
)

type {{ command_pascal }}UseCaseInputData interface{}
type {{ command_pascal }}UseCaseOutputData interface{}

{{ for subcommand in (subcommand_list | split ',' )-}}
{{ prefix := command_pascal + (subcommand | trim | pascal) -}}
type {{ prefix }}UseCaseInputData struct {
	{{ command_pascal }}UseCaseInputData
}
type {{ prefix }}UseCaseOutputData struct {
	{{ command_pascal }}UseCaseOutputData
	Message string
}
type {{ prefix }}UseCase interface {
	core.UseCase
	Handle(input *{{ prefix }}UseCaseInputData)
}
{{ end }}
type {{ command_pascal }}UseCaseBus interface {
	Handle(input {{ command_pascal }}UseCaseInputData)
}
type {{ command_camel }}UseCaseBus struct {
	{{ for subcommand in (subcommand_list | split ',') -}}
	{{ prefix := command_pascal + (subcommand | trim | pascal) -}}
	{{ subcommand | trim | camel }} {{ prefix }}UseCase `do:""`
	{{ end }}
}

func New{{ command_pascal }}UseCaseBus(i do.Injector) ({{ command_pascal }}UseCaseBus, error) {
	{{ for subcommand in (subcommand_list | split ',') -}}
	{{ prefix := command_pascal + (subcommand | trim | pascal) -}}
	{{ subcommand | trim | camel }}, err := do.Invoke[{{ prefix }}UseCase](i)
	if err != nil {
		return nil, err
	}
	{{ end }}return &{{ command_camel }}UseCaseBus{
		{{ for subcommand in (subcommand_list | split ',') -}}
		{{ subcommand | trim | camel }}: {{ subcommand | trim | camel }},
		{{ end }}}, nil
}

func (bus *{{ command_camel }}UseCaseBus) Handle(input {{ command_pascal }}UseCaseInputData) {
	switch data := input.(type) {
	{{ for subcommand in subcommand_list | split ',' -}}
	{{ prefix := command_pascal + (subcommand | trim | pascal) -}}
	case *{{ prefix }}UseCaseInputData:
		bus.{{ subcommand | trim | camel }}.Handle(data)
	{{ end }}default:
		panic(fmt.Errorf("handler for '%T' is not implemented", data))
	}
}

```

# `internal/usecase/interactor/{{ command_snake }}.go`

```go
package interactor

import (
	"{{ go_module }}/internal/adapter/presenter"
	"{{ go_module }}/internal/usecase/port"
	"github.com/samber/do/v2"
)

{{ for subcommand in (subcommand_list | split ',') -}}
{{ prefix_pascal := command_pascal + (subcommand | trim | pascal) -}}
{{ prefix_camel := prefix_pascal | camel }}
type {{ prefix_camel }}Interactor struct {
	presenter presenter.{{ prefix_pascal }}Presenter `do:""`
}

func New{{ prefix_pascal }}Interactor(i do.Injector) (port.{{ prefix_pascal }}UseCase, error) {
	p, err := do.Invoke[presenter.{{ prefix_pascal }}Presenter](i)
	if err != nil {
		return nil, err
	}
	return &{{ prefix_camel }}Interactor{
		presenter: p,
	}, nil
}

func (u *{{ prefix_camel }}Interactor) Handle(input *port.{{ prefix_pascal }}UseCaseInputData) {
	output := &port.{{ prefix_pascal }}UseCaseOutputData{}
	output.Message = "{{ command_snake }} {{ subcommand | snake }} called."
	u.presenter.Complete(output)
}
{{ end }}
```

# `internal/adapter/presenter/{{ command_snake }}.go`

```go
package presenter

import (
	"fmt"

	"{{ go_module }}/internal/usecase/port"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

{{ for subcommand in (subcommand_list | split ',') -}}
{{ prefix_pascal := command_pascal + (subcommand | trim | pascal) -}}
{{ prefix_camel := prefix_pascal | camel }}
type {{ prefix_pascal }}Presenter interface {
	Complete(output *port.{{ prefix_pascal }}UseCaseOutputData)
	Suspend(err error)
}

type {{ prefix_camel }}Presenter struct {
}

func New{{ prefix_pascal }}Presenter(i do.Injector) ({{ prefix_pascal }}Presenter, error) {
	return &{{ prefix_camel }}Presenter{}, nil
}

func (p *{{ prefix_camel }}Presenter) Complete(output *port.{{ prefix_pascal }}UseCaseOutputData) {
	fmt.Printf("%v\n", output)
}

func (p *{{ prefix_camel }}Presenter) Suspend(err error) {
	cobra.CheckErr(err)
}

{{ end }}
```

# `internal/inject/{{ command_snake }}.go`

```go
package inject

import (
	"{{ go_module }}/internal/adapter/controller"
	"{{ go_module }}/internal/adapter/presenter"
	"{{ go_module }}/internal/usecase/interactor"
	"{{ go_module }}/internal/usecase/port"
	"github.com/samber/do/v2"
)

// {{ command_pascal }}Package groups all services specific to the {{ command_snake }} command.
// Services are lazily initialized when first requested.
var {{ command_pascal }}Package = do.Package(
	// adapter/controller
	do.Lazy(controller.New{{ command_pascal }}Controller),

	// usecase/port
	do.Lazy(port.New{{ command_pascal }}UseCaseBus),

	// usecase/interactor
	{{ for subcommand in (subcommand_list | split ',') -}}
	{{ prefix_pascal := command_pascal + (subcommand | trim | pascal) -}}
	do.Lazy(interactor.New{{ prefix_pascal }}Interactor),
	{{ end }}
	// adapter/presenter
	{{ for subcommand in (subcommand_list | split ',') -}}
	{{ prefix_pascal := command_pascal + (subcommand | trim | pascal) -}}
	do.Lazy(presenter.New{{ prefix_pascal }}Presenter),
	{{ end }}
)

// Injector{{ command_pascal }} provides a DI container with both base and {{ command_snake }}-specific services.
var Injector{{ command_pascal }} = do.New(BasePackage, {{ command_pascal }}Package)

```
