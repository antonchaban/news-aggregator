package model

import (
	"fmt"
	"time"
)

type Article struct {
	Id          int
	Title       string
	Description string
	Link        string
	Source      string
	PubDate     time.Time
}

func (a Article) String() string {
	return fmt.Sprintf(
		"ID%d\n Title: %s\nDate: %s\nDescription: %s\nLink: %s\nSource: %s\n",
		a.Id,
		a.Title,
		a.PubDate.Format(time.RFC1123),
		a.Description,
		a.Link,
		a.Source,
	)
}
