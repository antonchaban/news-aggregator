package service

import (
	"fmt"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/repository"
	"time"
)

const (
	ABCNewsSource         = "ABC News: International"
	BBCNewsSource         = "BBC News"
	WashingtonTimesSource = "The Washington Times stories: World"
	NBCNewsSource         = "NBC News"
	USATodaySource        = "USA TODAY"
)

type ArticleService struct {
	articleRepoInMem repository.Article
}

type Article interface {
	GetAll() ([]model.Article, error)
	GetById(id int) (model.Article, error)
	Create(article model.Article) (model.Article, error)
	Delete(id int) error
	GetBySource(source string) ([]model.Article, error)
	GetByKeyword(keyword string) ([]model.Article, error)
	GetByDateInRange(startDate, endDate string) ([]model.Article, error)
}

func NewArticleService(articleRepo repository.Article) *ArticleService {
	return &ArticleService{articleRepoInMem: articleRepo}
}

func (a *ArticleService) GetAll() ([]model.Article, error) {
	return a.articleRepoInMem.GetAll()
}

func (a *ArticleService) GetById(id int) (model.Article, error) {
	return a.articleRepoInMem.GetById(id)
}

func (a *ArticleService) Create(article model.Article) (model.Article, error) {
	return a.articleRepoInMem.Create(article)
}

func (a *ArticleService) Delete(id int) error {
	return a.articleRepoInMem.Delete(id)
}

func (a *ArticleService) GetBySource(source string) ([]model.Article, error) {
	switch source {
	case "abcnews":
		return a.articleRepoInMem.GetBySource(ABCNewsSource)
	case "bbc":
		return a.articleRepoInMem.GetBySource(BBCNewsSource)
	case "washingtontimes":
		return a.articleRepoInMem.GetBySource(WashingtonTimesSource)
	case "nbc":
		return a.articleRepoInMem.GetBySource(NBCNewsSource)
	case "usatoday":
		return a.articleRepoInMem.GetBySource(USATodaySource)
	default:
		return nil, fmt.Errorf("source not found")
	}
}

func (a *ArticleService) GetByKeyword(keyword string) ([]model.Article, error) {
	return a.articleRepoInMem.GetByKeyword(keyword)
}

func (a *ArticleService) GetByDateInRange(startDate, endDate string) ([]model.Article, error) {
	var startDateObj, endDateObj time.Time
	var err error

	if startDate != "" {
		startDateObj, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse start date: %v", err)
		}
	} else {
		startDateObj = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

	}

	if endDate != "" {
		endDateObj, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse end date: %v", err)
		}
	} else {
		endDateObj = time.Now()
	}
	return a.articleRepoInMem.GetByDateInRange(startDateObj, endDateObj)
}
