package inmemory

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
)

// MemoryArticleStorage is a struct that contains the in-memory database for articles.
type memoryArticleStorage struct {
	Articles []model.Article
	nextID   int
}

func New() storage.ArticleStorage {
	return &memoryArticleStorage{
		Articles: []model.Article{},
		nextID:   1, // Initializing IDs for inmemory storage, and then auto-incrementing it after saving an article
	}
}

// GetAll returns all articles in the database.
func (a *memoryArticleStorage) GetAll() ([]model.Article, error) {
	return a.Articles, nil
}

// Save adds a new article to the database.
func (a *memoryArticleStorage) Save(article model.Article) (model.Article, error) {
	article.Id = a.nextID
	a.nextID++
	a.Articles = append(a.Articles, article)
	return article, nil
}

// Delete removes the article with the given ID from the database.
func (a *memoryArticleStorage) Delete(id int) error {
	for i, article := range a.Articles {
		if article.Id == id {
			a.Articles = append(a.Articles[:i], a.Articles[i+1:]...)
			return nil
		}
	}
	return errors.New("article not found")
}

func (a *memoryArticleStorage) SaveAll(articles []model.Article) error {
	for _, article := range articles {
		_, err := a.Save(article)
		if err != nil {
			return err
		}
	}
	return nil
}
