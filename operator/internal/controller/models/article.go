package models

import "time"

type Article struct {
	Id          int       `json:"Id"`          // Unique identifier for the article.
	Title       string    `json:"Title"`       // Full title of the article.
	Description string    `json:"Description"` // Full description of the article.
	Link        string    `json:"Link"`        // URL link to the article.
	Source      Source    `json:"Source"`      // Source of the article.
	PubDate     time.Time `json:"PubDate"`     // Date of publication of the article.
}
