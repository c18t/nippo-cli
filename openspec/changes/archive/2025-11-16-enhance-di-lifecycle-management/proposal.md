# Change: Enhance DI lifecycle management with shutdown, aliasing, and health checks

## Why

The current DI implementation uses samber/do v2.0.0 but doesn't leverage several
important features that would improve application reliability, testability, and
maintainability:

1. **Missing Graceful Shutdown**: The application doesn't implement proper
   cleanup when terminating. Resources like Google Drive API clients, file
   handles, and other services are not explicitly closed, potentially leading to
   resource leaks or incomplete operations. Additionally, there is a **bug** in
   the current Drive API client implementation: the `FilesService` is recreated
   on every file operation instead of being reused, which is inefficient and
   prevents proper lifecycle management.

2. **No Service Aliasing**: Services are registered and resolved by concrete
   types rather than interfaces. This creates tight coupling between layers and
   makes testing difficult. The project has well-defined port interfaces in
   `internal/usecase/port/` but doesn't use them for dependency injection.

3. **No Health Checks**: There's no mechanism to verify service health at
   startup or during operation. This makes it difficult to detect configuration
   or dependency issues early.

**Note on Configuration Management**: During implementation, we evaluated moving
configuration from global variables to DI container management using
`do.ProvideValue`. However, this approach would require all services to invoke
config in their constructors, which causes immediate evaluation before
`InitConfig()` runs. Attempting to defer config access to service methods
introduces a service locator anti-pattern that reduces maintainability. The
global `core.Cfg` variable pattern remains the most pragmatic solution for
application-wide singleton configuration.

These gaps were identified by reviewing samber/do v2 documentation
(<https://do.samber.dev> and <https://pkg.go.dev/github.com/samber/do/v2>),
which provides comprehensive lifecycle management features that are currently
unused.

## What Changes

- **Implement Graceful Shutdown**: Fix the Drive API client bug by storing and
  reusing the `FilesService` instance instead of recreating it. Add `Shutdowner`
  interface implementation to properly close the Drive API client. Use
  `do.Shutdown` to cleanly close resources on application termination. Note:
  Drive API client initialization must be lazy (deferred until config is loaded)
  to access credential file paths.
- **Add Service Aliasing**: Use `do.As` to bind concrete implementations to
  port interfaces, enabling interface-based dependency injection. Update all
  service registrations to support both concrete and interface-based resolution.
- **Add Health Checks**: Implement `Healthchecker` interface for critical
  services (Drive API client, file providers) and add startup validation using
  `do.HealthCheck`.
- **Update Command Lifecycle**: Modify command controllers to trigger shutdown
  on completion or interruption (SIGTERM, SIGINT).

## Impact

**Benefits:**

- **Improved Reliability**: Graceful shutdown ensures resources are properly
  released
- **Better Testability**: Interface-based injection enables easy mocking and
  testing
- **Clearer Dependencies**: Service aliasing makes dependency graphs more
  explicit
- **Early Failure Detection**: Health checks catch configuration issues at
  startup

**Breaking Changes:**

- BREAKING: Service shutdown now required for clean application termination

**Migration:**

- Add shutdown handlers to main command execution
- Update tests to use new DI patterns

## Related Work

Builds upon the DI refactoring completed in
`openspec/changes/archive/2025-11-15-refactor-di-implementation`, which
established lazy initialization (via `GetInjector()` singleton) and scope
isolation patterns (via `.Clone()` in command injectors).
