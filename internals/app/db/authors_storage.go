package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"dclibrary.com/internals/app/models"
)

type AuthorsStorage struct {
	DB *sql.DB
}

type GetBooksList struct {
	gbl *BooksStorage
}

func (m *AuthorsStorage) CreateNewAuthor(author_name, author_surname, author_patrynomic, author_photo string) (int, error) {

	stmt := `INSERT INTO dclibrary.authors (author_name, author_surname, author_patrynomic, author_photo) VALUES(?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, author_name, author_surname, author_patrynomic, author_photo)
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

func (m *AuthorsStorage) GetBooksByAuthorId(id int64) ([]*models.FinalBooks, error) {

	stmt := `SELECT book_id, bookname, author_id, count, photo, stars FROM dclibrary.books`
	sdmd := `SELECT author_id, author_name, author_photo, author_stars FROM dclibrary.authors WHERE author_id = ?`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ath bool
	var allbooks []*models.FinalBooks

	for rows.Next() {
		ath = false
		s := &models.Books{}
		b := &models.FinalBooks{}
		err = rows.Scan(&s.BookId, &s.BookName, &s.Authors, &s.Count, &s.BookPhoto, &s.Stars)
		if err != nil {
			return nil, err
		}

		authors := strings.Split(s.Authors, ", ")
		for _, value := range authors {
			a := &models.Authors{}
			val, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal(err)
			}

			auth := m.DB.QueryRow(sdmd, val)
			err = auth.Scan(&a.AuthorId, &a.AuthorName, &a.AuthorPhoto, &a.AuthorStars)
			if a.AuthorId == int(id) {
				ath = true
			}
			b.Authors = append(b.Authors, *a)
		}

		b.Book = *s
		if ath == true {
			allbooks = append(allbooks, b)
		}
	}
	return allbooks, nil
}
