package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/backuper"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	eventServerStart               = "server_start"
	eventLoadSourcesStart          = "load_sources_start"
	eventLoadSourcesError          = "load_sources_error"
	eventLoadSourcesComplete       = "load_sources_complete"
	eventInitializeSourcesStart    = "initialize_sources_start"
	eventInitializeSourcesComplete = "initialize_sources_complete"
	eventAddSourcesStart           = "add_sources_start"
	eventAddSourceError            = "add_source_error"
	eventAddSourceSuccess          = "add_source_success"
	eventAddSourcesComplete        = "add_sources_complete"
	eventLoadArticlesStart         = "load_articles_start"
	eventLoadArticlesError         = "load_articles_error"
	eventLoadArticlesComplete      = "load_articles_complete"
	eventSaveArticlesStart         = "save_articles_start"
	eventSaveArticlesError         = "save_articles_error"
	eventSaveArticlesComplete      = "save_articles_complete"
	eventServerListenStart         = "server_listen_start"
	eventServerListenError         = "server_listen_error"
	eventServerStarted             = "server_started"
)

// Server represents a web server with TLS support.
type Server struct {
	httpServer *http.Server
	certFile   string
	keyFile    string
}

// Run starts the HTTPS server on the specified port and initializes the sources.
// It also loads articles from a backup file and saves them using the provided
// article handler. The server listens for HTTPS requests using the specified
// certificate and key files.
//
// Parameters:
// - port: The port on which the server will listen for requests.
// - handler: The handler for handling requests.
// - artHandler: The web handler for managing articles.
//
// Returns an error if the server fails to start or if loading/saving articles fails.
func (s *Server) Run(port string, handler http.Handler, artHandler web.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	logrus.WithField("event_id", eventServerStart).Info("Starting server on port ", port)

	// Load sources from file
	logrus.WithField("event_id", eventLoadSourcesStart).Info("Loading sources from file")
	srcs, err := backuper.NewLoader(artHandler.SrcService()).LoadSrcsFromFile()
	if err != nil {
		logrus.WithField("event_id", eventLoadSourcesError).Error("Failed to load sources from file", err)
		return err
	}
	logrus.WithField("event_id", eventLoadSourcesComplete).Info("Sources loaded from file")

	// Initialize sources if none are found
	if len(srcs) == 0 {
		logrus.WithField("event_id", eventInitializeSourcesStart).Info("No sources found in file, initializing default sources")
		initializeSources(artHandler.SrcService())
		logrus.WithField("event_id", eventInitializeSourcesComplete).Info("Default sources initialized")
	} else {
		logrus.WithField("event_id", eventAddSourcesStart).Info("Adding sources from file")
		for _, src := range srcs {
			_, err := artHandler.SrcService().AddSource(src)
			if err != nil {
				logrus.WithField("event_id", eventAddSourceError).Errorf("Error occurred while adding source %s: %s", src.Name, err.Error())
			} else {
				logrus.WithField("event_id", eventAddSourceSuccess).Infof("Source %s added successfully", src.Name)
			}
		}
		logrus.WithField("event_id", eventAddSourcesComplete).Info("All sources added from file")
	}

	// Load articles from file
	logrus.WithField("event_id", eventLoadArticlesStart).Info("Loading articles from file")
	articles, err := backuper.NewLoader(artHandler.SrcService()).LoadAllFromFile()
	if err != nil {
		logrus.WithField("event_id", eventLoadArticlesError).Error("Failed to load articles from file", err)
		return err
	}
	logrus.WithField("event_id", eventLoadArticlesComplete).Info("Articles loaded from file")

	// Save all articles
	logrus.WithField("event_id", eventSaveArticlesStart).Info("Saving all articles")
	err = artHandler.ArticleService().SaveAll(articles)
	if err != nil {
		logrus.WithField("event_id", eventSaveArticlesError).Error("Failed to save articles", err)
		return err
	}
	logrus.WithField("event_id", eventSaveArticlesComplete).Info("All articles saved")

	logrus.WithField("event_id", eventServerListenStart).Info("Starting HTTPS server")
	if err := s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.WithField("event_id", eventServerListenError).Fatalf("Could not listen on %s: %v\n", port, err)
		return err
	}

	logrus.WithField("event_id", eventServerStarted).Info("Server started successfully")
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
func initializeSources(ssvc web.SourceService) {
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
