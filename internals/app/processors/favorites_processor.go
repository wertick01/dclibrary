package processors

import (
	"errors"
	"fmt"

	"dclibrary.com/internals/app/db"
	"dclibrary.com/internals/app/models"
)

type FavorietesProcessor struct {
	storage *db.FavorietesStorage
}

func NewFavorietesProcessor(storage *db.FavorietesStorage) *FavorietesProcessor {
	processor := new(FavorietesProcessor)
	processor.storage = storage
	return processor
}

func (processor *FavorietesProcessor) AddFavorieteBook(favoriete *models.FavorieteBooks) (int, error) {
	return processor.storage.AddToFavorieteBooks(favoriete)
}

func (processor *FavorietesProcessor) ListFavorieteBooks(user_id int64) ([]*models.Books, error) {
	return processor.storage.GetFavoriteBooksList(user_id)
}

func (processor *FavorietesProcessor) DeleteFromFavorieteBooks(id int64) (int64, error) {
	deleted, err := processor.storage.DeleteFavorieteBookById(id)
	if err != nil {
		return 0, errors.New("CANNOT DELETE BOOK")
	}
	fmt.Printf("Book %v has been deleted from favorietes.", deleted)
	return deleted, nil
}

func (processor *FavorietesProcessor) AddFavorieteAuthor(favoriete *models.FavorieteAuthors) (int, error) {
	return processor.storage.AddToFavorieteAuthors(favoriete)
}

func (processor *FavorietesProcessor) ListFavorieteAuthors(user_id int64) ([]*models.Authors, error) {
	return processor.storage.GetFavoriteAuthorsList(user_id)
}

func (processor *FavorietesProcessor) DeleteFromFavorieteAuthors(id int64) (int64, error) {
	deleted, err := processor.storage.DeleteFavorieteAuthorById(id)
	if err != nil {
		return 0, errors.New("CANNOT DELETE THE AUTHOR")
	}
	fmt.Printf("Author %v has been deleted from favorietes.", deleted)
	return deleted, nil
}
