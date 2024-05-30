package service

import (
	"fmt"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/storage"
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
	articleStorage storage.Article
}

// Article is an interface that defines the methods for interacting with the article storage.
type Article interface {
	GetAll() ([]model.Article, error)
	GetById(id int) (model.Article, error)
	Create(article model.Article) (model.Article, error)
	Delete(id int) error
	GetBySource(source string) ([]model.Article, error)
	GetByKeyword(keyword string) ([]model.Article, error)
	GetByDateInRange(startDate, endDate string) ([]model.Article, error)
}

func New(articleRepo storage.Article) *ArticleService {
	return &ArticleService{articleStorage: articleRepo}
}

// GetAll returns all articles in the database.
func (a *ArticleService) GetAll() ([]model.Article, error) {
	return a.articleStorage.GetAll()
}

// GetById returns the article with the given ID.
func (a *ArticleService) GetById(id int) (model.Article, error) {
	return a.articleStorage.GetById(id)
}

// Create adds a new article to the database.
func (a *ArticleService) Create(article model.Article) (model.Article, error) {
	return a.articleStorage.Create(article)
}

// Delete removes the article with the given ID from the database.
func (a *ArticleService) Delete(id int) error {
	return a.articleStorage.Delete(id)
}

// GetBySource returns all articles from the given source.
func (a *ArticleService) GetBySource(source string) ([]model.Article, error) {
	switch source {
	case "abcnews":
		return a.articleStorage.GetBySource(ABCNewsSource)
	case "bbc":
		return a.articleStorage.GetBySource(BBCNewsSource)
	case "washingtontimes":
		return a.articleStorage.GetBySource(WashingtonTimesSource)
	case "nbc":
		return a.articleStorage.GetBySource(NBCNewsSource)
	case "usatoday":
		return a.articleStorage.GetBySource(USATodaySource)
	default:
		return nil, fmt.Errorf("source not found")
	}
}

// GetByKeyword returns all articles that contain the given keyword.
func (a *ArticleService) GetByKeyword(keyword string) ([]model.Article, error) {
	return a.articleStorage.GetByKeyword(keyword)
}

// GetByDateInRange returns all articles published between the given start and end dates.
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
	return a.articleStorage.GetByDateInRange(startDateObj, endDateObj)
}
