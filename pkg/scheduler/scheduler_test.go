package scheduler

import (
	"fmt"
	service_mocks "github.com/antonchaban/news-aggregator/pkg/service/mocks"
	"github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewScheduler(t *testing.T) {
	asvc := new(service_mocks.MockArticleService)
	ssvc := new(service_mocks.MockSourceService)

	scheduler := NewScheduler(asvc, ssvc)

	assert.NotNil(t, scheduler)
	assert.NotNil(t, scheduler.scheduler)
	assert.Equal(t, asvc, scheduler.asvc)
	assert.Equal(t, ssvc, scheduler.ssvc)
}

func TestScheduler_Start(t *testing.T) {
	// Set up the log recorder
	recorder := &LogRecorder{}
	logrus.AddHook(recorder)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := service_mocks.NewMockArticleService(ctrl)
	mockSourceService := service_mocks.NewMockSourceService(ctrl)
	s := NewScheduler(mockArticleService, mockSourceService)

	mockSourceService.EXPECT().FetchFromAllSources().Return(nil).AnyTimes()
	s.Start()

	time.Sleep(2 * time.Second)

	s.Stop()

	// Check if the "Scheduler started successfully" log message is present
	found := false
	for _, entry := range recorder.Entries {
		if entry.Message == "Scheduler started successfully" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected log message 'Scheduler started successfully' not found")

	mockSourceService.EXPECT().FetchFromAllSources().Return(nil).AnyTimes()
}

func TestScheduler_updateArticles_Success(t *testing.T) {
	recorder := &LogRecorder{}
	logrus.AddHook(recorder)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleService := service_mocks.NewMockArticleService(ctrl)
	mockSourceService := service_mocks.NewMockSourceService(ctrl)

	s := NewScheduler(mockArticleService, mockSourceService)
	mockSourceService.EXPECT().FetchFromAllSources().Return(nil)
	recorder.Entries = nil

	// Call updateArticles directly
	s.updateArticles()

	// Assertions
	assert.Len(t, recorder.Entries, 2)
	assert.Equal(t, "Updating articles...", recorder.Entries[0].Message)
	assert.Equal(t, "Articles updated successfully", recorder.Entries[1].Message)
}

func TestScheduler_updateArticles_Error(t *testing.T) {
	recorder := &LogRecorder{}
	logrus.AddHook(recorder)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleService := service_mocks.NewMockArticleService(ctrl)
	mockSourceService := service_mocks.NewMockSourceService(ctrl)

	s := NewScheduler(mockArticleService, mockSourceService)

	mockSourceService.EXPECT().FetchFromAllSources().Return(fmt.Errorf("some error"))
	recorder.Entries = nil
	s.updateArticles()

	// Assertions
	assert.Len(t, recorder.Entries, 2)
	assert.Equal(t, "Updating articles...", recorder.Entries[0].Message)
	assert.Contains(t, recorder.Entries[1].Message, "Error occurred while fetching articles from sources: some error")
}

type LogRecorder struct {
	Entries []*logrus.Entry
}

func (r *LogRecorder) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (r *LogRecorder) Fire(entry *logrus.Entry) error {
	r.Entries = append(r.Entries, entry)
	return nil
}
