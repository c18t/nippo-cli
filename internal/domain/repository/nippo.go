package repository

import (
	"time"

	"github.com/c18t/nippo-cli/internal/domain/model"
)

type QueryListParam struct {
	Folders        []string
	UpdatedAt      time.Time
	FileExtensions []string
	OrderBy        string
	PageToken      string
}

type QueryListOption struct {
	WithContent bool
	Recursive   bool
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
