package service

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/domain/model"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	i "github.com/c18t/nippo-cli/internal/domain/service"
	"go.uber.org/dig"
)

type nippoFacade struct {
	remoteQuery  repository.RemoteNippoQuery
	localQuery   repository.LocalNippoQuery
	localCommand repository.LocalNippoCommand
}

type inNippoFacade struct {
	dig.In
	RemoteQuery  repository.RemoteNippoQuery
	LocalQuery   repository.LocalNippoQuery
	LocalCommand repository.LocalNippoCommand
}

func NewNippoFacade(serviceDeps inNippoFacade) i.NippoFacade {
	return &nippoFacade{
		remoteQuery:  serviceDeps.RemoteQuery,
		localQuery:   serviceDeps.LocalQuery,
		localCommand: serviceDeps.LocalCommand,
	}
}

func (s *nippoFacade) Send(request *i.NippoFacadeRequest, option *i.NippoFacadeOption) (*i.NippoFacadeReponse, error) {
	remoteFiles, err := s.remoteQuery.List(request.Query, request.Option)
	if err != nil {
		return nil, err
	}

	nippoList := make([]model.Nippo, len(remoteFiles))

	if request.Action&i.NippoFacadeActionDownload != 0 {
		for i, nippo := range remoteFiles {
			fmt.Printf("%s (%s)\n", nippo.RemoteFile.Name, nippo.RemoteFile.Id)
			err = s.remoteQuery.Download(&nippo)
			if err != nil {
				fmt.Printf("download failed: %v\n", err)
			}
			if len(request.Content) <= i {
				request.Content = append(request.Content, model.Nippo{})
			}
			request.Content[i] = model.Nippo(nippo)
		}
	}

	if request.Action&i.NippoFacadeActionCache != 0 {
		for i, nippo := range request.Content {
			err = s.localCommand.Create(&nippo)
			if err != nil {
				fmt.Printf("cache failed: %v\n", err)
			}
			nippoList[i] = nippo
		}
	} else {
		nippoList = request.Content
	}

	count := len(nippoList)
	return &i.NippoFacadeReponse{
		Result: &i.NippoFacadeResponseResult{
			Action:  request.Action,
			Count:   count,
			Message: fmt.Sprintf("%d files downloaded.", len(nippoList)),
		},
		Content: nippoList,
	}, nil
}
