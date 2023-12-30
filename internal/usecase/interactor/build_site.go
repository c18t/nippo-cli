package interactor

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"sort"
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
	nippoService    service.NippoFacade
	templateService service.TemplateService
	presenter       presenter.BuildSitePresenter
}

type inBuildSiteInteractor struct {
	dig.In
	AssetRepository repository.AssetRepository
	NippoService    service.NippoFacade
	TemplateService service.TemplateService
	Presenter       presenter.BuildSitePresenter
}

func NewBuildSiteInteractor(buildDeps inBuildSiteInteractor) port.BuildSiteUsecase {
	return &buildSiteInteractor{
		assetRepository: buildDeps.AssetRepository,
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
			Folder:        "1FZEaqRa8NmuRheHjTiW-_gUP3E5Ddw2T",
			FileExtension: "md",
			UpdatedAt:     core.Cfg.LastUpdateCheckTimestamp,
			OrderBy:       "name desc",
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
		nippo, err := model.NewNippo(filepath.Join(cacheDir, file.Name()))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		nippoHtml, err := nippo.GetHtml()
		if err != nil {
			fmt.Println(err)
			return nil
		}

		nippoFile := fmt.Sprintf("%v.html", nippo.Date.PathString())
		err = u.templateService.SaveTo(filepath.Join(outputDir, nippoFile), "nippo", Content{
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

func (u *buildSiteInteractor) buildArchivePage() error {
	cacheDir := filepath.Join(core.Cfg.GetCacheDir(), "md")
	outputDir := filepath.Join(core.Cfg.GetCacheDir(), "output", "archive")

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
		nippo, err := model.NewNippo(filepath.Join(cacheDir, file.Name()))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		nippoList = append(nippoList, *nippo)
	}

	calender, err := model.NewCalender(month, nippoList)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	archiveFile := fmt.Sprintf("%04d%02d.html", calender.YearMonth.Year, calender.YearMonth.Month)

	err = u.templateService.SaveTo(filepath.Join(outputDir, archiveFile), "calender", Archive{
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
