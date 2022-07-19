package processors

import (
	"errors"
	"fmt"

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

func (processor *BooksProcessor) CreateBook(book *models.Books) (*models.Books, error) {

	if book.BookName == "" {
		return nil, errors.New("name should not be empty")
	}

	return processor.storage.CreateNewBook(book)
}

func (processor *BooksProcessor) FindBook(id int64) (*models.Books, error) {
	book, err := processor.storage.GetBookById(id)

	if err != nil {
		return nil, errors.New("book not found")
	}

	return book, nil

}

func (processor *BooksProcessor) ListBooks() ([]*models.Books, error) {
	return processor.storage.GetBooksList()
}

func (processor *BooksProcessor) UpdateBook(id int64) (*models.Books, error) { //!!! ПРОВЕРИТЬ
	book, err := processor.FindBook(id)
	if err != nil {
		return nil, err
	}

	changedbook, err := processor.storage.ChangeBook(book)
	if err != nil {
		return nil, errors.New("SOMETHING IS WRONG")
	}

	return changedbook, nil
}

func (processor *BooksProcessor) DeleteBook(id int64) (int64, error) {
	deleted, err := processor.storage.DeleteBookById(id)
	if err != nil {
		return 0, errors.New("CANNOT DELETE BOOK")
	}
	fmt.Printf("Book %v has been deleted.", id)
	return deleted, nil
}

func (processor *BooksProcessor) StarTheBook(id int64) (*models.Books, error) {
	err := processor.storage.PutStarByBookId(id)
	if err != nil {
		return nil, errors.New("CANNOT DELETE BOOK")
	}
	fmt.Printf("Book %v has been stared.", id)
	book, _ := processor.FindBook(id)
	return book, nil
}
