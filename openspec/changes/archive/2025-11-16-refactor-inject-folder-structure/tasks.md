# Implementation Tasks

**Status**: Completed on 2025-11-16

## Summary of Completed Work

All phases successfully implemented:

1. **File Rename**: `000_inject.go` → `container.go` using `git mv` to preserve history
2. **do.Package Migration**: Refactored to use `BasePackage` with `do.Lazy()` wrappers for lazy initialization
3. **Documentation Updated**: Updated `openspec/project.md` and `AGENTS.md` to reference new filename and pattern
4. **Verification**: Build successful, `--help` and `--version` working correctly
5. **Godoc Enhanced**: Added comprehensive package and variable documentation

## 1. Rename File

- [ ] 1.1 Rename `internal/inject/000_inject.go` to `internal/inject/container.go`
  - [ ] 1.1.1 Use `git mv` to preserve history: `git mv internal/inject/000_inject.go internal/inject/container.go`
  - [ ] 1.1.2 Verify file contents remain unchanged (defer refactoring to step 2)

## 2. Migrate to do.Package Pattern

- [ ] 2.1 Create default config constructor
  - [ ] 2.1.1 Add `NewDefaultConfig() (*Config, error)` to `internal/core/config.go`
  - [ ] 2.1.2 Return empty `&Config{}` for DI initialization
  - [ ] 2.1.3 Document that this is overridden by `InitConfig()`

- [ ] 2.2 Refactor `container.go` to use `do.Package`
  - [ ] 2.2.1 Create `BasePackage` variable with `do.Package()`
  - [ ] 2.2.2 Wrap all service constructors with `do.Lazy()`
  - [ ] 2.2.3 Replace `addProvider()` implementation to use `do.New(BasePackage)`
  - [ ] 2.2.4 Maintain same service registration order for clarity
  - [ ] 2.2.5 Add grouping comments (core, adapter, domain, usecase)

- [ ] 2.3 Verify lazy initialization behavior
  - [ ] 2.3.1 Confirm services are not created until first `do.Invoke`
  - [ ] 2.3.2 Test that config override in `InitConfig()` still works
  - [ ] 2.3.3 Verify no circular dependency issues

## 3. Update Documentation

- [ ] 3.1 Update `openspec/project.md` references
  - [ ] 3.1.1 Change `internal/inject/000_inject.go` to `internal/inject/container.go`
  - [ ] 3.1.2 Add documentation about `do.Package` pattern
  - [ ] 3.1.3 Update "Base injector" description to mention `BasePackage`
  - [ ] 3.1.4 Document lazy initialization benefits

- [ ] 3.2 Update CLAUDE.md if it references the old filename
  - [ ] 3.2.1 Search for "000_inject.go" in CLAUDE.md
  - [ ] 3.2.2 Update to "container.go" if found

- [ ] 3.3 Update README.md if it references the old filename
  - [ ] 3.3.1 Search for "000_inject.go" in README.md
  - [ ] 3.3.2 Update to "container.go" if found

## 4. Verification

- [ ] 4.1 Verify compilation
  - [ ] 4.1.1 Run `mise run build` to ensure no build errors
  - [ ] 4.1.2 Verify all import paths still work correctly
  - [ ] 4.1.3 Check for any unused imports after refactoring

- [ ] 4.2 Verify functionality
  - [ ] 4.2.1 Run `./bin/nippo --help` to verify help displays
  - [ ] 4.2.2 Run `./bin/nippo --version` to verify version displays
  - [ ] 4.2.3 Test a command that uses Drive API (e.g., `nippo build --help`)
  - [ ] 4.2.4 Verify no runtime errors or panics occur

- [ ] 4.3 Verify lazy initialization
  - [ ] 4.3.1 Add temporary logging to service constructors
  - [ ] 4.3.2 Run `./bin/nippo --help` and verify services are NOT created
  - [ ] 4.3.3 Run actual command and verify services ARE created on demand
  - [ ] 4.3.4 Remove temporary logging

- [ ] 4.4 Verify git history preservation
  - [ ] 4.4.1 Run `git log --follow internal/inject/container.go` to confirm history is intact
  - [ ] 4.4.2 Verify the rename is properly tracked

## 5. Optional: Add Godoc Comments

- [ ] 5.1 Enhance package documentation in `container.go`
  - [ ] 5.1.1 Add package-level comment explaining the `do.Package` pattern
  - [ ] 5.1.2 Document `BasePackage` variable with usage examples
  - [ ] 5.1.3 Explain lazy initialization and when services are created
  - [ ] 5.1.4 Document the relationship between base and command-specific injectors
  - [ ] 5.1.5 Explain when to use `GetInjector()` vs `Clone()`
