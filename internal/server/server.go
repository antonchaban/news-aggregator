package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/backuper"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	certFile   string
	keyFile    string
}

func (s *Server) Run(port string, handler http.Handler, artHandler web.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	initializeSources(artHandler.SrcService())
	articles, err := backuper.NewLoader(artHandler.SrcService()).LoadAllFromFile()
	if err != nil {
		return err
	}
	err = artHandler.ArticleService().SaveAll(articles)

	if err := s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("Could not listen on %s: %v\n", port, err)
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context, articles []model.Article) error {
	fmt.Println("Shutting down the server...")
	err := backuper.NewSaver(articles).SaveAllToFile()
	if err != nil {
		return err
	}
	return s.httpServer.Shutdown(ctx)
}

func NewServer(certFile, keyFile string) *Server {
	return &Server{
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func initializeSources(ssvc service.SourceService) {
	sources := []model.Source{
		{Name: "BBC News", Link: "https://feeds.bbci.co.uk/news/rss.xml"},
		{Name: "ABC News: International", Link: "https://abcnews.go.com/abcnews/internationalheadlines"},
		{Name: "The Washington Times stories: World", Link: "https://www.washingtontimes.com/rss/headlines/news/world/"},
		{Name: "USA TODAY", Link: "https://www.usatoday.com/news/world/"},
	}

	for _, source := range sources {
		_, err := ssvc.AddSource(source)
		if err != nil {
			logrus.Errorf("error occurred while adding source %s: %s", source.Name, err.Error())
		}
	}
}
