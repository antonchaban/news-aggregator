package main

import (
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/antonchaban/news-aggregator/pkg/storage/postgres"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	pgHost    string
	pgUser    string
	pgPass    string
	pgDBName  string
	pgSSLMode string
)

func init() {
	pgHost = os.Getenv("POSTGRES_HOST")
	pgUser = os.Getenv("POSTGRES_USER")
	pgPass = os.Getenv("POSTGRES_PASSWORD")
	pgDBName = os.Getenv("POSTGRES_DB")
	pgSSLMode = "disable"
}

func main() {
	db, err := storage.NewPostgresDB(storage.Config{
		Host:     pgHost,
		Username: pgUser,
		Password: pgPass,
		DBName:   pgDBName,
		SSLMode:  pgSSLMode,
	})
	if err != nil {
		logrus.Fatal("error occurred while connecting to the database: ", err.Error())
	}

	srcDb := postgres.NewSrc(db)
	artDb := postgres.New(db)
	sourceService := service.NewSourceService(artDb, srcDb)
	err = sourceService.FetchFromAllSources()
	if err != nil {
		logrus.Fatal("error occurred while fetching articles from sources: ", err.Error())
		return
	}
}
