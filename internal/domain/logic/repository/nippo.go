package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/model"
	i "github.com/c18t/nippo-cli/internal/domain/repository"
)

type remoteNippoQuery struct {
	provider gateway.DriveFileProvider
}

type localNippoQuery struct {
	provider gateway.LocalFileProvider
}

func NewRemoteNippoQuery(p gateway.DriveFileProvider) i.RemoteNippoQuery {
	return &remoteNippoQuery{p}
}

func (r *remoteNippoQuery) List(param *i.QueryListParam, option *i.QueryListOption) ([]model.Nippo, error) {
	res, err := r.provider.List(param)
	if err != nil {
		return nil, err
	}

	nippoList := make([]model.Nippo, len(res.Files))
	for i, file := range res.Files {
		nippoList[i].Date = model.NewNippoDate(file.Name)
		nippoList[i].RemoteFile = file
	}
	return nippoList, nil
}

func (r *remoteNippoQuery) Download(nippo *model.Nippo) (err error) {
	nippo.Content, err = r.provider.Download(nippo.RemoteFile.Id)
	return
}

func NewLocalNippoQuery(p gateway.LocalFileProvider) i.LocalNippoQuery {
	return &localNippoQuery{p}
}

func (r *localNippoQuery) Exist(date *model.NippoDate) bool {
	return true
}

func (r *localNippoQuery) Find(date *model.NippoDate) (*model.Nippo, error) {
	return model.NewNippo("")
}

func (r *localNippoQuery) List(param *i.QueryListParam, option *i.QueryListOption) ([]model.Nippo, error) {
	files, err := r.provider.List(param)
	if err != nil {
		return nil, err
	}

	var nippoList []model.Nippo
	for _, file := range files {
		nippo, err := model.NewNippo(filepath.Join(param.Folder, file.Name()))
		if err != nil {
			return nil, err
		}
		if option.WithContent {
			r.Load(nippo)
		}
		nippoList = append(nippoList, *nippo)
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

func NewLocalNippoCommand(p gateway.LocalFileProvider) i.LocalNippoCommand {
	return &localNippoCommand{p}
}

func (r *localNippoCommand) Create(nippo *model.Nippo) error {
	cacheDir := filepath.Join(core.Cfg.GetCacheDir(), "md")
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	filePath := filepath.Join(cacheDir, fmt.Sprintf("%v.md", nippo.Date))
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
