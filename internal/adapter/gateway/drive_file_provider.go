package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type DriveFileTimestamp struct {
	t time.Time
}

func NewDriveFileTimestamp(t time.Time) DriveFileTimestamp {
	return DriveFileTimestamp{t}
}

func (t DriveFileTimestamp) String() string {
	return t.t.Format(time.RFC3339)
}

type DriveFileProvider interface {
	List(param *repository.QueryListParam) (*drive.FileList, error)
	Download(string) ([]byte, error)
}

type driveFileProvider struct {
	fs *drive.FilesService
}

func NewDriveFileProvider() DriveFileProvider {
	return &driveFileProvider{}
}

func (g *driveFileProvider) List(param *repository.QueryListParam) (*drive.FileList, error) {
	fileService, err := g.getFileService()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(
		"parents in '%v' and fileExtension = '%v' and modifiedTime >= '%v'",
		param.Folder, param.FileExtension, NewDriveFileTimestamp(param.UpdatedAt))
	r, err := fileService.List().
		Q(query).
		OrderBy(param.OrderBy).
		Fields("nextPageToken, files(id, name, fileExtension)").
		PageSize(50).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}
	return r, nil
}

func (g *driveFileProvider) Download(id string) ([]byte, error) {
	fileService, err := g.getFileService()
	if err != nil {
		return nil, err
	}
	res, err := fileService.Get(id).Download()
	if err != nil {
		return nil, err
	}
	defer func() { io.Copy(io.Discard, res.Body); res.Body.Close() }()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (g *driveFileProvider) getFileService() (*drive.FilesService, error) {
	if g.fs != nil {
		return g.fs, nil
	}

	dataDir := core.Cfg.GetDataDir()
	tok, err := g.tokenFromFile(filepath.Join(dataDir, "token.json"))
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(filepath.Join(core.Cfg.GetDataDir(), "credentials.json"))
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := config.Client(context.Background(), tok)
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	return srv.Files, nil
}

// Retrieves a token from a local file.
func (g *driveFileProvider) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
