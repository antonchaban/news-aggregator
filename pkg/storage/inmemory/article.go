package inmemory

import (
	"github.com/reiver/go-porterstemmer"
	"news-aggregator/pkg/model"
	"strings"
	"time"
)

// ArticleInMemory is a struct that contains the in-memory database for articles.
type ArticleInMemory struct {
	Articles []model.Article
	nextID   int
}

func New() *ArticleInMemory {
	return &ArticleInMemory{
		Articles: []model.Article{},
		nextID:   1,
	}
}

// GetAll returns all articles in the database.
func (a *ArticleInMemory) GetAll() ([]model.Article, error) {
	return a.Articles, nil
}

// GetById returns the article with the given ID.
func (a *ArticleInMemory) GetById(id int) (model.Article, error) {
	for _, article := range a.Articles {
		if article.Id == id {
			return article, nil
		}
	}
	return model.Article{}, nil
}

// Create adds a new article to the database.
func (a *ArticleInMemory) Create(article model.Article) (model.Article, error) {
	article.Id = a.nextID
	a.nextID++
	a.Articles = append(a.Articles, article)
	return article, nil
}

// Delete removes the article with the given ID from the database.
func (a *ArticleInMemory) Delete(id int) error {
	for i, article := range a.Articles {
		if article.Id == id {
			a.Articles = append(a.Articles[:i], a.Articles[i+1:]...)
			return nil
		}
	}
	return nil
}

// GetByKeyword returns all articles that contain the given keyword in their title or description.
func (a *ArticleInMemory) GetByKeyword(keyword string) ([]model.Article, error) {
	var articles []model.Article
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
func (a *ArticleInMemory) GetBySource(source string) ([]model.Article, error) {
	var articles []model.Article
	for _, article := range a.Articles {
		if article.Source == source {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

// GetByDateInRange returns all articles published between the given start and end dates.
func (a *ArticleInMemory) GetByDateInRange(startDate, endDate time.Time) ([]model.Article, error) {
	var articles []model.Article
	for _, article := range a.Articles {
		if article.PubDate.After(startDate) && article.PubDate.Before(endDate) {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func (a *ArticleInMemory) SaveAll(articles []model.Article) error {
	for _, article := range articles {
		_, err := a.Create(article)
		if err != nil {
			return err
		}
	}
	return nil
}
