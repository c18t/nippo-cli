# Implementation Tasks

## 0. Upgrade samber/do/v2 to Stable Version
- [x] 0.1 Update `go.mod` to use `v2.0.0` instead of `v2.0.0-beta.7`
- [ ] 0.2 Run `go mod download` to fetch the stable version
- [ ] 0.3 Run `go mod tidy` to update `go.sum`
- [ ] 0.4 Verify no API breaking changes between beta.7 and v2.0.0
- [ ] 0.5 Run `go build` to ensure compatibility

## 1. Implement Lazy Initialization
- [ ] 1.1 Update `internal/inject/000_inject.go` to implement lazy initialization pattern
  - [ ] 1.1.1 Replace `var Injector = AddProvider()` with private `injector` variable
  - [ ] 1.1.2 Add `sync.Once` for thread-safe initialization
  - [ ] 1.1.3 Create `GetInjector()` function that returns `*do.RootScope`
- [ ] 1.2 Update all command injectors to use `GetInjector().Clone()` instead of `Injector.Clone()`
  - [ ] 1.2.1 Update `internal/inject/init.go`
  - [ ] 1.2.2 Update `internal/inject/build.go`
  - [ ] 1.2.3 Update `internal/inject/deploy.go`
  - [ ] 1.2.4 Update `internal/inject/update.go`
  - [ ] 1.2.5 Update `internal/inject/root.go`
  - [ ] 1.2.6 Update `internal/inject/clean.go` (also fixes scope isolation)

## 2. Fix DI Scope Isolation
- [ ] 2.1 Ensure `internal/inject/clean.go` uses `GetInjector().Clone()` pattern (covered in 1.2.6)
- [ ] 2.2 Verify all command injectors return cloned scope, not base injector
- [ ] 2.3 Review base injector (`000_inject.go`) contains only shared services

## 3. Replace MustInvoke with Proper Error Handling
- [ ] 3.1 Update adapter/controller constructors (6 files)
  - [ ] 3.1.1 `internal/adapter/controller/init.go`
  - [ ] 3.1.2 `internal/adapter/controller/build.go`
  - [ ] 3.1.3 `internal/adapter/controller/deploy.go`
  - [ ] 3.1.4 `internal/adapter/controller/clean.go`
  - [ ] 3.1.5 `internal/adapter/controller/update.go`
  - [ ] 3.1.6 `internal/adapter/controller/root.go`
- [ ] 3.2 Update adapter/presenter constructors (6 files)
  - [ ] 3.2.1 `internal/adapter/presenter/init.go` (2 constructors)
  - [ ] 3.2.2 `internal/adapter/presenter/build.go`
  - [ ] 3.2.3 `internal/adapter/presenter/deploy.go`
  - [ ] 3.2.4 `internal/adapter/presenter/clean.go`
  - [ ] 3.2.5 `internal/adapter/presenter/update.go`
  - [ ] 3.2.6 `internal/adapter/presenter/root.go`
- [ ] 3.3 Update usecase/port constructors (6 files)
  - [ ] 3.3.1 `internal/usecase/port/init.go`
  - [ ] 3.3.2 `internal/usecase/port/build.go`
  - [ ] 3.3.3 `internal/usecase/port/deploy.go`
  - [ ] 3.3.4 `internal/usecase/port/clean.go`
  - [ ] 3.3.5 `internal/usecase/port/update.go`
  - [ ] 3.3.6 `internal/usecase/port/root.go`
- [ ] 3.4 Update usecase/interactor constructors (10+ files)
  - [ ] 3.4.1 All interactor files in `internal/usecase/interactor/`
- [ ] 3.5 Update domain/logic/repository constructors
  - [ ] 3.5.1 `internal/domain/logic/repository/nippo.go` (3 constructors)
  - [ ] 3.5.2 `internal/domain/logic/repository/asset.go`

## 4. Fix Unused Injector Parameters
- [ ] 4.1 Remove unused Injector from `internal/adapter/gateway/drive_file_provider.go`
  - [ ] 4.1.1 Change signature to `func NewDriveFileProvider() (DriveFileProvider, error)`
  - [ ] 4.1.2 Update registration in `internal/inject/000_inject.go`
- [ ] 4.2 Remove unused Injector from `internal/adapter/gateway/local_file_provider.go`
  - [ ] 4.2.1 Change signature to `func NewLocalFileProvider() (LocalFileProvider, error)`
  - [ ] 4.2.2 Update registration in `internal/inject/000_inject.go`
- [ ] 4.3 Remove unused Injector from `internal/adapter/presenter/view/init.go`
  - [ ] 4.3.1 Change signature for `NewConfigureProjectView`
  - [ ] 4.3.2 Update registration in `internal/inject/init.go`
- [ ] 4.4 Remove unused Injector from `internal/domain/logic/service/template_service.go`
  - [ ] 4.4.1 Change signature to `func NewTemplateService() (TemplateService, error)`
  - [ ] 4.4.2 Update registration in `internal/inject/000_inject.go`

## 5. Update Templates and Examples
- [ ] 5.1 Update `.scaffdog/command.md` template to use `GetInjector().Clone()` pattern
- [ ] 5.2 Update template to use `do.Invoke` instead of `do.MustInvoke`
- [ ] 5.3 Update template with proper error handling pattern
- [ ] 5.4 Add example error wrapping in template

## 6. Validation and Testing
- [ ] 6.1 Build the project to ensure no compilation errors
- [ ] 6.2 Run all commands to verify functionality (`init`, `build`, `deploy`, `clean`, `update`)
- [ ] 6.3 Run with race detector: `go run -race . --help`
- [ ] 6.4 Verify base injector is lazily initialized (only when accessed)
- [ ] 6.5 Verify base injector is not modified by command initialization
- [ ] 6.6 Test error handling by intentionally breaking a dependency
- [ ] 6.7 Verify error messages are descriptive and include context
- [ ] 6.8 Confirm no panics occur during dependency resolution failures

## 7. Documentation
- [ ] 7.1 Update AGENTS.md to document lazy initialization pattern
- [ ] 7.2 Add note about `GetInjector()` usage in DI documentation
- [ ] 7.3 Document error handling pattern for constructors
- [ ] 7.4 Document breaking changes and migration path
- [ ] 7.5 Add examples of proper `do.Invoke` usage vs `do.MustInvoke`
