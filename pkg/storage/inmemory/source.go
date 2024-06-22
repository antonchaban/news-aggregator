package inmemory

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
)

type memorySourceStorage struct {
	Sources []model.Source
	nextID  int
}

func (m *memorySourceStorage) GetAll() ([]model.Source, error) {
	return m.Sources, nil
}

func (m *memorySourceStorage) Save(src model.Source) (model.Source, error) {
	for _, s := range m.Sources {
		if s.Link == src.Link {
			return model.Source{}, errors.New("source already exists")
		}
	}
	src.Id = m.nextID
	m.nextID++
	m.Sources = append(m.Sources, src)
	return src, nil
}

func (m *memorySourceStorage) SaveAll(sources []model.Source) error {
	//TODO implement me
	panic("implement me")
}

func (m *memorySourceStorage) Delete(id int) error {
	for i, s := range m.Sources {
		if s.Id == id {
			m.Sources = append(m.Sources[:i], m.Sources[i+1:]...)
			return nil
		}
	}
	return errors.New("source not found")
}

func (m *memorySourceStorage) GetByID(id int) (model.Source, error) {
	for _, s := range m.Sources {
		if s.Id == id {
			return s, nil
		}
	}
	return model.Source{}, errors.New("source not found")
}

func NewSrc() storage.SourceStorage {
	return &memorySourceStorage{
		Sources: []model.Source{},
		nextID:  1, // Initializing IDs for inmemory storage, and then auto-incrementing it after saving an article
	}
}
