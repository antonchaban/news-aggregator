package main

import (
	"news-aggregator/pkg/handler/cli"
	"news-aggregator/pkg/service"
	"news-aggregator/pkg/storage/inmemory"
)

func main() {
	// Initialize storage and service
	db := inmemory.New()
	svc := service.New(db)

	// Initialize handler and execute CLI commands
	cli.NewHandler(svc)
}
