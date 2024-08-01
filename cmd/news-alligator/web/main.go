package main

import (
	"context"
	"fmt"
	_ "github.com/antonchaban/news-aggregator/cmd/news-alligator/web/docs"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/scheduler"
	"github.com/antonchaban/news-aggregator/pkg/server"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/antonchaban/news-aggregator/pkg/storage/postgres"
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
	certFileEnvVar  = "CERT_FILE"
	keyFileEnvVar   = "KEY_FILE"
	portEnvVar      = "PORT"
	dbHostEnvVar    = "DB_HOST"
	dbPortEnvVar    = "DB_PORT"
	dbUserEnvVar    = "DB_USERNAME"
	dbPassEnvVar    = "DB_PASSWORD"
	dbNameEnvVar    = "DB_NAME"
	dbSSLModeEnvVar = "DB_SSLMODE"
	storTypeEnvVar  = "STORAGE_TYPE"
)

func main() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5436")
	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "qwerty")
	os.Setenv("DB_NAME", "postgres")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("STORAGE_TYPE", "postgres")

	if err := checkEnvVars(
		certFileEnvVar, keyFileEnvVar, portEnvVar, dbHostEnvVar, dbPortEnvVar, dbUserEnvVar,
		dbPassEnvVar, dbNameEnvVar, dbSSLModeEnvVar, storTypeEnvVar,
	); err != nil {
		logrus.Fatal(err)
	}

	// Initialize in-memory databases
	db, err := storage.NewPostgresDB(storage.Config{
		Host:     os.Getenv(dbHostEnvVar),
		Port:     os.Getenv(dbPortEnvVar),
		Username: os.Getenv(dbUserEnvVar),
		Password: os.Getenv(dbPassEnvVar),
		DBName:   os.Getenv(dbNameEnvVar),
		SSLMode:  os.Getenv(dbSSLModeEnvVar),
	})
	if err != nil {
		logrus.Fatalf("error occurred while initializing database: %s", err.Error())
		panic(err)
	}

	//artDb := inmemory.New()
	//srcDb := inmemory.NewSrc()
	artDb := postgres.New(db)
	srcDb := postgres.NewSrc(db)
	asvc := service.New(artDb)
	ssvc := service.NewSourceService(artDb, srcDb)

	// Initialize web handler
	h := web.NewHandler(asvc, ssvc)

	// Create a new HTTPS server
	srv := server.NewServer(os.Getenv(certFileEnvVar), os.Getenv(keyFileEnvVar))

	// Start the server in a goroutine
	go func() {
		if err := srv.Run(os.Getenv(portEnvVar), h.InitRoutes()); err != nil {
			logrus.Fatal("error occurred while running server: ", err.Error())
		}
	}()

	logrus.Print("news-alligator 🐊 started")

	// Start the scheduler for updating articles
	newScheduler := scheduler.NewScheduler(asvc, ssvc)
	newScheduler.Start()

	// Wait for a signal to quit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("news-alligator 🐊 shutting down")

	// Stop the scheduler
	newScheduler.Stop()

	// Retrieve all articles before shutting down
	articles, err := asvc.GetAll()
	if err != nil {
		logrus.Errorf("error occurred on getting all articles: %s", err.Error())
	}

	sources, err := ssvc.GetAll()
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
