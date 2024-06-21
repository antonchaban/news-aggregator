package server

import (
	"context"
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/backuper"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler, artHandler web.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	articles, err := backuper.NewLoader(artHandler.ArticleService()).LoadAllFromFile()
	for _, feed := range getSupportedFeeds() {
		articlesFeed, err := backuper.NewLoader(artHandler.ArticleService()).UpdateFromFeed(feed)
		if err != nil {
			logrus.Fatalf("error occurred while updating articles from feed: %s", err.Error())
		}
		articles = append(articles, articlesFeed...)

	}
	if err != nil {
		return nil
	}
	err = artHandler.ArticleService().SaveAll(articles)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context, articles []model.Article) error {
	fmt.Println("Shutting down the server...")
	err := backuper.NewSaver(articles).SaveAllToFile()
	if err != nil {
		return err
	}
	return s.httpServer.Shutdown(ctx)
}

func getSupportedFeeds() []string {
	return []string{
		"https://feeds.bbci.co.uk/news/rss.xml",
		"https://abcnews.go.com/abcnews/internationalheadlines",
		"https://www.washingtontimes.com/rss/headlines/news/world/",
		"https://www.usatoday.com/news/world/",
	}
}
