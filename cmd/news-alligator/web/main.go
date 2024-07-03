package main

import (
	"context"
	_ "github.com/antonchaban/news-aggregator/cmd/news-alligator/web/docs"
	"github.com/antonchaban/news-aggregator/internal/server"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/scheduler"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/antonchaban/news-aggregator/pkg/storage/postgres"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// @title News Alligator API
// @version 1
// @description This is a News Alligator API server.
// @host https://localhost:8080
// @BasePath /articles

func main() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5436")
	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "qwerty")
	os.Setenv("DB_NAME", "postgres")
	os.Setenv("DB_SSLMODE", "disable")
	// Initialize in-memory databases
	db, err := storage.NewPostgresDB(storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})

	//artDb := inmemory.New()
	//srcDb := inmemory.NewSrc()
	artDb := postgres.New(db)
	srcDb := postgres.NewSrc(db)
	asvc := service.New(artDb)
	ssvc := service.NewSourceService(artDb, srcDb)

	// Initialize web handler
	h := web.NewHandler(asvc, ssvc)

	// Create a new HTTPS server
	srv := server.NewServer(os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"))

	// Start the server in a goroutine
	go func() {
		if err := srv.Run(os.Getenv("PORT"), h.InitRoutes(), *h); err != nil {
			logrus.Fatal("error occurred while running http server: ", err.Error())
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
