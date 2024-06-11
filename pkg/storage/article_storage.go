package storage

import (
	"news-aggregator/pkg/model"
	"time"
)

//go:generate mockgen -destination=mocks/mock_article.go -package=mocks news-aggregator/pkg/storage ArticleStorage

// ArticleStorage is an interface that defines the methods for interacting with the article storage.
//
// Key Responsibilities:
// 1. Retrieve all stored articles.
// 2. Save a single article to the storage.
// 3. Save multiple articles to the storage.
// 4. Delete an article from the storage by its ID.
// 5. Retrieve articles that match a given keyword.
// 6. Retrieve articles from a specific source.
// 7. Retrieve articles within a specified date range.
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
// - Use transactions where necessary to ensure data consistency, for example when saving multiple articles with SaveAll.
// - Ensure proper error handling and logging to facilitate debugging and monitoring of storage operations.
// - Validate input data thoroughly before performing any operations to prevent injection attacks or corrupt data.
type ArticleStorage interface {
	GetAll() ([]model.Article, error)
	Save(article model.Article) (model.Article, error)
	SaveAll(articles []model.Article) error
	Delete(id int) error
	GetByKeyword(keyword string) ([]model.Article, error)
	GetBySource(source string) ([]model.Article, error)
	GetByDateInRange(startDate, endDate time.Time) ([]model.Article, error)
}
