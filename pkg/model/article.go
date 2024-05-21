package model

import (
	"fmt"
)

type Article struct {
	Id          int
	Title       string
	Description string
	Link        string
	Source      string
	PubDate     string
}

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
