package main

import (
	"github.com/antonchaban/news-aggregator/pkg/handler/cli"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/antonchaban/news-aggregator/pkg/storage/inmemory"
)

// The main function in CLI package initializes the storage and service and creates a new CLI handler
// for News Aggregator application.
func main() {
	// Initialize storage and service
	db := inmemory.New()
	svc := service.New(db)
	srcDb := inmemory.NewSrc()
	srcSvc := service.NewSourceService(db, srcDb)

	// Initialize handler and execute CLI commands
	_, err := cli.NewHandler(svc, srcSvc)
	if err != nil {
		panic(err)
	}
}
