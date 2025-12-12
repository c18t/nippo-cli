package service

import (
	"errors"

	"github.com/c18t/nippo-cli/internal/domain/model"
	"github.com/c18t/nippo-cli/internal/domain/repository"
)

// ErrCancelled is returned when the operation is cancelled by the user
var ErrCancelled = errors.New("operation cancelled")

type NippoFacade interface {
	Send(request *NippoFacadeRequest, option *NippoFacadeOption) (*NippoFacadeReponse, error)
}

type NippoFacadeAction int

const (
	NippoFacadeActionSearch = 1 << iota
	NippoFacadeActionClean
	NippoFacadeActionDownload
	NippoFacadeActionCache
)

type NippoFacadeRequest struct {
	Action  NippoFacadeAction
	Query   *repository.QueryListParam
	Option  *repository.QueryListOption
	Content []model.Nippo
}
// ProgressCallback is called for each file processed during download/cache operations.
// Returns true to continue, false to cancel the operation.
type ProgressCallback func(filename string, fileId string, current int, total int) bool

type NippoFacadeOption struct {
	OnProgress ProgressCallback
}
type NippoFacadeReponse struct {
	Result  *NippoFacadeResponseResult
	Content []model.Nippo
}
type NippoFacadeResponseResult struct {
	Action  NippoFacadeAction
	Count   int
	Message string
}
