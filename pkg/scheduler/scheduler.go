package scheduler

import (
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"time"
)

type Scheduler struct {
	scheduler *gocron.Scheduler
	asvc      service.ArticleService
	ssvc      service.SourceService
}

func NewScheduler(asvc service.ArticleService, ssvc service.SourceService) *Scheduler {
	logrus.WithField("event_id", "scheduler_initialized").Info("Initializing Scheduler")
	return &Scheduler{
		scheduler: gocron.NewScheduler(time.UTC),
		asvc:      asvc,
		ssvc:      ssvc,
	}
}

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

func (s *Scheduler) Stop() {
	logrus.WithField("event_id", "scheduler_stop").Info("Stopping scheduler")
	s.scheduler.Stop()
	logrus.WithField("event_id", "scheduler_stopped").Info("Scheduler stopped successfully")
}

func (s *Scheduler) updateArticles() {
	logrus.WithField("event_id", "update_articles_start").Info("Updating articles...")
	err := s.ssvc.FetchFromAllSources()
	if err != nil {
		logrus.WithField("event_id", "update_articles_error").Errorf("Error occurred while fetching articles from sources: %s", err.Error())
		return
	}
	logrus.WithField("event_id", "update_articles_success").Info("Articles updated successfully")
}
