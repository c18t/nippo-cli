# dependency-injection Specification

## Purpose

TBD - created by archiving change refactor-di-implementation. Update Purpose
after archive.

## Requirements

### Requirement: Lazy Initialization Pattern

The base DI container SHALL use lazy initialization with `sync.Once` to defer
resource allocation until the injector is first accessed, improving startup
performance and ensuring thread safety.

#### Scenario: Injector initialized on first access

- **GIVEN** the package is loaded
- **WHEN** `GetInjector()` is called for the first time
- **THEN** it SHALL initialize the DI container using `sync.Once`
- **AND** subsequent calls SHALL return the same instance
- **AND** no initialization SHALL occur during package load time

#### Scenario: Thread-safe concurrent initialization

- **GIVEN** multiple goroutines call `GetInjector()` concurrently
- **WHEN** the injector has not been initialized
- **THEN** `sync.Once` SHALL ensure initialization happens exactly once
- **AND** all goroutines SHALL receive the same injector instance
- **AND** no race conditions SHALL occur

#### Scenario: Deferred resource allocation

- **GIVEN** the application imports the inject package
- **WHEN** the injector is never accessed
- **THEN** no DI container SHALL be created
- **AND** no memory SHALL be allocated for services
- **AND** startup time SHALL not be impacted by unused dependencies

### Requirement: Command Scope Isolation

Each command-specific injector SHALL create an isolated dependency scope by
cloning the base injector, preventing pollution of shared dependencies.

#### Scenario: Clean command uses isolated scope

- **GIVEN** the base injector is accessible via `GetInjector()`
- **WHEN** `AddCleanProvider()` is called
- **THEN** it SHALL clone the base injector using `GetInjector().Clone()`
- **AND** register command-specific dependencies in the cloned scope
- **AND** return the cloned scope without modifying the base injector

#### Scenario: All command injectors follow consistent pattern

- **GIVEN** multiple command injectors (`init`, `build`, `deploy`, `clean`,
  `update`, `root`)
- **WHEN** any command injector is initialized
- **THEN** each SHALL use `GetInjector().Clone()` to create an isolated scope
- **AND** follow the same provider registration pattern
- **AND** return the isolated scope

### Requirement: Base Injector Immutability

The base `Injector` SHALL remain unmodified after initialization to ensure
consistent shared service state across all command scopes.

#### Scenario: Base injector provides shared services only

- **GIVEN** `GetInjector()` is called for the first time
- **WHEN** the base injector is lazily initialized via `AddProvider()`
- **THEN** only gateway, repository, and domain services SHALL be registered
- **AND** no command-specific services SHALL be registered
- **AND** the injector SHALL be available for cloning by command-specific providers

#### Scenario: Command-specific services do not pollute base scope

- **GIVEN** a command-specific injector clones the base injector
- **WHEN** command-specific services are registered
- **THEN** these services SHALL only exist in the cloned scope
- **AND** the base `Injector` SHALL not contain command-specific services
- **AND** other command scopes SHALL not have access to these services

### Requirement: Proper Error Handling in Constructors

All service constructors SHALL use `do.Invoke` instead of `do.MustInvoke` and
properly propagate errors with context, ensuring graceful failure instead of
panics.

#### Scenario: Constructor uses Invoke with error handling

- **GIVEN** a constructor needs to resolve a dependency
- **WHEN** the constructor calls `do.Invoke[T](injector)`
- **THEN** it SHALL check the returned error
- **AND** if error is not nil, SHALL return nil and wrapped error with context
- **AND** if successful, SHALL use the resolved dependency
- **AND** SHALL NOT use `do.MustInvoke` which can panic

#### Scenario: Dependency resolution failure is gracefully handled

- **GIVEN** a service dependency cannot be resolved
- **WHEN** `do.Invoke` is called in a constructor
- **THEN** an error SHALL be returned with descriptive context
- **AND** the error SHALL be propagated up the call stack
- **AND** the application SHALL NOT panic
- **AND** the error message SHALL identify which dependency failed

#### Scenario: Error messages provide context

- **GIVEN** a constructor fails to resolve a dependency
- **WHEN** an error is returned
- **THEN** the error message SHALL include the constructor name
- **AND** SHALL include the dependency type that failed
- **AND** SHALL wrap the original error using `fmt.Errorf` with `%w`
- **AND** SHALL enable error chain inspection

### Requirement: Constructor Dependency Declaration

Constructors SHALL only accept `do.Injector` parameter if they actually use it
to resolve dependencies, ensuring clear dependency graphs.

#### Scenario: Constructor uses Injector parameter

- **GIVEN** a constructor accepts `do.Injector` parameter
- **WHEN** the constructor is implemented
- **THEN** it SHALL call `do.Invoke` at least once on the injector
- **AND** SHALL resolve all dependencies through the injector
- **AND** SHALL NOT leave the parameter unused

#### Scenario: Constructor without dependencies

- **GIVEN** a service has no dependencies
- **WHEN** defining the constructor
- **THEN** it SHALL NOT accept `do.Injector` parameter
- **AND** SHALL return the service instance directly
- **AND** the signature SHALL be `func NewService() (Service, error)`

#### Scenario: Unused Injector parameters are removed

- **GIVEN** existing constructors with unused `do.Injector` parameters
- **WHEN** refactoring the DI implementation
- **THEN** the parameter SHALL be removed if not used
- **AND** the constructor signature SHALL be simplified
- **AND** registration in `do.Provide` SHALL be updated accordingly
