package models

type FavorieteBooks struct {
	UserId         int `json: "user_id"`
	FavoriteBookId int `json: "favorites_books"`
}

type FavorieteAuthors struct {
	UserId           int `json: "user_id"`
	FavoriteAuthorId int `json: "favorites_books"`
}
