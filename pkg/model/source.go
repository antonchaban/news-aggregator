package model

import "fmt"

type Source struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

func (s Source) String() string {
	return fmt.Sprintf(
		"ID: %s\nName: %s\nLink: %s\n", s.Id, s.Name, s.Link)
}
