package repository

type AssetRepository interface {
	CleanNippoCache() error
	CleanBuildCache() error
}
