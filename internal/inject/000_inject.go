package inject

import (
	"sync"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/domain/logic/repository"
	"github.com/c18t/nippo-cli/internal/domain/logic/service"
	"github.com/samber/do/v2"
)

var (
	injector *do.RootScope
	once     sync.Once
)

// GetInjector returns the singleton DI container with lazy initialization.
// It is thread-safe and initializes the container only once.
func GetInjector() *do.RootScope {
	once.Do(func() {
		injector = addProvider()
	})
	return injector
}

func addProvider() *do.RootScope {
	var i = do.New()

	// adapter/gateway
	do.Provide(i, gateway.NewDriveFileProvider)
	do.Provide(i, gateway.NewLocalFileProvider)

	// domain/repository
	do.Provide(i, repository.NewRemoteNippoQuery)
	do.Provide(i, repository.NewLocalNippoQuery)
	do.Provide(i, repository.NewLocalNippoCommand)
	do.Provide(i, repository.NewAssetRepository)

	// domain/service
	do.Provide(i, service.NewNippoFacade)
	do.Provide(i, service.NewTemplateService)
	return i
}
