package postgres

import (
	"context"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresArticleStorage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) storage.ArticleStorage {
	return &postgresArticleStorage{db: db}
}

func (pa *postgresArticleStorage) GetAll() ([]model.Article, error) {
	rows, err := pa.db.Query(context.Background(), `
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
	return articles, nil
}

func (pa *postgresArticleStorage) Save(article model.Article) (model.Article, error) {
	var id int
	err := pa.db.QueryRow(context.Background(), "INSERT INTO articles (title, description, link, source_id, pub_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		article.Title, article.Description, article.Link, article.Source.Id, article.PubDate).Scan(&id)
	if err != nil {
		return model.Article{}, err
	}
	article.Id = id
	return article, nil
}

func (pa *postgresArticleStorage) SaveAll(articles []model.Article) error {
	tx, err := pa.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, article := range articles {
		_, err := tx.Exec(context.Background(), "INSERT INTO articles (title, description, link, source_id, pub_date) VALUES ($1, $2, $3, $4, $5)",
			article.Title, article.Description, article.Link, article.Source.Id, article.PubDate)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (pa *postgresArticleStorage) Delete(id int) error {
	_, err := pa.db.Exec(context.Background(), "DELETE FROM articles WHERE id = $1", id)
	return err
}

func (pa *postgresArticleStorage) DeleteBySourceID(id int) error {
	_, err := pa.db.Exec(context.Background(), "DELETE FROM articles WHERE source_id = $1", id)
	return err
}
