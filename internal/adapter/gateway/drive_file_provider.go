package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const DriveFolderMimeType = "application/vnd.google-apps.folder"

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

	query := g.queryBuilder(param)
	listCall := fileService.List().
		Fields("nextPageToken, files(id, name, fileExtension, mimeType)").
		PageSize(100).
		Q(query)
	if param.OrderBy != "" {
		listCall = listCall.OrderBy(param.OrderBy)
	}
	if param.PageToken != "" {
		listCall = listCall.PageToken(param.PageToken)
	}

	res, err := listCall.Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}
	return res, nil
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

func (g *driveFileProvider) queryBuilder(param *repository.QueryListParam) string {
	var sb strings.Builder
	folderQuery := fmt.Sprintf("mimeType = '%s'", DriveFolderMimeType)
	sb.WriteString(folderQuery)

	var parentValue string
	if len(param.Folders) > 0 {
		var parents = make([]string, len(param.Folders))
		for i, v := range param.Folders {
			parents[i] = fmt.Sprintf("'%s'", v)
		}
		parentValue = strings.Join(parents, ", ")
		if sb.Len() > 0 {
			sb.WriteString(" and ")
		}
		sb.WriteString(fmt.Sprintf("parents in %s", parentValue))
	}
	folderQuery = sb.String()

	sb.Reset()
	if parentValue != "" {
		if sb.Len() > 0 {
			sb.WriteString(" and ")
		}
		sb.WriteString(fmt.Sprintf("parents in %s", parentValue))
	}
	if len(param.FileExtensions) != 0 {
		var exts = make([]string, len(param.FileExtensions))
		for i, v := range param.FileExtensions {
			exts[i] = fmt.Sprintf("'%s'", v)
		}
		extValue := strings.Join(exts, ", ")
		if sb.Len() > 0 {
			sb.WriteString(" and ")
		}
		sb.WriteString(fmt.Sprintf("fileExtension in %s", extValue))
	}
	if !param.UpdatedAt.IsZero() {
		modifiedTimeQuery := fmt.Sprintf("modifiedTime >= '%v'", NewDriveFileTimestamp(param.UpdatedAt))
		if sb.Len() > 0 {
			sb.WriteString(" and ")
		}
		sb.WriteString(modifiedTimeQuery)
	}
	fileQuery := sb.String()

	if fileQuery == "" {
		return fmt.Sprintf("%s and trashed != true", fileQuery)
	} else {
		return fmt.Sprintf("((%s) or (%s)) and trashed != true", folderQuery, fileQuery)
	}
}
