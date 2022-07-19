package processors

import (
	"errors"

	"dclibrary.com/internals/app/db"
	"dclibrary.com/internals/app/models"
)

type AuthorsProcessor struct {
	storage *db.AuthorsStorage
}

func (processor *AuthorsProcessor) CreateAuthor(author models.Authors) (int, error) {

	if author.AuthorName.Name == "" && author.AuthorName.Surname == "" {
		return 0, errors.New("name should not be empty")
	}

	return processor.storage.CreateNewAuthor(
		author.AuthorName.Name,
		author.AuthorName.Surname,
		author.AuthorName.Patronymic,
		author.AuthorPhoto,
	)
}

func (processor *AuthorsProcessor) ListAuthors() ([]*models.Authors, error) {
	return processor.storage.GetAuthorsList()
}

func (processor *AuthorsProcessor) AuthorsBooks(id int64) ([]*models.FinalBooks, error) {
	book, err := processor.storage.GetBooksByAuthorId(id)

	if err != nil {
		return book, errors.New("Author not found")
	}

	return book, nil

}

func (processor *AuthorsProcessor) FindAuthor(id int64) (*models.Authors, error) {
	author, err := processor.storage.GetAuthorById(id)

	if err != nil {
		return author, errors.New("user not found")
	}

	return author, nil

}
