// scheduler/scheduler.go
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
	s := &Scheduler{
		scheduler: gocron.NewScheduler(time.UTC),
		asvc:      asvc,
		ssvc:      ssvc,
	}
	return s
}

func (s *Scheduler) Start() {
	// Schedule the updateArticles method to run every minute
	_, err := s.scheduler.Every(1).Minute().Do(s.updateArticles)
	if err != nil {
		return
	}
	s.scheduler.StartAsync()
}

func (s *Scheduler) Stop() {
	s.scheduler.Stop()
}

func (s *Scheduler) updateArticles() {
	logrus.Print("Updating articles...")
	for _, feed := range getSupportedFeeds() {
		articlesFeed, err := s.asvc.LoadFromFeed(feed)
		if err != nil {
			logrus.Errorf("error occurred while updating articles from feed: %s", err.Error())
			continue
		}
		err = s.asvc.SaveAll(articlesFeed)
		if err != nil {
			logrus.Errorf("error occurred while saving updated articles: %s", err.Error())
		}
	}
	err := s.ssvc.FetchFromAllSources()
	if err != nil {
		logrus.Errorf("error occurred while fetching articles from sources: %s", err.Error())
	}
	logrus.Print("Articles updated successfully")
}

func getSupportedFeeds() []string {
	return []string{
		"https://feeds.bbci.co.uk/news/rss.xml",
		"https://abcnews.go.com/abcnews/internationalheadlines",
		"https://www.washingtontimes.com/rss/headlines/news/world/",
		"https://www.usatoday.com/news/world/",
	}
}
