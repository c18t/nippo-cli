package interactor

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/model"
	"github.com/c18t/nippo-cli/internal/domain/service"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type buildSiteInteractor struct {
	templateService service.TemplateService
	presenter       presenter.BuildSitePresenter
}

type inBuildSiteInteractor struct {
	dig.In
	TemplateService service.TemplateService
	Presenter       presenter.BuildSitePresenter
}

func NewBuildSiteInteractor(buildDeps inBuildSiteInteractor) port.BuildSiteUsecase {
	return &buildSiteInteractor{
		templateService: buildDeps.TemplateService,
		presenter:       buildDeps.Presenter,
	}
}

func (u *buildSiteInteractor) Handle(input *port.BuildSiteUsecaseInputData) {
	output := &port.BuildSiteUsecaseOutputData{}

	err := downloadNippoData()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	err = clearBuildCache()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	err = buildIndexPage(u.templateService)
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	err = buildNippoPage(u.templateService)
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	err = buildArchivePage(u.templateService)
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	output.Message = "ok. "
	u.presenter.Complete(output)
}

type OpenGraph struct {
	Url         string
	Title       string
	Description string
	ImageUrl    string
}

// page content
type Content struct {
	PageTitle string
	Date      string
	Og        OpenGraph
	Content   template.HTML
}

type Archive struct {
	PageTitle string
	Date      string
	Og        OpenGraph
	Calender  *model.Calender
}

func buildIndexPage(ts service.TemplateService) error {
	cacheDir := path.Join(core.Cfg.GetCacheDir(), "md")
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() > files[j].Name() })
	var fileName string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName = file.Name()
		break
	}
	nippo, err := model.NewNippo(path.Join(cacheDir, fileName))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	nippoHtml, err := nippo.GetHtml()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	outputDir := path.Join(core.Cfg.GetCacheDir(), "output")
	err = os.MkdirAll(outputDir, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		return nil
	}

	f, err := os.Create(path.Join(outputDir, "index.html"))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	err = ts.SaveTo(f, "index", Content{
		Date: nippo.Date.TitleString(),
		Og: OpenGraph{
			Url:         "https://nippo.c18t.net/",
			Title:       "日報 - nippo.c18t.net",
			Description: "ɯ̹t͡ɕʲi's daily reports.",
			ImageUrl:    "https://nippo.c18t.net/nippo_ogp.png",
		},
		Content: template.HTML(nippoHtml),
	})
	return err
}

func buildNippoPage(ts service.TemplateService) error {
	cacheDir := path.Join(core.Cfg.GetCacheDir(), "md")
	outputDir := path.Join(core.Cfg.GetCacheDir(), "output")
	err := os.MkdirAll(outputDir, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		return nil
	}

	files, err := os.ReadDir(cacheDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() > files[j].Name() })
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		nippo, err := model.NewNippo(path.Join(cacheDir, file.Name()))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		nippoHtml, err := nippo.GetHtml()
		if err != nil {
			fmt.Println(err)
			return nil
		}

		f, err := os.Create(path.Join(outputDir, fmt.Sprintf("%v.html", nippo.Date.PathString())))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer f.Close()

		err = ts.SaveTo(f, "nippo", Content{
			PageTitle: nippo.Date.PathString(),
			Date:      nippo.Date.TitleString(),
			Og: OpenGraph{
				Url:         "https://nippo.c18t.net/" + nippo.Date.PathString(),
				Title:       nippo.Date.PathString() + " / 日報 - nippo.c18t.net",
				Description: "ɯ̹t͡ɕʲi's daily report for " + nippo.Date.PathString() + ".",
				ImageUrl:    "https://nippo.c18t.net/nippo_ogp.png",
			},
			Content: template.HTML(nippoHtml),
		})
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}
	return nil
}

func buildArchivePage(ts service.TemplateService) error {
	cacheDir := path.Join(core.Cfg.GetCacheDir(), "md")
	outputDir := path.Join(core.Cfg.GetCacheDir(), "output", "archive")
	err := os.MkdirAll(outputDir, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		return nil
	}

	files, err := os.ReadDir(cacheDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	month, err := time.Parse("20060102.md", files[0].Name())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	nippoList := []model.Nippo{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		nippo, err := model.NewNippo(path.Join(cacheDir, file.Name()))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		nippoList = append(nippoList, nippo)
	}

	calender, err := model.NewCalender(month, nippoList)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	archiveFile := fmt.Sprintf("%04d%02d.html", calender.YearMonth.Year, calender.YearMonth.Month)
	f, err := os.Create(path.Join(outputDir, archiveFile))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	err = ts.SaveTo(f, "calender", Archive{
		PageTitle: calender.YearMonth.PathString(),
		Date:      calender.YearMonth.TitleString(),
		Og: OpenGraph{
			Url:         "https://nippo.c18t.net/archive/" + calender.YearMonth.PathString(),
			Title:       calender.YearMonth.PathString() + " / 日報 - nippo.c18t.net",
			Description: "ɯ̹t͡ɕʲi's daily reports for " + calender.YearMonth.PathString() + ".",
			ImageUrl:    "https://nippo.c18t.net/nippo_ogp.png",
		},
		Calender: calender,
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return nil
}

// download nippo data in google drive
func downloadNippoData() error {
	b, err := os.ReadFile(path.Join(core.Cfg.GetDataDir(), "credentials.json"))
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v\n", err)
		return nil
	}

	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v\n", err)
		return nil
	}
	client := getClient(config)

	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve Drive client: %v\n", err)
	}

	r, err := srv.Files.List().
		Q("parents in '1FZEaqRa8NmuRheHjTiW-_gUP3E5Ddw2T' and fileExtension = 'md'").
		// OrderBy("modifiedTime desc").
		OrderBy("name desc").
		Fields("nextPageToken, files(id, name, fileExtension)").
		PageSize(30).Do()
	if err != nil {
		fmt.Printf("Unable to retrieve files: %v\n", err)
	}
	fmt.Println("Files:")

	cacheDir := path.Join(core.Cfg.GetCacheDir(), "md")
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		return nil
	}

	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
			err = downloadFile(srv.Files, i, cacheDir)
			if err != nil {
				fmt.Printf("download failed: %v\n", err)
				continue
			}
		}
	}

	return nil
}

func downloadFile(fsrv *drive.FilesService, i *drive.File, folderPath string) error {
	g, err := fsrv.Get(i.Id).Download()
	if err != nil {
		return err
	}
	defer g.Body.Close()

	f, err := os.Create(filepath.Join(folderPath, i.Name))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	_, err = io.Copy(f, g.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	defaultDataDir := path.Join(home, ".local", "share")
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" || !path.IsAbs(dataDir) {
		dataDir = defaultDataDir
	}
	dataDir = path.Join(dataDir, "nippo")
	tok, err := tokenFromFile(path.Join(dataDir, "token.json"))
	if err != nil {
		tok = getTokenFromWeb1(config)
		saveToken1(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb1(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken1(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
