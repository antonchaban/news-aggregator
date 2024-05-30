package main

import (
	"news-aggregator/pkg/handler/cli"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/repository"
	"news-aggregator/pkg/service"
)

func main() {
	// Initialize repository and service
	var articles []model.Article
	db := repository.NewInMemory(articles)
	svc := service.New(db)

	// Initialize handler and execute CLI commands
	handler := cli.NewHandler(svc)
	handler.Execute()
}
