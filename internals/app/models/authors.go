package models

type Authors struct {
	AuthorId    int         `json: "author_id"`
	AuthorName  string      `json: "author_name"`
	AuthorPhoto interface{} `json: "author_photo"`
}
