package db

import (
	"database/sql"
	"errors"
	"fmt"

	"dclibrary.com/internals/app/models"
)

type AuthorsStorage struct {
	BooksStorage
}

func NewAuthorsStorage(db *sql.DB) *AuthorsStorage {
	storage := new(AuthorsStorage)
	storage.DB = db
	return storage

}

func (m *AuthorsStorage) CreateNewAuthor(author *models.Authors) (int, error) {

	stmt := `INSERT INTO dclibrary.authors (author_name, author_surname, author_patrynomic, author_photo, author_stars) VALUES(?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, author.AuthorName.Name, author.AuthorName.Surname, author.AuthorName.Patronymic, author.AuthorPhoto, 0)
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

	stmt := `SELECT author_id, author_name, author_surname, author_patrynomic, author_photo, author_stars FROM dclibrary.authors`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var authors []*models.Authors

	for rows.Next() {
		s := &models.Authors{}
		err = rows.Scan(&s.AuthorId, &s.AuthorName, &s.AuthorPhoto, &s.AuthorStars)
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

	stmt := `SELECT author_id, author_name, author_photo, author_stars FROM dclibrary.authors WHERE author_id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Authors{}

	err := row.Scan(&s.AuthorId, &s.AuthorName, &s.AuthorPhoto, &s.AuthorStars)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *AuthorsStorage) GetBooksByAuthorId(id int64) ([]*models.Books, *models.Authors, error) {
	stmt := `SELECT book_id FROM dclibrary.books_authors WHERE author_id = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, nil, err
	}
	var book_id int64
	var books []*models.Books

	for rows.Next() {
		err = rows.Scan(&book_id)
		book, err := m.GetBookById(book_id)
		if err != nil {
			return nil, nil, err
		}
		books = append(books, book)
	}
	author, err := m.GetAuthorById(id)
	if err != nil {
		return books, nil, err
	}

	return books, author, nil
}

func (m *AuthorsStorage) PutStarByAuthorId(id int64) error { //!!!

	stmt := `UPDATE dclibrary.authors SET author_stars = ? WHERE author_id = ?`

	author, err := m.GetAuthorById(id)
	if err != nil {
		return err
	}

	author.AuthorStars += 1

	putstar, err := m.DB.Exec(stmt, author.AuthorStars, id)
	if err != nil {
		return err
	}

	id, err = putstar.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (m *AuthorsStorage) ChangeAuthor(old *models.Authors) (*models.Authors, error) {

	stmt := `UPDATE dclibrary.authors SET author_name = ?, author_surname = ?, author_patrynomic = ?, author_photo = ?, WHERE author_id = ?`
	sdmd := `DELETE FROM dclibrary.books_authors WHERE author_id = ?`

	change, err := m.DB.Exec(stmt, old.AuthorName.Name, old.AuthorName.Surname, old.AuthorName.Patronymic, old.AuthorPhoto)
	if err != nil {
		return nil, err
	}
	id, err := change.LastInsertId()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Author %v has been changed.", id)

	deleted, err := m.DB.Exec(sdmd, old.AuthorId)
	if err != nil {
		return nil, err
	}
	id, err = deleted.LastInsertId()
	if err != nil {
		return nil, err
	}

	connected, err := m.AuthorsBooksConnection(old)
	if err != nil {
		return nil, err
	}

	return connected, nil
}

func (m *AuthorsStorage) AuthorsBooksConnection(author *models.Authors) (*models.Authors, error) {
	stmt := `INSERT INTO dclibrary.books_authors (book_id, author_id) VALUES(?, ?)`

	books, author, err := m.GetBooksByAuthorId(int64(author.AuthorId))
	if err != nil {
		return nil, err
	}
	for _, val := range books {
		row, err := m.DB.Exec(stmt, val.BookId, author.AuthorId)
		if err != nil {
			return nil, err
		}
		id, err := row.LastInsertId()
		if err != nil {
			return nil, err
		}
		fmt.Printf("Id %v has beed added to DB.", id)
	}
	return author, nil
}

func (m *AuthorsStorage) DeleteAuthorById(id int64) (int64, error) {
	stmt := `DELETE FROM dclibrary.authors WHERE author_id = ?`

	deleted, err := m.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	res, err := deleted.LastInsertId()
	if err != nil {
		return 0, err
	}
	return res, nil
}