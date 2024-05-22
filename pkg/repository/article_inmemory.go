package repository

import (
	"github.com/reiver/go-porterstemmer"
	"news-aggregator/pkg/model"
	"strings"
	"time"
)

type ArticleInMemory struct {
	Articles []model.Article
}

type Article interface {
	GetAll() ([]model.Article, error)
	GetById(id int) (model.Article, error)
	Create(article model.Article) (model.Article, error)
	Delete(id int) error
	GetByKeyword(keyword string) ([]model.Article, error)
	GetBySource(source string) ([]model.Article, error)
	GetByDateInRange(startDate, endDate time.Time) ([]model.Article, error)
}

func NewArticleInMemory(db []model.Article) *ArticleInMemory {
	return &ArticleInMemory{Articles: db}
}

func (a *ArticleInMemory) GetAll() ([]model.Article, error) {
	return a.Articles, nil
}

func (a *ArticleInMemory) GetById(id int) (model.Article, error) {
	for _, article := range a.Articles {
		if article.Id == id {
			return article, nil
		}
	}
	return model.Article{}, nil
}

func (a *ArticleInMemory) Create(article model.Article) (model.Article, error) {
	a.Articles = append(a.Articles, article)
	return article, nil
}

func (a *ArticleInMemory) Delete(id int) error {
	for i, article := range a.Articles {
		if article.Id == id {
			a.Articles = append(a.Articles[:i], a.Articles[i+1:]...)
			return nil
		}
	}
	return nil
}

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

func (a *ArticleInMemory) GetBySource(source string) ([]model.Article, error) {
	var articles []model.Article
	for _, article := range a.Articles {
		if article.Source == source {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func (a *ArticleInMemory) GetByDateInRange(startDate, endDate time.Time) ([]model.Article, error) {
	var articles []model.Article
	for _, article := range a.Articles {
		if article.PubDate.After(startDate) && article.PubDate.Before(endDate) {
			articles = append(articles, article)
		}
	}
	return articles, nil
}
