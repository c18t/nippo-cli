package service

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/domain/model"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	ds "github.com/c18t/nippo-cli/internal/domain/service"
	"github.com/samber/do/v2"
)

type nippoFacade struct {
	remoteQuery  repository.RemoteNippoQuery  `do:""`
	localQuery   repository.LocalNippoQuery   `do:""`
	localCommand repository.LocalNippoCommand `do:""`
}

func NewNippoFacade(injector do.Injector) (ds.NippoFacade, error) {
	remoteQuery, err := do.Invoke[repository.RemoteNippoQuery](injector)
	if err != nil {
		return nil, err
	}
	localQuery, err := do.Invoke[repository.LocalNippoQuery](injector)
	if err != nil {
		return nil, err
	}
	localCommand, err := do.Invoke[repository.LocalNippoCommand](injector)
	if err != nil {
		return nil, err
	}
	return &nippoFacade{
		remoteQuery:  remoteQuery,
		localQuery:   localQuery,
		localCommand: localCommand,
	}, nil
}

func (s *nippoFacade) Send(request *ds.NippoFacadeRequest, option *ds.NippoFacadeOption) (*ds.NippoFacadeReponse, error) {
	remoteFiles, err := s.remoteQuery.List(request.Query, request.Option)
	if err != nil {
		return nil, err
	}

	nippoList := make([]model.Nippo, len(remoteFiles))
	total := len(remoteFiles)

	if request.Action&ds.NippoFacadeActionDownload != 0 {
		for i, nippo := range remoteFiles {
			// Report progress via callback if provided
			if option != nil && option.OnProgress != nil {
				if !option.OnProgress(nippo.RemoteFile.Name, nippo.RemoteFile.Id, i+1, total) {
					return nil, ds.ErrCancelled
				}
			}
			err = s.remoteQuery.Download(&nippo)
			if err != nil {
				// Progress callback can handle error display if needed
				// For now, continue processing other files
				_ = err
			}
			if len(request.Content) <= i {
				request.Content = append(request.Content, model.Nippo{})
			}
			request.Content[i] = model.Nippo(nippo)
		}
	}

	if request.Action&ds.NippoFacadeActionCache != 0 {
		for i, nippo := range request.Content {
			err = s.localCommand.Create(&nippo)
			if err != nil {
				// Continue processing other files on cache error
				_ = err
			}
			nippoList[i] = nippo
		}
	} else {
		nippoList = request.Content
	}

	count := len(nippoList)
	return &ds.NippoFacadeReponse{
		Result: &ds.NippoFacadeResponseResult{
			Action:  request.Action,
			Count:   count,
			Message: fmt.Sprintf("%d files downloaded.", len(nippoList)),
		},
		Content: nippoList,
	}, nil
}
