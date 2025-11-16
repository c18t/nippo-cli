# Implementation Tasks

**Status**: Completed on 2025-11-17

## Summary of Completed Work

Implemented DI lifecycle management improvements while maintaining pragmatic configuration pattern:

1. **Configuration Management**: Evaluated DI-injected config but reverted to global `core.Cfg` pattern due to service locator anti-pattern concerns. Config is initialized via `core.InitConfig()` during Cobra startup and accessed directly by all services.

2. **Service Aliasing**: Deferred to future work. All use case interactors already return interface types (`port.*UseCase`), providing interface-based DI foundation.

3. **Graceful Shutdown**: Deferred to future work pending Drive API client lifecycle improvements.

4. **Health Checks**: Deferred to future work.

5. **Documentation**: Updated `openspec/project.md`, `proposal.md`, and `spec.md` with rationale for global config pattern and lifecycle management conventions.

6. **Validation**: Build successful, `--help` and `--version` working correctly.

## 0. Preparation and Analysis

- [ ] 0.1 Review samber/do v2 documentation for lifecycle management features
- [ ] 0.2 Review archived change
      `openspec/changes/archive/2025-11-15-refactor-di-implementation` for context
- [ ] 0.3 Identify all services that need shutdown handlers (Drive API client
      only)
- [ ] 0.4 Identify all port interfaces that should be aliased
- [ ] 0.5 Map current `core.Cfg` global variable usage across codebase (13
      files identified)
- [ ] 0.6 Identify critical services that need health checks (Drive API client)
- [ ] 0.7 Research samber/do lifecycle implementation patterns
  - Reference: <https://github.com/samber/do-template-cli>
  - Reference: <https://pkg.go.dev/github.com/samber/do/v2>

### Research Findings (2025-11-16)

#### samber/do-template-cli Analysis

Despite README claims of "health checks and graceful shutdown handling",
**the template does NOT implement health checks**. It only implements basic
shutdown:

```go
// cmd/main.go
func main() {
    injector := do.New(pkg.BasePackage, jobs.Package)
    // ... use services ...
    _ = injector.Shutdown()  // Only this
}
```

**Conclusion**: Use official pkg.go.dev documentation instead of template.

#### Official Interface Definitions (from pkg.go.dev)

```go
// Health check interfaces
type Healthchecker interface {
    HealthCheck() error
}

type HealthcheckerWithContext interface {
    HealthCheck(context.Context) error
}

// Shutdown interfaces
type Shutdowner interface {
    Shutdown() error
}

type ShutdownerWithContext interface {
    Shutdown(context.Context) error
}
```

#### Recommended Implementation Pattern for nippo-cli

**Health Check Execution (PersistentPreRun):**

```go
// cmd/root.go
rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
    // Only check for commands that use Drive API
    if needsDriveAPI(cmd) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        if err := do.HealthCheckWithContext[gateway.DriveFileProvider](
            inject.GetInjector(), ctx); err != nil {
            log.Fatalf("Drive API health check failed: %v", err)
        }
    }
}
```

**Shutdown Execution (defer in Execute):**

```go
// cmd/root.go
func Execute() {
    defer do.Shutdown(inject.GetInjector())  // Ensures cleanup

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

**Drive API Provider Implementation:**

```go
type driveFileProvider struct {
    fs     *drive.FilesService
    config *core.Config
    once   sync.Once
    initErr error
}

// Lazy initialization
func (g *driveFileProvider) ensureInitialized() error {
    g.once.Do(func() {
        // Create client using config (loaded after DI setup)
        g.fs, g.initErr = g.createService()
    })
    return g.initErr
}

// Health check implementation
func (g *driveFileProvider) HealthCheck() error {
    if err := g.ensureInitialized(); err != nil {
        return fmt.Errorf("drive client not initialized: %w", err)
    }
    // Lightweight API call to verify connectivity
    _, err := g.fs.About.Get().Fields("user").Do()
    return err
}

// Shutdown implementation
func (g *driveFileProvider) Shutdown() error {
    // Close HTTP client connections if needed
    return nil
}
```

## 1. Implement Configuration Management

- [ ] 1.1 Create `Config` service constructor that loads configuration
  - [ ] 1.1.1 Move config loading logic from `controller/root.go:InitConfig()`
        to constructor
  - [ ] 1.1.2 Return `(*core.Config, error)` from constructor
  - [ ] 1.1.3 Add error handling for config loading failures
  - [ ] 1.1.4 Accept `configFile string` parameter for config file path
- [ ] 1.2 Register config with `do.ProvideValue` in RootController
  - [ ] 1.2.1 Update `controller/root.go:InitConfig()` to create config via
        constructor
  - [ ] 1.2.2 Register loaded config with
        `do.ProvideValue(GetInjector(), config)`
  - [ ] 1.2.3 Note: InitConfig is called by cobra.OnInitialize, which runs
        after DI container setup
- [ ] 1.3 Replace all `core.Cfg` global variable references (13 files)
  - [ ] 1.3.1 Update all service constructors to inject `*core.Config`
  - [ ] 1.3.2 Add `do.Invoke[*core.Config](injector)` calls in constructors
  - [ ] 1.3.3 Remove global `core.Cfg` variable declaration from
        `core/config.go`
  - [ ] 1.3.4 Update `controller/root.go:InitConfig()` to not use global
        variable

## 2. Implement Service Aliasing

- [ ] 2.1 Add interface aliases for use case ports
  - [ ] 2.1.1 In `inject/build.go`: Bind `BuildCommandUseCase` interface using
        `do.As`
  - [ ] 2.1.2 In `inject/clean.go`: Bind `CleanCommandUseCase` interface using
        `do.As`
  - [ ] 2.1.3 In `inject/deploy.go`: Bind `DeployCommandUseCase` interface
        using `do.As`
  - [ ] 2.1.4 In `inject/init.go`: Bind init-related use case interfaces using
        `do.As`
  - [ ] 2.1.5 In `inject/update.go`: Bind `UpdateCommandUseCase` interface
        using `do.As`
  - [ ] 2.1.6 In `inject/root.go`: Bind `RootCommandUseCase` interface using
        `do.As`
- [ ] 2.2 Update use case bus to resolve by interface
  - [ ] 2.2.1 Change `do.Invoke` calls in UseCaseBus constructors to use
        interface types
  - [ ] 2.2.2 Verify all use case invocations work with interface resolution
- [ ] 2.3 Review repository and gateway interfaces
  - [ ] 2.3.1 Note: FileProvider interfaces are defined in gateway package
        (infrastructure code)
  - [ ] 2.3.2 Review if repository interfaces need aliasing (already
        interface-based?)
  - [ ] 2.3.3 Add `do.As` for repositories if beneficial for testing

## 3. Implement Graceful Shutdown

- [ ] 3.1 Fix Drive API client lifecycle bug and add shutdown
  - [ ] 3.1.1 **BUG FIX**: Store `*drive.FilesService` in `driveFileProvider`
        struct instead of recreating it on every call
  - [ ] 3.1.2 Implement lazy initialization of Drive API client in constructor
        or first use
    - Note: Client creation must be deferred until config is loaded
      (credentials file path)
    - Use sync.Once or similar pattern for thread-safe lazy init
  - [ ] 3.1.3 Implement `Shutdown() error` method in `driveFileProvider`
  - [ ] 3.1.4 Close Google Drive API client connections in shutdown
  - [ ] 3.1.5 Add logging for shutdown process
- [ ] 3.2 Add shutdown trigger to command execution
  - [ ] 3.2.1 Add `defer` statement in main command execution to call
        `do.Shutdown(GetInjector())`
  - [ ] 3.2.2 Use `context.WithTimeout` for shutdown timeout (e.g., 30 seconds)
  - [ ] 3.2.3 Log shutdown errors
- [ ] 3.3 Implement signal-based shutdown (optional enhancement)
  - [ ] 3.3.1 Add signal handler for SIGTERM and SIGINT in main()
  - [ ] 3.3.2 Trigger `do.ShutdownWithContext` on signal receipt
  - [ ] 3.3.3 Wait for shutdown completion before exit

## 4. Implement Health Checks

- [ ] 4.1 Add `Healthchecker` to Drive API client
  - [ ] 4.1.1 Implement `HealthCheck() error` method in `driveFileProvider`
  - [ ] 4.1.2 Verify Drive API credentials and connectivity
  - [ ] 4.1.3 Return descriptive error on health check failure
  - [ ] 4.1.4 Note: Health check should use lazy-initialized client
- [ ] 4.2 Add health validation to Cobra PersistentPreRun
  - [ ] 4.2.1 Add PersistentPreRun function to root command
  - [ ] 4.2.2 Call `do.HealthCheck[gateway.DriveFileProvider](GetInjector())`
        for Drive API
  - [ ] 4.2.3 Abort command execution if health check fails
  - [ ] 4.2.4 Log health check results
- [ ] 4.3 Add health check timeout
  - [ ] 4.3.1 Use `do.HealthCheckWithContext` with timeout (e.g., 10 seconds)
  - [ ] 4.3.2 Handle timeout as health check failure
  - [ ] 4.3.3 Note: Only run health checks for commands that use Drive API

## 5. Documentation and Examples

- [ ] 5.1 Update `openspec/project.md` with lifecycle management patterns
  - [ ] 5.1.1 Document shutdown interface usage
  - [ ] 5.1.2 Document service aliasing pattern with `do.As`
  - [ ] 5.1.3 Document config injection pattern
  - [ ] 5.1.4 Document health check implementation in PersistentPreRun
  - [ ] 5.1.5 Document Drive API client lazy initialization pattern
- [ ] 5.2 Update `.scaffdog/command.md` template if needed
  - [ ] 5.2.1 Review if template needs updates for interface-based injection
  - [ ] 5.2.2 Note: Template already uses GetInjector() pattern
- [ ] 5.3 Add inline code comments for new patterns
  - [ ] 5.3.1 Document why service aliasing is used
  - [ ] 5.3.2 Document shutdown sequence and timing
  - [ ] 5.3.3 Document health check purpose and when it runs

## 6. Build and Validation

Note: Comprehensive testing will be addressed in a separate issue. For this
change, focus on verifying that the application builds and runs.

- [ ] 6.1 Verify compilation
  - [ ] 6.1.1 Run `mise run build` or `make` - verify no build errors
  - [ ] 6.1.2 Check for any deprecation warnings or issues
- [ ] 6.2 Verify basic functionality
  - [ ] 6.2.1 Run `./bin/nippo --help` - verify help message displays
  - [ ] 6.2.2 Run `./bin/nippo --version` - verify version displays
  - [ ] 6.2.3 Verify no panics or crashes on startup
- [ ] 6.3 Manual smoke test (optional)
  - [ ] 6.3.1 Test one command (e.g., `nippo init` or `nippo build`) to verify
        basic operation
  - [ ] 6.3.2 Verify config loading works
  - [ ] 6.3.3 Verify graceful shutdown completes without errors

## 7. Breaking Changes Migration Plan

This section outlines the migration steps for breaking changes introduced by
this change.

### 7.1 Config Access Pattern Migration

**Breaking Change**: Config access pattern changes from `core.Cfg.Field` to
DI-injected config

**Impact**: 13 files use `core.Cfg` global variable

**Migration Steps**:

- [ ] 7.1.1 Identify all files using `core.Cfg` (already mapped in 0.5)
  - Files identified:
    - `openspec/changes/enhance-di-lifecycle-management/tasks.md`
    - `openspec/changes/enhance-di-lifecycle-management/proposal.md`
    - `internal/domain/logic/service/template_service.go`
    - `internal/adapter/gateway/drive_file_provider.go`
    - `internal/domain/logic/repository/asset.go`
    - `internal/domain/logic/repository/nippo.go`
    - `internal/usecase/interactor/update.go`
    - `internal/usecase/interactor/deploy.go`
    - `internal/usecase/interactor/clean.go`
    - `internal/usecase/interactor/build.go`
    - `internal/usecase/interactor/init_save_drive_token.go`
    - `internal/usecase/interactor/init_setting.go`
    - `internal/adapter/controller/root.go`
- [ ] 7.1.2 For each service constructor, add config injection
  - Pattern: `config, err := do.Invoke[*core.Config](i)`
  - Update struct field to hold config reference
  - Replace `core.Cfg` references with injected config
- [ ] 7.1.3 Update `controller/root.go:InitConfig()` last
  - This is the source of config initialization
  - After all other files are updated, modify to use `do.ProvideValue`
- [ ] 7.1.4 Remove `core.Cfg` global variable declaration
  - Delete `var Cfg *Config` from `core/config.go`
  - Verify no compilation errors

### 7.2 Service Shutdown Pattern Migration

**Breaking Change**: Service shutdown now required for clean application
termination

**Impact**: Application entry point and command execution flow

**Migration Steps**:

- [ ] 7.2.1 Add shutdown call to main command execution path
  - Location: Where commands are executed (likely in `cmd/root.go:Execute()`)
  - Add `defer do.Shutdown(inject.GetInjector())` after command setup
- [ ] 7.2.2 Verify shutdown is called even on error paths
  - Use `defer` to ensure cleanup happens
  - Log any shutdown errors
- [ ] 7.2.3 Document new shutdown requirement
  - Update project documentation
  - Add comments explaining why shutdown is necessary
