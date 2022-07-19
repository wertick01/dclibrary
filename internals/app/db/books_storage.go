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

func (m *BooksStorage) CreateNewBook(book *models.Books) (*models.Books, error) {

	stmt := `INSERT INTO dclibrary.books (bookname, count, photo, stars) VALUES(?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, book.BookName, book.Count, book.BookPhoto, 0)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	connected, err := m.BooksAuthorsConnection(book)
	if err != nil {
		return nil, err
	}

	fmt.Printf("---> Book %v has been added to DB", id)

	return connected, nil
}

func (m *BooksStorage) BooksAuthorsConnection(books *models.Books) (*models.Books, error) {
	stmt := `INSERT INTO dclibrary.books_authors (book_id, author_id) VALUES(?, ?)`

	authors := books.Authors
	for _, val := range authors {
		row, err := m.DB.Exec(stmt, books.BookId, val.AuthorId)
		if err != nil {
			return nil, err
		}
		id, err := row.LastInsertId()
		if err != nil {
			return nil, err
		}
		fmt.Printf("Id %v has beed added to DB.", id)
	}
	return books, nil
}

func (m *BooksStorage) GetBooksList() ([]*models.Books, error) {

	stmt := `SELECT book_id, bookname, count, photo, stars FROM dclibrary.books`
	skmk := `SELECT author_id FROM dclibrary.books_authors WHERE book_id = ?`
	sdmd := `SELECT author_name, author_surname, author_patrynomic, author_photo, author_stars FROM dclibrary.authors WHERE author_id = ?`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var allbooks []*models.Books

	for rows.Next() {
		s := &models.Books{}
		err = rows.Scan(&s.BookId, &s.BookName, &s.Count, &s.BookPhoto, &s.Stars)
		if err != nil {
			return nil, err
		}

		connection, err := m.DB.Query(skmk, s.BookId)
		if err != nil {
			return nil, err
		}

		defer connection.Close()

		for connection.Next() {
			a := &models.Authors{}

			err = connection.Scan(&a.AuthorId)
			if err != nil {
				return nil, err
			}

			s.Authors = append(s.Authors, *a)
		}

		for _, val := range s.Authors {
			authors := m.DB.QueryRow(sdmd, val.AuthorId)

			err = authors.Scan(
				&val.AuthorName.Name,
				&val.AuthorName.Surname,
				&val.AuthorName.Patronymic,
				&val.AuthorPhoto,
				&val.AuthorStars,
			)

			if err != nil {
				return nil, err
			}
		}

		allbooks = append(allbooks, s)
	}

	return allbooks, nil
}

func (m *BooksStorage) GetBookById(id int64) (*models.Books, error) {

	stmt := `SELECT book_id, bookname, count, photo, stars FROM dclibrary.books WHERE book_id = ?`
	skmk := `SELECT author_id FROM dclibrary.books_authors WHERE book_id = ?`
	sdmd := `SELECT author_name, author_photo FROM dclibrary.authors WHERE author_id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Books{}

	err := row.Scan(&s.BookId, &s.BookName, &s.Count, &s.BookPhoto, &s.Stars)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	connection, err := m.DB.Query(skmk, s.BookId)
	if err != nil {
		return nil, err
	}

	defer connection.Close()

	for connection.Next() {
		a := &models.Authors{}

		err = connection.Scan(&a.AuthorId)
		if err != nil {
			return nil, err
		}

		s.Authors = append(s.Authors, *a)
	}

	for _, val := range s.Authors {
		authors := m.DB.QueryRow(sdmd, val.AuthorId)

		err = authors.Scan(
			&val.AuthorName.Name,
			&val.AuthorName.Surname,
			&val.AuthorName.Patronymic,
			&val.AuthorPhoto,
			&val.AuthorStars,
		)

		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (m *BooksStorage) ChangeBook(old *models.Books) (*models.Books, error) { //доделать с учётом нескольких авторов

	stmt := `UPDATE dclibrary.books SET book_name = ?, count = ?, photo = ? WHERE book_id = ?`
	sdmd := `DELETE FROM dclibrary.books_authors WHERE book_id = ?`

	change, err := m.DB.Exec(stmt, old.BookName, old.Count, old.BookPhoto)
	if err != nil {
		return nil, err
	}
	id, err := change.LastInsertId()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Book %v has been changed.", id)

	deleted, err := m.DB.Exec(sdmd, old.BookId)
	if err != nil {
		return nil, err
	}
	id, err = deleted.LastInsertId()
	if err != nil {
		return nil, err
	}

	connected, err := m.BooksAuthorsConnection(old)
	if err != nil {
		return nil, err
	}

	return connected, nil
}

func (m *BooksStorage) DeleteBookById(id int64) (int64, error) {
	stmt := `DELETE FROM dclibrary.books WHERE book_id = ?`
	sdmd := `DELETE FROM dclibrary.books_authors WHERE book_id = ?`

	deleted, err := m.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	res, err := deleted.LastInsertId()
	if err != nil {
		return 0, err
	}

	delet, err := m.DB.Exec(sdmd, id)
	if err != nil {
		return 0, err
	}

	result, err := delet.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Printf("Book %v has been deleted", result)
	return res, nil
}

func (m *BooksStorage) PutStarByBookId(id int64) error { //!!!

	stmt := `UPDATE dclibrary.books SET stars = ? WHERE book_id = ?`

	book, err := m.GetBookById(id)
	if err != nil {
		return err
	}

	book.Stars += 1

	putstar, err := m.DB.Exec(stmt, book.Stars, id)
	if err != nil {
		return err
	}

	id, err = putstar.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}
