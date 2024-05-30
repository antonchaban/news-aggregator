package main

import (
	"news-aggregator/pkg/handler/cli"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/service"
	"news-aggregator/pkg/storage/inmemory"
)

func main() {
	// Initialize storage and service
	var articles []model.Article
	db := inmemory.New(articles)
	svc := service.New(db)

	// Initialize handler and execute CLI commands
	handler := cli.NewHandler(svc)
	handler.InitCommands()
}
