package main

import (
	"context"
	"fmt"
	_ "github.com/antonchaban/news-aggregator/cmd/news-alligator/web/docs"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/server"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/antonchaban/news-aggregator/pkg/storage/inmemory"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	_ "go.uber.org/mock/mockgen/model"
	"os"
	"os/signal"
	"syscall"
)

// @title News Alligator API
// @version 1
// @description This is a News Alligator API server.
// @host https://localhost:443
// @BasePath /articles

const (
	certFileEnvVar = "CERT_FILE"
	keyFileEnvVar  = "KEY_FILE"
	portEnvVar     = "PORT"
)

var (
	dbHost  string
	dbUser  string
	dbPass  string
	dbName  string
	sslMode string
)

func init() {
	dbHost = os.Getenv("POSTGRES_HOST")
	dbUser = os.Getenv("POSTGRES_USER")
	dbPass = os.Getenv("POSTGRES_PASSWORD")
	dbName = os.Getenv("POSTGRES_DB")
	sslMode = "disable"
}

func main() {
	//os.Setenv(certFileEnvVar, "pkg/server/server.crt")
	//os.Setenv(keyFileEnvVar, "pkg/server/server.key")
	//os.Setenv(portEnvVar, "8443")
	//os.Setenv("SAVES_DIR", "backups")
	if err := checkEnvVars(
		certFileEnvVar, keyFileEnvVar, portEnvVar,
	); err != nil {
		logrus.Fatal(err)
	}

	// Initialize in-memory databases

	//db, err := storage.NewDB(storage.Config{
	//	Host:     dbHost,
	//	Username: dbUser,
	//	Password: dbPass,
	//	DBName:   dbName,
	//	SSLMode:  sslMode,
	//})
	//if err != nil {
	//	logrus.Fatal("error occurred while connecting to the database: ", err.Error())
	//}

	//srcDb := postgres.NewSrc(db)
	//artDb := postgres.New(db)

	artDb := inmemory.New()
	srcDb := inmemory.NewSrc()
	articleService := service.New(artDb)
	sourceService := service.NewSourceService(artDb, srcDb)

	// Initialize web handler
	h := web.NewHandler(articleService, sourceService)

	// Create a new HTTPS server
	srv := server.NewServer(os.Getenv(certFileEnvVar), os.Getenv(keyFileEnvVar))

	// Start the server in a goroutine
	go func() {
		if err := srv.RunWithFiles(os.Getenv(portEnvVar), h.InitRoutes(), *h); err != nil {
			logrus.Fatal("error occurred while running http server: ", err.Error())
		}
	}()

	logrus.Print("news-alligator 🐊 started")

	//newScheduler := scheduler.NewScheduler(articleService, sourceService)
	//newScheduler.Start()

	// Wait for a signal to quit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("news-alligator 🐊 shutting down")

	//newScheduler.Stop()

	// Retrieve all articles before shutting down
	articles, err := articleService.GetAll()
	if err != nil {
		logrus.Errorf("error occurred on getting all articles: %s", err.Error())
	}

	sources, err := sourceService.GetAll()
	if err != nil {
		logrus.Errorf("error occurred on getting all sources: %s", err.Error())
	}

	// Shutdown the server
	if err := srv.Shutdown(context.Background(), articles, sources); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}
}

// checkEnvVars checks if the required environment variables are set and returns an error if any are missing
func checkEnvVars(vars ...string) error {
	for _, v := range vars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("environment variable %s not set", v)
		}
	}
	return nil
}
