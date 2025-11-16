# Change: Refactor DI implementation with lazy initialization, scope isolation, and

proper error handling

## Why

Issue #37 reports that the DI implementation has multiple significant problems:

1. **Scope Isolation Issue**: `internal/inject/clean.go` directly modifies the
   base `Injector` instead of cloning it, which pollutes the shared dependency
   container. Other command injectors (`build`, `deploy`, `init`, `update`,
   `root`) correctly use `Injector.Clone()` to create isolated scopes.

2. **Eager Initialization Issue**: All injectors use eager initialization
   (`var Injector = AddProvider()`), creating DI containers at package load
   time regardless of whether they're needed. This wastes resources and impacts
   startup performance.

3. **Error Handling Issue**: 39 locations use `do.MustInvoke` which panics on
   error, despite constructors returning `error`. This violates Go error
   handling conventions and makes debugging difficult.

4. **Unused Dependencies Issue**: 4 constructors accept `do.Injector` parameter
   but never use it, violating samber/do conventions and creating confusion
   about actual dependencies.

5. **Outdated Dependency Issue**: The project uses `samber/do/v2` beta version
   (`v2.0.0-beta.7`) instead of the stable release (`v2.0.0` released
   Sep 21, 2024), missing bug fixes and stability improvements.

These issues violate best practices demonstrated in samber/do-template-cli and
the boilerplate-go-cli project (see c18t/boilerplate-go-cli#4).

## What Changes

- **BREAKING**: Upgrade `github.com/samber/do/v2` from `v2.0.0-beta.7` to
  `v2.0.0` (stable release)
- **BREAKING**: Replace `var Injector = AddProvider()` with lazy
  initialization using `sync.Once` pattern
- **BREAKING**: Change all references from `inject.Injector` to
  `inject.GetInjector()`
- **BREAKING**: Replace all `do.MustInvoke` with `do.Invoke` and proper error
  handling (39 locations)
- Fix `internal/inject/clean.go` to use `Injector.Clone()` pattern
- Update all command-specific injectors to use `GetInjector().Clone()` instead
  of `Injector.Clone()`
- Fix 4 constructors that accept unused `do.Injector` parameters:
  - `internal/adapter/gateway/drive_file_provider.go`
  - `internal/adapter/gateway/local_file_provider.go`
  - `internal/adapter/presenter/view/init.go`
  - `internal/domain/logic/service/template_service.go`
- Update scaffdog template (`.scaffdog/command.md`) to use proper error
  handling pattern
- Ensure thread-safe, lazy initialization of DI containers
- Implement consistent error wrapping with context in all constructors

## Impact

- Affected specs: `dependency-injection`
- Affected code:
  - `internal/inject/000_inject.go` (lazy initialization)
  - `internal/inject/*.go` (6 files: GetInjector usage)
  - `internal/adapter/controller/*.go` (6 files: error handling)
  - `internal/adapter/presenter/*.go` (6 files: error handling)
  - `internal/adapter/gateway/*.go` (2 files: unused Injector fix + error handling)
  - `internal/adapter/presenter/view/init.go` (unused Injector fix)
  - `internal/usecase/interactor/*.go` (10+ files: error handling)
  - `internal/usecase/port/*.go` (6 files: error handling)
  - `internal/domain/logic/service/template_service.go` (unused Injector fix)
  - `internal/domain/logic/repository/*.go` (error handling)
  - `.scaffdog/command.md` (template update for all patterns)
- Risk: Medium-High - Breaking changes with extensive code modifications
- Breaking: **YES**
  - `inject.Injector` → `inject.GetInjector()` migration required
  - Constructor signatures remain same, but error returns are now meaningful
- Benefits:
  - Improved startup performance by deferring DI container creation
  - Thread-safe initialization with `sync.Once`
  - Consistent scope isolation across all commands
  - Reduced memory footprint during package initialization
  - **Proper error handling** - No more silent panics, clear error messages
  - **Better debuggability** - Error propagation with context
  - **Cleaner code** - Unused parameters removed
  - **samber/do v2 compliance** - Follows official best practices
