package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/model"
	i "github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/samber/do/v2"
	"google.golang.org/api/drive/v3"
)

type remoteNippoQuery struct {
	provider gateway.DriveFileProvider `do:""`
}

type localNippoQuery struct {
	provider gateway.LocalFileProvider `do:""`
}

func NewRemoteNippoQuery(injector do.Injector) (i.RemoteNippoQuery, error) {
	provider, err := do.Invoke[gateway.DriveFileProvider](injector)
	if err != nil {
		return nil, err
	}
	return &remoteNippoQuery{
		provider: provider,
	}, nil
}

func (r *remoteNippoQuery) List(param *i.QueryListParam, option *i.QueryListOption) ([]model.Nippo, error) {
	tempParam := *param
	res, nippoList, folderList, err := r.list(&tempParam, option)
	if err != nil {
		return nil, err
	}
	for res.NextPageToken != "" {
		tempParam.PageToken = res.NextPageToken
		var pageNippoList []model.Nippo
		var pageFolderList []drive.File
		res, pageNippoList, pageFolderList, err = r.list(&tempParam, option)
		if err != nil {
			return nil, err
		}
		nippoList = append(nippoList, pageNippoList...)
		folderList = append(folderList, pageFolderList...)
	}

	if option.Recursive && len(folderList) > 0 {
		folderIds := make([]string, len(folderList))
		for i, folder := range folderList {
			folderIds[i] = folder.Id
		}
		tempParam.Folders = folderIds
		childNippoList, err := r.List(&tempParam, option)
		if err != nil {
			return nil, err
		}
		nippoList = append(nippoList, childNippoList...)
	}
	return nippoList, nil
}

func (r *remoteNippoQuery) list(param *i.QueryListParam, option *i.QueryListOption) (*drive.FileList, []model.Nippo, []drive.File, error) {
	res, err := r.provider.List(param)
	if err != nil {
		return nil, nil, nil, err
	}

	nippoList := []model.Nippo{}
	folderList := []drive.File{}
	for _, file := range res.Files {
		if file.MimeType == gateway.DriveFolderMimeType {
			folderList = append(folderList, *file)
		} else {
			nippo := &model.Nippo{}
			nippo.Date = model.NewNippoDate(file.Name)
			nippo.RemoteFile = file
			if option.WithContent {
				r.Download(nippo)
			}
			nippoList = append(nippoList, *nippo)
		}
	}
	return res, nippoList, folderList, nil
}

func (r *remoteNippoQuery) Download(nippo *model.Nippo) (err error) {
	nippo.Content, err = r.provider.Download(nippo.RemoteFile.Id)
	return
}

func NewLocalNippoQuery(injector do.Injector) (i.LocalNippoQuery, error) {
	provider, err := do.Invoke[gateway.LocalFileProvider](injector)
	if err != nil {
		return nil, err
	}
	return &localNippoQuery{
		provider: provider,
	}, nil
}

func (r *localNippoQuery) Exist(date *model.NippoDate) bool {
	return true
}

func (r *localNippoQuery) Find(date *model.NippoDate) (*model.Nippo, error) {
	return model.NewNippo("")
}

func (r *localNippoQuery) List(param *i.QueryListParam, option *i.QueryListOption) ([]model.Nippo, error) {
	var nippoList []model.Nippo
	for _, folder := range param.Folders {
		localParam := &i.QueryListParam{
			Folders:        []string{folder},
			FileExtensions: param.FileExtensions,
			UpdatedAt:      param.UpdatedAt,
			OrderBy:        param.OrderBy,
		}
		files, err := r.provider.List(localParam)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			nippo, err := model.NewNippo(filepath.Join(folder, file.Name()))
			if err != nil {
				return nil, err
			}
			if option.WithContent {
				r.Load(nippo)
			}
			nippoList = append(nippoList, *nippo)
		}
	}
	return nippoList, nil
}

func (r *localNippoQuery) Load(nippo *model.Nippo) (err error) {
	nippo.Content, err = r.provider.Read(nippo.FilePath)
	return
}

type localNippoCommand struct {
	provider gateway.LocalFileProvider
}

func NewLocalNippoCommand(injector do.Injector) (i.LocalNippoCommand, error) {
	provider, err := do.Invoke[gateway.LocalFileProvider](injector)
	if err != nil {
		return nil, err
	}
	return &localNippoCommand{
		provider: provider,
	}, nil
}

func (r *localNippoCommand) Create(nippo *model.Nippo) error {
	cacheDir := filepath.Join(core.Cfg.GetCacheDir(), "md")
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	filePath := filepath.Join(cacheDir, fmt.Sprintf("%v.md", nippo.Date.FileString()))
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(nippo.Content)
	if err != nil {
		return err
	}
	nippo.FilePath = filePath
	return nil
}

func (r *localNippoCommand) Delete(nippo *model.Nippo) error {
	return nil
}
