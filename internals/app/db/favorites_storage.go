package db

import (
	"database/sql"
	"errors"
	"fmt"

	"dclibrary.com/internals/app/models"
)

type FavorietesStorage struct {
	DB *sql.DB
}

func (m *FavorietesStorage) AddToFavorietes(book_id, user_id int64) (int, error) {

	stmt := `UPDATE dclibrary.favorietes SET book_id = ? WHERE user_id = ?`

	var favorietes *models.Favorietes

	result, err := m.DB.Exec(stmt, favorietes.FavoriteBooks, favorietes.UserId)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Printf("---> book %v has been added to Favorietes", id)

	return int(id), nil
}

func (m *FavorietesStorage) GetFavoritesList(user_id int64) ([]*models.Books, error) {

	stmt := `SELECT user_id, favirietes FROM dclibrary.favorietes WHERE user_id = ?`
	sbmb := `SELECT book_id, bookname, author_id, count, photo FROM dclibrary.books WHERE id = ?`
	sdmd := `SELECT author_name, author_photo FROM dclibrary.authors WHERE author_id = ?`

	var favorietes *models.Favorietes

	row := m.DB.QueryRow(stmt, user_id)
	err := row.Scan(&favorietes.UserId, &favorietes.FavoriteBooks)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	for _, value := range favorietes.FavoriteBooks {

	}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []*models.Books

	for rows.Next() {
		s := &models.Books{}
		err = rows.Scan(&s.BookId, &s.BookName, &s.Author.AuthorId, &s.Count, &s.BookPhoto)
		if err != nil {
			return nil, err
		}
		auth := m.DB.QueryRow(sdmd, s.Author.AuthorId)
		err = auth.Scan(&s.Author.AuthorName, &s.Author.AuthorPhoto)
		books = append(books, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
