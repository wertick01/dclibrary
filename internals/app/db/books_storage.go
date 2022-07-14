package db

import (
	"database/sql"
	"errors"
	"fmt"

	"dclibrary.com/internals/app/models"
)

type BooksStorage struct {
	DB *sql.DB
}

func (m *BooksStorage) CreateNewBook(bookname, photo string, author_id, count int) (int, error) {

	stmt := `INSERT INTO dclibrary.books (bookname, author_id, count, photo) VALUES(?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, bookname, author_id, count, photo)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Printf("---> Book %v has been added to DB", id)

	return int(id), nil
}

func (m *BooksStorage) GetBooksList() ([]*models.Books, error) {

	stmt := `SELECT book_id, bookname, author_id, count, photo FROM dclibrary.books`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []*models.Books

	for rows.Next() {
		s := &models.Books{}
		err = rows.Scan(&s.BookId, &s.BookName, &s.AuthorId, &s.Count, &s.BookPhoto)
		if err != nil {
			return nil, err
		}
		books = append(books, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (m *BooksStorage) GetBookById(id int) (*models.Books, error) {

	stmt := `SELECT book_id, bookname, author_id, count, photo FROM dclibrary.books WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Books{}

	err := row.Scan(&s.BookId, &s.BookName, &s.AuthorId, &s.Count, &s.BookPhoto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}
