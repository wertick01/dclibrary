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

func (m *UsersStorage) CreateNewBook(bookname, photo string, author_id, count int) (int, error) {

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

func (m *UsersStorage) GetUsersList() ([]*models.User, error) {

	stmt := `SELECT user_id, name, phone, mail, hash, role_id FROM dclibrary.users`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []*models.User

	for rows.Next() {
		s := &models.User{}
		err = rows.Scan(&s.UserId, &s.Name, &s.Phone, &s.Mail, &s.Role)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *UsersStorage) GetUserById(id int64) (*models.User, error) {

	stmt := `SELECT user_id, name, phone, mail, hash, role_id FROM dclibrary.users WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.User{}

	err := row.Scan(&s.UserId, &s.Name, &s.Phone, &s.Mail, &s.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}
