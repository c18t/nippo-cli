# Project Context

## Purpose

nippo-cli is a command-line tool that powers the author's "nippo" (daily report)
workflow. It provides functionality to initialize project settings, build static
sites from markdown content, and deploy them to Google Drive.

## Tech Stack

- **Language**: Go 1.25.3
- **CLI Framework**: [spf13/cobra](https://github.com/spf13/cobra)
- **Configuration**: [spf13/viper](https://github.com/spf13/viper)
- **Dependency Injection**: [samber/do/v2](https://github.com/samber/do) v2.0.0
- **Markdown Processing**: [gomarkdown/markdown](https://github.com/gomarkdown/markdown)
- **RSS Generation**: [gorilla/feeds](https://github.com/gorilla/feeds)
- **External APIs**: Google OAuth2, Google Drive API
- **Development Tools**: mise, pre-commit, Docker/devcontainer

## Project Conventions

### Architecture Patterns

The project follows **Clean Architecture** with three distinct layers:

1. **Adapter Layer** (`internal/adapter/`)
   - `controller/` - CLI command handlers (Cobra commands)
   - `gateway/` - External service integrations (file providers)
   - `presenter/` - Output formatting and user interaction

2. **Domain Layer** (`internal/domain/`)
   - `logic/service/` - Business logic services
   - `logic/repository/` - Data access interfaces

3. **Usecase Layer** (`internal/usecase/`)
   - `interactor/` - Application use cases
   - `port/` - Interface definitions for use cases

### Dependency Injection

- Uses `samber/do/v2` for DI container management
- **Base injector** (`internal/inject/container.go`): Defines `BasePackage`
  using `do.Package()` pattern with lazily initialized services. The singleton
  injector is created via `GetInjector()` with `sync.Once` for thread-safety.
- **Package pattern**: Services are grouped with `do.Lazy()` wrappers for lazy
  initialization. Services are created only when first requested, improving
  startup performance.
- **Command injectors** (`internal/inject/{command}.go`): Each command defines
  a `{Command}Package` using `do.Package()` with command-specific services.
  The command injector is created via `do.New(BasePackage, CommandPackage)` to
  compose base and command services into a single container.
- **Constructor pattern**: All constructors use `do.Invoke` (not `MustInvoke`)
  and return `(T, error)` for proper error handling

#### Lifecycle Management

- **Configuration Management**: Config is loaded and initialized via
  `core.InitConfig(configFile)` during Cobra initialization hook. The global
  `core.Cfg` variable is set once at startup and accessed directly by all
  services. This global singleton pattern is used because:
  - Config initialization must occur after flag parsing but before services run
  - DI-managed config would require immediate evaluation in constructors (before
    InitConfig runs)
  - Deferring config access to methods introduces service locator anti-pattern
  - Application-wide singleton config has no per-request variation
- **Graceful Shutdown**: Services implementing resource cleanup provide a
  `Shutdown() error` method. The DI container's `Shutdown()` method is called in
  `cmd.Execute()` via defer to ensure resources are properly released on
  application exit.
- **Health Checks**: Critical services (e.g., `DriveFileProvider`) implement
  `HealthCheck() error` to verify initialization and connectivity. Health checks
  can be invoked using `do.HealthCheck[ServiceType](injector)`.
- **Interface-based DI**: Use case interactors return interface types
  (`port.*UseCase`) to enable easy mocking and testing.

### Code Style

- Follow Go conventions and idioms
- Use `gofmt` for formatting
- Constructor naming: `New{ServiceName}(injector do.Injector) (Service, error)`
- Error wrapping: Use `fmt.Errorf` with `%w` for error chain inspection
- Unused parameters: Use blank identifier `_` if parameter is required by
  interface but not used

### Testing Strategy

- Run tests with race detector: `go run -race . --help`
- Verify DI container initialization and scope isolation
- Test error handling paths in constructors
- Integration tests for each command (`init`, `build`, `deploy`, `clean`,
  `update`)

### Git Workflow

- **Branching**: Feature branches following `feature/{issue-number}_{description}`
  pattern
- **Commits**: Conventional Commits specification
  - Format: `type(scope): description` (title in English, ~80 chars)
  - Body: Japanese, bulleted list
  - Types: `feat`, `fix`, `refactor`, `chore`, `docs`, `test`, etc.
- **Pre-commit hooks**: Enforces formatting, linting (markdownlint, prettier,
  golangci-lint, shellcheck)
- **Pull Requests**: Use `.github/pull_request_template.md` template

## Domain Context

### Commands

The CLI provides five main commands:

1. **init**: Initialize project settings and Google Drive authentication
2. **build**: Build static site from markdown content
3. **deploy**: Publish built content to Google Drive
4. **clean**: Clean build artifacts
5. **update**: Update existing content

### Development Commands

- **Setup**: `mise run setup` - Install all dependencies (mise tools, go modules,
  pnpm packages, MCP servers, pre-commit hooks)
- **Build**: `make` or `mise run build` - Compile the CLI application
- **Run**: `go run nippo/nippo.go {command}` - Execute during development
- **Release**: `mise run release` - Build release binaries with goreleaser

### OpenSpec Change Management

The project uses OpenSpec for managing significant changes:

- **Proposals**: Create in `openspec/changes/{change-name}/proposal.md`
- **Specs**: Define requirements in
  `openspec/changes/{change-name}/specs/{spec-name}/spec.md`
- **Archive**: Completed changes move to
  `openspec/changes/archive/{date}-{change-name}/`
- **Authoritative specs**: Published to `openspec/specs/{spec-name}/spec.md`

Refer to `@openspec/AGENTS.md` for detailed workflow.

## Important Constraints

- **Go version**: Must use Go 1.25.3 or compatible
- **DI pattern**: Always use `GetInjector().Clone()` for command-specific
  injectors to maintain scope isolation
- **Error handling**: Never use `do.MustInvoke` - always use `do.Invoke` with
  proper error propagation
- **Markdown line length**: 80 characters (enforced by markdownlint for OpenSpec
  files, excluding code blocks)
- **Google Drive integration**: Requires OAuth2 authentication flow during `init`

## External Dependencies

- **Google Drive API**: Used for storing and deploying generated content
- **Google OAuth2**: Handles user authentication for Drive access
- **Container runtime**: Docker for devcontainer development environment
- **mise**: Tool version management (replaces asdf/rtx)
- **pre-commit**: Git hook management for code quality checks
