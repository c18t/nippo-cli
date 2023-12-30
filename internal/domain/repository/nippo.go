package repository

import (
	"time"

	"github.com/c18t/nippo-cli/internal/domain/model"
)

type QueryListParam struct {
	Folder        string
	UpdatedAt     time.Time
	FileExtension string
	OrderBy       string
}

type QueryListOption struct {
	WithContent bool
}

type RemoteNippoQuery interface {
	List(param *QueryListParam, option *QueryListOption) ([]model.Nippo, error)
	Download(nippo *model.Nippo) error
}

type LocalNippoQuery interface {
	Exist(*model.NippoDate) bool
	Find(*model.NippoDate) (*model.Nippo, error)
	List(param *QueryListParam, option *QueryListOption) ([]model.Nippo, error)
	Load(nippo *model.Nippo) error
}

type LocalNippoCommand interface {
	Create(*model.Nippo) error
	Delete(*model.Nippo) error
}
