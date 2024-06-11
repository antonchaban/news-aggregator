package inmemory

import (
	"errors"
	"github.com/reiver/go-porterstemmer"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/storage"
	"strings"
	"time"
)

// MemoryArticleStorage is a struct that contains the in-memory database for articles.
type memoryArticleStorage struct {
	Articles []model.Article
	nextID   int
}

func New() storage.ArticleStorage {
	return &memoryArticleStorage{
		Articles: []model.Article{},
		nextID:   1,
	}
}

// GetAll returns all articles in the database.
func (a *memoryArticleStorage) GetAll() ([]model.Article, error) {
	return a.Articles, nil
}

// Create adds a new article to the database.
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

// GetByKeyword returns all articles that contain the given keyword in their title or description.
func (a *memoryArticleStorage) GetByKeyword(keyword string) ([]model.Article, error) {
	articles := []model.Article{}
	for _, article := range a.Articles {
		normalizedTitle := strings.ToLower(article.Title)
		normalizedDesc := strings.ToLower(article.Description)

		stemmedKeyword := porterstemmer.StemString(keyword)
		stemmedTitle := porterstemmer.StemString(normalizedTitle)
		stemmedDesc := porterstemmer.StemString(normalizedDesc)

		if strings.Contains(stemmedTitle, stemmedKeyword) ||
			strings.Contains(stemmedDesc, stemmedKeyword) {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

// GetBySource returns all articles from the given source.
func (a *memoryArticleStorage) GetBySource(source string) ([]model.Article, error) {
	var articles []model.Article
	for _, article := range a.Articles {
		if article.Source == source {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

// GetByDateInRange returns all articles published between the given start and end dates.
func (a *memoryArticleStorage) GetByDateInRange(startDate, endDate time.Time) ([]model.Article, error) {
	articles := []model.Article{}
	for _, article := range a.Articles {
		if article.PubDate.After(startDate) && article.PubDate.Before(endDate) {
			articles = append(articles, article)
		}
	}
	return articles, nil
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
