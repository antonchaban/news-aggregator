package inmemory

import (
	"github.com/stretchr/testify/assert"
	"news-aggregator/pkg/model"
	"testing"
	"time"
)

func TestArticleInMemory_Create(t *testing.T) {
	storage := New()
	tests := []struct {
		name      string
		articles  model.Article
		expected  []model.Article
		expectErr bool
	}{
		{
			name: "Save 1 article successfully",
			articles: model.Article{
				Title:       "Article 1",
				Description: "Description 1",
				Link:        "http://link1.com",
				Source:      "Source 1",
				PubDate:     time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			},
			expected: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expectErr: false,
		},
		{
			name: "Save 2 articles successfully",
			articles: model.Article{
				Title:       "Article 2",
				Description: "Description 2",
				Link:        "http://link2.com",
				Source:      "Source 2",
				PubDate:     time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			},
			expected: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 2, Title: "Article 2", Description: "Description 2", Link: "http://link2.com", Source: "Source 2", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := storage.Save(tt.articles)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, storage.Articles)
			}
		})
	}
}

func TestArticleInMemory_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		articles  []model.Article
		expected  []model.Article
		expectErr bool
	}{
		{
			name: "Delete article successfully",
			id:   1,
			articles: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 2, Title: "Article 2", Description: "Description 2", Link: "http://link2.com", Source: "Source 2", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 3, Title: "Article 3", Description: "Description 3", Link: "http://link3.com", Source: "Source 3", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expected: []model.Article{
				{Id: 2, Title: "Article 2", Description: "Description 2", Link: "http://link2.com", Source: "Source 2", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 3, Title: "Article 3", Description: "Description 3", Link: "http://link3.com", Source: "Source 3", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expectErr: false,
		},
		{
			name:      "Delete non-existing article",
			id:        2,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemoryArticleStorage{Articles: tt.articles}
			err := storage.Delete(tt.id)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, storage.Articles)
			}
		})
	}
}

func TestArticleInMemory_GetAll(t *testing.T) {
	tests := []struct {
		name      string
		articles  []model.Article
		expected  []model.Article
		expectErr bool
	}{
		{
			name: "Get all articles successfully",
			articles: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 2, Title: "Article 2", Description: "Description 2", Link: "http://link2.com", Source: "Source 2", PubDate: time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)},
			},
			expected: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 2, Title: "Article 2", Description: "Description 2", Link: "http://link2.com", Source: "Source 2", PubDate: time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)},
			},
			expectErr: false,
		},
		{
			name:      "Get empty list of articles",
			articles:  []model.Article{},
			expected:  []model.Article{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage := &MemoryArticleStorage{Articles: tt.articles}
			articles, err := storage.GetAll()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, articles)
			}
		})
	}
}

func TestArticleInMemory_GetByDateInRange(t *testing.T) {
	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		articles  []model.Article
		expected  []model.Article
		expectErr bool
	}{
		{
			name:      "Get articles in date range successfully",
			startDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC),
			articles: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 2, Title: "Article 2", Description: "Description 2", Link: "http://link2.com", Source: "Source 2", PubDate: time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)},
			},
			expected: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expectErr: false,
		},
		{
			name:      "No articles in date range",
			startDate: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
			articles: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expected:  []model.Article{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemoryArticleStorage{Articles: tt.articles}
			articles, err := storage.GetByDateInRange(tt.startDate, tt.endDate)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, articles)
			}
		})
	}
}

func TestArticleInMemory_GetByKeyword(t *testing.T) {
	tests := []struct {
		name      string
		keyword   string
		articles  []model.Article
		expected  []model.Article
		expectErr bool
	}{
		{
			name:    "Get articles by keyword successfully",
			keyword: "Article 1",
			articles: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 2, Title: "Another Article", Description: "Description 2", Link: "http://link2.com", Source: "Source 2", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expected: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expectErr: false,
		},
		{
			name:    "No articles found with keyword",
			keyword: "NonExistingKeyword",
			articles: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expected:  []model.Article{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemoryArticleStorage{Articles: tt.articles}
			articles, err := storage.GetByKeyword(tt.keyword)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, articles)
			}
		})
	}
}

func TestArticleInMemory_GetBySource(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		articles  []model.Article
		expected  []model.Article
		expectErr bool
	}{
		{
			name:   "Get articles by source successfully",
			source: "Source 1",
			articles: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 2, Title: "Another Article", Description: "Description 2", Link: "http://link2.com", Source: "Source 2", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expected: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemoryArticleStorage{Articles: tt.articles}
			articles, err := storage.GetBySource(tt.source)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, articles)
			}
		})
	}
}

func TestArticleInMemory_SaveAll(t *testing.T) {
	tests := []struct {
		name      string
		articles  []model.Article
		expected  []model.Article
		expectErr bool
	}{
		{
			name: "Save all articles successfully",
			articles: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 2, Title: "Article 2", Description: "Description 2", Link: "http://link2.com", Source: "Source 2", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expected: []model.Article{
				{Id: 1, Title: "Article 1", Description: "Description 1", Link: "http://link1.com", Source: "Source 1", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
				{Id: 2, Title: "Article 2", Description: "Description 2", Link: "http://link2.com", Source: "Source 2", PubDate: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			},
			expectErr: false,
		},
		{
			name:      "Save empty list of articles",
			articles:  []model.Article{},
			expected:  []model.Article{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := New()
			err := storage.SaveAll(tt.articles)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, storage.Articles)
			}
		})
	}
}
