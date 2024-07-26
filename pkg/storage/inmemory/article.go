package inmemory

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/sirupsen/logrus"
)

const (
	eventArticleStorageInitialized = "article_storage_initialized"
	eventDeleteArticlesBySourceID  = "delete_articles_by_source_id"
	eventGetAllArticles            = "get_all_articles"
	eventSaveArticle               = "save_article"
	eventSaveArticleError          = "save_article_error"
	eventArticleSaved              = "article_saved"
	eventDeleteArticle             = "delete_article"
	eventArticleDeleted            = "article_deleted"
	eventDeleteArticleError        = "delete_article_error"
	eventSaveAllArticles           = "save_all_articles"
	eventSaveAllArticlesSkip       = "save_all_articles_skip"
	eventAllArticlesSaved          = "all_articles_saved"
)

// MemoryArticleStorage is a struct that contains the in-memory database for articles.
type memoryArticleStorage struct {
	Articles []model.Article
	nextID   int
}

func New() service.ArticleStorage {
	logrus.WithField("event_id", eventArticleStorageInitialized).Info("Initializing Article Storage")
	return &memoryArticleStorage{
		Articles: []model.Article{},
		nextID:   1, // Initializing IDs for in-memory storage, and then auto-incrementing it after saving an article
	}
}

// DeleteBySourceID removes all articles with the given source ID from the database.
func (a *memoryArticleStorage) DeleteBySourceID(id int) error {
	logrus.WithField("event_id", eventDeleteArticlesBySourceID).Info("Deleting articles by source ID", id)

	// New slice to store articles that don't match the source ID
	newArticles := a.Articles[:0] // create a zero-length slice with the same capacity

	for _, article := range a.Articles {
		if article.Source.Id != id {
			newArticles = append(newArticles, article)
		}
	}
	a.Articles = newArticles
	return nil
}

// GetAll returns all articles in the database.
func (a *memoryArticleStorage) GetAll() ([]model.Article, error) {
	logrus.WithField("event_id", eventGetAllArticles).Info("Fetching all articles")
	return a.Articles, nil
}

// Save adds a new article to the database.
func (a *memoryArticleStorage) Save(article model.Article) (model.Article, error) {
	logrus.WithField("event_id", eventSaveArticle).Info("Saving new article", article.Link)
	for _, art := range a.Articles {
		if art.Link == article.Link {
			logrus.WithField("event_id", eventSaveArticleError).Error("Article already exists", article.Link)
			return model.Article{}, errors.New("article already exists")
		}
	}
	article.Id = a.nextID
	a.nextID++
	a.Articles = append(a.Articles, article)
	logrus.WithFields(logrus.Fields{
		"event_id":   eventArticleSaved,
		"article_id": article.Id,
	}).Info("Article saved successfully")
	return article, nil
}

// Delete removes the article with the given ID from the database.
func (a *memoryArticleStorage) Delete(id int) error {
	logrus.WithField("event_id", eventDeleteArticle).Info("Deleting article", id)
	for i, article := range a.Articles {
		if article.Id == id {
			a.Articles = append(a.Articles[:i], a.Articles[i+1:]...)
			logrus.WithField("event_id", eventArticleDeleted).Info("Article deleted successfully", id)
			return nil
		}
	}
	logrus.WithField("event_id", eventDeleteArticleError).Error("Article not found", id)
	return errors.New("article not found")
}

func (a *memoryArticleStorage) SaveAll(articles []model.Article) error {
	logrus.WithField("event_id", eventSaveAllArticles).Info("Saving multiple articles")
	for _, article := range articles {
		_, err := a.Save(article)
		if err != nil {
			if errors.Is(err, errors.New("article already exists")) {
				logrus.WithField("event_id", eventSaveAllArticlesSkip).Warn("Article already exists, skipping", article.Link)
				continue
			}
		}
	}
	logrus.WithField("event_id", eventAllArticlesSaved).Info("All articles processed")
	return nil
}

func (a *memoryArticleStorage) GetByFilter(query string, args []interface{}) ([]model.Article, error) {
	logrus.WithField("event_id", "get_by_filter_not_supported").Warn("GetByFilter operation is not supported in in-memory storage")
	return nil, errors.New("GetByFilter operation is not supported in in-memory storage")
}
