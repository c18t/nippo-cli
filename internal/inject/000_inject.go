package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/domain/logic/repository"
	"github.com/c18t/nippo-cli/internal/domain/logic/service"
	"github.com/samber/do/v2"
)

var Injector = AddProvider()

func AddProvider() *do.RootScope {
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
