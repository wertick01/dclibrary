package models

type Books struct {
	BookId    int    `json: "book_id"`
	Count     int    `json: "count"`
	Stars     int    `json: "stars"`
	BookPhoto string `json: "book_photo"`
	BookName  string `json: "book_name"`
	Authors   string `json: "author"`
}

type FinalBooks struct {
	Book    Books     `json: "books"`
	Authors []Authors `json: "authors"`
}
