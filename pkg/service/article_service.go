package service

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=mocks/mock_article.go -package=mocks github.com/antonchaban/news-aggregator/pkg/storage ArticleStorage

// ArticleStorage is an interface that defines the methods for interacting with the article storage.
//
// Expected Behaviors or Guarantees:
// - The GetAll method should return all articles available in the storage or an error if the operation fails.
// - The Save method should store the provided article and return the saved article (with updated fields such as ID) or an error if the operation fails.
// - The SaveAll method should store all provided articles and return an error if the operation fails for any reason.
// - The Delete method should remove the article with the specified ID from the storage and return an error if the operation fails.
// - The GetByKeyword method should return all articles that contain the specified keyword in their content or metadata, or an error if the operation fails.
// - The GetBySource method should return all articles from the specified source, or an error if the operation fails.
// - The GetByDateInRange method should return all articles within the specified date range, or an error if the operation fails.
//
// Common Errors or Exceptions and Handling:
// - `error`: This general error can occur in any of the methods. It should be handled by logging the error and returning an appropriate message to the user or retrying the operation if possible.
// - Data validation errors: Validate input data before attempting to save or retrieve articles. For example, check if the `id` in Delete method is a valid positive integer.
//
// Known Limitations or Restrictions:
// - The Delete method does not specify what happens if the article does not exist. It should be clarified whether it returns an error or silently succeeds.
// - The methods do not define any constraints on the size or format of the articles being saved.
//
// Usage Guidelines or Best Practices:
// - Use transactions where necessary to ensure data consistency, for example, when saving multiple articles with SaveAll.
// - Ensure proper error handling and logging to facilitate debugging and monitoring of storage operations.
// - Validate input data thoroughly before performing any operations to prevent injection attacks or corrupt data.
type ArticleStorage interface {
	GetAll() ([]model.Article, error)
	Save(article model.Article) (model.Article, error)
	SaveAll(articles []model.Article) error
	Delete(id int) error
	DeleteBySourceID(id int) error
}

type articleService struct {
	articleStorage ArticleStorage
}

func New(articleRepo ArticleStorage) web.ArticleService {
	return &articleService{articleStorage: articleRepo}
}

// SaveAll saves multiple articles to the database.
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

// GetByFilter returns all articles that match the given filters.
func (a *articleService) GetByFilter(f filter.Filters) ([]model.Article, error) {
	logrus.WithField("event_id", "get_by_filter_start").Info("Fetching articles with filter")

	// Fetch all articles initially
	articles, err := a.GetAll()
	if err != nil {
		logrus.WithField("event_id", "get_all_articles_error").Error("Error fetching all articles", err)
		return nil, err
	}
	logrus.WithField("event_id", "all_articles_fetched").Info("All articles fetched successfully")

	// Create filter handlers
	sourceFilter := &filter.SourceFilter{}
	keywordFilter := &filter.KeywordFilter{}
	dateRangeFilter := &filter.DateRangeFilter{}

	logrus.WithField("event_id", "filters_created").Info("Filter handlers created")

	// Create the chain
	sourceFilter.SetNext(keywordFilter).SetNext(dateRangeFilter)
	logrus.WithField("event_id", "filters_chained").Info("Filters chained together")

	// Start filtering
	filteredArticles, err := sourceFilter.Filter(articles, f)
	if err != nil {
		logrus.WithField("event_id", "filtering_error").Error("Error during filtering", err)
		return nil, err
	}
	logrus.WithField("event_id", "filtering_complete").Info("Filtering completed successfully")

	return filteredArticles, nil
}
