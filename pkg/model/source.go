package model

import "fmt"

type Source struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Link      string `json:"link" db:"link"`
	ShortName string `json:"short_name" db:"short_name"`
}

func (s Source) String() string {
	return fmt.Sprintf("Source{Id: %d,"+
		" Name: %s,"+
		" Link: %s}", s.Id, s.Name, s.Link)
}
