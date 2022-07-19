package models

type Authors struct {
	AuthorId    int        `json: "author_id"`
	AuthorStars int        `json: "author_stars"`
	AuthorName  AuthorName `json: "author_name"`
	AuthorPhoto string     `json: "author_photo"`
}

type AuthorName struct {
	Name       string `json: "name"`
	Surname    string `json: "surname"`
	Patronymic string `json: "patronymic"`
}
