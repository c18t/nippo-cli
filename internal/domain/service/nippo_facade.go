package service

import (
	"github.com/c18t/nippo-cli/internal/domain/model"
	"github.com/c18t/nippo-cli/internal/domain/repository"
)

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
type NippoFacadeOption struct {
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
