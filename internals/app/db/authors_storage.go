package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"dclibrary.com/internals/app/models"
)

type AuthorsStorage struct {
	DB *sql.DB
}

func (m *AuthorsStorage) CreateNewAuthor(author_name, author_photo string) (int, error) {

	stmt := `INSERT INTO dclibrary.authors (author_name, author_photo) VALUES(?, ?)`

	result, err := m.DB.Exec(stmt, author_name, author_photo)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Printf("---> Author %v has been added to DB", id)

	return int(id), nil
}

func (m *AuthorsStorage) GetAuthorsList() ([]*models.Authors, error) {

	stmt := `SELECT author_id, author_name, author_photo FROM dclibrary.authors`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var authors []*models.Authors

	for rows.Next() {
		s := &models.Authors{}
		err = rows.Scan(&s.AuthorId, &s.AuthorName, &s.AuthorPhoto)
		if err != nil {
			return nil, err
		}
		authors = append(authors, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (m *AuthorsStorage) GetAuthorById(id int64) (*models.Authors, error) {

	stmt := `SELECT author_id, author_name, author_photo FROM dclibrary.authors WHERE author_id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Authors{}

	err := row.Scan(&s.AuthorId, &s.AuthorName, &s.AuthorPhoto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *AuthorsStorage) GetBooksByAuthorId(author_id int64) ([]*models.Books, error) { //???

	stmt := `SELECT book_id, bookname, authors, count, photo FROM dclibrary.books WHERE author_id = ?`
	//lst := *BooksStorage.GetBooksList()
	sdmd := `SELECT author_name, author_photo FROM dclibrary.authors WHERE author_id = ?`

	rows, err := m.DB.Query(stmt, author_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []*models.Books
	var authors []string

	for rows.Next() {
		s := &models.Books{}
		err = rows.Scan(&s.BookId, &s.BookName, &s.Authors, &s.Count, &s.BookPhoto)
		if err != nil {
			return nil, err
		}

		authors = strings.Split(s.Authors, ", ")
		auth := m.DB.QueryRow(sdmd, s.Author.AuthorId)
		err = auth.Scan(&s.Author.AuthorName, &s.Author.AuthorPhoto)
		books = append(books, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
