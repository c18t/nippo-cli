# Change: Refactor inject folder structure for better clarity

## Why

The current `internal/inject/` folder structure uses `000_inject.go` to define
the base injector that provides shared services. The `000_` prefix was added to
avoid naming conflicts with command-specific injector files (e.g., `build.go`,
`clean.go`), but this approach has several issues:

1. **Unclear naming convention**: The `000_` prefix is non-standard and doesn't
   clearly communicate the file's purpose. New developers might not understand
   why this prefix exists.

2. **File sorting dependence**: The naming relies on alphabetical sorting to
   ensure the base injector is visually first in the directory listing, which
   is a fragile convention.

3. **Poor discoverability**: The file name doesn't indicate that it contains
   the singleton base injector factory, making it harder to navigate the
   codebase.

4. **Not following Go conventions**: Go packages typically use descriptive
   names rather than numeric prefixes.

After reviewing the
[samber/do-template-cli](https://github.com/samber/do-template-cli)
implementation, they use a clearer pattern:

- **`pkg/base.go`**: Contains the base package provider with shared services
- **`pkg/jobs/package.go`**: Contains job-specific providers
- **Main initialization**: `do.New(pkg.BasePackage, jobs.Package)`

This pattern clearly separates base/shared services from feature-specific services.

## What Changes

Refactor the `internal/inject/` folder with two improvements:

### 1. Rename `000_inject.go` → `container.go`

```text
internal/inject/
├── container.go          # Base injector with shared services
├── build.go             # Build command injector
├── clean.go             # Clean command injector
├── deploy.go            # Deploy command injector
├── init.go              # Init command injector
├── root.go              # Root command injector
└── update.go            # Update command injector
```

**Naming rationale**:

- **Clear intent**: The name "container" explicitly indicates it contains the
  DI container factory
- **No conflicts**: Commands don't use "container" as a name
- **Discoverable**: Easy to find when looking for DI container setup
- **Standard Go naming**: Uses descriptive nouns without numeric prefixes

### 2. Migrate from `do.Provide` to `do.Package` pattern

**Current pattern (do.Provide)**:

```go
func addProvider() *do.RootScope {
    var i = do.New()
    do.Provide(i, gateway.NewDriveFileProvider)
    do.Provide(i, gateway.NewLocalFileProvider)
    // ... more individual registrations
    return i
}
```

**New pattern (do.Package)**:

```go
// BasePackage groups all shared services
var BasePackage = do.Package(
    // core
    do.Lazy(core.NewDefaultConfig),

    // adapter/gateway
    do.Lazy(gateway.NewDriveFileProvider),
    do.Lazy(gateway.NewLocalFileProvider),

    // domain/repository
    do.Lazy(repository.NewRemoteNippoQuery),
    do.Lazy(repository.NewLocalNippoQuery),
    do.Lazy(repository.NewLocalNippoCommand),
    do.Lazy(repository.NewAssetRepository),

    // domain/service
    do.Lazy(service.NewNippoFacade),
    do.Lazy(service.NewTemplateService),
)

func GetInjector() *do.RootScope {
    once.Do(func() {
        injector = do.New(BasePackage)
    })
    return injector
}
```

**Benefits of do.Package**:

- **Modularity**: Services are grouped as a reusable package
- **Composability**: Easy to combine with other packages: `do.New(BasePackage, TestPackage)`
- **Lazy initialization**: Services are created only when first requested
- **Better testing**: Can create alternative packages for testing
- **Follows best practices**: Aligns with samber/do-template-cli pattern

## Impact

**Benefits:**

- **Improved readability**: Clear, self-documenting file names
- **Better maintainability**: New developers can easily understand the structure
- **Follows best practices**: Aligns with Go conventions and samber/do patterns
- **Better modularity**: Package pattern enables cleaner service grouping
- **Improved testability**: Can easily substitute entire service packages in tests
- **Lazy loading**: Services initialize on first use, improving startup time

**Breaking Changes:**

- None (internal refactoring only)

**Migration:**

1. Rename `internal/inject/000_inject.go` → `internal/inject/container.go`
2. Refactor `addProvider()` to use `do.Package` with `do.Lazy` wrappers
3. Create a default config constructor for initial DI setup
4. Update documentation referencing the old patterns

## Related Work

Builds upon the DI refactoring completed in:

- `openspec/changes/archive/2025-11-15-refactor-di-implementation`
- `openspec/changes/enhance-di-lifecycle-management`
