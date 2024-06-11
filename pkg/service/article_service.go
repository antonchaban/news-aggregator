package service

import (
	"errors"
	"fmt"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/storage"
	"time"
)

const (
	abcNewsSource         = "ABC News: International"
	bbcNewsSource         = "BBC News"
	washingtonTimesSource = "The Washington Times stories: World"
	nbcNewsSource         = "NBC News"
	usaTodaySource        = "USA TODAY"
)

//go:generate mockgen -destination=../service/mocks/mock_article_service.go -package=mocks news-aggregator/pkg/service ArticleService

// ArticleService is an interface that defines the methods for interacting with the article storage.
type ArticleService interface {
	GetAll() ([]model.Article, error)
	Create(article model.Article) (model.Article, error)
	Delete(id int) error
	GetBySource(source string) ([]model.Article, error)
	GetByKeyword(keyword string) ([]model.Article, error)
	GetByDateInRange(startDate, endDate string) ([]model.Article, error)
	SaveAll(articles []model.Article) error
}

type articleService struct {
	articleStorage storage.ArticleStorage
}

func New(articleRepo storage.ArticleStorage) ArticleService {
	return &articleService{articleStorage: articleRepo}
}

func (a *articleService) SaveAll(articles []model.Article) error {
	err := a.articleStorage.SaveAll(articles)
	if err != nil {
		return errors.New("failed to save articles")
	}
	return err
}

// GetAll returns all articles in the database.
func (a *articleService) GetAll() ([]model.Article, error) {
	return a.articleStorage.GetAll()
}

// Create adds a new article to the database.
func (a *articleService) Create(article model.Article) (model.Article, error) {
	return a.articleStorage.Save(article)
}

// Delete removes the article with the given ID from the database.
func (a *articleService) Delete(id int) error {
	return a.articleStorage.Delete(id)
}

// GetBySource returns all articles from the given source.
func (a *articleService) GetBySource(source string) ([]model.Article, error) {
	switch source {
	case "abcnews":
		return a.articleStorage.GetBySource(abcNewsSource)
	case "bbc":
		return a.articleStorage.GetBySource(bbcNewsSource)
	case "washingtontimes":
		return a.articleStorage.GetBySource(washingtonTimesSource)
	case "nbc":
		return a.articleStorage.GetBySource(nbcNewsSource)
	case "usatoday":
		return a.articleStorage.GetBySource(usaTodaySource)
	default:
		return nil, fmt.Errorf("source not found")
	}
}

// GetByKeyword returns all articles that contain the given keyword.
func (a *articleService) GetByKeyword(keyword string) ([]model.Article, error) {
	return a.articleStorage.GetByKeyword(keyword)
}

// GetByDateInRange returns all articles published between the given start and end dates.
func (a *articleService) GetByDateInRange(startDate, endDate string) ([]model.Article, error) {
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
