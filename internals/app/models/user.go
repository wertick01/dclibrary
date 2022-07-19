package models

import "errors"

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type User struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
	Phone  string `json: "phone"`
	Mail   string `json: "mail"`
	Hash   string `json: "hash"`
	RoleId int    `json:"Role"`
}

type Role struct {
	RoleId int    `json: "role_id"`
	Role   string `json: "role"`
}

type FinalUser struct {
	User User `json: "user"`
	Role Role `json: "user_role"`
}

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
