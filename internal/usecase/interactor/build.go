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
	return &buildCommandInteractor{
		assetRepository: do.MustInvoke[repository.AssetRepository](i),
		localNippoQuery: do.MustInvoke[repository.LocalNippoQuery](i),
		nippoService:    do.MustInvoke[service.NippoFacade](i),
		templateService: do.MustInvoke[service.TemplateService](i),
		fileProvider:    do.MustInvoke[gateway.LocalFileProvider](i),
		presenter:       do.MustInvoke[presenter.BuildCommandPresenter](i),
	}, nil
}

func (u *buildCommandInteractor) Handle(input *port.BuildCommandUseCaseInputData) {
	output := &port.BuildCommandUseCaseOutputData{}

	if err := u.downloadNippo(); err != nil {
		u.presenter.Suspend(err)
		return
	}

	if err := u.assetRepository.CleanBuildCache(); err != nil {
		u.presenter.Suspend(err)
		return
	}

	if err := u.buildIndexPage(); err != nil {
		u.presenter.Suspend(err)
		return
	}

	if err := u.buildNippoPage(); err != nil {
		u.presenter.Suspend(err)
		return
	}

	if err := u.buildArchivePage(); err != nil {
		u.presenter.Suspend(err)
		return
	}

	if err := u.buildFeed(); err != nil {
		u.presenter.Suspend(err)
		return
	}

	if err := u.buildSiteMap(); err != nil {
		u.presenter.Suspend(err)
		return
	}

	output.Message = "ok. "
	u.presenter.Complete(output)
}

func (u *buildCommandInteractor) downloadNippo() error {
	_, err := u.nippoService.Send(&service.NippoFacadeRequest{
		Action: service.NippoFacadeActionSearch | service.NippoFacadeActionDownload | service.NippoFacadeActionCache,
		Query: &repository.QueryListParam{
			Folders:        []string{"1HNSRS2tJI2t7DKP_8XQJ2NTleSH-rs4y"},
			FileExtensions: []string{"md"},
			UpdatedAt:      core.Cfg.LastUpdateCheckTimestamp,
			OrderBy:        "name",
		},
		Option: &repository.QueryListOption{
			Recursive: true,
		},
	}, &service.NippoFacadeOption{})
	if err != nil {
		return err
	}
	core.Cfg.LastUpdateCheckTimestamp = time.Now()
	return core.Cfg.SaveConfig()
}

type OpenGraph struct {
	Url         string
	Title       string
	Description string
	ImageUrl    string
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

	err = u.templateService.SaveTo(filepath.Join(outputDir, "index.html"), "index", Content{
		Url:         "https://nippo.c18t.net/",
		Date:        nippo.Date.TitleString(),
		Description: "ɯ̹t͡ɕʲi's daily reports.",
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
	for _, nippo := range nippoList {
		nippoHtml, err := nippo.GetHtml()
		if err != nil {
			return err
		}

		nippoFile := fmt.Sprintf("%v.html", nippo.Date.PathString())
		err = u.templateService.SaveTo(filepath.Join(outputDir, nippoFile), "nippo", Content{
			Url:         "https://nippo.c18t.net/" + nippo.Date.PathString(),
			PageTitle:   nippo.Date.FileString(),
			Description: "ɯ̹t͡ɕʲi's daily report for " + nippo.Date.FileString() + ".",
			Date:        nippo.Date.TitleString(),
			Og: OpenGraph{
				Url:         "https://nippo.c18t.net/" + nippo.Date.PathString(),
				Title:       nippo.Date.FileString() + " / 日報 - nippo.c18t.net",
				Description: "ɯ̹t͡ɕʲi's daily report for " + nippo.Date.FileString() + ".",
				ImageUrl:    "https://nippo.c18t.net/nippo_ogp.png",
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
			Url:         "https://nippo.c18t.net/" + calender.YearMonth.PathString(),
			PageTitle:   calender.YearMonth.FileString(),
			Description: "ɯ̹t͡ɕʲi's daily reports for " + calender.YearMonth.FileString() + ".",
			Date:        calender.YearMonth.TitleString(),
			Og: OpenGraph{
				Url:         "https://nippo.c18t.net/" + calender.YearMonth.PathString(),
				Title:       calender.YearMonth.FileString() + " / 日報 - nippo.c18t.net",
				Description: "ɯ̹t͡ɕʲi's daily reports for " + calender.YearMonth.FileString() + ".",
				ImageUrl:    "https://nippo.c18t.net/nippo_ogp.png",
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

	author := &feeds.Author{Name: "ɯ̹t͡ɕʲi"}

	feed := &feeds.Feed{
		Title:       "日報 - nippo.c18t.net",
		Link:        &feeds.Link{Href: "https://nippo.c18t.net"},
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

	for _, nippo := range nippoList[len(nippoList)-20:] {
		nippoHtml, err := nippo.GetHtml()
		if err != nil {
			return err
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Title:       nippo.Date.FileString() + " / 日報 - nippo.c18t.net",
			Link:        &feeds.Link{Href: "https://nippo.c18t.net/" + nippo.Date.PathString()},
			Id:          "https://nippo.c18t.net/" + nippo.Date.PathString(),
			Description: "ɯ̹t͡ɕʲi's daily report for " + nippo.Date.FileString() + ".",
			Author:      author,
			Created:     time.Date(nippo.Date.Year(), nippo.Date.Month(), nippo.Date.Day(), 0, 0, 0, 0, time.Local),
			Content:     string(nippoHtml),
		})
	}

	feed.Sort(func(i, j *feeds.Item) bool {
		return i.Created.After(j.Created)
	})

	rss, err := feed.ToAtom()
	if err != nil {
		return err
	}
	u.fileProvider.Write(filepath.Join(outputDir, "feed.xml"), []byte(rss))
	return nil
}

func (u *buildCommandInteractor) buildSiteMap() error {
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "output")

	files, err := u.fileProvider.List(&repository.QueryListParam{
		Folders:        []string{outputDir},
		FileExtensions: []string{"html"},
	})
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
			data.AddItem("https://nippo.c18t.net/", now, "daily", 0.5)
		} else {
			data.AddItem("https://nippo.c18t.net/"+fileName, now, "monthly", 0.5)
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

		sitemapIndex.AddItem("https://nippo.c18t.net/"+sitemapFileName, now)
	}

	xmlString, err := sitemapIndex.ToXMLString()
	if err != nil {
		return nil
	}
	u.fileProvider.Write(filepath.Join(outputDir, "sitemap_index.xml"), []byte(xmlString))
	return nil
}
