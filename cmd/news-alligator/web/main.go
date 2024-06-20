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
	db := inmemory.New()
	svc := service.New(db)

	if err := godotenv.Load(dotenvPath); err != nil {
		logrus.Fatal("error occurred while loading env variables: ", err.Error())
	}

	h := web.NewHandler(svc)
	srv := new(server.Server)
	go func() {
		if err := srv.Run(os.Getenv("PORT"), h.InitRoutes()); err != nil {
			logrus.Fatal("error occurred while running http server: ", err.Error())
		}
	}()

	logrus.Print("news-alligator 🐊 started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("news-alligator 🐊 shutting down")

	articles, err := svc.GetAll()
	if err != nil {
		logrus.Errorf("error occurred on getting all articles: %s", err.Error())
	}
	if err := srv.Shutdown(context.Background(), articles); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}
}
