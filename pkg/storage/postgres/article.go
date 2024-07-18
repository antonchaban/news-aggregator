package postgres

import (
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/jmoiron/sqlx"
)

type postgresArticleStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) storage.ArticleStorage {
	return &postgresArticleStorage{db: db}
}

func (pa *postgresArticleStorage) GetAll() ([]model.Article, error) {
	var articles []model.Article
	query := fmt.Sprintf(`SELECT a.id, a.title, a.description, a.link, a.pub_date,
			       s.id, s.name, s.link
			FROM articles a
			JOIN sources s ON a.source_id = s.id`)
	err := pa.db.Select(&articles, query)
	if err != nil {
		return nil, err
	}
	return articles, nil
	/*	rows, err := pa.db.Query(context.Background(), `
			SELECT a.id, a.title, a.description, a.link, a.pub_date,
			       s.id, s.name, s.link
			FROM articles a
			JOIN sources s ON a.source_id = s.id
		`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var articles []model.Article
		for rows.Next() {
			var article model.Article
			err := rows.Scan(
				&article.Id, &article.Title, &article.Description, &article.Link, &article.PubDate,
				&article.Source.Id, &article.Source.Name, &article.Source.Link,
			)
			if err != nil {
				return nil, err
			}
			articles = append(articles, article)
		}
		return articles, nil*/
}

func (pa *postgresArticleStorage) Save(article model.Article) (model.Article, error) {
	var id int
	createQuery := fmt.Sprintf(`INSERT INTO articles (title, description, link, source_id, pub_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`)

	err := pa.db.QueryRow(createQuery, article.Title, article.Description, article.Link, article.Source.Id, article.PubDate).Scan(&id)
	if err != nil {
		return model.Article{}, err
	}
	article.Id = id
	return article, nil
	/*var id int
	err := pa.db.QueryRow(context.Background(), "INSERT INTO articles (title, description, link, source_id, pub_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		article.Title, article.Description, article.Link, article.Source.Id, article.PubDate).Scan(&id)
	if err != nil {
		return model.Article{}, err
	}
	article.Id = id
	return article, nil*/
}

func (pa *postgresArticleStorage) SaveAll(articles []model.Article) error {
	for _, article := range articles {
		_, err := pa.Save(article)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pa *postgresArticleStorage) Delete(id int) error {
	query := fmt.Sprintf(`DELETE FROM articles WHERE id = $1`)
	_, err := pa.db.Exec(query, id)
	return err
	/*	_, err := pa.db.Exec(context.Background(), "DELETE FROM articles WHERE id = $1", id)
		return err*/
}

func (pa *postgresArticleStorage) DeleteBySourceID(id int) error {
	query := fmt.Sprintf(`DELETE FROM articles WHERE source_id = $1`)
	_, err := pa.db.Exec(query, id)
	if err != nil {
		return err
	}
	return err
	/*_, err := pa.db.Exec(context.Background(), "DELETE FROM articles WHERE source_id = $1", id)
	return err*/
}

func (pa *postgresArticleStorage) GetByFilter(query string, args []interface{}) ([]model.Article, error) {
	var articles []model.Article
	err := pa.db.Select(&articles, query, args...)
	if err != nil {
		return nil, err
	}
	return articles, nil

	/*rows, err := pa.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		err := rows.Scan(
			&article.Id, &article.Title, &article.Description, &article.Link, &article.PubDate,
			&article.Source.Id, &article.Source.Name, &article.Source.Link,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil*/
}
