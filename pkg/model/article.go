package model

import (
	"fmt"
	"time"
)

// Article is a struct that represents a news article
// It has the following fields:
// - Id: an integer that represents the unique identifier of the article
// - Title: a string that represents the title of the article
// - Description: a string that represents the description of the article
// - Link: a string that represents the link to the original article
// - Source: a string that represents the source of the article
// - PubDate: a time.Time that represents the publication date of the article
type Article struct {
	Id          int       `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Link        string    `db:"link"`
	Source      Source    `db:"source_id"`
	PubDate     time.Time `db:"pub_date"`
}

// String method returns a string representation of the Article struct
func (a Article) String() string {
	return fmt.Sprintf(
		"ID: %d\nTitle: %s\nDate: %s\nDescription: %s\nLink: %s\nSource: %s\n",
		a.Id,
		a.Title,
		a.PubDate,
		a.Description,
		a.Link,
		a.Source,
	)
}
