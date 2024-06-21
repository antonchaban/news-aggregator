package storage

import "github.com/antonchaban/news-aggregator/pkg/model"

type SourceStorage interface {
	GetAll() ([]model.Source, error)
	Save(src model.Source) (model.Source, error)
	SaveAll(sources []model.Source) error
	Delete(id int) error
	GetByID(id int) (model.Source, error)
}
