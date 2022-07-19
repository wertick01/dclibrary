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
func NewUsersStorage() *UsersStorage {
	storage := new(UsersStorage)
	storage.databasePool = pool
	return storage

}
*/

func (m *UsersStorage) CreateNewUser(user *models.User) (int, error) {

	stmt := `INSERT INTO dclibrary.users (name, surname, patrynomyc, phone, mail, hash, role_id) VALUES(?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(
		stmt,
		user.Name,
		user.Surname,
		user.Patrynomic,
		user.Phone,
		user.Mail,
		user.Hash,
		user.Role.RoleId,
	)

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

	stmt := `SELECT user_id, name, surname, patrynomic, phone, mail, role_id FROM dclibrary.users`
	sdmd := `SELECT role FROM dclibrary.roles WHERE role_id = ?`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*models.User
	//var rol m

	for rows.Next() {
		s := &models.User{}
		err = rows.Scan(&s.UserId, &s.Name, &s.Surname, &s.Patrynomic, &s.Phone, &s.Mail, &s.Role.RoleId)
		if err != nil {
			return nil, err
		}
		rol := m.DB.QueryRow(sdmd, s.Role.RoleId)
		err = rol.Scan(&s.Role.Role)
		users = append(users, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UsersStorage) GetUserById(id int64) (*models.User, error) {

	stmt := `SELECT user_id, name, surname, patrynomic, phone, mail, hash, role_id FROM dclibrary.users WHERE user_id = ?`
	sdmd := `SELECT role_id, role FROM dclibrary.roles WHERE role_id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.User{}

	err := row.Scan(&s.UserId, &s.Name, &s.Surname, &s.Patrynomic, &s.Phone, &s.Mail, &s.Role.RoleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	rol := m.DB.QueryRow(sdmd, s.Role.RoleId)
	err = rol.Scan(&s.Role.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *UsersStorage) ChangeUserById(old *models.User) (*models.User, error) {

	stmt := `UPDATE dclibrary.users SET name = ?, surname = ?, patrynomic = ?, phone = ?, mail = ?, hash = ?, role = ? WHERE user_id = ?`

	change, err := m.DB.Exec(stmt, old.Name, old.Surname, old.Patrynomic, old.Phone, old.Mail, old.Hash, old.Role.RoleId)
	if err != nil {
		return nil, err
	}
	id, err := change.LastInsertId()
	if err != nil {
		return nil, err
	}
	fmt.Printf("User %v has been changed.", id)

	return old, nil
}

func (m *UsersStorage) DeleteUserById(id int64) (int, error) {
	stmt := `DELETE FROM dclibrary.users WHERE user_id = ?`
	deleted, err := m.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	res, err := deleted.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(res), nil
}
