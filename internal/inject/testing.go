package inject

import (
	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/c18t/nippo-cli/internal/domain/service"
	"github.com/samber/do/v2"
)

// TestBasePackageOptions allows selective replacement of services in the base package.
// Any non-nil field will override the default implementation.
type TestBasePackageOptions struct {
	// core
	Config *core.Config

	// adapter/gateway
	DriveFileProvider gateway.DriveFileProvider
	LocalFileProvider gateway.LocalFileProvider

	// domain/repository
	RemoteNippoQuery  repository.RemoteNippoQuery
	LocalNippoQuery   repository.LocalNippoQuery
	LocalNippoCommand repository.LocalNippoCommand
	AssetRepository   repository.AssetRepository

	// domain/service
	NippoFacade     service.NippoFacade
	TemplateService service.TemplateService
}

// NewTestInjector creates a test injector with optional service replacements.
//
// Usage:
//
//	// Use all default services
//	injector := inject.NewTestInjector(nil)
//
//	// Replace specific services with mocks
//	injector := inject.NewTestInjector(&inject.TestBasePackageOptions{
//	    DriveFileProvider: mockDriveProvider,
//	    Config: testConfig,
//	})
//
//	// Invoke services as normal
//	controller, _ := do.Invoke[controller.BuildController](injector)
func NewTestInjector(opts *TestBasePackageOptions) *do.RootScope {
	// Start with base package
	injector := do.New(BasePackage)

	if opts == nil {
		return injector
	}

	// Override services with test implementations
	if opts.DriveFileProvider != nil {
		do.Override(injector, func(do.Injector) (gateway.DriveFileProvider, error) {
			return opts.DriveFileProvider, nil
		})
	}

	if opts.LocalFileProvider != nil {
		do.Override(injector, func(do.Injector) (gateway.LocalFileProvider, error) {
			return opts.LocalFileProvider, nil
		})
	}

	if opts.RemoteNippoQuery != nil {
		do.Override(injector, func(do.Injector) (repository.RemoteNippoQuery, error) {
			return opts.RemoteNippoQuery, nil
		})
	}

	if opts.LocalNippoQuery != nil {
		do.Override(injector, func(do.Injector) (repository.LocalNippoQuery, error) {
			return opts.LocalNippoQuery, nil
		})
	}

	if opts.LocalNippoCommand != nil {
		do.Override(injector, func(do.Injector) (repository.LocalNippoCommand, error) {
			return opts.LocalNippoCommand, nil
		})
	}

	if opts.AssetRepository != nil {
		do.Override(injector, func(do.Injector) (repository.AssetRepository, error) {
			return opts.AssetRepository, nil
		})
	}

	if opts.NippoFacade != nil {
		do.Override(injector, func(do.Injector) (service.NippoFacade, error) {
			return opts.NippoFacade, nil
		})
	}

	if opts.TemplateService != nil {
		do.Override(injector, func(do.Injector) (service.TemplateService, error) {
			return opts.TemplateService, nil
		})
	}

	return injector
}
