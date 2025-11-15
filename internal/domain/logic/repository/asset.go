package repository

import (
	"os"
	"path/filepath"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/core"
	i "github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/samber/do/v2"
)

type assetRepository struct {
	provider gateway.LocalFileProvider `do:""`
}

func NewAssetRepository(injector do.Injector) (i.AssetRepository, error) {
	provider, err := do.Invoke[gateway.LocalFileProvider](injector)
	if err != nil {
		return nil, err
	}
	return &assetRepository{
		provider: provider,
	}, nil
}

func (r *assetRepository) CleanNippoCache() error {
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "md")
	return r.clean(&i.QueryListParam{
		Folders:        []string{outputDir},
		FileExtensions: []string{"md"},
	})
}

func (r *assetRepository) CleanBuildCache() error {
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "output")
	return r.clean(&i.QueryListParam{
		Folders:        []string{outputDir},
		FileExtensions: []string{"html"},
	})
}

func (r *assetRepository) clean(query *i.QueryListParam) error {
	files, err := r.provider.List(query)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.Remove(filepath.Join(query.Folders[0], file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
