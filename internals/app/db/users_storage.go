package db

import (
	"bufio"
	"errors"
	"fmt"
	"os"

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

	stmt := `SELECT user_id, name, phone, mail, role_id FROM dclibrary.users`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []*models.User

	for rows.Next() {
		s := &models.User{}
		err = rows.Scan(&s.UserId, &s.Name, &s.Phone, &s.Mail, &s.RoleId)
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

	err := row.Scan(&s.UserId, &s.Name, &s.Phone, &s.Mail, &s.RoleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *UsersStorage) ChangeUserById(id int64, old *models.User) (*models.User, error) {

	params := m.CheckParams()
	if search("name", params) {
		fmt.Println("Введите новое имя: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		old.Name = scanner.Text()
	} else if search("phone", params) {
		fmt.Println("Введите новый номер телефона: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		old.Phone = scanner.Text()
	} else if search("mail", params) {
		fmt.Println("Введите новый почтовый ящик: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		old.Mail = scanner.Text()
	} else if search("password", params) {
		fmt.Println("Введите новый ПАРОЛЬ: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		old.Hash = scanner.Text()
	}

	stmt := `UPDATE dclibrary.users SET name = ?, phone = ?, mail = ?, hash = ? WHERE user_id = ?`
	change, err := m.DB.Exec(stmt, old.Name, old.Phone, old.Mail, old.Hash, id)
	if err != nil {
		return old, err
	}
	id, err = change.LastInsertId()
	if err != nil {
		return old, err
	}
	fmt.Printf("---> User %v has been added to DB", id)

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
