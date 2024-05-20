package service

import (
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/repository"
)

type ArticleService struct {
	articleRepoInMem repository.Article
}

type Article interface {
	GetAll() ([]model.Article, error)
	GetById(id int) (model.Article, error)
	Create(article model.Article) (model.Article, error)
	Delete(id int) error
}

func NewArticleService(articleRepo repository.Article) *ArticleService {
	return &ArticleService{articleRepoInMem: articleRepo}
}

func (a *ArticleService) GetAll() ([]model.Article, error) {
	return a.articleRepoInMem.GetAll()
}

func (a *ArticleService) GetById(id int) (model.Article, error) {
	return a.articleRepoInMem.GetById(id)
}

func (a *ArticleService) Create(article model.Article) (model.Article, error) {
	return a.articleRepoInMem.Create(article)
}

func (a *ArticleService) Delete(id int) error {
	return a.articleRepoInMem.Delete(id)
}
