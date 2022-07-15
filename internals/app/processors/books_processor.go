package processors

import (
	"errors"

	"dclibrary.com/internals/app/db"
	"dclibrary.com/internals/app/models"
)

type BooksProcessor struct {
	storage *db.BooksStorage
}

func NewBooksProcessor(storage *db.BooksStorage) *BooksProcessor {
	processor := new(BooksProcessor)
	processor.storage = storage
	return processor
}

func (processor *BooksProcessor) CreateBook(book models.Books) (int, error) {

	if book.BookName == "" {
		return 0, errors.New("name should not be empty")
	}

	return processor.storage.CreateNewBook(book.BookName, book.BookPhoto, book.AuthorId, book.Count)
}

func (processor *BooksProcessor) FindBook(id int) (*models.Books, error) {
	book, err := processor.storage.GetBookById(id)

	if err != nil {
		return book, errors.New("user not found")
	}

	return book, nil

}

func (processor *BooksProcessor) ListBooks() ([]*models.Books, error) {
	return processor.storage.GetBooksList()
}
