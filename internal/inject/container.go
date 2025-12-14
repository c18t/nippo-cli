// Package inject provides dependency injection container setup using samber/do/v2.
//
// The package uses the do.Package pattern to group shared services into a reusable
// module. Services are lazily initialized - they are created only when first requested
// via do.Invoke, improving application startup time.
//
// Base injector usage:
//
//	injector := inject.GetInjector()           // Get singleton base injector
//	service, _ := do.Invoke[MyService](injector)  // Resolve service (creates if needed)
//
// Command-specific injectors:
//
//	cmdInjector := inject.GetInjector().Clone()  // Create isolated scope
//	do.Provide(cmdInjector, NewCommandService)   // Add command-specific services
//
// The base injector is initialized once using sync.Once and provides shared services
// (gateways, repositories, domain services) to all commands. Each command creates an
// isolated scope by cloning the base injector to avoid cross-command pollution.
package inject

import (
	"sync"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/domain/logic/repository"
	"github.com/c18t/nippo-cli/internal/domain/logic/service"
	"github.com/samber/do/v2"
)

var (
	injector *do.RootScope
	once     sync.Once
)

// BasePackage groups all shared services for the application.
//
// Services are registered with do.Lazy() wrappers, enabling lazy initialization.
// Each service is created only when first requested via do.Invoke, not during
// package initialization. This improves startup performance by deferring expensive
// initializations until actually needed.
//
// The package includes:
//   - adapter/gateway: File providers (Drive API, local filesystem)
//   - domain/repository: Data access (nippo queries, commands, assets)
//   - domain/service: Business logic (nippo facade, template service)
//
// Note: Configuration is managed via the global core.Cfg variable initialized
// by core.InitConfig() at application startup, not through dependency injection.
//
// Command-specific services (controllers, presenters, use cases) are registered
// separately in each command's injector file (e.g., build.go, clean.go).
var BasePackage = do.Package(
	// adapter/gateway
	do.Lazy(gateway.NewDriveFileProvider),
	do.Lazy(gateway.NewLocalFileProvider),

	// adapter/presenter
	do.Lazy(presenter.NewConsolePresenter),

	// domain/repository
	do.Lazy(repository.NewRemoteNippoQuery),
	do.Lazy(repository.NewLocalNippoQuery),
	do.Lazy(repository.NewLocalNippoCommand),
	do.Lazy(repository.NewAssetRepository),

	// domain/service
	do.Lazy(service.NewNippoFacade),
	do.Lazy(service.NewTemplateService),
)

// GetInjector returns the singleton DI container with lazy initialization.
// It is thread-safe and initializes the container only once.
func GetInjector() *do.RootScope {
	once.Do(func() {
		injector = do.New(BasePackage)
	})
	return injector
}
