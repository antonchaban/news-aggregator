package models

type Source struct {
	Id        int    `json:"id"`         // Unique identifier for the source.
	Name      string `json:"name"`       // Full name of the source, as it appears in the news aggregator.
	Link      string `json:"link"`       // URL link to the source.
	ShortName string `json:"short_name"` // Shortened name of the source for search purposes.
}
