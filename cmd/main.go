package main

import (
	"fmt"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/repository"
	"news-aggregator/pkg/service"
)

func main() {
	var parser service.Parser
	parser = &service.RssParser{}
	var articles []model.Article
	db := repository.NewArticleInMemory(articles)
	svc := service.NewArticleService(db)
	err := startParsing(parser, db, svc, "data/abcnews-international-category-19-05-24.xml")
	err = startParsing(parser, db, svc, "data/bbc-world-category-19-05-24.xml")
	if err != nil {
		return
	}
	fmt.Println("Parsing and saving articles completed!")
	fmt.Println("Articles in the database:")
	articlesInDb, err := svc.GetAll()
	for _, article := range articlesInDb {
		fmt.Println(article)
	}
}

func startParsing(parser service.Parser, db *repository.ArticleInMemory, svc *service.ArticleService, filePath string) error {
	file, err := parser.ParseFile(filePath)
	if err != nil {
		panic(err)
	}
	for _, item := range file.Items {
		article := model.Article{
			Id:          len(db.Articles) + 1,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
		}
		_, err := svc.Create(article)
		if err != nil {
			panic(err)
		}

	}
	return err
}
