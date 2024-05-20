package repository

import "news-aggregator/pkg/model"

type ArticleInMemory struct {
	Articles []model.Article
}

type Article interface {
	GetAll() ([]model.Article, error)
	GetById(id int) (model.Article, error)
	Create(article model.Article) (model.Article, error)
	Delete(id int) error
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
