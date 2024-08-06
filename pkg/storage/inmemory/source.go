package inmemory

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/sirupsen/logrus"
)

const (
	eventSourceStorageInitialized = "source_storage_initialized"
	eventGetAllSources            = "get_all_sources"
	eventSaveSource               = "save_source"
	eventSaveSourceError          = "save_source_error"
	eventSourceSaved              = "source_saved"
	eventSaveAllSources           = "save_all_sources"
	eventSaveAllSourcesSkip       = "save_all_sources_skip"
	eventSaveAllSourcesError      = "save_all_sources_error"
	eventAllSourcesSaved          = "all_sources_saved"
	eventDeleteSource             = "delete_source"
	eventSourceDeleted            = "source_deleted"
	eventDeleteSourceError        = "delete_source_error"
	eventGetSourceByID            = "get_source_by_id"
	eventGetSourceByIDError       = "get_source_by_id_error"
	eventUpdateSource             = "update_source"
	eventSourceUpdated            = "source_updated"
	eventUpdateSourceError        = "update_source_error"
)

// memorySourceStorage represents an in-memory storage for sources.
type memorySourceStorage struct {
	Sources []model.Source
	nextID  int
}

// NewSrc creates a new instance of the in-memory source storage.
func NewSrc() service.SourceStorage {
	logrus.WithField("event_id", eventSourceStorageInitialized).Info("Initializing Source Storage")
	return &memorySourceStorage{
		Sources: []model.Source{},
		nextID:  1, // Initializing IDs for in-memory storage, and then auto-incrementing it after saving a source
	}
}

// GetAll returns all sources from the in-memory storage.
func (m *memorySourceStorage) GetAll() ([]model.Source, error) {
	logrus.WithField("event_id", eventGetAllSources).Info("Fetching all sources")
	return m.Sources, nil
}

// Save saves a new source to the in-memory storage, if it is not a duplicate.
func (m *memorySourceStorage) Save(src model.Source) (model.Source, error) {
	logrus.WithField("event_id", eventSaveSource).Info("Saving new source", src.Link)
	for _, s := range m.Sources {
		if s.Link == src.Link {
			logrus.WithField("event_id", eventSaveSourceError).Error("Source already exists", src.Link)
			return model.Source{}, errors.New("source already exists")
		}
	}
	src.Id = m.nextID
	m.nextID++
	m.Sources = append(m.Sources, src)
	logrus.WithFields(logrus.Fields{
		"event_id":  eventSourceSaved,
		"source_id": src.Id,
	}).Info("Source saved successfully")
	return src, nil
}

// SaveAll saves multiple sources to the in-memory storage.
func (m *memorySourceStorage) SaveAll(sources []model.Source) error {
	logrus.WithField("event_id", eventSaveAllSources).Info("Saving multiple sources")
	for _, src := range sources {
		_, err := m.Save(src)
		if err != nil {
			if errors.Is(err, errors.New("source already exists")) {
				logrus.WithField("event_id", eventSaveAllSourcesSkip).Warn("Source already exists, skipping", src.Link)
				continue
			}
			logrus.WithField("event_id", eventSaveAllSourcesError).Error("Error saving source", err)
			return err
		}
	}
	logrus.WithField("event_id", eventAllSourcesSaved).Info("All sources processed")
	return nil
}

// Delete removes a source from the in-memory storage by its ID.
func (m *memorySourceStorage) Delete(id int) error {
	logrus.WithField("event_id", eventDeleteSource).Info("Deleting source", id)
	for i, s := range m.Sources {
		if s.Id == id {
			m.Sources = append(m.Sources[:i], m.Sources[i+1:]...)
			logrus.WithField("event_id", eventSourceDeleted).Info("Source deleted successfully", id)
			return nil
		}
	}
	logrus.WithField("event_id", eventDeleteSourceError).Error("Source not found", id)
	return errors.New("source not found")
}

// GetByID retrieves a source from the in-memory storage by its ID.
func (m *memorySourceStorage) GetByID(id int) (model.Source, error) {
	logrus.WithField("event_id", eventGetSourceByID).Info("Fetching source by ID", id)
	for _, s := range m.Sources {
		if s.Id == id {
			return s, nil
		}
	}
	logrus.WithField("event_id", eventGetSourceByIDError).Error("Source not found", id)
	return model.Source{}, errors.New("source not found")
}

// Update updates a source in the in-memory storage by its ID.
func (m *memorySourceStorage) Update(id int, src model.Source) (model.Source, error) {
	logrus.WithField("event_id", eventUpdateSource).Info("Updating source", id)
	for i, s := range m.Sources {
		if s.Id == id {
			m.Sources[i] = src
			logrus.WithField("event_id", eventSourceUpdated).Info("Source updated successfully", id)
			return src, nil
		}
	}
	logrus.WithField("event_id", eventUpdateSourceError).Error("Source not found", id)
	return model.Source{}, errors.New("source not found")
}

func (m *memorySourceStorage) GetByShortName(shortName string) (model.Source, error) {
	logrus.WithField("event_id", eventGetSourceByID).Info("Fetching source by short name", shortName)
	for _, s := range m.Sources {
		if s.ShortName == shortName {
			return s, nil
		}
	}
	logrus.WithField("event_id", eventGetSourceByIDError).Error("Source not found", shortName)
	return model.Source{}, errors.New("source not found")
}
