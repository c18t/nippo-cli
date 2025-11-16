# dependency-injection Specification

## ADDED Requirements

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
