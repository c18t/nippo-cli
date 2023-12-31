package interactor

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/domain/model"
	"github.com/c18t/nippo-cli/internal/domain/repository"
	"github.com/c18t/nippo-cli/internal/domain/service"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"go.uber.org/dig"
)

type buildSiteInteractor struct {
	assetRepository repository.AssetRepository
	localNippoQuery repository.LocalNippoQuery
	nippoService    service.NippoFacade
	templateService service.TemplateService
	presenter       presenter.BuildSitePresenter
}

type inBuildSiteInteractor struct {
	dig.In
	AssetRepository repository.AssetRepository
	LocalNippoQuery repository.LocalNippoQuery
	NippoService    service.NippoFacade
	TemplateService service.TemplateService
	Presenter       presenter.BuildSitePresenter
}

func NewBuildSiteInteractor(buildDeps inBuildSiteInteractor) port.BuildSiteUsecase {
	return &buildSiteInteractor{
		assetRepository: buildDeps.AssetRepository,
		localNippoQuery: buildDeps.LocalNippoQuery,
		nippoService:    buildDeps.NippoService,
		templateService: buildDeps.TemplateService,
		presenter:       buildDeps.Presenter,
	}
}

func (u *buildSiteInteractor) Handle(input *port.BuildSiteUsecaseInputData) {
	output := &port.BuildSiteUsecaseOutputData{}

	err := u.downloadNippo()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	err = u.assetRepository.CleanBuildCache()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	err = u.buildIndexPage()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	err = u.buildNippoPage()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	err = u.buildArchivePage()
	if err != nil {
		u.presenter.Suspend(err)
		return
	}

	output.Message = "ok. "
	u.presenter.Complete(output)
}

func (u *buildSiteInteractor) downloadNippo() error {
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

func (u *buildSiteInteractor) buildIndexPage() error {
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

func (u *buildSiteInteractor) buildNippoPage() error {
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
			PageTitle: nippo.Date.FileString(),
			Date:      nippo.Date.TitleString(),
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

func (u *buildSiteInteractor) buildArchivePage() error {
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
			PageTitle: calender.YearMonth.FileString(),
			Date:      calender.YearMonth.TitleString(),
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
