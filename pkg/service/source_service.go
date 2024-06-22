package service

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/parser"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"net/url"
)

type SourceService interface {
	FetchFromAllSources() error
	FetchSourceByID(id int) ([]model.Article, error)
	AddSource(source model.Source) (model.Source, error)
	DeleteSource(id int) error
}

type sourceService struct {
	articleStorage storage.ArticleStorage
	srcStorage     storage.SourceStorage
}

func (s *sourceService) DeleteSource(id int) error {
	return s.srcStorage.Delete(id)
}

func (s *sourceService) AddSource(source model.Source) (model.Source, error) {
	save, err := s.srcStorage.Save(source)
	if err != nil {
		return model.Source{}, err
	}
	return save, nil
}

func (s *sourceService) FetchFromAllSources() error {
	allSrcs, err := s.srcStorage.GetAll()
	if err != nil {
		return err
	}

	for _, src := range allSrcs {
		urlParsed, err := url.Parse(src.Link)
		if err != nil {
			return err
		}
		articles, err := parser.ParseArticlesFromFeed(*urlParsed)
		if err != nil {
			return err
		}
		err = s.articleStorage.SaveAll(articles)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sourceService) FetchSourceByID(id int) ([]model.Article, error) {
	src, err := s.srcStorage.GetByID(id)
	if err != nil {
		return nil, err
	}
	urlParsed, err := url.Parse(src.Link)
	if err != nil {
		return nil, err
	}
	articles, err := parser.ParseArticlesFromFeed(*urlParsed)
	if err != nil {
		return nil, err
	}
	err = s.articleStorage.SaveAll(articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func NewSourceService(articleRepo storage.ArticleStorage, srcRepo storage.SourceStorage) SourceService {
	return &sourceService{articleStorage: articleRepo, srcStorage: srcRepo}
}
