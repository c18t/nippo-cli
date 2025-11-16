# code-organization Specification

## Purpose

Define code organization conventions for the dependency injection layer, ensuring clear file naming and structure that improves discoverability and maintainability.

## Requirements

### Requirement: DI Container File Naming

The base DI container file SHALL use a descriptive name that clearly indicates its purpose, avoiding numeric prefixes and following Go naming conventions.

#### Scenario: Base injector file is named container.go

- **GIVEN** the DI container needs to be initialized with shared services
- **WHEN** a developer navigates to the `internal/inject/` directory
- **THEN** they SHALL find a file named `container.go` that contains the base injector factory
- **AND** the file SHALL export a `GetInjector()` function that returns the singleton DI container
- **AND** the filename SHALL clearly indicate it contains DI container setup without needing numeric prefixes

#### Scenario: File naming avoids conflicts with command names

- **GIVEN** command-specific injectors are named after their commands (e.g., `build.go`, `clean.go`)
- **WHEN** the base injector file is named `container.go`
- **THEN** there SHALL be no naming conflicts with existing or future command files
- **AND** the name SHALL not require numeric prefixes like `000_` for disambiguation

#### Scenario: Documentation references correct filename

- **GIVEN** project documentation describes the DI architecture
- **WHEN** developers read `openspec/project.md` or other documentation
- **THEN** references to the base injector SHALL use `internal/inject/container.go`
- **AND** SHALL NOT reference legacy filenames like `000_inject.go`
