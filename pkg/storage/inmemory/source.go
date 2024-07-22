package inmemory

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/sirupsen/logrus"
)

// memorySourceStorage represents an in-memory storage for sources.
type memorySourceStorage struct {
	Sources []model.Source
	nextID  int
}

// NewSrc creates a new instance of the in-memory source storage.
func NewSrc() service.SourceStorage {
	logrus.WithField("event_id", "source_storage_initialized").Info("Initializing Source Storage")
	return &memorySourceStorage{
		Sources: []model.Source{},
		nextID:  1, // Initializing IDs for in-memory storage, and then auto-incrementing it after saving a source
	}
}

// GetAll returns all sources from the in-memory storage.
func (m *memorySourceStorage) GetAll() ([]model.Source, error) {
	logrus.WithField("event_id", "get_all_sources").Info("Fetching all sources")
	return m.Sources, nil
}

// Save saves a new source to the in-memory storage, if it is not a duplicate.
func (m *memorySourceStorage) Save(src model.Source) (model.Source, error) {
	logrus.WithField("event_id", "save_source").Info("Saving new source", src.Link)
	for _, s := range m.Sources {
		if s.Link == src.Link {
			logrus.WithField("event_id", "save_source_error").Error("Source already exists", src.Link)
			return model.Source{}, errors.New("source already exists")
		}
	}
	src.Id = m.nextID
	m.nextID++
	m.Sources = append(m.Sources, src)
	logrus.WithFields(logrus.Fields{
		"event_id":  "source_saved",
		"source_id": src.Id,
	}).Info("Source saved successfully")
	return src, nil
}

// SaveAll saves multiple sources to the in-memory storage.
func (m *memorySourceStorage) SaveAll(sources []model.Source) error {
	logrus.WithField("event_id", "save_all_sources").Info("Saving multiple sources")
	for _, src := range sources {
		_, err := m.Save(src)
		if err != nil {
			if errors.Is(err, errors.New("source already exists")) {
				logrus.WithField("event_id", "save_all_sources_skip").Warn("Source already exists, skipping", src.Link)
				continue
			}
			logrus.WithField("event_id", "save_all_sources_error").Error("Error saving source", err)
			return err
		}
	}
	logrus.WithField("event_id", "all_sources_saved").Info("All sources processed")
	return nil
}

// Delete removes a source from the in-memory storage by its ID.
func (m *memorySourceStorage) Delete(id int) error {
	logrus.WithField("event_id", "delete_source").Info("Deleting source", id)
	for i, s := range m.Sources {
		if s.Id == id {
			m.Sources = append(m.Sources[:i], m.Sources[i+1:]...)
			logrus.WithField("event_id", "source_deleted").Info("Source deleted successfully", id)
			return nil
		}
	}
	logrus.WithField("event_id", "delete_source_error").Error("Source not found", id)
	return errors.New("source not found")
}

// GetByID retrieves a source from the in-memory storage by its ID.
func (m *memorySourceStorage) GetByID(id int) (model.Source, error) {
	logrus.WithField("event_id", "get_source_by_id").Info("Fetching source by ID", id)
	for _, s := range m.Sources {
		if s.Id == id {
			return s, nil
		}
	}
	logrus.WithField("event_id", "get_source_by_id_error").Error("Source not found", id)
	return model.Source{}, errors.New("source not found")
}

// Update updates a source in the in-memory storage by its ID.
func (m *memorySourceStorage) Update(id int, src model.Source) (model.Source, error) {
	logrus.WithField("event_id", "update_source").Info("Updating source", id)
	for i, s := range m.Sources {
		if s.Id == id {
			m.Sources[i] = src
			logrus.WithField("event_id", "source_updated").Info("Source updated successfully", id)
			return src, nil
		}
	}
	logrus.WithField("event_id", "update_source_error").Error("Source not found", id)
	return model.Source{}, errors.New("source not found")
}
