package postgres

import (
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/jmoiron/sqlx"
)

type postgresArticleStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) service.ArticleStorage {
	return &postgresArticleStorage{db: db}
}

func (pa *postgresArticleStorage) GetAll() ([]model.Article, error) {
	var articles []model.Article
	query := fmt.Sprintf(`SELECT a.id, a.title, a.description, a.link, a.pub_date,
			       s.id AS source_id, s.name AS source_name, s.link AS source_link
			FROM articles a
			JOIN sources s ON a.source_id = s.id`)
	rows, err := pa.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		var source model.Source

		err := rows.Scan(&article.Id, &article.Title, &article.Description, &article.Link, &article.PubDate,
			&source.Id, &source.Name, &source.Link)
		if err != nil {
			return nil, err
		}
		article.Source = source
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
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
}

func (pa *postgresArticleStorage) DeleteBySourceID(id int) error {
	query := fmt.Sprintf(`DELETE FROM articles WHERE source_id = $1`)
	_, err := pa.db.Exec(query, id)
	if err != nil {
		return err
	}
	return err
}

func (pa *postgresArticleStorage) GetByFilter(query string, args []interface{}) ([]model.Article, error) {
	rows, err := pa.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var articles []model.Article

	for rows.Next() {
		var article model.Article
		var source model.Source

		err := rows.Scan(
			&article.Id, &article.Title, &article.Description, &article.Link, &article.PubDate,
			&source.Id, &source.Name, &source.Link,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		article.Source = source
		articles = append(articles, article)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return articles, nil
}
