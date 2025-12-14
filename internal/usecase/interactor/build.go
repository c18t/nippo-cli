package interactor

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/c18t/nippo-cli/internal/adapter/gateway"
	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/model"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/c18t/nippo-cli/internal/domain/service"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/carlosstrand/go-sitemap"
	"github.com/gorilla/feeds"
	"github.com/samber/do/v2"
)

type buildCommandInteractor struct {
	assetRepository repository.AssetRepository      `do:""`
	localNippoQuery repository.LocalNippoQuery      `do:""`
	nippoService    service.NippoFacade             `do:""`
	templateService service.TemplateService         `do:""`
	fileProvider    gateway.LocalFileProvider       `do:""`
	presenter       presenter.BuildCommandPresenter `do:""`
}

func NewBuildCommandInteractor(i do.Injector) (port.BuildCommandUseCase, error) {
	assetRepository, err := do.Invoke[repository.AssetRepository](i)
	if err != nil {
		return nil, err
	}
	localNippoQuery, err := do.Invoke[repository.LocalNippoQuery](i)
	if err != nil {
		return nil, err
	}
	nippoService, err := do.Invoke[service.NippoFacade](i)
	if err != nil {
		return nil, err
	}
	templateService, err := do.Invoke[service.TemplateService](i)
	if err != nil {
		return nil, err
	}
	fileProvider, err := do.Invoke[gateway.LocalFileProvider](i)
	if err != nil {
		return nil, err
	}
	p, err := do.Invoke[presenter.BuildCommandPresenter](i)
	if err != nil {
		return nil, err
	}
	return &buildCommandInteractor{
		assetRepository: assetRepository,
		localNippoQuery: localNippoQuery,
		nippoService:    nippoService,
		templateService: templateService,
		fileProvider:    fileProvider,
		presenter:       p,
	}, nil
}

func (u *buildCommandInteractor) Handle(input *port.BuildCommandUseCaseInputData) {
	downloadedFiles, err := u.downloadNippo()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	var buildError error

	if err := u.assetRepository.CleanBuildCache(); err != nil {
		buildError = err
	}

	if buildError == nil {
		if err := u.buildIndexPage(); err != nil {
			buildError = err
		}
	}

	if buildError == nil {
		if err := u.buildNippoPage(); err != nil {
			buildError = err
		}
	}

	if buildError == nil {
		if err := u.buildArchivePage(); err != nil {
			buildError = err
		}
	}

	if buildError == nil {
		if err := u.buildFeed(); err != nil {
			buildError = err
		}
	}

	if buildError == nil {
		if err := u.buildSiteMap(); err != nil {
			buildError = err
		}
	}

	// Show summary (downloaded files and any build errors)
	u.presenter.Summary(downloadedFiles, nil, buildError)

	if buildError != nil {
		return
	}
}

func (u *buildCommandInteractor) downloadNippo() ([]presenter.FileInfo, error) {
	// Show spinner while fetching file list
	u.presenter.Progress(&port.BuildCommandUseCaseOutputData{Message: "Fetching file list from Google Drive..."})

	started := false
	var downloadedFiles []presenter.FileInfo

	// Use configured drive folder ID
	driveFolderId := core.Cfg.Project.DriveFolderId
	if driveFolderId == "" {
		return nil, fmt.Errorf("drive folder ID is not configured. Run `nippo init` to configure")
	}

	_, err := u.nippoService.Send(&service.NippoFacadeRequest{
		Action: service.NippoFacadeActionSearch | service.NippoFacadeActionDownload | service.NippoFacadeActionCache,
		Query: &repository.QueryListParam{
			Folders:        []string{driveFolderId},
			FileExtensions: []string{"md"},
			UpdatedAt:      core.Cfg.LastUpdateCheckTimestamp,
			OrderBy:        "name",
		},
		Option: &repository.QueryListOption{
			Recursive: true,
		},
	}, &service.NippoFacadeOption{
		OnProgress: func(filename string, fileId string, current int, total int) bool {
			if !started {
				// Stop the "fetching" spinner and start build progress
				u.presenter.StopProgress()
				u.presenter.StartBuildProgress(total)
				started = true
			}
			u.presenter.UpdateBuildProgress(filename, fileId)
			downloadedFiles = append(downloadedFiles, presenter.FileInfo{Name: filename, Id: fileId})
			// Return false if user cancelled
			return !u.presenter.IsBuildCancelled()
		},
	})
	if started {
		u.presenter.StopBuildProgress()
	} else {
		// No files to download, stop the spinner
		u.presenter.StopProgress()
	}
	if err != nil {
		if err == service.ErrCancelled {
			// User cancelled, exit silently
			return downloadedFiles, err
		}
		return downloadedFiles, err
	}
	core.Cfg.LastUpdateCheckTimestamp = time.Now()
	return downloadedFiles, core.Cfg.SaveConfig()
}

type OpenGraph struct {
	Url         string
	Title       string
	Description string
	ImageUrl    string
}

// getSiteUrl returns the configured site URL or error if not configured
func getSiteUrl() (string, error) {
	if core.Cfg.Project.SiteUrl == "" {
		return "", fmt.Errorf("site URL is not configured. Run `nippo init` to configure")
	}
	return strings.TrimSuffix(core.Cfg.Project.SiteUrl, "/"), nil
}

// page content
type Content struct {
	Url         string
	PageTitle   string
	Description string
	Date        string
	Og          OpenGraph
	Content     template.HTML
}

type Archive struct {
	Url         string
	PageTitle   string
	Description string
	Date        string
	Og          OpenGraph
	Calender    *model.Calender
}

func (u *buildCommandInteractor) buildIndexPage() error {
	cacheDir := filepath.Join(core.Cfg.GetCacheDir(), "md")
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "output")

	nippoList, err := u.localNippoQuery.List(&repository.QueryListParam{
		Folders: []string{cacheDir},
	}, &repository.QueryListOption{})
	if err != nil {
		return err
	}
	nippo := nippoList[len(nippoList)-1]
	nippoHtml, err := nippo.GetHtml()
	if err != nil {
		return err
	}

	siteUrl, err := getSiteUrl()
	if err != nil {
		return err
	}
	err = u.templateService.SaveTo(filepath.Join(outputDir, "index.html"), "index", Content{
		Url:         siteUrl + "/",
		Date:        nippo.Date.TitleString(),
		Description: "ɯ̹t͡ɕʲi's daily reports.",
		Og: OpenGraph{
			Url:         siteUrl + "/",
			Title:       "日報 - nippo.c18t.me",
			Description: "ɯ̹t͡ɕʲi's daily reports.",
			ImageUrl:    siteUrl + "/nippo_ogp.png",
		},
		Content: template.HTML(nippoHtml),
	})
	return err
}

func (u *buildCommandInteractor) buildNippoPage() error {
	cacheDir := filepath.Join(core.Cfg.GetCacheDir(), "md")
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "output")

	nippoList, err := u.localNippoQuery.List(&repository.QueryListParam{
		Folders: []string{cacheDir},
	}, &repository.QueryListOption{
		WithContent: true,
	})
	if err != nil {
		return err
	}
	siteUrl, err := getSiteUrl()
	if err != nil {
		return err
	}
	for _, nippo := range nippoList {
		nippoHtml, err := nippo.GetHtml()
		if err != nil {
			return err
		}

		nippoFile := fmt.Sprintf("%v.html", nippo.Date.PathString())
		err = u.templateService.SaveTo(filepath.Join(outputDir, nippoFile), "nippo", Content{
			Url:         siteUrl + "/" + nippo.Date.PathString(),
			PageTitle:   nippo.Date.FileString(),
			Description: "ɯ̹t͡ɕʲi's daily report for " + nippo.Date.FileString() + ".",
			Date:        nippo.Date.TitleString(),
			Og: OpenGraph{
				Url:         siteUrl + "/" + nippo.Date.PathString(),
				Title:       nippo.Date.FileString() + " / 日報 - nippo.c18t.me",
				Description: "ɯ̹t͡ɕʲi's daily report for " + nippo.Date.FileString() + ".",
				ImageUrl:    siteUrl + "/nippo_ogp.png",
			},
			Content: template.HTML(nippoHtml),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *buildCommandInteractor) buildArchivePage() error {
	cacheDir := filepath.Join(core.Cfg.GetCacheDir(), "md")
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "output")

	nippoList, err := u.localNippoQuery.List(&repository.QueryListParam{
		Folders: []string{cacheDir},
	}, &repository.QueryListOption{
		WithContent: true,
	})
	if err != nil {
		return err
	}
	var monthMap = map[string]bool{}
	for _, nippo := range nippoList {
		month := nippo.Date.FileString()[:7]
		monthMap[month] = true
	}

	siteUrl, err := getSiteUrl()
	if err != nil {
		return err
	}
	for key := range monthMap {
		month, err := model.NewCalenderYearMonth(key)
		if err != nil {
			return err
		}

		calender, err := model.NewCalender(month, nippoList)
		if err != nil {
			return err
		}

		archiveFile := fmt.Sprintf("%04d%02d.html", calender.YearMonth.Year, calender.YearMonth.Month)

		err = u.templateService.SaveTo(filepath.Join(outputDir, archiveFile), "calender", Archive{
			Url:         siteUrl + "/" + calender.YearMonth.PathString(),
			PageTitle:   calender.YearMonth.FileString(),
			Description: "ɯ̹t͡ɕʲi's daily reports for " + calender.YearMonth.FileString() + ".",
			Date:        calender.YearMonth.TitleString(),
			Og: OpenGraph{
				Url:         siteUrl + "/" + calender.YearMonth.PathString(),
				Title:       calender.YearMonth.FileString() + " / 日報 - nippo.c18t.me",
				Description: "ɯ̹t͡ɕʲi's daily reports for " + calender.YearMonth.FileString() + ".",
				ImageUrl:    siteUrl + "/nippo_ogp.png",
			},
			Calender: calender,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *buildCommandInteractor) buildFeed() error {
	cacheDir := filepath.Join(core.Cfg.GetCacheDir(), "md")
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "output")

	siteUrl, err := getSiteUrl()
	if err != nil {
		return err
	}
	author := &feeds.Author{Name: "ɯ̹t͡ɕʲi"}

	feed := &feeds.Feed{
		Title:       "日報 - nippo.c18t.me",
		Link:        &feeds.Link{Href: siteUrl},
		Description: "ɯ̹t͡ɕʲi's daily reports.",
		Author:      author,
		Created:     time.Now(),
	}

	nippoList, err := u.localNippoQuery.List(&repository.QueryListParam{
		Folders: []string{cacheDir},
	}, &repository.QueryListOption{})
	if err != nil {
		return err
	}

	// Get the last 20 nippo entries (or all if less than 20)
	startIdx := len(nippoList) - 20
	if startIdx < 0 {
		startIdx = 0
	}
	for _, nippo := range nippoList[startIdx:] {
		nippoHtml, err := nippo.GetHtml()
		if err != nil {
			return err
		}

		// Use front-matter created time if available, fallback to filename-derived date
		createdTime := nippo.GetCreatedTime()

		item := &feeds.Item{
			Title:       nippo.Date.FileString() + " / 日報 - nippo.c18t.me",
			Link:        &feeds.Link{Href: siteUrl + "/" + nippo.Date.PathString()},
			Id:          siteUrl + "/" + nippo.Date.PathString(),
			Description: "ɯ̹t͡ɕʲi's daily report for " + nippo.Date.FileString() + ".",
			Author:      author,
			Created:     createdTime,
			Content:     string(nippoHtml),
		}

		// Set updated time if available from front-matter
		updatedTime := nippo.GetUpdatedTime()
		if !updatedTime.IsZero() {
			item.Updated = updatedTime
		}

		feed.Items = append(feed.Items, item)
	}

	feed.Sort(func(i, j *feeds.Item) bool {
		return i.Created.After(j.Created)
	})

	rss, err := feed.ToAtom()
	if err != nil {
		return err
	}
	return u.fileProvider.Write(filepath.Join(outputDir, "feed.xml"), []byte(rss))
}

func (u *buildCommandInteractor) buildSiteMap() error {
	cacheDir := filepath.Join(core.Cfg.GetCacheDir(), "md")
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "output")

	// Get nippo list to extract last modified times from front-matter
	nippoList, err := u.localNippoQuery.List(&repository.QueryListParam{
		Folders: []string{cacheDir},
	}, &repository.QueryListOption{
		WithContent: true,
	})
	if err != nil {
		return err
	}

	// Build a map of pathString -> last modified time
	lastModifiedMap := make(map[string]time.Time)
	for _, nippo := range nippoList {
		// GetMarkdown() parses front-matter, so call it to populate FrontMatter
		_, _ = nippo.GetMarkdown()
		pathStr := nippo.Date.PathString()
		updatedTime := nippo.GetUpdatedTime()
		if !updatedTime.IsZero() {
			lastModifiedMap[pathStr] = updatedTime
		} else {
			// Fallback to created time
			lastModifiedMap[pathStr] = nippo.GetCreatedTime()
		}
	}

	files, err := u.fileProvider.List(&repository.QueryListParam{
		Folders:        []string{outputDir},
		FileExtensions: []string{"html"},
	})
	if err != nil {
		return err
	}

	siteUrl, err := getSiteUrl()
	if err != nil {
		return err
	}
	now := time.Now()
	sitemaps := []sitemap.Sitemap{}

	count := 0
	data := sitemap.NewSitemap([]*sitemap.SitemapItem{}, nil)
	for _, file := range files {
		count++
		if count >= 10000 {
			sitemaps = append(sitemaps, *data)
			data = sitemap.NewSitemap([]*sitemap.SitemapItem{}, nil)
			count = 0
		}

		fileName := strings.TrimSuffix(file.Name(), ".html")
		if fileName == "index" {
			data.AddItem(siteUrl+"/", now, "daily", 0.5)
		} else {
			// Use last modified time from front-matter if available
			lastMod := now
			if t, ok := lastModifiedMap[fileName]; ok {
				lastMod = t
			}
			data.AddItem(siteUrl+"/"+fileName, lastMod, "monthly", 0.5)
		}
	}
	if count > 0 {
		sitemaps = append(sitemaps, *data)
	}

	sitemapIndex := sitemap.NewSitemapIndex([]*sitemap.SitemapIndexItem{}, nil)
	for i, sitemap := range sitemaps {
		xmlString, err := sitemap.ToXMLString()
		if err != nil {
			return nil
		}
		sitemapFileName := fmt.Sprintf("sitemap_%d.xml", i+1)
		err = u.fileProvider.Write(filepath.Join(outputDir, sitemapFileName), []byte(xmlString))
		if err != nil {
			return nil
		}

		sitemapIndex.AddItem(siteUrl+"/"+sitemapFileName, now)
	}

	xmlString, err := sitemapIndex.ToXMLString()
	if err != nil {
		return nil
	}
	return u.fileProvider.Write(filepath.Join(outputDir, "sitemap_index.xml"), []byte(xmlString))
}
