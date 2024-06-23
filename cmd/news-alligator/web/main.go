package main

import (
	"context"
	"github.com/antonchaban/news-aggregator/internal/server"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/antonchaban/news-aggregator/pkg/storage/inmemory"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

const dotenvPath = "../../../.env"

func main() {
	// Initialize in-memory databases
	db := inmemory.New()
	srcDb := inmemory.NewSrc()
	asvc := service.New(db)
	ssvc := service.NewSourceService(db, srcDb)

	// Load environment variables
	if err := godotenv.Load(dotenvPath); err != nil {
		logrus.Fatal("error occurred while loading env variables: ", err.Error())
	}

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

	// Wait for a signal to quit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("news-alligator 🐊 shutting down")

	// Retrieve all articles before shutting down
	articles, err := asvc.GetAll()
	if err != nil {
		logrus.Errorf("error occurred on getting all articles: %s", err.Error())
	}

	// Shutdown the server
	if err := srv.Shutdown(context.Background(), articles); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}
}
