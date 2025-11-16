# dependency-injection Specification

## Purpose

Define dependency injection patterns and lifecycle management for the application using samber/do v2, ensuring services are properly initialized, scoped, and cleaned up while maintaining testability and modularity.

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

### Requirement: Graceful Shutdown

Services that manage external resources SHALL implement the `Shutdowner`
interface and be cleanly shut down when the application terminates, ensuring
proper resource cleanup.

#### Scenario: Service implements shutdown interface

- **GIVEN** a service manages external resources (e.g., Google Drive API
  client, file handles)
- **WHEN** the service is registered in the DI container
- **THEN** it SHALL implement one of the shutdown interfaces (`Shutdowner`,
  `ShutdownerWithContext`, `ShutdownerWithError`,
  `ShutdownerWithContextAndError`)
- **AND** the shutdown method SHALL release all held resources
- **AND** SHALL return any cleanup errors to the caller

#### Scenario: Application shutdown triggers service cleanup

- **GIVEN** services are registered with shutdown implementations
- **WHEN** the application receives a termination signal (SIGTERM, SIGINT) or
  completes execution
- **THEN** `do.Shutdown` or `do.ShutdownWithContext` SHALL be called on the DI
  container
- **AND** all services implementing shutdown interfaces SHALL have their cleanup
  methods invoked
- **AND** shutdown SHALL occur in reverse dependency order
- **AND** shutdown errors SHALL be logged and reported

#### Scenario: Shutdown respects context timeout

- **GIVEN** services with `ShutdownerWithContext` implementations
- **WHEN** shutdown is initiated with a context timeout
- **THEN** each service SHALL have the opportunity to clean up within the
  timeout
- **AND** if timeout expires, remaining services SHALL be forcefully terminated
- **AND** timeout expiration SHALL be logged as a warning

### Requirement: Service Aliasing

Services SHALL be bound to their port interfaces using `do.As`, enabling
interface-based dependency injection and improving testability.

#### Scenario: Concrete service registered with interface alias

- **GIVEN** a concrete service implementation and corresponding port interface
- **WHEN** the service is registered in the DI container
- **THEN** `do.As[ConcreteType, InterfaceType]` SHALL be called to create the
  binding
- **AND** the service SHALL be resolvable by both concrete type and interface
  type
- **AND** `do.Invoke[InterfaceType]` SHALL return the registered concrete
  instance

#### Scenario: Multiple implementations of same interface

- **GIVEN** multiple concrete implementations of the same interface
- **WHEN** services are registered with named aliases
- **THEN** `do.AsNamed[Concrete, Interface](injector, "concrete-name",
"interface-name")` SHALL be used
- **AND** each implementation SHALL be resolvable by its named interface alias
- **AND** consumers SHALL specify which implementation to use via the name

#### Scenario: Interface-based dependency resolution

- **GIVEN** a service depends on an interface rather than a concrete type
- **WHEN** the service constructor uses `do.Invoke[InterfaceType](injector)`
- **THEN** the bound concrete implementation SHALL be injected
- **AND** the consuming service SHALL not have direct knowledge of the concrete
  type
- **AND** tests SHALL be able to substitute mock implementations easily

### Requirement: Configuration Management

Application configuration SHALL be managed through a global singleton variable
initialized at application startup, providing pragmatic access to
application-wide configuration.

#### Rationale: Why not DI-managed config

During implementation, we evaluated managing configuration through DI container
using `do.ProvideValue`. However, this approach has fundamental timing issues:

- Services invoking config in constructors causes immediate evaluation before
  `InitConfig()` runs
- Deferring config access to service methods introduces service locator
  anti-pattern
- The resulting code has worse maintainability than global variable pattern

For application-wide singleton configuration with no variation per request/scope,
a global variable is the most pragmatic solution in Go.

#### Scenario: Global config initialization

- **GIVEN** the application starts up
- **WHEN** `core.InitConfig(configFile)` is called during Cobra initialization
- **THEN** the global `core.Cfg` variable SHALL be set
- **AND** all services SHALL access config via `core.Cfg.Field`
- **AND** config SHALL remain immutable after initialization

#### Scenario: Services access global config

- **GIVEN** a service needs application configuration
- **WHEN** the service method executes
- **THEN** it SHALL access config via `core.Cfg.Field` directly
- **AND** SHALL NOT store config in service fields (to avoid stale references)
- **AND** MAY call config methods like `core.Cfg.GetDataDir()`

### Requirement: Health Checks

Critical services SHALL implement health checks using the `Healthchecker`
interface, enabling startup validation and operational monitoring.

#### Scenario: Service implements health check interface

- **GIVEN** a service depends on external resources or configuration
- **WHEN** the service is created
- **THEN** it SHALL implement `Healthchecker` or `HealthcheckerWithContext`
  interface
- **AND** the `HealthCheck()` method SHALL verify the service is ready to
  operate
- **AND** SHALL return an error if health check fails

#### Scenario: Startup health validation

- **GIVEN** all services are registered with health check implementations
- **WHEN** the application starts up
- **THEN** `do.HealthCheck[ServiceType](injector)` SHALL be called for critical
  services
- **AND** if any health check fails, application startup SHALL abort
- **AND** health check errors SHALL be logged with service name and error
  details
- **AND** the application SHALL exit with non-zero status code

#### Scenario: Health check with timeout

- **GIVEN** a service implements `HealthcheckerWithContext`
- **WHEN** health check is performed with a timeout context
- **THEN** the service SHALL complete its health check within the timeout
- **AND** SHALL return error if check cannot complete in time
- **AND** timeout expiration SHALL be treated as health check failure
