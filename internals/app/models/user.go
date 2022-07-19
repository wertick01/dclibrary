package models

import "errors"

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type User struct {
	UserId     int64  `json:"user_id"`
	Name       string `json:"name"`
	Surname    string `json: "surname"`
	Patrynomic string `json: "patrynomic"`
	Phone      string `json: "phone"`
	Mail       string `json: "mail"`
	Hash       string `json: "hash"`
	Role       Role   `json:"Role"`
}

type Role struct {
	RoleId int    `json: "role_id"`
	Role   string `json: "role"`
}

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
