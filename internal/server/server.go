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

// Server represents a web server with TLS support.
type Server struct {
	httpServer *http.Server
	certFile   string
	keyFile    string
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	if err := s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("Could not listen on %s: %v\n", port, err)
		return err
	}

	return nil
}

// RunWithFiles starts the HTTP server on the specified port and initializes the sources.
// It also loads articles from a backup file and saves them using the provided
// article handler.
// The server listens for HTTPS requests using the specified
// certificate and key files.
//
// Parameters:
// - port: The port on which the server will listen for requests.
// - handler: The HTTP handler to use for handling requests.
// - artHandler: The web handler for managing articles.
//
// Returns an error if the server fails to start or if loading/saving articles fails.
func (s *Server) RunWithFiles(port string, handler http.Handler, artHandler web.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	srcs, err := backuper.NewLoader(artHandler.SrcService()).LoadSrcsFromFile()
	if err != nil {
		return err
	}
	if len(srcs) == 0 {
		initializeSources(artHandler.SrcService())
	} else {
		for _, src := range srcs {
			_, err := artHandler.SrcService().AddSource(src)
			if err != nil {
				logrus.Errorf("error occurred while adding source %s: %s", src.Name, err.Error())
			}
		}
	}
	articles, err := backuper.NewLoader(artHandler.SrcService()).LoadAllFromFile()
	if err != nil {
		return err
	}
	err = artHandler.ArticleService().SaveAll(articles)
	if err != nil {
		return err
	}

	if err := s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("Could not listen on %s: %v\n", port, err)
		return err
	}

	return nil
}

// Shutdown gracefully shuts down the server and saves articles to a backup file.
// It also logs the shutdown process.
//
// Parameters:
// - ctx: The context to use for the shutdown process.
// - articles: The list of articles to save to the backup file.
//
// Returns an error if saving the articles or shutting down the server fails.
func (s *Server) Shutdown(ctx context.Context, articles []model.Article, sources []model.Source) error {
	fmt.Println("Shutting down the server...")
	err := backuper.NewSaver(articles, sources).SaveAllToFile()
	if err != nil {
		return err
	}
	err = backuper.NewSaver(articles, sources).SaveSrcsToFile()
	if err != nil {
		return err
	}
	return s.httpServer.Shutdown(ctx)
}

// NewServer creates a new Server instance with the specified certificate and key files.
//
// Parameters:
// - certFile: The path to the certificate file for TLS.
// - keyFile: The path to the key file for TLS.
//
// Returns a new Server instance.
func NewServer(certFile, keyFile string) *Server {
	return &Server{
		certFile: certFile,
		keyFile:  keyFile,
	}
}

// initializeSources initializes the sources for the provided SourceService.
// It adds a predefined list of sources to the service.
//
// Parameters:
// - ssvc: The SourceService to use for adding sources.
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
