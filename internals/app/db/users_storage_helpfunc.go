package db

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (m *UsersStorage) CheckParams() []string {
	var params string
	var mass, res []string

	fmt.Println("Введите параметры:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	params = scanner.Text()
	mass = strings.Split(params, " ")

	if search("name", mass) {
		res = append(res, "name")
	} else if search("phone", mass) {
		res = append(res, "phone")
	} else if search("mail", mass) {
		res = append(res, "mail")
	} else if search("password", mass) {
		res = append(res, "password")
	}

	return res
}

func (m *BooksStorage) CheckParams() []string {
	var params string
	var mass, res []string

	fmt.Println("Введите параметры:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	params = scanner.Text()
	mass = strings.Split(params, " ")

	if search("book_name", mass) {
		res = append(res, "book_name")
	} else if search("author", mass) {
		res = append(res, "author")
	} else if search("count", mass) {
		res = append(res, "count")
	} else if search("book_photo", mass) {
		res = append(res, "book_photo")
	}

	return res
}

func search(key string, arr []string) bool {
	for _, value := range arr {
		if value == key {
			return true
		}
	}
	return false
}

/*
func (m *BooksStorage) CheckAuthors(authors string) {
	auth := strings.Split(authors, ", ")
}
*/
