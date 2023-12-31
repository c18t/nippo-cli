package repository

import (
	"os"
	"path/filepath"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/core"
	i "github.com/c18t/nippo-cli/internal/domain/repository"
)

type assetRepository struct {
	provider gateway.LocalFileProvider
}

func NewAssetRepository(p gateway.LocalFileProvider) i.AssetRepository {
	return &assetRepository{p}
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
