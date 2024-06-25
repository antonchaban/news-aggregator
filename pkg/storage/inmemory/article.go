package inmemory

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/sirupsen/logrus"
)

// MemoryArticleStorage is a struct that contains the in-memory database for articles.
type memoryArticleStorage struct {
	Articles []model.Article
	nextID   int
}

func New() storage.ArticleStorage {
	logrus.WithField("event_id", "article_storage_initialized").Info("Initializing Article Storage")
	return &memoryArticleStorage{
		Articles: []model.Article{},
		nextID:   1, // Initializing IDs for in-memory storage, and then auto-incrementing it after saving an article
	}
}

// GetAll returns all articles in the database.
func (a *memoryArticleStorage) GetAll() ([]model.Article, error) {
	logrus.WithField("event_id", "get_all_articles").Info("Fetching all articles")
	return a.Articles, nil
}

// Save adds a new article to the database.
func (a *memoryArticleStorage) Save(article model.Article) (model.Article, error) {
	logrus.WithField("event_id", "save_article").Info("Saving new article", article.Link)
	for _, art := range a.Articles {
		if art.Link == article.Link {
			logrus.WithField("event_id", "save_article_error").Error("Article already exists", article.Link)
			return model.Article{}, errors.New("article already exists")
		}
	}
	article.Id = a.nextID
	a.nextID++
	a.Articles = append(a.Articles, article)
	logrus.WithFields(logrus.Fields{
		"event_id":   "article_saved",
		"article_id": article.Id,
	}).Info("Article saved successfully")
	return article, nil
}

// Delete removes the article with the given ID from the database.
func (a *memoryArticleStorage) Delete(id int) error {
	logrus.WithField("event_id", "delete_article").Info("Deleting article", id)
	for i, article := range a.Articles {
		if article.Id == id {
			a.Articles = append(a.Articles[:i], a.Articles[i+1:]...)
			logrus.WithField("event_id", "article_deleted").Info("Article deleted successfully", id)
			return nil
		}
	}
	logrus.WithField("event_id", "delete_article_error").Error("Article not found", id)
	return errors.New("article not found")
}

func (a *memoryArticleStorage) SaveAll(articles []model.Article) error {
	logrus.WithField("event_id", "save_all_articles").Info("Saving multiple articles")
	for _, article := range articles {
		_, err := a.Save(article)
		if err != nil {
			if errors.Is(err, errors.New("article already exists")) {
				logrus.WithField("event_id", "save_all_articles_skip").Warn("Article already exists, skipping", article.Link)
				continue
			}
			logrus.WithField("event_id", "save_all_articles_error").Error("Error saving article", err)
			return err
		}
	}
	logrus.WithField("event_id", "all_articles_saved").Info("All articles processed")
	return nil
}
