package models

type Favorietes struct {
	UserId        int   `json: "user_id"`
	FavoriteBooks []int `json: "favorites_books"`
}
