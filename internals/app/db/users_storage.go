package db

import (
	"errors"
	"fmt"

	"dclibrary.com/internals/app/models"

	//"github.com/georgysavva/scany/pgxscan"
	//"github.com/jackc/pgx/v4/pgxpool"
	//log "github.com/sirupsen/logrus"

	"database/sql"
)

type UsersStorage struct {
	DB *sql.DB
}

/*
func NewUsersStorage(pool *pgxpool.Pool) *UsersStorage {
	storage := new(UsersStorage)
	storage.databasePool = pool
	return storage

}
*/

func (m *UsersStorage) CreateNewUser(name, phone, mail, hash string, role_id int) (int, error) {

	stmt := `INSERT INTO dclibrary.users (name, phone, mail, hash, role_id) VALUES(?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, name, phone, mail, hash, role_id)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Printf("---> User %v has been added to DB", id)

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

func (m *UsersStorage) GetUserById(id int) (*models.User, error) {

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
