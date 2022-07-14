package models

type Books struct {
	BookId    int    `json: "book_id"`
	BookName  string `json: "book_name"`
	AuthorId  int    `json: "author"`
	Count     int    `json: "count"`
	BookPhoto string `json: "book_photo"`
}
