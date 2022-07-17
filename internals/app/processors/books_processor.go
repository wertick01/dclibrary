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

func (processor *BooksProcessor) CreateBook(book models.Books) (int, error) {

	if book.BookName == "" {
		return 0, errors.New("name should not be empty")
	}

	return processor.storage.CreateNewBook(book.BookName, book.BookPhoto, book.Authors, book.Count)
}

func (processor *BooksProcessor) FindBook(id int64) (*models.FinalBooks, error) {
	book, err := processor.storage.GetBookById(id)

	if err != nil {
		return book, errors.New("user not found")
	}

	return book, nil

}

func (processor *BooksProcessor) ListBooks() ([]*models.FinalBooks, error) {
	return processor.storage.GetBooksList()
}

func (processor *BooksProcessor) UpdateBook(id int64) (*models.Books, error) { //!!! ПРОВЕРИТЬ
	book, _ := processor.FindBook(id)

	changedbook, err := processor.storage.ChangeBookById(id, &book.Book)
	if err != nil {
		return &book.Book, errors.New("SOMETHING IS WRONG")
	}

	return changedbook, nil
}

func (processor *BooksProcessor) DeleteBook(id int64) (*models.FinalBooks, error) {
	book, _ := processor.FindBook(id)
	_, err := processor.storage.DeleteBookById(id)
	if err != nil {
		return book, errors.New("CANNOT DELETE BOOK")
	}
	fmt.Printf("Book %v has been deleted.", id)
	return book, nil
}
