package scheduler

import (
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"time"
)

// Scheduler is a struct that holds the gocron.Scheduler instance and the services required for the scheduler to work
type Scheduler struct {
	scheduler *gocron.Scheduler
	asvc      web.ArticleService
	ssvc      web.SourceService
}

// NewScheduler initializes a new Scheduler instance with the provided article and source services.
func NewScheduler(asvc web.ArticleService, ssvc web.SourceService) *Scheduler {
	logrus.WithField("event_id", "scheduler_initialized").Info("Initializing Scheduler")
	return &Scheduler{
		scheduler: gocron.NewScheduler(time.UTC),
		asvc:      asvc,
		ssvc:      ssvc,
	}
}

// Start schedules the updateArticles task to run every minute and starts the scheduler asynchronously.
func (s *Scheduler) Start() {
	logrus.WithField("event_id", "scheduler_start").Info("Starting scheduler")
	_, err := s.scheduler.Every(1).Minute().Do(s.updateArticles)
	if err != nil {
		logrus.WithField("event_id", "schedule_error").Errorf("Error scheduling updateArticles job: %s", err.Error())
		return
	}
	s.scheduler.StartAsync()
	logrus.WithField("event_id", "scheduler_started").Info("Scheduler started successfully")
}

// Stop stops the scheduler and logs the stop event.
func (s *Scheduler) Stop() {
	logrus.WithField("event_id", "scheduler_stop").Info("Stopping scheduler")
	s.scheduler.Stop()
	logrus.WithField("event_id", "scheduler_stopped").Info("Scheduler stopped successfully")
}

// updateArticles fetches articles from all sources using the source service.
func (s *Scheduler) updateArticles() {
	logrus.WithField("event_id", "update_articles_start").Info("Updating articles...")
	err := s.ssvc.FetchFromAllSources()
	if err != nil {
		logrus.WithField("event_id", "update_articles_error").Errorf("Error occurred while fetching articles from sources: %s", err.Error())
		return
	}
	logrus.WithField("event_id", "update_articles_success").Info("Articles updated successfully")
}
