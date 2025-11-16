# AGENTS.md

This file provides guidance to AI assistants (Claude Code and other agents) when
working with code in this repository.

> **For comprehensive project context**, see `openspec/project.md` which
> contains detailed information about:
>
> - Project purpose and goals
> - Complete tech stack and dependencies
> - Architecture patterns and conventions
> - Testing strategy and constraints
> - Domain context and external dependencies

This file focuses on **practical commands and development workflows**.

<!-- OPENSPEC:START -->

## OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:

- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big
  performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:

- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

## Development Environment

This project uses Dev Containers (Codespaces or VS Code). All development tools
are managed via [mise](https://mise.jdx.dev/).

### container-use

ALWAYS use ONLY Environments for ANY and ALL file, code, or shell operations—NO
EXCEPTIONS—even for simple or generic requests.

DO NOT install or use the git cli with the environment_run_cmd tool. All
environment tools will handle git operations for you. Changing ".git" yourself
will compromise the integrity of your environment.

You MUST inform the user how to view your work using
`container-use log <env_id>` AND `container-use checkout <env_id>`.
Failure to do this will make your work inaccessible to others.

### Initial Setup

```bash
mise trust
mise run setup    # Installs tools, go modules, and pre-commit hooks
```

### Essential Commands

**Building:**

```bash
mise run build    # or: mise run b
make              # Direct make invocation
make debug        # Build with debug symbols to /tmp
```

**Testing:**

```bash
# Run all tests
go test ./...

# Run specific package tests with verbose output
go test -v ./path/to/package

# Run tests with coverage
go test -cover ./...

# Verify build before committing
mise run build
```

**Code Quality:**

**IMPORTANT**: When you modify files, you MUST run the appropriate linters to
check for errors before committing:

**For Go files** (`**/*.go`):

```bash
staticcheck ./...
```

This ensures that Go code follows best practices and maintains consistency
across the project.

**For Markdown files** (any `.md` file):

```bash
markdownlint-cli2 "**/*.md"
```

This validates Markdown formatting and ensures documentation consistency.

**Run all linters:**

```bash
pre-commit run --all-files  # Run all linters and formatters
```

**Release:**

```bash
mise run release  # Build release binaries with goreleaser (snapshot mode)
```

## Architecture

This is a Clean Architecture CLI boilerplate using Cobra for command handling
and samber/do for dependency injection.

### Layer Structure

```text
cmd/                    # Cobra command definitions
  └── *_invoker.go      # Command invokers that wire DI to Cobra RunE
internal/
  ├── adapter/
  │   ├── controller/   # Adapts Cobra commands to use cases
  │   └── presenter/    # Formats use case output
  ├── usecase/
  │   ├── port/         # Use case interfaces and DTOs
  │   └── interactor/   # Use case implementations
  ├── inject/           # DI container setup per command
  └── core/             # Shared types (Controller, UseCase interfaces)
```

### Data Flow

1. **cmd/\*\_invoker.go** creates controller via DI
2. **Controller** receives Cobra command/args, calls UseCase via Bus
3. **UseCaseBus** routes input to appropriate Interactor
4. **Interactor** processes logic, sends output to Presenter
5. **Presenter** formats and displays result

### Adding New Commands

Use scaffdog to generate boilerplate:

```bash
cobra-cli add <command-name>
scaffdog generate command --answer "name:<command-name>" --answer "usecase:command"
```

This generates all layers:

- `cmd/<command>_invoker.go` - DI wiring
- `internal/adapter/controller/<command>.go` - Controller
- `internal/usecase/port/<command>.go` - Port interfaces and DTOs
- `internal/usecase/interactor/<command>.go` - Use case implementation
- `internal/adapter/presenter/<command>.go` - Output presenter
- `internal/inject/<command>.go` - DI configuration

Wire the command in `cmd/<command>.go`:

```go
func init() {
    commandCmd.RunE = createCommandCommand()
    rootCmd.AddCommand(commandCmd)
}
```

### Dependency Injection

Each command has its own package that composes with the base package:

```go
// internal/inject/command.go
var CommandPackage = do.Package(
    do.Lazy(controller.NewCommandController),
    do.Lazy(port.NewCommandUseCaseBus),
    do.Lazy(interactor.NewCommandInteractor),
    do.Lazy(presenter.NewCommandPresenter),
)

var InjectorCommand = do.New(BasePackage, CommandPackage)
```

Base services go in `internal/inject/container.go` using the `BasePackage`
pattern. All services use `do.Lazy()` for lazy initialization.

For testing, use `inject.NewTestInjector()` to create test injectors with mock
services:

```go
// Test with mocked services
injector := inject.NewTestInjector(&inject.TestBasePackageOptions{
    DriveFileProvider: mockDriveProvider,
    Config: testConfig,
})
```

### Key Patterns

- **UseCaseBus**: Type-switches input DTOs to route to correct use case
- **Presenter**: Handles both success (Complete) and error (Suspend) cases
- **Controller Params**: Struct to hold command flags/parameters
- **Cobra Integration**: Controllers implement `core.Controller` interface
  with `Exec(cmd, args) error` method

## Build Configuration

- **Makefile**: Uses git describe for version, builds to `bin/` directory
- **.goreleaser.yaml**: Cross-platform release builds (Linux, Windows, macOS)
- **Version injection**: `main.version` is set via ldflags during build
